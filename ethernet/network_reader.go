package ethernet

import (
	"github.com/hsheth2/logs"
)

// A FrameHeader is used when returning data from Reader.Read() calls
type FrameHeader struct {
	//Rmac, Lmac *MAC_Address
	Packet []byte
}

type ethernetReader struct {
	ethertype EtherType
	input     chan []byte
	processed chan *FrameHeader
}

func newEthernetReader(etp EtherType) (*ethernetReader, error) {
	ethr := &ethernetReader{
		ethertype: etp,
		input:     make(chan []byte, ethProtocolBufferSize),
		processed: make(chan *FrameHeader, ethProtocolBufferSize),
	}

	go ethr.readAll()

	return ethr, nil
}

func (ethr *ethernetReader) readAll() {
	for {
		//logs.Trace.Println("Ethernet reader attempting to get work")
		data := <-ethr.input
		//logs.Trace.Println("Ethernet reader received packet")

		ethHead := &FrameHeader{
			//Rmac: &MAC_Address{Data: data[ETH_MAC_ADDR_SZ : 2*ETH_MAC_ADDR_SZ]},
			//Lmac:   &MAC_Address{Data: data[0:ETH_MAC_ADDR_SZ]},
			Packet: data[ethHeaderSize:],
		}
		//			/*logs*/logs.Info.Println("Beginning to forward ethernet packet")
		select {
		case ethr.processed <- ethHead:
			//logs.Trace.Println("Forwarding ethernet packet")
		default:
			logs.Warn.Println("Dropping Ethernet packet: buffer full")
		}
	}
}

// blocking read call
func (ethr *ethernetReader) Read() (*FrameHeader, error) {
	return <-ethr.processed, nil
}

func (ethr *ethernetReader) Close() error {
	err := Unbind(ethr.ethertype)
	if err != nil {
		return err
	}
	close(ethr.input)
	close(ethr.processed)
	return nil
}
