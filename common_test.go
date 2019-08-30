package sms_test

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"

	"github.com/fkgi/sms"
)

func randTransactionID() byte {
	return byte(rand.Int31n(16))
}

func randBool() bool {
	return rand.Int31n(2) == 0
}

func randByte() byte {
	return byte(rand.Int31n(256))
}

func randDate() time.Time {
	z := int(rand.Int31n(105) - 48)
	return time.Date(
		2000+rand.Int()%100, time.Month(rand.Int()%12+1),
		rand.Int()%32, rand.Int()%24, rand.Int()%60, rand.Int()%60, 0,
		time.FixedZone("", z*15*60))
}

func randVP() sms.ValidityPeriod {
	switch rand.Int31n(4) {
	case 1:
		return sms.VPRelative(randByte())
	case 2:
		t := time.Now()
		mx := time.Date(
			2099, time.December, 31, 23, 59, 54, 0, t.Location())
		t = t.Add(time.Duration(rand.Int63n(int64(mx.Sub(t)))) + time.Second*5)
		var r [7]byte
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
		return sms.VPAbsolute(r)
	case 3:
		return sms.VPEnhanced{
			randByte(), randByte(), randByte(), randByte(), randByte(), randByte(), randByte()}
	}
	return nil
}

func int2SemiOctet(i int) (b byte) {
	b = byte(i % 10)
	b = (b << 4) | byte((i/10)%10)
	return
}

func randText(len int) string {
	var b bytes.Buffer
	for i := 0; i < len; i++ {
		c := 0x1b
		for code[c] == '\x00' || code[c] == '\x1b' {
			c = rand.Int() % (128 + 16)
		}
		switch c {
		case 0x80:
			i++
		case 0x84:
			i++
		case 0x85:
			i++
		case 0x88:
			i++
		case 0x89:
			i++
		case 0x8a:
			i++
		case 0x8c:
			i++
		case 0x8d:
			i++
		case 0x8e:
			i++
		case 0x8f:
			i++
		}
		b.WriteRune(code[c])
	}
	return b.String()
}

func randDigit(len int) string {
	var b bytes.Buffer
	for i := 0; i < len; i++ {
		b.WriteString(strconv.Itoa(rand.Int() % 10))
	}
	return b.String()
}
