package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// SubmitReport is TPDU message from SC to MS
type SubmitReport struct {
	cpData

	RMR  byte  `json:"rp-mr"`          // M / Message Reference
	CS   byte  `json:"rp-cs"`          // M / Cause
	DIAG *byte `json:"diag,omitempty"` // O / Diagnostics

	FCS  byte       `json:"tp-fcs,omitempty"` // C / Failure Cause
	SCTS time.Time  `json:"tp-scts"`          // M / Service Centre Time Stamp
	PID  *byte      `json:"tp-pid,omitempty"` // O / Protocol Identifier
	DCS  DataCoding `json:"tp-dcs,omitempty"` // O / Data Coding Scheme
	UD   UserData   `json:"tp-uid,omitempty"` // O / User Data
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

// MarshalRP output byte data of this RPDU
func (d SubmitReport) MarshalRP() []byte {
	if d.FCS&0x80 == 0x80 {
		rp := RpError{RMR: d.RMR, CS: d.CS, DIAG: d.DIAG}
		return rp.marshalRP(false, d.MarshalTP())
	}
	rp := RpAck{RMR: d.RMR}
	return rp.marshalRP(false, d.MarshalTP())
}

// MarshalCP output byte data of this CPDU
func (d SubmitReport) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
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
		e = ErrExtraData
	}
	return
}

// UnmarshalRP get data of this RPDU
func (d *SubmitReport) UnmarshalRP(b []byte) (e error) {
	if len(b) == 0 {
		return io.EOF
	}
	switch b[0] & 0x07 {
	case 0x03:
		rp := RpAck{}
		b, e = rp.unmarshalRP(false, b)
		if e == nil {
			d.RMR = rp.RMR
		}
	case 0x05:
		rp := RpError{}
		b, e = rp.unmarshalRP(false, b)
		if e == nil {
			d.RMR = rp.RMR
			d.CS = rp.CS
			d.DIAG = rp.DIAG
		}
	}
	if b == nil {
		e = io.EOF
	} else {
		e = d.UnmarshalTP(b)
	}
	return
}

// UnmarshalCP get data of this CPDU
func (d *SubmitReport) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

// MarshalJSON provide custom marshaller
func (d SubmitReport) MarshalJSON() ([]byte, error) {
	type alias SubmitReport
	al := struct {
		*alias
		Fcs *byte     `json:"tp-fcs,omitempty"`
		Dcs *byte     `json:"tp-dcs,omitempty"`
		Ud  *UserData `json:"tp-ud,omitempty"`
	}{alias: (*alias)(&d)}
	if d.FCS&0x80 == 0x80 {
		al.Fcs = &d.FCS
	}
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
	type alias SubmitReport
	al := struct {
		*alias
		Fcs *byte     `json:"tp-fcs,omitempty"`
		Dcs *byte     `json:"tp-dcs,omitempty"`
		Ud  *UserData `json:"tp-ud,omitempty"`
	}{alias: (*alias)(d)}
	if e := json.Unmarshal(b, &al); e != nil {
		return e
	}
	if al.Fcs != nil && *al.Fcs&0x80 == 0x80 {
		d.FCS = *al.Fcs
	}
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

	fmt.Fprintf(w, "TP-SubmitReport")
	if d.FCS&0x80 == 0x80 {
		fmt.Fprintf(w, " for RP-ERROR\n")
		fmt.Fprintf(w, "%sCP-TI:   %s\n", Indent, cpTIStat(d.TI))
		fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.RMR)
		fmt.Fprintf(w, "%sRP-CS:   cause=%s",
			Indent, rpCauseStat(d.CS))
		if d.DIAG != nil {
			fmt.Fprintf(w, ", diagnostic=%d\n", *d.DIAG)
		} else {
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "%sTP-FCS:  %s\n", Indent, fcsStat(d.FCS))
	} else {
		fmt.Fprintf(w, " for RP-ACK\n")
		fmt.Fprintf(w, "%sCP-TI:   %s\n", Indent, cpTIStat(d.TI))
		fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.RMR)
	}

	fmt.Fprintf(w, "%sTP-SCTS: %s\n", Indent, d.SCTS)
	if d.PID != nil {
		fmt.Fprintf(w, "%sTP-PID:  %s\n", Indent, pidStat(*d.PID))
	}
	if d.DCS != nil {
		fmt.Fprintf(w, "%sTP-DCS:  %s\n", Indent, d.DCS)
	}
	if !d.UD.isEmpty() {
		fmt.Fprintf(w, "%sTP-UD:\n", Indent)
		fmt.Fprintf(w, "%s", d.UD.String())
	}

	return w.String()
}
