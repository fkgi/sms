package sms

import (
	"bytes"
	"fmt"
)

// Data is RP-DATA RPDU
type Data struct {
	MR byte    // M / Message Reference
	OA Address // C / Originator Address
	DA Address // C / Recipient Address
	UD TPDU    // M / User Data
}

// EncodeMO returns binary data
func (d *Data) EncodeMO() []byte {
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

// EncodeMT returns binary data
func (d *Data) EncodeMT() []byte {
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

// DecodeMO reads binary data
func (d *Data) DecodeMO(b []byte) error {
	if d == nil {
		return fmt.Errorf("nil data")
	}
	r := bytes.NewReader(b)
	var e error
	if tmp, e := r.ReadByte(); e != nil {
		return e
	} else if tmp != 0 {
		return fmt.Errorf("invalid data")
	}
	if d.MR, e = r.ReadByte(); e != nil {
		return e
	}
	if tmp, e := r.ReadByte(); e != nil {
		return e
	} else if tmp != 0 {
		return fmt.Errorf("invalid data")
	}
	if d.DA, e = readAddr(r); e != nil {
		return e
	}
	if l, e := r.ReadByte(); e == nil {
		b = make([]byte, int(l))
	} else {
		return e
	}
	if _, e := r.Read(b); e != nil {
		return e
	}
	d.UD, e = DecodeAsSC(b)
	return e
}

// DecodeMT reads binary data
func (d *Data) DecodeMT(b []byte) error {
	if d == nil {
		return fmt.Errorf("nil data")
	}
	r := bytes.NewReader(b)
	var e error
	if tmp, e := r.ReadByte(); e != nil {
		return e
	} else if tmp != 1 {
		return fmt.Errorf("invalid data")
	}
	if d.MR, e = r.ReadByte(); e != nil {
		return e
	}
	if d.OA, e = readAddr(r); e != nil {
		return e
	}
	if tmp, e := r.ReadByte(); e != nil {
		return e
	} else if tmp != 0 {
		return fmt.Errorf("invalid data")
	}
	if l, e := r.ReadByte(); e == nil {
		b = make([]byte, int(l))
	} else {
		return e
	}
	if _, e = r.Read(b); e != nil {
		return e
	}
	d.UD, e = DecodeAsMS(b)
	return e
}

func (d *Data) String() string {
	if d == nil {
		return "<nil>"
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "SMS message stack: Data\n")
	fmt.Fprintf(w, "%sRP-MR:   %d\n", Indent, d.MR)
	fmt.Fprintf(w, "%sRP-OA:   %s\n", Indent, d.OA)
	fmt.Fprintf(w, "%sRP-DA:   %s\n", Indent, d.DA)
	fmt.Fprintf(w, "%sRP-UD:   %s\n", Indent, d.UD)
	return w.String()
}
