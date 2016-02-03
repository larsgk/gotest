package comm

import (
	serial "github.com/facchinm/go-serial"
)

type CommPort struct {
	Name string
}

func GetSerialPortList() ([]CommPort, error) {
	serPorts, err := serial.GetPortsList()

	commPorts := []CommPort{}
	for _, item := range serPorts {
		commPorts = append(commPorts, CommPort{Name: item})
	}

	return commPorts, err
}
