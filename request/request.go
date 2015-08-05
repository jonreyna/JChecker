package request

import (
	"github.com/Juniper/go-netconf/netconf"
)

type Request interface {
	Run(int, chan<- *netconf.RPCReply)
}
