package sms

import (
	"fmt"
	"io"
	"time"
)

var (
	msgRef byte
	fcsStr map[byte]string
	stStr  map[byte]string
)

func init() {
	msgRef = byte(time.Now().Nanosecond())
	fcsStr = map[byte]string{
		0x80: "Telematic interworking not supported",
		0x81: "Short message Type 0 not supported",
		0x82: "Cannot replace short message",
		0x8F: "Unspecified TP-PID error",
		0x90: "Data coding scheme (alphabet) not supported",
		0x91: "Message class not supported",
		0x9F: "Unspecified TP-DCS error",
		0xA0: "Command cannot be actioned",
		0xA1: "Command unsupported",
		0xAF: "Unspecified TP-Command error",
		0xB0: "TPDU not supported",
		0xC0: "SC busy",
		0xC1: "No SC subscription",
		0xC2: "SC system failure",
		0xC3: "Invalid SME address",
		0xC4: "Destination SME barred",
		0xC5: "SM Rejected-Duplicate SM",
		0xC6: "TP-VPF not supported",
		0xC7: "TP-VP not supported",
		0xD0: "(U)SIM SMS storage full",
		0xD1: "No SMS storage capability in (U)SIM",
		0xD2: "Error in MS",
		0xD3: "Memory Capacity Exceeded",
		0xD4: "(U)SIM Application Toolkit Busy",
		0xD5: "(U)SIM data download error",
		0xFF: "Unspecified error cause"}
	stStr = map[byte]string{
		0x00: "Short message received by the SME",
		0x01: "Short message forwarded by the SC to the SME but the SC is unable to confirm delivery",
		0x02: "Short message replaced by the SC",
		0x20: "Congestion",
		0x21: "SME busy",
		0x22: "No response from SME",
		0x23: "Service rejected",
		0x24: "Quality of service not available",
		0x25: "Error in SME",
		0x40: "Remote procedure error",
		0x41: "Incompatible destination",
		0x42: "Connection rejected by SME",
		0x43: "Not obtainable",
		0x44: "Quality of service not available",
		0x45: "No interworking available",
		0x46: "SM Validity Period Expired",
		0x47: "SM Deleted by originating SME",
		0x48: "SM Deleted by SC Administration",
		0x49: "SM does not exist",
		0x60: "Congestion",
		0x61: "SME busy",
		0x62: "No response from SME",
		0x63: "Service rejected",
		0x64: "Quality of service not available",
		0x65: "Error in SME"}
}

// TPDU represents a SMS PDU
type TPDU interface {
	WriteTo(w io.Writer) (n int64, e error)
	readFrom(h byte, r io.Reader) (n int64, e error)
	PrintStack(w io.Writer)
}

// Read parse byte data to TPDU.
// r is input byte stream.
// sc is true when decode data as SC, false when decode as MS,
func Read(r io.Reader, sc bool) (t TPDU, n int64, e error) {
	h := make([]byte, 1)
	if n, e = readBytes(r, n, h); e != nil {
		return
	}

	switch h[0] & 0x03 {
	case 0x00:
		if !sc {
			t = &Deliver{}
		} else {
			t = &DeliverReport{}
		}
	case 0x01:
		if sc {
			t = &Submit{}
		} else {
			t = &SubmitReport{}
		}
	case 0x02:
		if sc {
			// t = &Command{}
		} else {
			t = &StatusReport{}
		}
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

func writeBytes(w io.Writer, n int64, b []byte) (int64, error) {
	i, e := w.Write(b)
	n += int64(i)
	return n, e
}

func readBytes(r io.Reader, n int64, b []byte) (int64, error) {
	i, e := r.Read(b)
	n += int64(i)
	if e == nil && i != len(b) {
		e = fmt.Errorf("more data required")
	}
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
func srqStat(b bool) string {
	if b {
		return "This is result of an SMS-COMMAND"
	}
	return "This is result of a SMS-SUBMIT"
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
