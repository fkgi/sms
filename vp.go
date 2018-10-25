package sms

import (
	"bytes"
	"fmt"
	"time"
)

type vp interface {
	fmt.Stringer
	ExpireTime(t time.Time) time.Time
	Duration() time.Duration
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

// ExpireTime return expire time
func (f VPRelative) ExpireTime(t time.Time) time.Time {
	return t.Add(relativeFormatDuration(f[0]))
}

// Duration return duration to expire time
func (f VPRelative) Duration() time.Duration {
	return relativeFormatDuration(f[0])
}

func relativeFormatDuration(b byte) time.Duration {
	if b < 144 {
		return time.Duration(b+1) * 5 * time.Minute
	}
	if b < 168 {
		return time.Duration(b-143) * 30 * time.Minute
	}
	if b < 197 {
		return time.Duration(b-166) * time.Hour * 24
	}
	return time.Duration(b-192) * time.Hour * 24 * 7
}

// VPAbsolute is absolute format VP value
type VPAbsolute [7]byte

func (f VPAbsolute) String() string {
	return fmt.Sprintf("absolute, %s", decodeSCTimeStamp(f))
}

// ExpireTime return expire time
func (f VPAbsolute) ExpireTime(t time.Time) time.Time {
	return decodeSCTimeStamp(f)
}

// Duration return duration to expire time
func (f VPAbsolute) Duration() time.Duration {
	return time.Until(decodeSCTimeStamp(f))
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
			semiOctet2Int(f[1]),
			semiOctet2Int(f[2]),
			semiOctet2Int(f[3])))
	default:
		s.WriteString(", invalid format")
	}
	return s.String()
}

// ExpireTime return expire time
func (f VPEnhanced) ExpireTime(t time.Time) time.Time {
	switch f[0] & 0x03 {
	case 0x00:
		return time.Time{}
	case 0x01:
		return t.Add(relativeFormatDuration(f[1]))
	case 0x02:
		return t.Add(time.Duration(f[1]) * time.Second)
	case 0x03:
		i := time.Duration(semiOctet2Int(f[1])) * time.Hour
		i += time.Duration(semiOctet2Int(f[2])) * time.Minute
		i += time.Duration(semiOctet2Int(f[3])) * time.Second
		return t.Add(i)
	}
	return time.Time{}
}

// Duration return duration to expire time
func (f VPEnhanced) Duration() time.Duration {
	switch f[0] & 0x03 {
	case 0x00:
		return time.Duration(0)
	case 0x01:
		return relativeFormatDuration(f[1])
	case 0x02:
		return time.Duration(f[1]) * time.Second
	case 0x03:
		i := time.Duration(semiOctet2Int(f[1])) * time.Hour
		i += time.Duration(semiOctet2Int(f[2])) * time.Minute
		i += time.Duration(semiOctet2Int(f[3])) * time.Second
		return i
	}
	return time.Duration(0)
}
