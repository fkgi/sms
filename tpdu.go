package sms

import (
	"fmt"
	"io"
	"time"
)

var (
	msgRef = byte(0)
)

func init() {
	msgRef = byte(time.Now().Nanosecond())
}

// TPDU represents a SMS PDU
type TPDU interface {
	WriteTo(w io.Writer) (n int64, e error)
	ReadFrom(h byte, r io.Reader) (n int64, e error)
}

// ParseAsSM parse byte data to TPDU as SM side
func ParseAsSM(r io.Reader) (t TPDU, n int64, e error) {
	h := make([]byte, 1)
	i := 0
	if i, e = r.Read(h); e != nil {
		return
	} else if i != 1 {
		e = fmt.Errorf("no data")
		return
	}

	switch h[0] & 0x03 {
	case 0x00:
		t = &Deliver{}
	case 0x01:
		//		t = &SubmitReport{}
	case 0x02:
		//		t = &StatusReport{}
	case 0x03:
		e = fmt.Errorf("invalid data: reserved TPDU type")
		return
	}

	if n, e = t.ReadFrom(h[0], r); e != nil {
		return
	}
	n++

	return
}

// ParseAsSC parse byte data to TPDU as SC side
func ParseAsSC(r io.Reader) (t TPDU, n int64, e error) {
	h := make([]byte, 1)
	i := 0
	if i, e = r.Read(h); e != nil {
		return
	} else if i != 1 {
		e = fmt.Errorf("no data")
		return
	}

	switch h[0] & 0x03 {
	case 0x00:
		//		t = &DeliverReport{}
	case 0x01:
		//		t = &Submit{}
	case 0x02:
		//		t = &Command{}
	case 0x03:
		e = fmt.Errorf("invalid data: reserved TPDU type")
		return
	}

	if n, e = t.ReadFrom(h[0], r); e != nil {
		return
	}
	n++

	return
}

/*
	Req bool

	MTI  byte    // Message Type Indicator
	MMS  bool    // More Messages to Send
	RD   bool    // Reject Duplicates
	LP   bool    // Loop Prevention
	VPF  byte    // Validity Period Format
	SRI  bool    // Status Report Indication
	SRR  bool    // Status Report Request
	SRQ  bool    // Status Report Qualifier
	UDHI bool    // User Data Header Indicator
	RP   bool    // Reply Path
	FCS  byte    // Failure Cause
	MR   byte    // Message Reference
	DA   TPAddr  // Destination Address
	OA   TPAddr  // Originating Address
	RA   TPAddr  // Recipient Address
	SCTS [7]byte // Service Centre Time Stamp
	DT   [7]byte // Discharge Time
	ST   byte    // Status
	PI   byte    // Parameter Indicator
	PID  byte    // Protocol Identifier
	DCS  byte    // Data Coding Scheme
	VP   []byte  // Validity Period
	CDL  byte    // User Data Length
	UD   string  // User Data
	CT   byte    // Command Type
	MN   byte    // Message Number
	CD   []byte  // Command Data

	// User Data
	UDh []byte
}

func (p TPDU) encodeUD() []byte {
	s := p.UD
	if len(s)*2 > 140-len(p.UDh) {
		s = s[0 : (140-len(p.UDh))/2]
	} else if len(s) == 0 {
		return make([]byte, 0)
	}

	u := utf16.Encode([]rune(s))
	b := make([]byte, len(u)*2)
	for i, c := range u {
		b[i*2] = byte((c >> 8) & 0xff)
		b[i*2+1] = byte(c & 0xff)
	}
	return b
}
*/

// Address is SMS originator/destination address
type Address struct {
	EXT   bool
	TON   byte
	NPI   byte
	Digit addrValue
}

type addrValue interface {
	Length() int
	String() string
	WriteTo(w io.Writer) (n int64, e error)
}

// WriteTo wite binary data to io.Writer
func (a Address) WriteTo(w io.Writer) (n int64, e error) {
	i := 0
	b := []byte{byte(a.Digit.Length()), 0x00}

	if a.EXT {
		b[1] = 0x80
	} else {
		b[1] = 0x00
	}
	b[1] = b[1] | (a.TON&0x07)<<4
	b[1] = b[1] | (a.NPI & 0x0f)
	if i, e = w.Write(b); e != nil {
		return
	}

	if n, e = a.Digit.WriteTo(w); e != nil {
		return
	}
	n += int64(i)

	return
}

// ReadFrom read byte data and set parameter of the Address
func (a *Address) ReadFrom(r io.Reader) (n int64, e error) {
	i := 0
	b := make([]byte, 2)
	if i, e = r.Read(b); e != nil {
		return
	} else if i != 2 {
		e = fmt.Errorf("more data required")
		return
	}

	l := int(b[0])
	a.EXT = b[1]&0x80 == 0x80
	a.TON = (b[1] >> 4) & 0x07
	a.NPI = b[1] & 0x0f

	if l%2 == 1 {
		l++
	}
	l = l / 2
	b = make([]byte, l)
	if i, e = r.Read(b); e != nil {
		return
	} else if i != l {
		e = fmt.Errorf("more data required")
		return
	}
	n = int64(i + 2)

	return
}

// DateTime is time data for TP-SCTS, TP-DT and in Absolute format of TP-VP
type DateTime [7]byte

// EncodeTime create DateTime
func EncodeTime(t time.Time) DateTime {
	var r [7]byte

	r[0] = byte(t.Year() % 10)
	r[0] = (r[0] << 4) | byte((t.Year()/10)%10)
	r[1] = byte(t.Month() % 10)
	r[1] = (r[1] << 4) | byte((t.Month()/10)%10)
	r[2] = byte(t.Day() % 10)
	r[2] = (r[2] << 4) | byte((t.Day()/10)%10)
	r[3] = byte(t.Hour() % 10)
	r[3] = (r[3] << 4) | byte((t.Hour()/10)%10)
	r[4] = byte(t.Minute() % 10)
	r[4] = (r[4] << 4) | byte((t.Minute()/10)%10)
	r[5] = byte(t.Second() % 10)
	r[5] = (r[5] << 4) | byte((t.Second()/10)%10)

	_, z := t.Zone()
	z /= 900
	r[6] = byte((z % 10) & 0x0f)
	r[6] = (r[6] << 4) | byte(((z/10)%10)&0x0f)
	if z < 0 {
		r[6] = r[6] | 0x08
	}
	return r
}
