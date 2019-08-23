package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// DeliverReport is TPDU message from MS to SC
type DeliverReport struct {
	cpData

	RMR  byte  `json:"rmr"`            // M / Message Reference for RP
	CS   byte  `json:"cs"`             // M / Cause
	DIAG *byte `json:"diag,omitempty"` // O / Diagnostics

	FCS byte       `json:"fcs,omitempty"` // C / Failure Cause
	PID *byte      `json:"pid,omitempty"` // O / Protocol Identifier
	DCS DataCoding `json:"dcs,omitempty"` // O / Data Coding Scheme
	UD  UserData   `json:"ud,omitempty"`  // O / User Data
}

// MarshalTP output byte data of this TPDU
func (d DeliverReport) MarshalTP() []byte {
	w := new(bytes.Buffer)

	b := byte(0x00)
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
func (d DeliverReport) MarshalRP() []byte {
	w := new(bytes.Buffer)

	if d.FCS&0x80 == 0x80 {
		w.WriteByte(4) // MTI
		w.WriteByte(d.RMR)
		if d.DIAG != nil {
			w.WriteByte(2)
			w.WriteByte(d.CS)
			w.WriteByte(*d.DIAG)
		} else {
			w.WriteByte(1)
			w.WriteByte(d.CS)
		}
	} else {
		w.WriteByte(2) // MTI
		w.WriteByte(d.RMR)
	}
	b := d.MarshalTP()
	w.WriteByte(0x41)
	w.WriteByte(byte(len(b)))
	w.Write(b)

	return w.Bytes()
}

// MarshalCP output byte data of this CPDU
func (d DeliverReport) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// UnmarshalDeliverReport decode DeliverReport from bytes
func UnmarshalDeliverReport(b []byte) (d DeliverReport, e error) {
	e = d.UnmarshalTP(b)
	return
}

// UnmarshalTP get data of this TPDU
func (d *DeliverReport) UnmarshalTP(b []byte) (e error) {
	if len(b) == 0 {
		return io.EOF
	}
	if b[0]&0x03 != 0x00 {
		return UnexpectedMessageTypeError{
			Expected: 0x00, Actual: b[0] & 0x03}
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

// UnmarshalRP get data of this TPDU
func (d *DeliverReport) UnmarshalRP(b []byte) (e error) {
	r := bytes.NewReader(b)

	if mti, e := r.ReadByte(); e != nil {
		return e
	} else if mti == 4 {
		if d.RMR, e = r.ReadByte(); e != nil {
			return e
		}
	} else if mti == 6 {
		if d.RMR, e = r.ReadByte(); e != nil {
			return e
		}
		if l, e := r.ReadByte(); e != nil {
			return e
		} else if l == 0 || l > 2 {
			return InvalidLengthError{}
		} else if d.CS, e = r.ReadByte(); e != nil {
			return e
		} else if l == 2 {
			if l, e = r.ReadByte(); e != nil {
				return e
			}
			d.DIAG = &l
		}
	} else {
		return UnexpectedMessageTypeError{
			Expected: 0, Actual: mti}
	}

	if iei, e := r.ReadByte(); e != nil {
		return e
	} else if iei != 0x41 {
		return UnexpectedInformationElementError{
			Expected: 0x41, Actual: iei}
	}
	if l, e := r.ReadByte(); e == nil {
		b = make([]byte, int(l))
	} else {
		return e
	}
	if n, e := r.Read(b); e != nil {
		return e
	} else if n != len(b) {
		return io.EOF
	}
	if r.Len() != 0 {
		return InvalidLengthError{}
	}
	return d.UnmarshalTP(b)
}

// UnmarshalCP get data of this CPDU
func (d *DeliverReport) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

// UnmarshalJSON provide custom marshaller
func (d *DeliverReport) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	al := struct {
		Fcs *byte     `json:"fcs,omitempty"`
		Pid *byte     `json:"pid,omitempty"`
		Dcs *byte     `json:"dcs,omitempty"`
		Ud  *UserData `json:"ud,omitempty"`
	}{}
	if e := json.Unmarshal(b, &al); e != nil {
		return e
	}
	if al.Fcs != nil && *al.Fcs&0x80 == 0x80 {
		d.FCS = *al.Fcs
	}
	d.PID = al.Pid
	if al.Dcs != nil {
		d.DCS = UnmarshalDataCoding(*al.Dcs)
	}
	if al.Ud != nil {
		d.UD = *al.Ud
	}

	return nil
}

// MarshalJSON provide custom marshaller
func (d DeliverReport) MarshalJSON() ([]byte, error) {
	al := struct {
		Fcs *byte     `json:"fcs,omitempty"`
		Pid *byte     `json:"pid,omitempty"`
		Dcs *byte     `json:"dcs,omitempty"`
		Ud  *UserData `json:"ud,omitempty"`
	}{}
	if d.FCS&0x80 == 0x80 {
		al.Fcs = &d.FCS
	}
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

func (d DeliverReport) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "SMS message stack: Deliver Report")
	if d.FCS&0x80 == 0x80 {
		fmt.Fprintf(w, " for RP-ERROR\n")
		fmt.Fprintf(w, "%sCP-TI:   %d\n", Indent, d.TI)
		fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.RMR)
		fmt.Fprintf(w, "%sRP-CS:   cause=%s",
			Indent, rpCauseStat(d.CS))
		if d.DIAG != nil {
			fmt.Fprintf(w, "diagnostic=%d\n", *d.DIAG)
		} else {
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "%sTP-FCS:  %s\n", Indent, fcsStat(d.FCS))
	} else {
		fmt.Fprintf(w, " for RP-ACK\n")
		fmt.Fprintf(w, "%sCP-TI:   %d\n", Indent, d.TI)
		fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.RMR)
	}

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
