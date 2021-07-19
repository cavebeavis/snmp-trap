package main

import (
	"context"
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
	logger := log.New(os.Stdout, "", 0)
	logger.Println("Starting trap crap...")

	goSnmpCFG := &GoSnmpCFG{
		Target:             "127.0.0.1",
		Port:               1162,
		Transport:          "udp",
		Community:          "public",
		Version:            gosnmp.Version2c,
		Timeout:            time.Second * 3,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            20,
		MaxRepetitions:     100,
	}

	g := NewGoSNMP(context.Background(), goSnmpCFG, logger)
	err := g.Connect()
	if err != nil {
		logger.Fatalf("Connect() err: %v", err)
	}
	defer g.Conn.Close()

	vars := []gosnmp.SnmpPDU{
		{Name: "1.3.6.1.2.1.1.1.0", Type: gosnmp.Integer, Value: 1},
		{Name: "1.3.6.1.2.1.1.2.0", Type: gosnmp.OctetString, Value: "TRAPTEST1234"},
	}
	
	trap := gosnmp.SnmpTrap{
		IsInform: false,
		Variables: vars,
	}

	result, err := g.SendTrap(trap)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Print(*result)

	logger.Println("Finished")
}

type GoSnmpCFG struct {
	Target             string
	Port               uint16
	Transport          string
	Community          string
	Version            gosnmp.SnmpVersion
	Timeout            time.Duration
	Retries            int
	ExponentialTimeout bool
	MaxOids            int
	MaxRepetitions     uint32
}

func NewGoSNMP(ctx context.Context, cfg *GoSnmpCFG, logger gosnmp.LoggerInterface) *gosnmp.GoSNMP {
	return &gosnmp.GoSNMP{
		// Conn is net connection to use, typically established using GoSNMP.Connect().
		// Conn net.Conn

		// Target is an ipv4 address.
		// Target string
		Target: cfg.Target,

		// Port is a port.
		// Port uint16
		Port: cfg.Port,

		// Transport is the transport protocol to use ("udp" or "tcp"); if unset "udp" will be used.
		// Transport string
		Transport: cfg.Transport,

		// Community is an SNMP Community string.
		// Community string
		Community: cfg.Community,

		// Version is an SNMP Version.
		// Version SnmpVersion
		Version: cfg.Version,

		// Context allows for overall deadlines and cancellation.
		// Context context.Context
		Context: ctx,

		// Timeout is the timeout for one SNMP request/response.
		// Timeout time.Duration
		Timeout: cfg.Timeout,

		// Set the number of retries to attempt.
		// Retries int
		Retries: cfg.Retries,

		// Double timeout in each retry.
		// ExponentialTimeout bool
		ExponentialTimeout: cfg.ExponentialTimeout,

		// Logger is the GoSNMP.Logger to use for debugging.
		// For verbose logging to stdout:
		// x.Logger = NewLogger(log.New(os.Stdout, "", 0))
		// For Release builds, you can turn off logging entirely by using the go build tag "gosnmp_nodebug" even if the logger was installed.
		// Logger Logger
		Logger: gosnmp.NewLogger(logger),

		// Message hook methods allow passing in a functions at various points in the packet handling.
		// For example, this can be used to collect packet timing, add metrics, or implement tracing.
		/*

		*/
		// PreSend is called before a packet is sent.
		// PreSend func(*GoSNMP)
		PreSend: func(gs *gosnmp.GoSNMP) {
			logger.Print("gosnmp preparing packet")
		},

		// OnSent is called when a packet is sent.
		// OnSent func(*GoSNMP)
		OnSent: func(gs *gosnmp.GoSNMP) {
			logger.Print("gosnmp packet sent")
		},

		// OnRecv is called when a packet is received.
		// OnRecv func(*GoSNMP)
		OnRecv: func(gs *gosnmp.GoSNMP) {
			logger.Print("gosnmp packet received")
		},

		// OnRetry is called when a retry attempt is done.
		// OnRetry func(*GoSNMP)
		OnRetry: func(gs *gosnmp.GoSNMP) {
			logger.Print("gosnmp packet retried")
		},

		// OnFinish is called when the request completed.
		// OnFinish func(*GoSNMP)
		OnFinish: func(gs *gosnmp.GoSNMP) {
			logger.Print("gosnmp request complete")
		},

		// MaxOids is the maximum number of oids allowed in a Get().
		// (default: MaxOids)
		// MaxOids int
		MaxOids: cfg.MaxOids,

		// MaxRepetitions sets the GETBULK max-repetitions used by BulkWalk*
		// Unless MaxRepetitions is specified it will use defaultMaxRepetitions (50)
		// This may cause issues with some devices, if so set MaxRepetitions lower.
		// See comments in https://github.com/gosnmp/gosnmp/issues/100
		// MaxRepetitions uint32
		MaxRepetitions: cfg.MaxRepetitions,

		// NonRepeaters sets the GETBULK max-repeaters used by BulkWalk*.
		// (default: 0 as per RFC 1905)
		// NonRepeaters int

		// UseUnconnectedUDPSocket if set, changes net.Conn to be unconnected UDP socket.
		// Some multi-homed network gear isn't smart enough to send SNMP responses
		// from the address it received the requests on. To work around that,
		// we open unconnected UDP socket and use sendto/recvfrom.
		// UseUnconnectedUDPSocket bool

		// netsnmp has '-C APPOPTS - set various application specific behaviours'
		//
		// - 'c: do not check returned OIDs are increasing' - use AppOpts = map[string]interface{"c":true} with
		//   Walk() or BulkWalk(). The library user needs to implement their own policy for terminating walks.
		// - 'p,i,I,t,E' -> pull requests welcome
		// AppOpts map[string]interface{}

		// MsgFlags is an SNMPV3 MsgFlags.
		// MsgFlags SnmpV3MsgFlags

		// SecurityModel is an SNMPV3 Security Model.
		// SecurityModel SnmpV3SecurityModel

		// SecurityParameters is an SNMPV3 Security Model parameters struct.
		// SecurityParameters SnmpV3SecurityParameters

		// ContextEngineID is SNMPV3 ContextEngineID in ScopedPDU.
		// ContextEngineID string

		// ContextName is SNMPV3 ContextName in ScopedPDU
		// ContextName string
	}
}