package main

import (
	"config"
	"fmt"
	"cms2000"
	"currentcost"
	"limitlessled"
)

func main() {
	fmt.Println("listening")

	msgReader, err := currentcost.NewSerialMessageReader(config.GetString("currentcost.device"), config.GetInt("currentcost.baudrate"))
	if err != nil {
		fmt.Println(err)
		return
	}
	var led limitlessled.LimitlessLed
	var cms2000 cms2000.Cms2000
	ledUdp := config.GetString("limitlessled.udp")
	cms2000Tcp := config.GetString("cms2000.tcp")
	for {
		if msg, err := msgReader.ReadMessage(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(msg.ToString())
			reading, err := cms2000.GetReading(cms2000Tcp)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(reading.ToString())
			if reading.ETotal == 0 {
				led.Command(limitlessled.YELLOW, ledUdp)
			} else if float64(msg.Channel1.Watts) > reading.ENow {
				led.Command(limitlessled.RED, ledUdp)
			} else {
				led.Command(limitlessled.GREEN, ledUdp)
			}
		}
	}

}
