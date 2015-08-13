package reply

import (
	"encoding/xml"
)

type Zone struct {
	Name           string `xml:"name,omitempty"             json:"name,omitempty"`
	DrivingFRUName string `xml:"driving-fru-name,omitempty" json:"driving-fru-name,omitempty"`
	Temperature    string `xml:"temperature,omitempty"      json:"temperature,omitempty"`
	ZoneStatus     string `xml:"zone-status,omitempty"      json:"zone-status,omitempty"`
	FanMissingCnt  int    `xml:"fan-missing-cnt,omitempty"  json:"fan-missing-cnt,omitempty"`
	FanFailedCnt   int    `xml:"fan-failed-cnt,omitempty"   json:"fan-failed-cnt,omitempty"`
	FanDutyCycle   int    `xml:"fan-dutycycle,omitempty"    json:"fan-dutycycle,omitempty"`
}
type ChassisZonesReply struct {
	XMLName xml.Name `xml:"chassis-zones-information,omitempty" json:"-"`
	Zones   []Zone
}

func NewZonesResponse(ncReply *string) (*ChassisZonesReply, error) {

	var (
		replyNoNewlines      = newLine.Replace(*ncReply)
		chassisZonesResponse = new(ChassisZonesReply)
	)

	if err := xml.Unmarshal([]byte(replyNoNewlines), chassisZonesResponse); err != nil {
		return nil, err
	} else {
		return chassisZonesResponse, nil
	}
}
