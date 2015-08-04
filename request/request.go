package request

import (
	"github.com/Juniper/go-netconf/netconf"
)

type Request interface {
	Method() netconf.RawMethod
}