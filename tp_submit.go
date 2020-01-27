package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Submit is TPDU message from MS to SC
type Submit struct {
	rpRequest

	RD  bool `json:"rd"`  // M / Reject Duplicates
	SRR bool `json:"srr"` // O / Status Report Request
	RP  bool `json:"rp"`  // M / Reply Path

	TMR byte           `json:"tmr"`          // M / Message Reference for TP
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
	w.WriteByte(d.TMR)
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

// MarshalRP output byte data of this RPDU
func (d Submit) MarshalRP() []byte {
	return d.rpRequest.marshal(true, d.MarshalTP())
}

// MarshalCP output byte data of this CPDU
func (d Submit) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
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

	if d.TMR, e = r.ReadByte(); e != nil {
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
		e = ErrExtraData
	}
	return
}

// UnmarshalRP get data of this RPDU
func (d *Submit) UnmarshalRP(b []byte) (e error) {
	if b, e = d.unmarshal(true, b); e == nil {
		e = d.UnmarshalTP(b)
	}
	return
}

// UnmarshalCP get data of this CPDU
func (d *Submit) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
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

	fmt.Fprintf(w, "TP-Submit\n")
	fmt.Fprintf(w, "%sCP-TI:   %d\n", Indent, d.TI)
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.RMR)
	fmt.Fprintf(w, "%sRP-OA:   <nil>\n", Indent)
	fmt.Fprintf(w, "%sRP-DA:   %s\n", Indent, d.SCA)
	fmt.Fprintf(w, "%sTP-RD:   %s\n", Indent, rdStat(d.RD))
	fmt.Fprintf(w, "%sTP-SRR:  %s\n", Indent, srrStat(d.SRR))
	fmt.Fprintf(w, "%sTP-RP:   %s\n", Indent, rpStat(d.RP))
	fmt.Fprintf(w, "%sTP-MR:   %d\n", Indent, d.TMR)
	fmt.Fprintf(w, "%sTP-DA:   %s\n", Indent, d.DA)
	fmt.Fprintf(w, "%sTP-PID:  %s\n", Indent, pidStat(d.PID))
	fmt.Fprintf(w, "%sTP-DCS:  %s\n", Indent, d.DCS)
	if d.VP != nil {
		fmt.Fprintf(w, "%sTP-VP:   %s\n", Indent, d.VP)
	}
	if !d.UD.isEmpty() {
		fmt.Fprintf(w, "%sTP-UD:\n", Indent)
		fmt.Fprintf(w, "%s", d.UD.String())
	}

	return w.String()
}
