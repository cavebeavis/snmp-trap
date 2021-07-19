package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/gosnmp/gosnmp"
)

type Response struct {
	SnmpPacket *gosnmp.SnmpPacket
	UdpAddress *net.UDPAddr
	Timestamp  time.Time
}

func main() {
	// Allowing tcpdump to start first.
	time.Sleep(time.Millisecond * 1000)

	logger := log.New(os.Stdout, "", 0)
	logger.Println("Starting trap listener...")

	//address := "0.0.0.0:1162"
	addresses, err := net.LookupIP("listener")
	if err != nil {
		logger.Fatal("looking up listener", err)
	}
	if len(addresses) < 1 {
		logger.Fatal("net.LookupIP('listener') returned empty addresses")
	}

	listenerAddress := addresses[0]

	/*addresses, err = net.LookupIP("sender")
	if err != nil {
		logger.Fatal("looking up sender", err)
	}
	if len(addresses) < 1 {
		logger.Fatal("net.LookupIP('sender') returned empty addresses")
	}

	senderAddress := addresses[0]*/

	listener := gosnmp.NewTrapListener()
	defer listener.Close()

	trapCh := make(chan Response, 1)
	listener.OnNewTrap = func(s *gosnmp.SnmpPacket, u *net.UDPAddr) {
		logger.Println(*u, *s)
		//if u.IP.Equal(senderAddress) {
			trapCh<- Response{
				SnmpPacket: s,
				UdpAddress: u,
				Timestamp: time.Now(),
			}
		//}
	}

	// listener goroutine
	errch := make(chan error)
	go func() {
		// defer close(errch)
		err := listener.Listen(listenerAddress.String() + ":1162")
		if err != nil {
			errch <- err
		}
	}()

	select {
	case <-listener.Listening():

	case err := <-errch:
		logger.Fatal(err)
	}

	vars := []gosnmp.SnmpPDU{
		{Name: ".1.3.6.1.2.1.1.1.0", Type: gosnmp.Integer, Value: 1},
		{Name: ".1.3.6.1.2.1.1.2.0", Type: gosnmp.OctetString, Value: "TRAPTEST1234"},
	}
	
	select {
	case <-time.After(time.Second * 10):
		logger.Fatal("ran out of time")
	case t := <-trapCh:
		logger.Println("trap received ", t.Timestamp, *t.SnmpPacket, *t.UdpAddress)

		for _, v1 := range vars {
			var found bool
			for _, v2 := range t.SnmpPacket.Variables {
				if v1.Type.String() == v2.Type.String() &&
				   v1.Name == v2.Name {

						v2Val8, ok := v2.Value.([]uint8)
						if v1.Value == v2.Value || (ok &&  v1.Value == string(v2Val8)) {
							found = true
							continue
						}
				}

				logger.Printf("not found v1: %#v\tv2:%#v", v1, v2)
			}

			if !found {
				logger.Fatalf("want: %#v, got: %#v", v1, t.SnmpPacket.Variables)
			}
		}
	}

	logger.Println("Trap OK")
}