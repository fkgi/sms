package sms

import (
	"bytes"
	"time"
	"unicode/utf16"
)

var (
	msgRef = byte(0)
)

func init() {
	msgRef = byte(time.Now().Nanosecond())
}

type TPDU struct {
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

func (p TPDU) Encode() []byte {
	switch p.MTI {
	case 0x00:
		if p.Req {
			return p.encodeDeliver()
		} else {
			return p.encodeDeliverReport()
		}
	case 0x01:
		if p.Req {
			return p.encodeSubmit()
		} else {
			return p.encodeSubmitReport()
		}
	}
	return nil
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

type TPAddr struct {
	EXT   bool
	TON   byte
	NPI   byte
	Digit string
}

func (p TPAddr) encode() []byte {
	var buf bytes.Buffer

	b := byte(len(p.Digit))
	buf.WriteByte(b)

	if p.EXT {
		b = 0x80
	} else {
		b = 0x00
	}
	b = b | (p.TON&0x07)<<4
	b = b | (p.NPI & 0x0f)
	buf.WriteByte(b)

	buf.Write(stotbcd(p.Digit))

	return buf.Bytes()
}

func getTime(t time.Time) []byte {
	r := make([]byte, 7)

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

func stotbcd(s string) []byte {
	if len(s)%2 != 0 {
		s = s + " "
	}
	r := make([]byte, len(s)/2)
	for i, c := range s {
		v := byte(0x0f)
		switch c {
		case '0':
			v = 0x00
		case '1':
			v = 0x01
		case '2':
			v = 0x02
		case '3':
			v = 0x03
		case '4':
			v = 0x04
		case '5':
			v = 0x05
		case '6':
			v = 0x06
		case '7':
			v = 0x07
		case '8':
			v = 0x08
		case '9':
			v = 0x09
		case '*':
			v = 0x09
		case '#':
			v = 0x09
		case 'a', 'A':
			v = 0x09
		case 'b', 'B':
			v = 0x09
		case 'c', 'C':
			v = 0x09
		default:
			v = 0x0f
		}
		if i%2 == 1 {
			v = v << 4
		}
		r[i/2] = r[i/2] | v
	}
	return r
}
