package sms

import (
	"bytes"
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
	Encode() []byte
	decode([]byte) error
	PrintStack(io.Writer)
}

// DecodeAsSC parse byte data to TPDU as SC.
func DecodeAsSC(b []byte) (t TPDU, e error) {
	return decode(b, true)
}

// DecodeAsMS parse byte data to TPDU as MS.
func DecodeAsMS(b []byte) (t TPDU, e error) {
	return decode(b, false)
}

func decode(b []byte, sc bool) (t TPDU, e error) {
	if len(b) == 0 {
		e = fmt.Errorf("invalid data")
	} else if sc {
		switch b[0] & 0x03 {
		case 0x00:
			t = &DeliverReport{}
		case 0x01:
			t = &Submit{}
		case 0x02:
			// t = &Command{}
		case 0x03:
			e = fmt.Errorf("invalid data: reserved TPDU type")
		}
	} else {
		switch b[0] & 0x03 {
		case 0x00:
			t = &Deliver{}
		case 0x01:
			t = &SubmitReport{}
		case 0x02:
			t = &StatusReport{}
		case 0x03:
			e = fmt.Errorf("invalid data: reserved TPDU type")
		}
	}

	if e == nil {
		e = t.decode(b)
	}
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

func readDCS(r *bytes.Reader) (dcs, error) {
	p, e := r.ReadByte()
	if e != nil {
		return nil, e
	}
	d := decodeDCS(p)
	if d == nil {
		return nil, fmt.Errorf("invalid TP-DCS data: % x", p)
	}
	return d, nil
}

func read7Bytes(r *bytes.Reader) ([7]byte, error) {
	if r.Len() < 7 {
		return [7]byte{}, io.EOF
	}
	b := make([]byte, 7)
	r.Read(b)
	return [7]byte{
		b[0], b[1], b[2], b[3], b[4], b[5], b[6]}, nil
}

func readUD(r *bytes.Reader, d dcs, h bool) ([]byte, []udh, error) {
	p, e := r.ReadByte()
	if e != nil {
		return nil, nil, e
	}
	l := d.unitSize()
	l *= int(p)
	if l%8 != 0 {
		l += 8 - l%8
	}

	ud := make([]byte, l/8)
	if r.Len() < len(ud) {
		return nil, nil, io.EOF
	}
	r.Read(ud)

	if h {
		return ud[ud[0]+1:], decodeUDH(ud[0 : ud[0]+1]), nil
	}
	return ud, nil, nil
}

func writeUD(w *bytes.Buffer, ud []byte, h []udh, d dcs) {
	udh := encodeUDH(h)

	u := d.unitSize()
	l := len(udh) + len(ud)
	l = ((l * 8) - (l * 8 % u)) / u

	w.WriteByte(byte(l))
	w.Write(udh)
	w.Write(ud)
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
	return time.Date(2000+d[0],
		time.Month(d[1]), d[2], d[3], d[4], d[5], 0,
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
func fcsStat(b byte) string {
	switch b {
	case 0x80:
		return "Telematic interworking not supported"
	case 0x81:
		return "Short message Type 0 not supported"
	case 0x82:
		return "Cannot replace short message"
	case 0x8F:
		return "Unspecified TP-PID error"
	case 0x90:
		return "Data coding scheme (alphabet) not supported"
	case 0x91:
		return "Message class not supported"
	case 0x9F:
		return "Unspecified TP-DCS error"
	case 0xA0:
		return "Command cannot be actioned"
	case 0xA1:
		return "Command unsupported"
	case 0xAF:
		return "Unspecified TP-Command error"
	case 0xB0:
		return "TPDU not supported"
	case 0xC0:
		return "SC busy"
	case 0xC1:
		return "No SC subscription"
	case 0xC2:
		return "SC system failure"
	case 0xC3:
		return "Invalid SME address"
	case 0xC4:
		return "Destination SME barred"
	case 0xC5:
		return "SM Rejected-Duplicate SM"
	case 0xC6:
		return "TP-VPF not supported"
	case 0xC7:
		return "TP-VP not supported"
	case 0xD0:
		return "(U)SIM SMS storage full"
	case 0xD1:
		return "No SMS storage capability in (U)SIM"
	case 0xD2:
		return "Error in MS"
	case 0xD3:
		return "Memory Capacity Exceeded"
	case 0xD4:
		return "(U)SIM Application Toolkit Busy"
	case 0xD5:
		return "(U)SIM data download error"
	case 0xFF:
		return "Unspecified error cause"
	}
	return fmt.Sprintf("Reserved(%d)", b)
}
func stStat(b byte) string {
	switch b {
	case 0x00:
		return "Short message received by the SME"
	case 0x01:
		return "Short message forwarded by the SC to the SME" +
			" but the SC is unable to confirm delivery"
	case 0x02:
		return "Short message replaced by the SC"
	case 0x20:
		return "Congestion"
	case 0x21:
		return "SME busy"
	case 0x22:
		return "No response from SME"
	case 0x23:
		return "Service rejected"
	case 0x24:
		return "Quality of service not available"
	case 0x25:
		return "Error in SME"
	case 0x40:
		return "Remote procedure error"
	case 0x41:
		return "Incompatible destination"
	case 0x42:
		return "Connection rejected by SME"
	case 0x43:
		return "Not obtainable"
	case 0x44:
		return "Quality of service not available"
	case 0x45:
		return "No interworking available"
	case 0x46:
		return "SM Validity Period Expired"
	case 0x47:
		return "SM Deleted by originating SME"
	case 0x48:
		return "SM Deleted by SC Administration"
	case 0x49:
		return "SM does not exist"
	case 0x60:
		return "Congestion"
	case 0x61:
		return "SME busy"
	case 0x62:
		return "No response from SME"
	case 0x63:
		return "Service rejected"
	case 0x64:
		return "Quality of service not available"
	case 0x65:
		return "Error in SME"
	}
	return fmt.Sprintf("Reserved(%d)", b)
}
func pidStat(b byte) string {
	switch b {
	case 0:
		return "Default store and forward short message"
	case 32:
		return "implicit telemetic device"
	case 33:
		return "Telex or teletex reduced to telex format"
	case 34:
		return "Group 3 telefax"
	case 35:
		return "Group 4 telefax"
	case 36:
		return "Voice telephone"
	case 37:
		return "ERMES (European Radio Messaging System)"
	case 38:
		return "National Paging system (known to the SC)"
	case 39:
		return "Videotex (T.100 [20] /T.101 [21])"
	case 40:
		return "Teletex, carrier unspecified"
	case 41:
		return "Teletex, in PSPDN"
	case 42:
		return "Teletex, in CSPDN"
	case 43:
		return "Teletex, in analog PSTN"
	case 44:
		return "Teletex, in digital ISDN"
	case 45:
		return "UCI (Universal Computer Interface, ETSI DE/PS 3 01 3)"
	case 48:
		return "A message handling facility (known to the SC)"
	case 49:
		return "Any public X.400 based message handling system"
	case 50:
		return "Internet Electronic Mail"
	case 63:
		return "A GSM/UMTS mobile station"
	case 64:
		return "Short Message Type 0"
	case 65:
		return "Replace Short Message Type 1"
	case 66:
		return "Replace Short Message Type 2"
	case 67:
		return "Replace Short Message Type 3"
	case 68:
		return "Replace Short Message Type 4"
	case 69:
		return "Replace Short Message Type 5"
	case 70:
		return "Replace Short Message Type 6"
	case 71:
		return "Replace Short Message Type 7"
	case 72:
		return "Device Triggering Short Message"
	case 94:
		return "Enhanced Message Service"
	case 95:
		return "Return Call Message"
	case 124:
		return "ANSI-136 R-DATA"
	case 125:
		return "ME Data download"
	case 126:
		return "ME De personalization Short Message"
	case 127:
		return "(U)SIM Data download"
	}
	if b > 0 && b < 32 {
		return "no telematic interworking, but SME to SME protocol"
	}
	if b > 55 && b < 63 {
		return "SC specific use"
	}
	if b > 191 && b <= 255 {
		return "SC specific use"
	}

	return fmt.Sprintf("Reserved(%d)", b)
}
