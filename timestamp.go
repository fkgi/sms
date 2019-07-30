package sms

import (
	"time"
)

func marshalSCTimeStamp(t time.Time) (r []byte) {
	r = make([]byte, 7)
	r[0] = int2SemiOctet(t.Year())
	r[1] = int2SemiOctet(int(t.Month()))
	r[2] = int2SemiOctet(t.Day())
	r[3] = int2SemiOctet(t.Hour())
	r[4] = int2SemiOctet(t.Minute())
	r[5] = int2SemiOctet(t.Second())

	_, z := t.Zone()
	z /= 900
	if z < 0 {
		z = -z
		r[6] = byte(z % 10)
		r[6] = (r[6] << 4) | byte(((z/10)%10)&0x07)
		r[6] = r[6] | 0x08
	} else {
		r[6] = byte(z % 10)
		r[6] = (r[6] << 4) | byte(((z/10)%10)&0x07)
	}
	return
}

func unmarshalSCTimeStamp(t [7]byte) time.Time {
	d := [6]int{}
	for i := range d {
		d[i] = semiOctet2Int(t[i])
	}
	l := semiOctet2Int(t[6] & 0xf7)
	if t[6]&0x08 == 0x08 {
		l = -l
	}
	return time.Date(2000+d[0],
		time.Month(d[1]), d[2], d[3], d[4], d[5], 0,
		time.FixedZone("", l*900))
}
