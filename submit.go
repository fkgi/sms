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
	if d.UDH != nil && len(d.UDH) != 0 {
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

	i := 0
	b := make([]byte, 1)
	if i, e = r.Read(b); e != nil {
		return
	} else if i != len(b) {
		e = fmt.Errorf("more data required")
		return
	}
	d.MR = b[0]

	d.DA = Address{}
	if n, e = d.DA.ReadFrom(r); e != nil {
		return
	}

	b = make([]byte, 2)
	if i, e = r.Read(b); e != nil {
		return
	} else if i != len(b) {
		e = fmt.Errorf("more data required")
		return
	}
	n += int64(i) + 1
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
		if i, e = r.Read(b); e != nil {
			return
		} else if i != len(b) {
			e = fmt.Errorf("more data required")
			return
		}
		n++
		d.VP = VPRelative([1]byte{b[0]})
	case 0x08:
		b = make([]byte, 7)
		if i, e = r.Read(b); e != nil {
			return
		} else if i != len(b) {
			e = fmt.Errorf("more data required")
			return
		}
		n += int64(i) + 1
		d.VP = VPEnhanced([7]byte{b[0], b[1], b[2], b[3], b[4], b[5], b[6]})
	case 0x18:
		b = make([]byte, 7)
		if i, e = r.Read(b); e != nil {
			return
		} else if i != len(b) {
			e = fmt.Errorf("more data required")
			return
		}
		n += int64(i) + 1
		d.VP = VPAbsolute([7]byte{b[0], b[1], b[2], b[3], b[4], b[5], b[6]})
	}

	b = make([]byte, 1)
	if i, e = r.Read(b); e != nil {
		return
	} else if i != len(b) {
		e = fmt.Errorf("more data required")
		return
	}
	l := d.DCS.unitSize()
	l *= int(b[0])
	if l%8 != 0 {
		l += 8 - l%8
	}

	d.UD = make([]byte, l/8)
	if i, e = r.Read(d.UD); e != nil {
		return
	} else if i != len(d.UD) {
		e = fmt.Errorf("more data required")
		return
	}
	n += int64(i)

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

	fmt.Fprintf(w, "TP-UD:\n")
	for _, h := range d.UDH {
		fmt.Fprintf(w, " %s\n", h)
	}
	if d.UD != nil && len(d.UD) != 0 {
		fmt.Fprintf(w, "%s\n", d.DCS.decodeData(d.UD))
	}
}

// SubmitReport is TPDU message from SC to MS
type SubmitReport struct {
	FCS  byte            // Failure Cause
	PI   byte            // Parameter Indicator
	SCTS time.Time       // Service Centre Time Stamp
	PID  byte            // Protocol Identifier
	DCS  dcs             // Data Coding Scheme
	UDH  map[byte][]byte // User Data Header
	UD   []byte          // User Data
}
