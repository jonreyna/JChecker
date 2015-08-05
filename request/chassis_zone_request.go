package request

import (
	"time"

	"github.com/Juniper/go-netconf/netconf"
	log "github.com/Sirupsen/logrus"
	"net"
)

const (
	CHASSIS_ZONES_REQUEST_TMPL = "get-chassis-zones-information></get-chassis-zones-information>"
)

type ChassisZonesRequest struct {
	LoopBackIP net.IP
	UserName   string
	Password   string
	Interval   time.Duration
}

func (*ChassisZonesRequest) Method() netconf.RawMethod {
	return netconf.RawMethod(CHASSIS_ZONES_REQUEST_TMPL)
}

func NewChassisZonesRequest(loopBackIP net.IP, interval time.Duration,
	userName, password string) *ChassisZonesRequest {

	return &ChassisZonesRequest{
		LoopBackIP: loopBackIP,
		UserName:   userName,
		Password:   password,
		Interval:   interval,
	}
}

func (zr *ChassisZonesRequest) Run(count int, replyChan chan<- *netconf.RPCReply) {

	go func() {
		ticker := time.NewTicker(zr.Interval)

		for {
			select {
			case <-ticker.C:
				// Get SSH/NETCONF session going for use throughout the life of this request.
				if s, err := netconf.DialSSH(zr.LoopBackIP.String(),
					netconf.SSHConfigPassword(zr.UserName, zr.Password)); err != nil {
					log.Errorln(err)
					s.Close()
					close(replyChan)
				} else if ncReply, err := s.Exec(zr.Method()); err != nil {
					log.Errorln(err)
					s.Close()
					close(replyChan)
				} else {
					s.Close()
					replyChan <- ncReply
				}

				count--
				if count < 0 {
					break
				}
			}
		}
	}()

	/*		cfg := config.GetConfig()

			file, err := os.Create(cfg.ChassisZonesResultsFile)

			if err != nil {
				log.Fatalln(err)
			}

			for _ = range ticker.C {
				if ncReply, err := s.Exec(zr.Method()); err != nil {
					log.Error(err)
				} else if chassisEnvResp, err := response.NewChassisEnvResponse(&ncReply.Data); err != nil {
					log.Error(err)
				} else {
					chassisEnvResp.WriteCSV(file)
				}
			}*/
}
