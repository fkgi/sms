package sms

import (
	"fmt"
	"io"
	"time"
)

// StatusReport is TPDU message from SC to MS
type StatusReport struct {
	MMS bool // More Messages to Send (true=more messages)
	LP  bool // Loop Prevention
	SRQ bool // Status Report Qualifier (true=status report shall be returned)

	MR   byte      // Message Reference
	RA   Address   // Destination Address
	SCTS time.Time // Service Centre Time Stamp
	DT   time.Time // Discharge Time
	ST   byte      // Status
	PID  *byte     // Protocol Identifier
	DCS  dcs       // Data Coding Scheme
	UDH  []udh     // User Data Header
	UD   []byte    // User Data
}

// WriteTo output byte data of this TPDU
func (d *StatusReport) WriteTo(w io.Writer) (n int64, e error) {
	b := []byte{0x02, d.MR}
	if !d.MMS {
		b[0] = b[0] | 0x04
	}
	if d.LP {
		b[0] = b[0] | 0x08
	}
	if d.SRQ {
		b[0] = b[0] | 0x20
	}
	if d.UDH != nil && len(d.UDH) != 0 {
		b[0] = b[0] | 0x40
	}
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}

	var nn int64
	nn, e = d.RA.WriteTo(w)
	n += nn
	if e != nil {
		return
	}

	b = encodeSCTimeStamp(d.SCTS)
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}

	b = encodeSCTimeStamp(d.DT)
	if n, e = writeBytes(w, n, b); e != nil {
		return
	}

	b = []byte{d.ST}
	if n, e = writeBytes(w, n, b); e != nil {
		return
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

func (d *StatusReport) readFrom(h byte, r io.Reader) (n int64, e error) {
	d.MMS = h&0x04 != 0x04
	d.LP = h&0x08 == 0x08
	d.SRQ = h&0x20 == 0x20

	b := make([]byte, 1)
	if n, e = readBytes(r, n, b); e != nil {
		return
	}
	d.MR = b[0]

	d.RA = Address{}
	if n, e = d.RA.ReadFrom(r); e != nil {
		return
	}

	b = make([]byte, 16)
	if n, e = readBytes(r, n, b); e != nil {
		return
	}
	d.SCTS = decodeSCTimeStamp(
		[7]byte{b[0], b[1], b[2], b[3], b[4], b[5], b[6]})
	d.DT = decodeSCTimeStamp(
		[7]byte{b[7], b[8], b[9], b[10], b[11], b[12], b[13]})
	d.ST = b[14]
	pi := b[15]

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
func (d *StatusReport) PrintStack(w io.Writer) {
	fmt.Fprintf(w, "SMS message stack: Status Report\n")
	fmt.Fprintf(w, "TP-MMS:  %s\n", mmsStat(d.MMS))
	fmt.Fprintf(w, "TP-LP:   %s\n", lpStat(d.LP))
	fmt.Fprintf(w, "TP-SRQ:  %s\n", srqStat(d.SRQ))

	fmt.Fprintf(w, "TP-MR:   %d\n", d.MR)
	fmt.Fprintf(w, "TP-RA:   %s\n", d.RA)
	fmt.Fprintf(w, "TP-SCTS: %s\n", d.SCTS)
	fmt.Fprintf(w, "TP-DT:   %s\n", d.SCTS)
	fmt.Fprintf(w, "TP-ST:   %d\n", d.ST)
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
