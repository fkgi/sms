package sms

import (
	"bytes"
	"fmt"
	"time"
)

// ValidityPeriod is type of validity period
type ValidityPeriod interface {
	fmt.Stringer
	ExpireTime(t time.Time) time.Time
	Duration() time.Duration
	SingleAttempt() bool
	Equal(ValidityPeriod) bool
}

type jvp struct {
	T time.Duration `json:"duration"`
	S bool          `json:"single"`
}

// ValidityPeriodOf returns VP from deadend time and single-attempt flag
func ValidityPeriodOf(t time.Duration, s bool) ValidityPeriod {
	if t%time.Second != 0 {
		t = t.Truncate(time.Second)
	}

	if s {
		vp := VPEnhanced{}
		vp[0] = 0x40
		if t == 0 {
		} else if t%(time.Hour*24*7) == 0 && t >= time.Hour*24*7*5 && t <= time.Hour*24*7*63 {
			vp[0] |= 0x01
			vp[1] = byte(t/(time.Hour*24*7)) + 192
		} else if t%(time.Hour*24) == 0 && t >= time.Hour*24*2 && t <= time.Hour*24*30 {
			vp[0] |= 0x01
			vp[1] = byte(t/(time.Hour*24)) + 166
		} else if t%(time.Minute*30) == 0 && t <= time.Hour*24 && t >= time.Hour*12+time.Minute*30 {
			vp[0] |= 0x01
			vp[1] = byte((t-time.Hour*12)/(time.Minute*30)) + 143
		} else if t%(time.Minute*5) == 0 && t <= time.Hour*12 && t >= time.Minute*5 {
			vp[0] |= 0x01
			vp[1] = byte(t/(time.Minute*5)) - 1
		} else if t <= time.Second*255 {
			vp[0] |= 0x02
			vp[1] = byte(t / time.Second)
		} else if t <= time.Hour*99+time.Minute*59+time.Second*59 {
			vp[0] |= 0x03
			vp[1] = int2SemiOctet(int(t / time.Hour))
			vp[2] = int2SemiOctet(int((t % time.Hour) / time.Minute))
			vp[3] = int2SemiOctet(int((t % time.Minute) / time.Second))
			//} else if t <= time.Hour*24*30 {
			//	vp[0] |= 0x01
			//	vp[1] = byte(t/(time.Hour*24)) + 166
		}
		return vp
	}

	if t == 0 {
		return VPEnhanced{}
	} else if t%(time.Hour*24*7) == 0 && t >= time.Hour*24*7*5 && t <= time.Hour*24*7*63 {
		return VPRelative(byte(t/(time.Hour*24*7)) + 192)
	} else if t%(time.Hour*24) == 0 && t >= time.Hour*24*2 && t <= time.Hour*24*30 {
		return VPRelative(byte(t/(time.Hour*24)) + 166)
	} else if t%(time.Minute*30) == 0 && t <= time.Hour*24 && t >= time.Hour*12+time.Minute*30 {
		return VPRelative(byte((t-time.Hour*12)/(time.Minute*30)) + 143)
	} else if t%(time.Minute*5) == 0 && t <= time.Hour*12 && t >= time.Minute*5 {
		return VPRelative(byte(t/(time.Minute*5)) - 1)
	} else if t <= time.Second*255 {
		vp := VPEnhanced{}
		vp[0] = 0x02
		vp[1] = byte(t / time.Second)
		return vp
	} else if t <= time.Hour*99+time.Minute*59+time.Second*59 {
		vp := VPEnhanced{}
		vp[0] = 0x03
		vp[1] = int2SemiOctet(int(t / time.Hour))
		vp[2] = int2SemiOctet(int((t % time.Hour) / time.Minute))
		vp[3] = int2SemiOctet(int((t % time.Minute) / time.Second))
		return vp
	}

	vp := marshalSCTimeStamp(time.Now().Add(t))
	return VPAbsolute{vp[0], vp[1], vp[2], vp[3], vp[4], vp[5], vp[6]}
}

// VPRelative is relative format VP value
type VPRelative byte

func (f VPRelative) String() string {
	return "ralative, " + relativeFormatString(byte(f))
}

func relativeFormatString(b byte) string {
	if b < 144 {
		i := int(b+1) * 5
		return fmt.Sprintf("%d:%02d:00", (i-i%60)/60, i%60)
	}
	if b < 168 {
		i := int(b-143) * 30
		return fmt.Sprintf("%d:%02d:00", (i-i%60)/60+12, i%60)
	}
	if b < 197 {
		return fmt.Sprintf("%d days", b-166)
	}
	return fmt.Sprintf("%d weeks", b-192)
}

// ExpireTime return expire time
func (f VPRelative) ExpireTime(t time.Time) time.Time {
	return t.Add(relativeFormatDuration(byte(f)))
}

// Duration return duration to expire time
func (f VPRelative) Duration() time.Duration {
	return relativeFormatDuration(byte(f))
}

func relativeFormatDuration(b byte) time.Duration {
	if b < 144 {
		return time.Duration(b+1) * 5 * time.Minute
	}
	if b < 168 {
		return time.Duration(b-143)*30*time.Minute + time.Hour*12
	}
	if b < 197 {
		return time.Duration(b-166) * time.Hour * 24
	}
	return time.Duration(b-192) * time.Hour * 24 * 7
}

// SingleAttempt return single attempt is required or not
func (f VPRelative) SingleAttempt() bool {
	return false
}

// Equal reports a and b are same
func (f VPRelative) Equal(b ValidityPeriod) bool {
	a, ok := b.(VPRelative)
	if !ok {
		return false
	}
	return a == f
}

// VPAbsolute is absolute format VP value
type VPAbsolute [7]byte

func (f VPAbsolute) String() string {
	return fmt.Sprintf("absolute, %s", unmarshalSCTimeStamp(f))
}

// ExpireTime return expire time
func (f VPAbsolute) ExpireTime(t time.Time) time.Time {
	return unmarshalSCTimeStamp(f)
}

// Duration return duration to expire time
func (f VPAbsolute) Duration() time.Duration {
	return time.Until(unmarshalSCTimeStamp(f))
}

// SingleAttempt return single attempt is required or not
func (f VPAbsolute) SingleAttempt() bool {
	return false
}

// Equal reports a and b are same
func (f VPAbsolute) Equal(b ValidityPeriod) bool {
	a, ok := b.(VPAbsolute)
	if !ok {
		return false
	}
	for i := range f {
		if a[i] != f[i] {
			return false
		}
	}
	return true
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
		s.WriteString(fmt.Sprintf(", %d:%02d:%02d",
			semiOctet2Int(f[1]),
			semiOctet2Int(f[2])%60,
			semiOctet2Int(f[3])%60))
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
		i += time.Duration(semiOctet2Int(f[2])%60) * time.Minute
		i += time.Duration(semiOctet2Int(f[3])%60) * time.Second
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
		i += time.Duration(semiOctet2Int(f[2])%60) * time.Minute
		i += time.Duration(semiOctet2Int(f[3])%60) * time.Second
		return i
	}
	return time.Duration(0)
}

// SingleAttempt return single attempt is required or not
func (f VPEnhanced) SingleAttempt() bool {
	if f[0]&0x40 == 0x40 {
		return true
	}
	return false
}

// Equal reports a and b are same
func (f VPEnhanced) Equal(b ValidityPeriod) bool {
	a, ok := b.(VPEnhanced)
	if !ok {
		return false
	}
	for i := range f {
		if a[i] != f[i] {
			return false
		}
	}
	return true
}
