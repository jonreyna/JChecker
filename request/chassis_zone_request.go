package request

import (
	"github.com/Juniper/go-netconf/netconf"
	"time"
)

const (
	CHASSIS_ZONES_REQUEST_TMPL = "get-chassis-zones-information></get-chassis-zones-information>"
)

type ChassisZonesRequest struct {
	LoopBackIP string
	Interval   time.Duration
}

func (*ChassisZonesRequest) Method() netconf.RawMethod {
	return netconf.RawMethod(CHASSIS_ZONES_REQUEST_TMPL)
}

func NewChassisZonesRequest(loopBackIP string, interval time.Duration) *ChassisZonesRequest {
	return &ChassisZonesRequest{LoopBackIP:loopBackIP}
}