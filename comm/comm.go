package comm

import (
	"errors"
	"fmt"
	serial "github.com/facchinm/go-serial" // bugst is the original but facchinm has merged tarm's windows code in...
	"io"
	"log"
)

type CommPort struct {
	Path        string
	VendorId    uint16
	ProductId   uint16
	DisplayName string
	connected   bool
	port        *serial.SerialPort
}

func (cp *CommPort) Connect(baudRate int) (io.ReadWriter, error) {
	mode := &serial.Mode{
		BaudRate: baudRate,
		Parity:   serial.PARITY_NONE,
		DataBits: 8,
		StopBits: serial.STOPBITS_ONE,
	}

	if cp.connected {
		return nil, errors.New("Already connected")
	}

	var err error
	cp.port, err = serial.OpenPort(cp.Path, mode)
	if err != nil {
		log.Println(err)
		cp.port = nil
		return nil, err
	}

	cp.connected = true

	return cp.port, nil
}

func (cp *CommPort) Disconnect() {

	if !cp.connected {
		return
	}
	// cut handlers
	cp.connected = false
	cp.port.Close()
	cp.port = nil
}

func (cp *CommPort) Write(p []byte) (n int, err error) {
	return cp.port.Write(p)
}

func (cp *CommPort) IsConnected() bool {
	return cp.connected
}

func FindFirstMatch(cpl []CommPort, vid uint16, pid uint16) (*CommPort, error) {
	for _, port := range cpl {
		if port.VendorId == vid && port.ProductId == pid {
			fmt.Printf("Serial port with VID/PID = %x/%x found in list!\n", vid, pid)
			return &port, nil
		}
	}

	return nil, fmt.Errorf("Serial port with VID/PID = %x/%x not found in list.", vid, pid)
}

func FindAllMatches(cpl []CommPort, vid uint16, pid uint16) ([]CommPort, error) {
	var result []CommPort
	for _, port := range cpl {
		if port.VendorId == vid && port.ProductId == pid {
			//fmt.Printf("Serial port with VID/PID = %x/%x found in list!\n", vid, pid)
			result = append(result, port)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("Serial port with VID/PID = %x/%x not found in list.", vid, pid)
	}

	return result, nil
}
