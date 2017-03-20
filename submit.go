package sms

import (
	"fmt"
	"io"
	"time"
)

// Submit is TPDU message from MS to SC
type Submit struct {
	RD  bool // Reject Duplicates
	SRR bool // Status Report Request
	RP  bool // Reply Path

	MR  byte    // Message Reference
	DA  Address // Destination Address
	PID byte    // Protocol Identifier
	DCS dcs     // Data Coding Scheme
	VP  vp      // Validity Period
	UDH []udh   // User Data Header
	UD  []byte  // User Data
}

// WriteTo output byte data of this TPDU
func (d *Submit) WriteTo(w io.Writer) (n int64, e error) {
	b := []byte{0x01, d.MR}
	if d.RD {
		b[0] = b[0] | 0x04
	}
	var vp []byte
	switch v := d.VP.(type) {
	case VPRelative:
		b[0] = b[0] | 0x10
		vp = []byte{v[0]}
	case VPEnhanced:
		b[0] = b[0] | 0x08
		vp = []byte{v[0], v[1], v[2], v[3], v[4], v[5], v[6]}
	case VPAbsolute:
		b[0] = b[0] | 0x18
		vp = []byte{v[0], v[1], v[2], v[3], v[4], v[5], v[6]}
	}
	if d.SRR {
		b[0] = b[0] | 0x20
	}
	if len(d.UDH) != 0 {
		b[0] = b[0] | 0x40
	}
	if d.RP {
		b[0] = b[0] | 0x80
	}
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}

	var nn int64
	nn, e = d.DA.WriteTo(w)
	n += nn
	if e != nil {
		return
	}

	b = []byte{d.PID, d.DCS.encode()}
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}

	if n, e = writeBytes(w, n, vp); e != nil {
		return
	}

	udh := encodeUDH(d.UDH)
	u := d.DCS.unitSize()
	l := len(udh) + len(d.UD)
	l = ((l * 8) - (l * 8 % u)) / u
	b = []byte{byte(l)}
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}
	if n, e = writeBytes(w, n, udh); e != nil {
		return
	}
	n, e = writeBytes(w, n, d.UD)
	return
}

func (d *Submit) readFrom(h byte, r io.Reader) (n int64, e error) {
	d.RD = h&0x04 == 0x04
	d.SRR = h&0x20 == 0x20
	d.RP = h&0x80 == 0x80

	b := make([]byte, 1)
	if n, e = readBytes(r, n, b); e != nil {
		return
	}
	d.MR = b[0]

	d.DA = Address{}
	var nn int64
	if nn, e = d.DA.ReadFrom(r); e != nil {
		return
	}
	n += nn

	b = make([]byte, 2)
	if n, e = readBytes(r, n, b); e != nil {
		return
	}
	d.PID = b[0]
	d.DCS = decodeDCS(b[1])
	if d.DCS == nil {
		e = fmt.Errorf("invalid TP-DCS data: % x", b[1])
		return
	}

	switch h & 0x18 {
	case 0x00:
		d.VP = nil
	case 0x10:
		b = make([]byte, 1)
		if n, e = readBytes(r, n, b); e != nil {
			return
		}
		d.VP = VPRelative([1]byte{b[0]})
	case 0x08:
		b = make([]byte, 7)
		if n, e = readBytes(r, n, b); e != nil {
			return
		}
		d.VP = VPEnhanced([7]byte{
			b[0], b[1], b[2], b[3], b[4], b[5], b[6]})
	case 0x18:
		b = make([]byte, 7)
		if n, e = readBytes(r, n, b); e != nil {
			return
		}
		d.VP = VPAbsolute([7]byte{
			b[0], b[1], b[2], b[3], b[4], b[5], b[6]})
	}

	b = make([]byte, 1)
	if n, e = readBytes(r, n, b); e != nil {
		return
	}
	l := d.DCS.unitSize()
	l *= int(b[0])
	if l%8 != 0 {
		l += 8 - l%8
	}

	d.UD = make([]byte, l/8)
	if n, e = readBytes(r, n, b); e != nil {
		return
	}

	if h&0x40 == 0x40 {
		d.UDH = decodeUDH(d.UD[0 : d.UD[0]+1])
		d.UD = d.UD[d.UD[0]+1:]
	}

	return
}

// PrintStack show PDU parameter
func (d *Submit) PrintStack(w io.Writer) {
	fmt.Fprintf(w, "SMS message stack: Submit\n")
	fmt.Fprintf(w, "TP-RD:   %s\n", rdStat(d.RD))
	fmt.Fprintf(w, "TP-SRR:  %s\n", srrStat(d.SRR))
	fmt.Fprintf(w, "TP-RP:   %s\n", rpStat(d.RP))

	fmt.Fprintf(w, "TP-MR:   %d\n", d.MR)
	fmt.Fprintf(w, "TP-DA:   %s\n", d.DA)
	fmt.Fprintf(w, "TP-PID:  %d\n", d.PID)
	fmt.Fprintf(w, "TP-DCS:  %s\n", d.DCS)
	if d.VP != nil {
		fmt.Fprintf(w, "TP-VP:   %s\n", d.VP)
	}

	if len(d.UDH)+len(d.UD) != 0 {
		fmt.Fprintf(w, "TP-UD:\n")
		for _, h := range d.UDH {
			fmt.Fprintf(w, "%s\n", h)
		}
		if len(d.UD) != 0 {
			fmt.Fprintf(w, "%s\n", d.DCS.decodeData(d.UD))
		}
	}
}

// SubmitReport is TPDU message from SC to MS
type SubmitReport struct {
	FCS  *byte     // Failure Cause
	SCTS time.Time // Service Centre Time Stamp
	PID  *byte     // Protocol Identifier
	DCS  dcs       // Data Coding Scheme
	UDH  []udh     // User Data Header
	UD   []byte    // User Data
}

// WriteTo output byte data of this TPDU
func (d *SubmitReport) WriteTo(w io.Writer) (n int64, e error) {
	b := []byte{0x01}
	if len(d.UDH) != 0 {
		b[0] = b[0] | 0x40
	}
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}

	if d.FCS != nil {
		b = []byte{*d.FCS}
		if n, e = writeBytes(w, n, b); e != nil {
			return
		}
	}

	b = []byte{0x00}
	if d.PID != nil {
		b[0] = b[0] | 0x01
	}
	if d.DCS != nil {
		b[0] = b[0] | 0x02
	}
	if len(d.UDH)+len(d.UD) != 0 {
		b[0] = b[0] | 0x04
	}
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}

	b = encodeSCTimeStamp(d.SCTS)
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}

	if d.PID != nil {
		b = []byte{*d.PID}
		if n, e = writeBytes(w, n, b); e != nil {
			return
		}
	}
	if d.DCS != nil {
		b = []byte{d.DCS.encode()}
		if n, e = writeBytes(w, n, b); e != nil {
			return
		}
	}

	if len(d.UDH)+len(d.UD) != 0 {
		udh := encodeUDH(d.UDH)
		u := d.DCS.unitSize()
		l := len(udh) + len(d.UD)
		l = ((l * 8) - (l * 8 % u)) / u
		b = []byte{byte(l)}
		if n, e = writeBytes(w, n, b); e != nil {
			return
		}
		if n, e = writeBytes(w, n, udh); e != nil {
			return
		}
		n, e = writeBytes(w, n, d.UD)
	}
	return
}

func (d *SubmitReport) readFrom(h byte, r io.Reader) (n int64, e error) {
	b := make([]byte, 1)
	if n, e = readBytes(r, n, b); e != nil {
		return
	}
	if b[0]&0x80 == 0x80 {
		*d.FCS = b[0]
		if n, e = readBytes(r, n, b); e != nil {
			return
		}
	}
	pi := b[0]

	b = make([]byte, 7)
	if n, e = readBytes(r, n, b); e != nil {
		return
	}
	d.SCTS = decodeSCTimeStamp([7]byte{b[0], b[1], b[2], b[3], b[4], b[5], b[6]})

	b = make([]byte, 1)
	if pi&0x01 == 0x01 {
		if n, e = readBytes(r, n, b); e != nil {
			return
		}
		d.PID = &b[0]
	}
	if pi&0x02 == 0x02 {
		if n, e = readBytes(r, n, b); e != nil {
			return
		}
		d.DCS = decodeDCS(b[0])
		if d.DCS == nil {
			e = fmt.Errorf("invalid TP-DCS data: % x", b[0])
			return
		}
	}
	if pi&0x04 == 0x04 {
		if d.DCS == nil {
			d.DCS = &GeneralDataCoding{
				AutoDelete: false,
				Compressed: false,
				MsgClass:   NoMessageClass,
				Charset:    GSM7bitAlphabet}
		}
		if n, e = readBytes(r, n, b); e != nil {
			return
		}
		l := d.DCS.unitSize()
		l *= int(b[0])
		if l%8 != 0 {
			l += 8 - l%8
		}

		d.UD = make([]byte, l/8)
		if n, e = readBytes(r, n, d.UD); e != nil {
			return
		}

		if h&0x40 == 0x40 {
			d.UDH = decodeUDH(d.UD[0 : d.UD[0]+1])
			d.UD = d.UD[d.UD[0]+1:]
		}
	}
	return
}

// PrintStack show PDU parameter
func (d *SubmitReport) PrintStack(w io.Writer) {
	fmt.Fprintf(w, "SMS message stack: Submit Report")
	if d.FCS != nil {
		fmt.Fprintf(w, " for RP-ERROR\n")
		v, ok := fcsStr[*d.FCS]
		if !ok {
			v = fmt.Sprintf("Reserved(%d)", *d.FCS)
		}
		fmt.Fprintf(w, "TP-FCS:  %s\n", v)
	} else {
		fmt.Fprintf(w, " for RP-ACK\n")
	}

	fmt.Fprintf(w, "TP-SCTS: %s\n", d.SCTS)
	if d.PID != nil {
		fmt.Fprintf(w, "TP-PID:  %d\n", *d.PID)
	}
	if d.DCS != nil {
		fmt.Fprintf(w, "TP-DCS:  %s\n", d.DCS)
	}
	if len(d.UDH)+len(d.UD) != 0 {
		fmt.Fprintf(w, "TP-UD:\n")
		for _, h := range d.UDH {
			fmt.Fprintf(w, "%s\n", h)
		}
		if len(d.UD) != 0 {
			fmt.Fprintf(w, "%s\n", d.DCS.decodeData(d.UD))
		}
	}
}
