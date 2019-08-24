package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Deliver is TPDU message from SC to MS
type Deliver struct {
	rpRequest

	MMS bool `json:"mms"` // M / More Messages to Send (true=more messages)
	LP  bool `json:"lp"`  // O / Loop Prevention
	SRI bool `json:"sri"` // O / Status Report Indication (true=status report shall be returned)
	RP  bool `json:"rp"`  // M / Reply Path

	OA   Address    `json:"oa"`           // M / Originating Address
	PID  byte       `json:"pid"`          // M / Protocol Identifier
	DCS  DataCoding `json:"dcs"`          // M / Data Coding Scheme
	SCTS time.Time  `json:"scts"`         // M / Service Centre Time Stamp
	UD   UserData   `json:"ud,omitempty"` // O / User Data
}

// MarshalTP output byte data of this TPDU
func (d Deliver) MarshalTP() []byte {
	w := new(bytes.Buffer)

	b := byte(0x00)
	if !d.MMS {
		b |= 0x04
	}
	if d.LP {
		b |= 0x08
	}
	if d.SRI {
		b |= 0x20
	}
	if len(d.UD.UDH) != 0 {
		b |= 0x40
	}
	if d.RP {
		b |= 0x80
	}
	w.WriteByte(b)
	l, a := d.OA.Marshal()
	w.WriteByte(l)
	w.Write(a)
	w.WriteByte(d.PID)
	if d.DCS == nil {
		w.WriteByte(0x00)
	} else {
		w.WriteByte(d.DCS.Marshal())
	}
	w.Write(marshalSCTimeStamp(d.SCTS))
	d.UD.write(w, d.DCS)

	return w.Bytes()
}

// MarshalRP output byte data of this RPDU
func (d Deliver) MarshalRP() []byte {
	return d.rpRequest.marshal(false, d.MarshalTP())
}

// MarshalCP output byte data of this CPDU
func (d Deliver) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// UnmarshalDeliver decode Deliver from bytes
func UnmarshalDeliver(b []byte) (d Deliver, e error) {
	e = d.UnmarshalTP(b)
	return
}

// UnmarshalTP get data of this TPDU
func (d *Deliver) UnmarshalTP(b []byte) (e error) {
	if len(b) == 0 {
		return io.EOF
	}
	if b[0]&0x03 != 0x00 {
		return UnexpectedMessageTypeError{
			Expected: 0x00, Actual: b[0] & 0x03}
	}

	d.MMS = b[0]&0x04 != 0x04
	d.LP = b[0]&0x08 == 0x08
	d.SRI = b[0]&0x20 == 0x20
	d.RP = b[0]&0x80 == 0x80

	r := bytes.NewReader(b[1:])

	if d.OA, e = readTPAddr(r); e != nil {
		return
	}
	if d.PID, e = r.ReadByte(); e != nil {
		return
	}
	if d.DCS, e = readDataCoding(r); e != nil {
		return
	}
	var tmp [7]byte
	if tmp, e = read7Bytes(r); e != nil {
		return
	}
	d.SCTS = unmarshalSCTimeStamp(tmp)
	if e = d.UD.read(r, d.DCS, b[0]&0x40 == 0x40); e != nil {
		return
	}
	if r.Len() != 0 {
		e = InvalidLengthError{}
	}
	return
}

// UnmarshalRP get data of this TPDU
func (d *Deliver) UnmarshalRP(b []byte) (e error) {
	if b, e = d.unmarshal(false, b); e == nil {
		e = d.UnmarshalTP(b)
	}
	return
}

// UnmarshalCP get data of this CPDU
func (d *Deliver) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

// MarshalJSON provide custom marshaller
func (d Deliver) MarshalJSON() ([]byte, error) {
	type alias Deliver
	al := struct {
		*alias
		Dcs byte      `json:"dcs"`
		Ud  *UserData `json:"ud,omitempty"`
	}{
		alias: (*alias)(&d)}
	if d.DCS != nil {
		al.Dcs = d.DCS.Marshal()
	} else {
		al.Dcs = 0
	}
	if !d.UD.isEmpty() {
		al.Ud = &d.UD
	}
	return json.Marshal(al)
}

// UnmarshalJSON provide custom marshaller
func (d *Deliver) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	type alias Deliver
	al := struct {
		Dcs byte      `json:"dcs"`
		Ud  *UserData `json:"ud,omitempty"`
		*alias
	}{
		alias: (*alias)(d)}
	if e := json.Unmarshal(b, &al); e != nil {
		return e
	}
	d.DCS = UnmarshalDataCoding(al.Dcs)
	if al.Ud != nil {
		d.UD = *al.Ud
	}
	return nil
}

func (d Deliver) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "TP-Deliver\n")
	fmt.Fprintf(w, "%sCP-TI:   %d\n", Indent, d.TI)
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.RMR)
	fmt.Fprintf(w, "%sRP-OA:   %d\n", Indent, d.SCA)
	fmt.Fprintf(w, "%sRP-DA:   <nil>\n", Indent)
	fmt.Fprintf(w, "%sTP-MMS:  %s\n", Indent, mmsStat(d.MMS))
	fmt.Fprintf(w, "%sTP-LP:   %s\n", Indent, lpStat(d.LP))
	fmt.Fprintf(w, "%sTP-SRI:  %s\n", Indent, sriStat(d.SRI))
	fmt.Fprintf(w, "%sTP-RP:   %s\n", Indent, rpStat(d.RP))
	fmt.Fprintf(w, "%sTP-OA:   %s\n", Indent, d.OA)
	fmt.Fprintf(w, "%sTP-PID:  %s\n", Indent, pidStat(d.PID))
	if d.DCS != nil {
		fmt.Fprintf(w, "%sTP-DCS:  %s\n", Indent, d.DCS)
	} else {
		fmt.Fprintf(w, "%sTP-DCS:  %s\n", Indent, &GeneralDataCoding{})
	}
	fmt.Fprintf(w, "%sTP-SCTS: %s\n", Indent, d.SCTS)
	if !d.UD.isEmpty() {
		fmt.Fprintf(w, "%s", d.UD.String())
	}

	return w.String()
}
