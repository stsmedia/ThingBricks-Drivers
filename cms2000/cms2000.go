package cms2000

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	//	"time"
	"errors"
	"math/big"
	"strings"
)

type Cms2000 struct {
	Conn *net.TCPConn
}

type Reading struct {
	Temp   float64
	EToday float64
	ENow   float64
	ETotal float64
	Vdc    float64
	Iac    float64
	Vac    float64
	Fac    float64
	Pac    float64
	Zac    float64
	HTotal float64
}

func (r *Reading) ToString() string {
	return fmt.Sprintf("temp: %v, etoday: %v, enow: %v, etotal: %v, vdc: %v, iac: %v, vac: %v, fac: %v, fac: %v, pac: %v, zac: %v, htotal: %v",
		r.Temp, r.EToday, r.ENow, r.ETotal, r.Vdc, r.Iac, r.Vac, r.Fac, r.Pac, r.Zac, r.HTotal)
}

func (c *Cms2000) connect(tcp string) {
	addr, err := net.ResolveTCPAddr("tcp4", tcp)
	if err != nil {
		fmt.Println(err)
	}
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		fmt.Println(err)
	}
	c.Conn = conn
}

func (c *Cms2000) GetReading(tcp string) (Reading, error) {
	if c.Conn == nil {
		c.connect(tcp)
	}
	var r Reading

	// reset the serial port
	c.Conn.Write([]byte{0xAA, 0xAA, 0x01, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x01, 0x59})
	c.Conn.Write([]byte{0xAA, 0xAA, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x55})

	ch := make(chan []byte)
	eCh := make(chan error)

	// Start a goroutine to read from our net connection
	go func(ch chan []byte, eCh chan error) {
		for {
			// try to read the data
			data := make([]byte, 512)
			_, err := c.Conn.Read(data)
			if err != nil {
				// send an error if it's encountered
				eCh <- err
				return
			}
			// send data if we read some.
			ch <- data
		}
	}(ch, eCh)

	// continuously read from the connection
	for {
		select {
		// This case means we recieved data on the connection
		// CMS 2000 protocol:
		case data := <-ch:
			msg := hex.Dump(data)
			//			fmt.Println(msg)
			if strings.Contains(msg, "aa aa 00 00 01 00 00 80") {
				payload := append(append([]byte{0xAA, 0xAA, 0x01, 0x00, 0x00, 0x00, 0x00, 0x01, 0x0B}, data[9:19]...), byte(0x01))
				res := append(payload, c.checksum(payload)...)
				c.Conn.Write(res)
			}
			if strings.Contains(msg, "aa aa 00 01 01 00 00 81") {
				c.Conn.Write([]byte{0xAA, 0xAA, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x57})
			}
			if strings.Contains(msg, "aa aa 00 01 01 00 01 80") {
				c.Conn.Write([]byte{0xAA, 0xAA, 0x01, 0x00, 0x00, 0x01, 0x01, 0x02, 0x00, 0x01, 0x59})
			}
			if strings.Contains(msg, "aa aa 00 01 01 00 01 82") {
				payload := data[9 : len(data)-2]
				c.Conn.Close()
				c.Conn = nil
				return c.parseData(payload), nil
			}
		case err := <-eCh:
			fmt.Println(err)
			return r, err
		}
	}
	return r, errors.New("foo")
}

func (c *Cms2000) parseData(payload []byte) Reading {
	var reading Reading
	reading.Temp = c.convert(payload[0:2], 10)
	reading.EToday = c.convert(payload[2:4], 100)
	reading.ENow = c.convert(payload[12:14], 1)
	reading.ETotal = c.convert(payload[18:20], 10)
	reading.Vdc = c.convert(payload[4:6], 10)
	reading.Iac = c.convert(payload[6:8], 10)
	reading.Vac = c.convert(payload[8:10], 10)
	reading.Fac = c.convert(payload[10:12], 100)
	reading.Pac = c.convert(payload[12:14], 1)
	reading.Zac = c.convert(payload[14:16], 1)
	reading.HTotal = c.convert(payload[20:22], 1)
	return reading
}

func (c *Cms2000) convert(in []byte, divide float64) float64 {
	i := new(big.Int)
	num := i.SetBytes(in)
	return float64(num.Int64()) / divide
}

func (c *Cms2000) checksum(in []byte) []byte {
	sum := 0
	for _, byte := range in {
		sum += int(byte)
	}
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(sum))
	return b
}
