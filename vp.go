package sms

import (
	"bytes"
	"fmt"
	"time"
)

type vp interface {
	fmt.Stringer
	Period(t time.Time) time.Time
}

// VPRelative is relative format VP value
type VPRelative [1]byte

func (f VPRelative) String() string {
	return "ralative, " + relativeFormatString(f[0])
}
func relativeFormatString(b byte) string {
	if b < 144 {
		i := int(b+1) * 5
		return fmt.Sprintf("%d:%d", (i-i%60)/60, i%60)
	}
	if b < 168 {
		i := int(b-143) * 30
		return fmt.Sprintf("%d:%d", (i-i%60)/60+12, i%60)
	}
	if b < 197 {
		return fmt.Sprintf("%d days", b-166)
	}
	return fmt.Sprintf("%d weeks", b-192)
}

// Period return period time
func (f VPRelative) Period(t time.Time) time.Time {
	return relativeFormatPeriod(t, f[0])
}

func relativeFormatPeriod(t time.Time, b byte) time.Time {
	if b < 144 {
		i := time.Duration(b+1) * 5 * time.Minute
		return t.Add(i)
	}
	if b < 168 {
		i := time.Duration(b-143) * 30 * time.Minute
		return t.Add(i)
	}
	if b < 197 {
		return t.AddDate(0, 0, int(b-166))
	}
	return t.AddDate(0, 0, int(b-192)*7)
}

// VPAbsolute is absolute format VP value
type VPAbsolute [7]byte

func (f VPAbsolute) String() string {
	return fmt.Sprintf("absolute, %s", decodeSCTimeStamp(f))
}

// Period return period time
func (f VPAbsolute) Period(t time.Time) time.Time {
	return decodeSCTimeStamp(f)
}

// VPEnhanced is enhanced format VP value
type VPEnhanced [7]byte

func (f VPEnhanced) String() string {
	var s bytes.Buffer
	s.WriteString("enhanced")
	if f[0]&0x40 == 0x40 {
		s.WriteString(", single-shot")
	}
	switch f[0] & 0x03 {
	case 0x00:
		s.WriteString(", no validity period")
	case 0x01:
		s.WriteString(", ")
		s.WriteString(relativeFormatString(f[1]))
	case 0x02:
		s.WriteString(fmt.Sprintf(", %d sec", f[1]))
	case 0x03:
		s.WriteString(fmt.Sprintf(", %d:%d:%d",
			semiOctet2Int(f[1]), semiOctet2Int(f[2]), semiOctet2Int(f[3])))
	default:
		s.WriteString(", invalid format")
	}
	return s.String()
}

// Period return period time
func (f VPEnhanced) Period(t time.Time) time.Time {
	switch f[0] & 0x03 {
	case 0x00:
		return t
	case 0x01:
		return relativeFormatPeriod(t, f[1])
	case 0x02:
		return t.Add(time.Duration(f[1]) * time.Second)
	case 0x03:
		j := int(f[1]&0x0f)*10 + int((f[1]&0xf0)>>4)
		i := time.Duration(j) * time.Hour
		j = int(f[2]&0x0f)*10 + int((f[2]&0xf0)>>4)
		i += time.Duration(j) * time.Minute
		j = int(f[3]&0x0f)*10 + int((f[3]&0xf0)>>4)
		i += time.Duration(j) * time.Second
		return t.Add(i)
	}
	return t
}
