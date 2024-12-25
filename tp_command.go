package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Command is TPDU message from MS to SC
type Command struct {
	rpData

	SRR bool `json:"tp-srr"` // O / Status Report Request

	TMR byte     `json:"tp-mr"`           // M / Message Reference for TP
	PID byte     `json:"tp-pid"`          // M / Protocol Identifier
	CT  byte     `json:"tp-ct"`           // M / Command Type
	MN  byte     `json:"tp-mn"`           // M / Message Number
	DA  Address  `json:"tp-da"`           // M / Destination Address
	CD  UserData `json:"tp-cd,omitempty"` // O / Command Data
}

var binDCS = GeneralDataCoding{MsgCharset: Charset8bitData}

// MarshalTP output byte data of this TPDU
func (d Command) MarshalTP() []byte {
	w := new(bytes.Buffer)

	b := byte(0x02)
	if d.SRR {
		b |= 0x20
	}
	if len(d.CD.UDH) != 0 {
		b |= 0x40
	}
	w.WriteByte(b)
	w.WriteByte(d.TMR)
	w.WriteByte(d.PID)
	w.WriteByte(d.CT)
	w.WriteByte(d.MN)
	l, a := d.DA.Marshal()
	w.WriteByte(l)
	w.Write(a)
	d.CD.write(w, binDCS)

	return w.Bytes()
}

// MarshalRP output byte data of this RPDU
func (d Command) MarshalRP() []byte {
	return d.rpData.marshal(true, d.MarshalTP())
}

// MarshalCP output byte data of this CPDU
func (d Command) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// UnmarshalCommand decode Submit from bytes
func UnmarshalCommand(b []byte) (d Command, e error) {
	e = d.UnmarshalTP(b)
	return
}

// UnmarshalTP get data of this TPDU
func (d *Command) UnmarshalTP(b []byte) (e error) {
	if len(b) == 0 {
		return io.EOF
	}
	if b[0]&0x03 != 0x02 {
		return UnexpectedMessageTypeError{
			Expected: 0x02, Actual: b[0] & 0x03}
	}

	d.SRR = b[0]&0x20 == 0x20

	r := bytes.NewReader(b[1:])

	if d.TMR, e = r.ReadByte(); e != nil {
		return
	}
	if d.PID, e = r.ReadByte(); e != nil {
		return
	}
	if d.CT, e = r.ReadByte(); e != nil {
		return
	}
	if d.MN, e = r.ReadByte(); e != nil {
		return
	}
	if d.DA, e = readTPAddr(r); e != nil {
		return
	}

	if e = d.CD.read(r, binDCS, b[0]&0x40 == 0x40); e != nil {
		return
	}
	if r.Len() != 0 {
		e = ErrExtraData
	}
	return
}

// UnmarshalRP get data of this RPDU
func (d *Command) UnmarshalRP(b []byte) (e error) {
	if b, e = d.unmarshal(true, b); e == nil {
		e = d.UnmarshalTP(b)
	}
	return
}

// UnmarshalCP get data of this CPDU
func (d *Command) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

// MarshalJSON provide custom marshaller
func (d Command) MarshalJSON() ([]byte, error) {
	type alias Command
	al := struct {
		*alias
		Cd *UserData `json:"tp-cd,omitempty"`
	}{alias: (*alias)(&d)}
	if !d.CD.isEmpty() {
		al.Cd = &d.CD
	}
	return json.Marshal(al)
}

// UnmarshalJSON provide custom marshaller
func (d *Command) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	type alias Command
	al := struct {
		Cd *UserData `json:"tp-cd,omitempty"`
		*alias
	}{alias: (*alias)(d)}
	if e := json.Unmarshal(b, &al); e != nil {
		return e
	}
	if al.Cd != nil {
		d.CD = *al.Cd
	}
	return nil
}

func (d Command) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "TP-Command\n")
	fmt.Fprintf(w, "%sCP-TI:   %s\n", Indent, cpTIStat(d.TI))
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.RMR)
	fmt.Fprintf(w, "%sRP-OA:   <nil>\n", Indent)
	fmt.Fprintf(w, "%sRP-DA:   %s\n", Indent, d.SCA)
	fmt.Fprintf(w, "%sTP-SRR:  %s\n", Indent, srrStat(d.SRR))
	fmt.Fprintf(w, "%sTP-MR:   %d\n", Indent, d.TMR)
	fmt.Fprintf(w, "%sTP-PID:  %s\n", Indent, pidStat(d.PID))
	fmt.Fprintf(w, "%sTP-CT:   %d\n", Indent, d.CT)
	fmt.Fprintf(w, "%sTP-MN:   %d\n", Indent, d.MN)
	fmt.Fprintf(w, "%sTP-DA:   %s\n", Indent, d.DA)
	if !d.CD.isEmpty() {
		fmt.Fprintf(w, "%sTP-CD:\n", Indent)
		fmt.Fprintf(w, "%s", d.CD.String())
	}

	return w.String()
}
