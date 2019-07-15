package sms

import (
	"bytes"
	"fmt"
)

// DataMO is MO RP-DATA RPDU
type DataMO struct {
	MR byte
	DA Address
	UD TPDU
}

// Encode returns binary data
func (d *DataMO) Encode() []byte {
	if d == nil {
		return []byte{}
	}

	w := new(bytes.Buffer)
	w.WriteByte(0) // MTI
	w.WriteByte(d.MR)
	w.WriteByte(0) // OA
	l, a := d.DA.Encode()
	w.WriteByte(l)
	w.Write(a)
	b := d.UD.Encode()
	w.WriteByte(byte(len(b)))
	w.Write(b)

	return w.Bytes()
}

// Decode reads binary data
func (d *DataMO) Decode(b []byte) (e error) {
	if d == nil {
		return fmt.Errorf("nil data")
	}
	r := bytes.NewReader(b)
	var tmp byte
	if tmp, e = r.ReadByte(); e != nil {
		return
	}
	if tmp != 0 {
		return fmt.Errorf("invalid data")
	}
	if d.MR, e = r.ReadByte(); e != nil {
		return
	}
	if tmp, e = r.ReadByte(); e != nil {
		return
	}
	if tmp != 0 {
		return fmt.Errorf("invalid data")
	}
	if d.DA, e = readAddr(r); e != nil {
		return
	}
	var l byte
	if l, e = r.ReadByte(); e != nil {
		return
	}
	b = make([]byte, int(l))
	if _, e = r.Read(b); e != nil {
		return
	}
	d.UD, e = DecodeAsSC(b)
	return
}

func (d *DataMO) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: MO Data\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sRP-OA:   <nil>\n", Indent)
	fmt.Fprintf(w, "%sRP-DA:   %s\n", Indent, d.DA)
	fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	return w.String()
}

// DataMT is MT RP-DATA RPDU
type DataMT struct {
	MR byte
	OA Address
	UD TPDU
}

// Encode returns binary data
func (d *DataMT) Encode() []byte {
	if d == nil {
		return []byte{}
	}

	w := new(bytes.Buffer)
	w.WriteByte(1) // MTI
	w.WriteByte(d.MR)
	l, a := d.OA.Encode()
	w.WriteByte(l)
	w.Write(a)
	w.WriteByte(0) // DA
	b := d.UD.Encode()
	w.WriteByte(byte(len(b)))
	w.Write(b)

	return w.Bytes()
}

// Decode reads binary data
func (d *DataMT) Decode(b []byte) (e error) {
	if d == nil {
		return fmt.Errorf("nil data")
	}
	r := bytes.NewReader(b)
	var tmp byte
	if tmp, e = r.ReadByte(); e != nil {
		return
	}
	if tmp != 1 {
		return fmt.Errorf("invalid data")
	}
	if d.MR, e = r.ReadByte(); e != nil {
		return
	}
	if d.OA, e = readAddr(r); e != nil {
		return
	}
	if tmp, e = r.ReadByte(); e != nil {
		return
	}
	if tmp != 0 {
		return fmt.Errorf("invalid data")
	}
	var l byte
	if l, e = r.ReadByte(); e != nil {
		return
	}
	b = make([]byte, int(l))
	if _, e = r.Read(b); e != nil {
		return
	}
	d.UD, e = DecodeAsSC(b)
	return
}

func (d *DataMT) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: MT Data\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sRP-OA:   %s\n", Indent, d.OA)
	fmt.Fprintf(w, "%sRP-DA:   <nil>\n", Indent)
	fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	return w.String()
}
