package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Submit is TPDU message from MS to SC
type Submit struct {
	RD  bool `json:"rd"`  // M / Reject Duplicates
	SRR bool `json:"srr"` // O / Status Report Request
	RP  bool `json:"rp"`  // M / Reply Path

	MR  byte           `json:"mr"`           // M / Message Reference
	DA  Address        `json:"da"`           // M / Destination Address
	PID byte           `json:"pid"`          // M / Protocol Identifier
	DCS DataCoding     `json:"dcs"`          // M / Data Coding Scheme
	VP  ValidityPeriod `json:"vp,omitempty"` // O / Validity Period
	UD  UserData       `json:"ud,omitempty"` // O / User Data
}

// MarshalTP output byte data of this TPDU
func (d Submit) MarshalTP() []byte {
	w := new(bytes.Buffer)

	b := byte(0x01)
	if d.RD {
		b |= 0x04
	}
	var vp []byte
	switch v := d.VP.(type) {
	case VPRelative:
		b |= 0x10
		vp = []byte{byte(v)}
	case VPEnhanced:
		b |= 0x08
		vp = []byte{v[0], v[1], v[2], v[3], v[4], v[5], v[6]}
	case VPAbsolute:
		b |= 0x18
		vp = []byte{v[0], v[1], v[2], v[3], v[4], v[5], v[6]}
	default:
		// nil VP value
	}
	if d.SRR {
		b |= 0x20
	}
	if len(d.UD.UDH) != 0 {
		b |= 0x40
	}
	if d.RP {
		b |= 0x80
	}
	w.WriteByte(b)
	w.WriteByte(d.MR)
	l, a := d.DA.Marshal()
	w.WriteByte(l)
	w.Write(a)
	w.WriteByte(d.PID)
	if d.DCS == nil {
		w.WriteByte(0x00)
	} else {
		w.WriteByte(d.DCS.Marshal())
	}
	w.Write(vp)
	d.UD.write(w, d.DCS)

	return w.Bytes()
}

// UnmarshalSubmit decode Submit from bytes
func UnmarshalSubmit(b []byte) (d Submit, e error) {
	e = d.UnmarshalTP(b)
	return
}

// UnmarshalTP get data of this TPDU
func (d *Submit) UnmarshalTP(b []byte) (e error) {
	if len(b) == 0 {
		return io.EOF
	}
	if b[0]&0x03 != 0x01 {
		return UnexpectedMessageTypeError{
			Expected: 0x01, Actual: b[0] & 0x03}
	}

	d.RD = b[0]&0x04 == 0x04
	d.SRR = b[0]&0x20 == 0x20
	d.RP = b[0]&0x80 == 0x80

	r := bytes.NewReader(b[1:])

	if d.MR, e = r.ReadByte(); e != nil {
		return
	}
	if d.DA, e = readTPAddr(r); e != nil {
		return
	}
	if d.PID, e = r.ReadByte(); e != nil {
		return
	}
	if d.DCS, e = readDataCoding(r); e != nil {
		return
	}
	switch b[0] & 0x18 {
	case 0x00:
		d.VP = nil
	case 0x10:
		var p byte
		if p, e = r.ReadByte(); e == nil {
			d.VP = VPRelative(p)
		} else {
			return
		}
	case 0x08:
		var p [7]byte
		if p, e = read7Bytes(r); e == nil {
			d.VP = VPEnhanced(p)
		} else {
			return
		}
	case 0x18:
		var p [7]byte
		if p, e = read7Bytes(r); e == nil {
			d.VP = VPAbsolute(p)
		} else {
			return
		}
	}
	if e = d.UD.read(r, d.DCS, b[0]&0x40 == 0x40); e != nil {
		return
	}
	if r.Len() != 0 {
		e = InvalidLengthError{}
	}
	return
}

// MarshalJSON provide custom marshaller
func (d Submit) MarshalJSON() ([]byte, error) {
	type alias Submit
	al := struct {
		*alias
		Dcs byte      `json:"dcs"`
		Vp  *jvp      `json:"vp,omitempty"`
		Ud  *UserData `json:"ud,omitempty"`
	}{
		alias: (*alias)(&d)}
	if d.DCS != nil {
		al.Dcs = d.DCS.Marshal()
	} else {
		al.Dcs = 0
	}
	if d.VP != nil {
		al.Vp = &jvp{
			T: d.VP.Duration(),
			S: d.VP.SingleAttempt()}
	}
	if !d.UD.isEmpty() {
		al.Ud = &d.UD
	}
	return json.Marshal(al)
}

// UnmarshalJSON provide custom marshaller
func (d *Submit) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	type alias Submit
	al := struct {
		Dcs byte      `json:"dcs"`
		Vp  *jvp      `json:"vp,omitempty"`
		Ud  *UserData `json:"ud,omitempty"`
		*alias
	}{
		alias: (*alias)(d)}
	if e := json.Unmarshal(b, &al); e != nil {
		return e
	}
	d.DCS = UnmarshalDataCoding(al.Dcs)
	if al.Vp != nil {
		d.VP = ValidityPeriodOf(al.Vp.T, al.Vp.S)
	}
	if al.Ud != nil {
		d.UD = *al.Ud
	}
	return nil
}

func (d Submit) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "SMS message stack: Submit\n")
	fmt.Fprintf(w, "%sTP-RD:   %s\n", Indent, rdStat(d.RD))
	fmt.Fprintf(w, "%sTP-SRR:  %s\n", Indent, srrStat(d.SRR))
	fmt.Fprintf(w, "%sTP-RP:   %s\n", Indent, rpStat(d.RP))
	fmt.Fprintf(w, "%sTP-MR:   %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sTP-DA:   %s\n", Indent, d.DA)
	fmt.Fprintf(w, "%sTP-PID:  %s\n", Indent, pidStat(d.PID))
	fmt.Fprintf(w, "%sTP-DCS:  %s\n", Indent, d.DCS)
	if d.VP != nil {
		fmt.Fprintf(w, "%sTP-VP:   %s\n", Indent, d.VP)
	}
	if !d.UD.isEmpty() {
		fmt.Fprintf(w, "%s", d.UD.String())
	}

	return w.String()
}

// SubmitReport is TPDU message from SC to MS
type SubmitReport struct {
	FCS  byte       `json:"fcs,omitempty"` // C / Failure Cause
	SCTS time.Time  `json:"scts"`          // M / Service Centre Time Stamp
	PID  *byte      `json:"pid,omitempty"` // O / Protocol Identifier
	DCS  DataCoding `json:"dcs,omitempty"` // O / Data Coding Scheme
	UD   UserData   `json:"uid,omitempty"` // O / User Data
}

// MarshalTP output byte data of this TPDU
func (d SubmitReport) MarshalTP() []byte {
	w := new(bytes.Buffer)

	b := byte(0x01)
	if len(d.UD.UDH) != 0 {
		b |= 0x40
	}
	w.WriteByte(b)
	if d.FCS&0x80 == 0x80 {
		w.WriteByte(d.FCS)
	}
	b = byte(0x00)
	if d.PID != nil {
		b |= 0x01
	}
	if d.DCS != nil {
		b |= 0x02
	}
	if len(d.UD.Text) != 0 || len(d.UD.UDH) != 0 {
		b |= 0x04
	}
	w.WriteByte(b)
	w.Write(marshalSCTimeStamp(d.SCTS))
	if d.PID != nil {
		w.WriteByte(*d.PID)
	}
	if d.DCS != nil {
		w.WriteByte(d.DCS.Marshal())
	}
	if !d.UD.isEmpty() {
		d.UD.write(w, d.DCS)
	}

	return w.Bytes()
}

// UnmarshalSubmitReport decode SubmitReport from bytes
func UnmarshalSubmitReport(b []byte) (d SubmitReport, e error) {
	e = d.UnmarshalTP(b)
	return
}

// UnmarshalTP get data of this TPDU
func (d *SubmitReport) UnmarshalTP(b []byte) (e error) {
	if len(b) == 0 {
		return io.EOF
	}
	if b[0]&0x03 != 0x01 {
		return UnexpectedMessageTypeError{
			Expected: 0x01, Actual: b[0] & 0x03}
	}

	r := bytes.NewReader(b[1:])

	var pi byte
	if pi, e = r.ReadByte(); e == nil && pi&0x80 == 0x80 {
		d.FCS = pi
		pi, e = r.ReadByte()
	}
	if e != nil {
		return
	}
	var p [7]byte
	if p, e = read7Bytes(r); e != nil {
		return
	}
	d.SCTS = unmarshalSCTimeStamp(p)
	if pi&0x01 == 0x01 {
		var p byte
		if p, e = r.ReadByte(); e != nil {
			return
		}
		d.PID = &p
	}
	if pi&0x02 == 0x02 {
		if d.DCS, e = readDataCoding(r); e != nil {
			return
		}
	}
	if pi&0x04 == 0x04 {
		e = d.UD.read(r, d.DCS, b[0]&0x40 == 0x40)
		if e != nil {
			return
		}
	}
	if r.Len() != 0 {
		e = InvalidLengthError{}
	}
	return
}

// MarshalJSON provide custom marshaller
func (d SubmitReport) MarshalJSON() ([]byte, error) {
	al := struct {
		Fcs  *byte     `json:"fcs,omitempty"`
		Scts time.Time `json:"scts"`
		Pid  *byte     `json:"pid,omitempty"`
		Dcs  *byte     `json:"dcs,omitempty"`
		Ud   *UserData `json:"ud,omitempty"`
	}{}
	if d.FCS&0x80 == 0x80 {
		al.Fcs = &d.FCS
	}
	al.Scts = d.SCTS
	al.Pid = d.PID
	if d.DCS != nil {
		tmp := d.DCS.Marshal()
		al.Dcs = &tmp
	}
	if !d.UD.isEmpty() {
		al.Ud = &d.UD
	}
	return json.Marshal(al)
}

// UnmarshalJSON provide custom marshaller
func (d *SubmitReport) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	al := struct {
		Fcs  *byte     `json:"fcs,omitempty"`
		Scts time.Time `json:"scts"`
		Pid  *byte     `json:"pid,omitempty"`
		Dcs  *byte     `json:"dcs,omitempty"`
		Ud   *UserData `json:"ud,omitempty"`
	}{}
	if e := json.Unmarshal(b, &al); e != nil {
		return e
	}
	if al.Fcs != nil && *al.Fcs&0x80 == 0x80 {
		d.FCS = *al.Fcs
	}
	d.SCTS = al.Scts
	d.PID = al.Pid
	if al.Dcs != nil {
		d.DCS = UnmarshalDataCoding(*al.Dcs)
	}
	if al.Ud != nil {
		d.UD = *al.Ud
	}
	return nil
}

func (d SubmitReport) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "SMS message stack: Submit Report")
	if d.FCS&0x80 == 0x80 {
		fmt.Fprintf(w, " for RP-ERROR\n")
		fmt.Fprintf(w, "%sTP-FCS:  %s\n", Indent, fcsStat(d.FCS))
	} else {
		fmt.Fprintf(w, " for RP-ACK\n")
	}

	fmt.Fprintf(w, "%sTP-SCTS: %s\n", Indent, d.SCTS)
	if d.PID != nil {
		fmt.Fprintf(w, "%sTP-PID:  %s\n", Indent, pidStat(*d.PID))
	}
	if d.DCS != nil {
		fmt.Fprintf(w, "%sTP-DCS:  %s\n", Indent, d.DCS)
	}
	if !d.UD.isEmpty() {
		fmt.Fprintf(w, "%s", d.UD.String())
	}

	return w.String()
}
