package limitlessled

import (
	"fmt"
	"net"
	"time"
)

var (
	ALL_OFF = []byte{0x21, 0x00, 0x55}
	ALL_ON  = []byte{0x22, 0x00, 0x55}

	BRIGHT_HIGH = []byte{0x4E, 0x1B, 0x55}
	BRIGHT_LOW  = []byte{0x4E, 0x02, 0x55}
	BRIGHT_HALF = []byte{0x4E, 0xD0, 0x55}

	DISCO_ON     = []byte{0x4D, 0x00, 0x55}
	DISCO_FASTER = []byte{0x44, 0x00, 0x55}
	DISCO_SLOWER = []byte{0x43, 0x00, 0x55}
	DISCO_MODE   = []byte{0x4D, 0x00, 0x55}

	GREEN         = []byte{0x20, 0x60, 0x55}
	BLUE          = []byte{0x20, 0x10, 0x55}
	YELLOW        = []byte{0x20, 0x80, 0x55}
	ORANGE        = []byte{0x20, 0xC0, 0x55}
	PINK          = []byte{0x20, 0xB9, 0x55}
	RED           = []byte{0x20, 0xB0, 0x55}
	VIOLET        = []byte{0x20, 0xA0, 0x55}
	BABY_BLUE     = []byte{0x20, 0x20, 0x55}
	AQUA          = []byte{0x20, 0x30, 0x55}
	MINT          = []byte{0x20, 0x40, 0x55}
	SEAFORM_GREEN = []byte{0x20, 0x50, 0x55}
	YELLOW_ORANGE = []byte{0x20, 0x90, 0x55}
	FUSIA         = []byte{0x20, 0xD0, 0x55}
	LILAC         = []byte{0x20, 0xE0, 0x55}
	LAVENDAR      = []byte{0x20, 0xF0, 0x55}
)

type LimitlessLed struct {
}

func (l *LimitlessLed) Command(command []byte, udp string) {
	conn, err := l.connect(udp)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write(command)
}

func (l *LimitlessLed) TransitionColor(colorFrom []byte, colorTo []byte, udp string) {
	conn, err := l.connect(udp)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write(colorFrom)
	newColor := colorFrom
	for {
		time.Sleep(200 * time.Millisecond)
		newColor = []byte{0x20, newColor[1] + 0x01, 0x55}
		fmt.Printf("%v\n", newColor[1])
		conn.Write(newColor)
		if newColor[1] == colorTo[1] {
			return
		}
	}
}

func (l *LimitlessLed) connect(udp string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp4", udp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	conn, err2 := net.DialUDP("udp4", nil, addr)
	if err2 != nil {
		fmt.Println(err2)
		return nil, err
	}
	return conn, nil
}
