package request

import (
	"github.com/Juniper/go-netconf/netconf"
)

const (
	CHASSIS_ENV_REQUEST_TMPL = "<get-environment-information></get-environment-information>"
)

type ChassisEnvRequest struct {
	LoopBackIP string
}

func (*ChassisEnvRequest) Method() netconf.RawMethod {
	return netconf.RawMethod(CHASSIS_ENV_REQUEST_TMPL)
}