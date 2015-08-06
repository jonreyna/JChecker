package request

import (
	"net"
	"time"

	"github.com/Juniper/go-netconf/netconf"
	log "github.com/Sirupsen/logrus"
)

const (
	CHASSIS_ENV_REQUEST_TMPL = "<get-environment-information></get-environment-information>"
)

type ChassisEnvRequest struct {
	LoopBackIP net.IP
	Interval   time.Duration
	UserName   string
	Password   string
}

func (*ChassisEnvRequest) Method() netconf.RawMethod {
	return netconf.RawMethod(CHASSIS_ENV_REQUEST_TMPL)
}

func NewChassisEnvRequest(loopBackIP net.IP, interval time.Duration,
	userName, password string) *ChassisEnvRequest {
	return &ChassisEnvRequest{
		LoopBackIP: loopBackIP,
		UserName:   userName,
		Password:   password,
		Interval:   interval,
	}
}

func (cer *ChassisEnvRequest) Run(count int, replyChan chan<- *netconf.RPCReply) {

	go func() {
		ticker := time.NewTicker(cer.Interval)

		for {
			select {
			case <-ticker.C:
				// Get SSH/NETCONF session going for use throughout the life of this request.
				if s, err := netconf.DialSSH(cer.LoopBackIP.String(),
					netconf.SSHConfigPassword(cer.UserName, cer.Password)); err != nil {
					log.Errorln(err)
					s.Close()
					close(replyChan)
				} else if ncReply, err := s.Exec(cer.Method()); err != nil {
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
