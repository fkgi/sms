package sms

import (
	"bytes"
	"fmt"
)

// Data is RP-DATA RPDU
type Data struct {
	MR byte
	OA Address
	DA Address
	UD TPDU
}

func (d *Data) encode(mo bool) []byte {
	if d == nil {
		return []byte{}
	}

	w := new(bytes.Buffer)
	if mo {
		w.WriteByte(0)
	} else {
		w.WriteByte(1)
	}
	w.WriteByte(d.MR)

	if mo {
		// w.WriteByte(IEI)
		w.WriteByte(0)

		// w.WriteByte(IEI)
		l, a := d.DA.Encode()
		w.WriteByte(l)
		w.Write(a)
	} else {
		// w.WriteByte(IEI)
		l, a := d.OA.Encode()
		w.WriteByte(l)
		w.Write(a)

		// w.WriteByte(IEI)
		w.WriteByte(0)
	}

	b := d.UD.Encode()
	// w.WriteByte(IEI)
	w.WriteByte(byte(len(b)))
	w.Write(b)
	return w.Bytes()
}

func (d *Data) EncodeMO() []byte {
	return d.encode(true)
}

func (d *Data) EncodeMT() []byte {
	return d.encode(false)
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
