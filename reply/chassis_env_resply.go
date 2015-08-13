package reply

import (
	"encoding/csv"
	"encoding/xml"
	"io"
	"time"

	"github.com/Juniper/go-netconf/netconf"
	log "github.com/Sirupsen/logrus"
)

type EnvItem struct {
	Name        string `xml:"name,omitempty"        json:"name,omitempty"`
	Class       string `xml:"class,omitempty"       json:"class,omitempty"`
	Status      string `xml:"status,omitempty"      json:"status,omitempty"`
	Temperature string `xml:"temperature,omitempty" json:"temperature,omitempty"`
}

type ChassisEnvReply struct {
	XMLName xml.Name  `xml:"environment-information,omitempty" json:"-"`
	EnvItem []EnvItem `xml:"environment-item,omitempty"        json:"environment-item,omitempty"`
}

func ProcessChassisEnv(ncReplies <-chan *netconf.RPCReply, errChan <-chan error) (<-chan *ChassisEnvReply, <-chan error) {

	replyc, errc := make(chan *ChassisEnvReply), make(chan error)

	go func() {
		const thisFunc = "reply.ProcessesChassisEnv.func()"
		cntxLog := log.WithFields(log.Fields{
			"func": thisFunc,
		})
		done := false
		for !done {
			select {
			case ncReply, ok := <-ncReplies:
				if !ok {
					cntxLog.Debugln("ncReply channel is closed")
					cntxLog.Debugln("closing output reply and error channels")
					done = true
					close(replyc)
					close(errc)
					break
				} else if reply, err := NewChassisEnvReply(ncReply); err != nil {
					cntxLog.Debugln("NewChassisEnvReply returned error")
					errc <- err
				} else {
					cntxLog.Debugln("sending ChassisEnvReply over channel")
					replyc <- reply
				}
			case err, ok := <-errChan:
				if !ok {
					cntxLog.Debugln("error channel closed")
					done = true
					close(replyc)
					close(errc)
					break
				} else {
					cntxLog.Debugln("passing error down output channel")
					errc <- err
				}
			}
		}
	}()

	return replyc, errc
}
func NewChassisEnvReply(ncReply *netconf.RPCReply) (*ChassisEnvReply, error) {

	const thisFunc = "reply.NewChassisEnvReply()"
	cntxLog := log.WithFields(log.Fields{
		"method": thisFunc,
	})

	replyNoNewlines := newLine.Replace(ncReply.Data)
	chassisEnvReply := new(ChassisEnvReply)

	cntxLog.Debugln("unmarshalling XML into chassisEnvReply")
	if err := xml.Unmarshal([]byte(replyNoNewlines), chassisEnvReply); err != nil {
		cntxLog.Debugln("could not unmarshal XML into chassisEnvReply")
		cntxLog.Debugln(err)
		return nil, err
	} else {
		return chassisEnvReply, nil
	}
}

func (cer *ChassisEnvReply) WriteCSV(w io.Writer) error {

	const thisMethod = "*ChassisEnvReply.WriteCSV()"
	cntxLog := log.WithFields(log.Fields{
		"method": thisMethod,
	})

	records := make([][]string, 0, 32)
	tStamp := time.Now()

	cntxLog.Debugln("building records for CSV file")
	for _, envItem := range cer.EnvItem {
		record := make([]string, 0, 4)
		if envItem.Class == "Power" {
			continue
		}
		record = append(record, tStamp.Format(time.Stamp))
		record = append(record, envItem.Name)
		record = append(record, envItem.Status)
		record = append(record, envItem.Temperature)
		records = append(records, record)
	}

	csvWriter := csv.NewWriter(w)
	cntxLog.Debugln("writing records")
	if err := csvWriter.WriteAll(records); err != nil {
		cntxLog.Debugln("could not write all CSV records")
		cntxLog.Debugln(err)
		return err
	} else {
		return nil
	}
}
