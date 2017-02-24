package sms

import (
	"fmt"
	"io"
	"time"
)

var (
	msgRef byte
)

func init() {
	msgRef = byte(time.Now().Nanosecond())
}

// TPDU represents a SMS PDU
type TPDU interface {
	WriteTo(w io.Writer) (n int64, e error)
	readFrom(h byte, r io.Reader) (n int64, e error)
	PrintStack(w io.Writer)
}

// ReadAsSM parse byte data to TPDU as SM side
func ReadAsSM(r io.Reader) (t TPDU, n int64, e error) {
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

	if n, e = t.readFrom(h[0], r); e != nil {
		return
	}
	n++

	return
}

// ReadAsSC parse byte data to TPDU as SC side
func ReadAsSC(r io.Reader) (t TPDU, n int64, e error) {
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
		t = &Submit{}
	case 0x02:
		//		t = &Command{}
	case 0x03:
		e = fmt.Errorf("invalid data: reserved TPDU type")
		return
	}

	if n, e = t.readFrom(h[0], r); e != nil {
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

*/

func writeBytes(w io.Writer, n int64, b []byte) (int64, error) {
	i, e := w.Write(b)
	n += int64(i)
	return n, e
}

func encodeSCTimeStamp(t time.Time) (r []byte) {
	r = make([]byte, 7)
	r[0] = int2SemiOctet(t.Year())
	r[1] = int2SemiOctet(int(t.Month()))
	r[2] = int2SemiOctet(t.Day())
	r[3] = int2SemiOctet(t.Hour())
	r[4] = int2SemiOctet(t.Minute())
	r[5] = int2SemiOctet(t.Second())

	_, z := t.Zone()
	z /= 900
	r[6] = byte((z % 10) & 0x0f)
	r[6] = (r[6] << 4) | byte(((z/10)%10)&0x0f)
	if z < 0 {
		r[6] = r[6] | 0x08
	}
	return
}

func decodeSCTimeStamp(t [7]byte) time.Time {
	d := [6]int{}
	for i := range d {
		d[i] = semiOctet2Int(t[i])
	}
	l := semiOctet2Int(t[6] & 0x7f)
	if t[6]&0x80 == 0x80 {
		l = -l
	}
	return time.Date(2000+d[0], time.Month(d[1]), d[2], d[3], d[4], d[5], 0,
		time.FixedZone("unknown", l*15*60))
}

func int2SemiOctet(i int) (b byte) {
	b = byte(i % 10)
	b = (b << 4) | byte((i/10)%10)
	return
}

func semiOctet2Int(b byte) (i int) {
	i = int(b & 0x0f)
	i = (i * 10) + int((b&0xf0)>>4)
	return
}

func mmsStat(b bool) string {
	if b {
		return "More messages are waiting"
	}
	return "No more messages are waiting"
}
func lpStat(b bool) string {
	if b {
		return "Forwarded/spawned message"
	}
	return "Not forwarded/spawned message"
}
func sriStat(b bool) string {
	if b {
		return "Status report shall be returned"
	}
	return "Status report shall not be returned"
}
func rpStat(b bool) string {
	if b {
		return "Reply Path is set"
	}
	return "Reply Path is not set"
}
func rdStat(b bool) string {
	if b {
		return "Reject duplicate SUBMIT"
	}
	return "Accept duplicate SUBMIT"
}
func srrStat(b bool) string {
	if b {
		return "Status report is requested"
	}
	return "Status report is not requested"
}
