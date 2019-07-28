package sms

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf16"
)

// UD is TP-UD
type UD struct {
	Text string `json:"text,omitempty"`
	UDH  []UDH  `json:"hdr,omitempty"`
}

func (u UD) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "%sTP-UD:", Indent)

	for _, h := range u.UDH {
		fmt.Fprintf(w, "\n%s%s%s", Indent, Indent, h)
	}
	if len(u.Text) != 0 {
		fmt.Fprintf(w, "\n%s%s%s", Indent, Indent, u.Text)
	}
	return w.String()
}

// Equal reports a and b are same
func (u UD) Equal(b UD) bool {
	if len(u.UDH) != len(b.UDH) {
		return false
	}
	for i := range u.UDH {
		if !u.UDH[i].Equal(b.UDH[i]) {
			return false
		}
	}
	return u.Text == b.Text
}

type judh struct {
	Key byte   `json:"key"`
	Val []byte `json:"value"`
}

// MarshalJSON provide custom marshaller
func (u UD) MarshalJSON() ([]byte, error) {
	type alias UD
	ud := struct {
		UDH []judh `json:"hdr,omitempty"`
		*alias
	}{
		UDH:   make([]judh, len(u.UDH)),
		alias: (*alias)(&u)}
	for i, h := range u.UDH {
		ud.UDH[i] = judh{Key: h.Key(), Val: h.Value()}
	}
	return json.Marshal(ud)
}

// UnmarshalJSON provide custom marshaller
func (u *UD) UnmarshalJSON(b []byte) (e error) {
	type alias UD
	al := struct {
		UDH []judh `json:"hdr,omitempty"`
		*alias
	}{
		alias: (*alias)(u)}
	if e = json.Unmarshal(b, &al); e != nil {
		return
	}
	u.UDH = make([]UDH, 0, len(al.UDH))
	for _, h := range al.UDH {
		switch h.Key {
		case 0x00:
			t := UnmarshalConcatenatedSM(h.Val)
			u.UDH = append(u.UDH, t)
		default:
			t := UnmarshalGeneric(h.Val)
			t.K = h.Key
			u.UDH = append(u.UDH, t)
		}
	}
	return nil
}

// Set8bitData set binary data as UD
func (u *UD) Set8bitData(d []byte) {
	if u != nil && len(d) != 0 {
		u.Text = base64.StdEncoding.EncodeToString(d)
	}
}

// Get8bitData set binary data as UD
func (u UD) Get8bitData() ([]byte, error) {
	return base64.StdEncoding.DecodeString(u.Text)
}

func (u *UD) read(r *bytes.Reader, d DCS, h bool) error {
	p, e := r.ReadByte()
	if e != nil {
		return e
	}

	c := CharsetGSM7bit
	if d != nil {
		c = d.Charset()
	}
	l := int(p)
	if c == CharsetGSM7bit {
		l *= 7
		if l%8 != 0 {
			l = l/8 + 1
		} else {
			l = l / 8
		}
	}
	ud := make([]byte, l)
	if r.Len() < len(ud) {
		return io.EOF
	}
	r.Read(ud)

	l = int(p)
	o := 0
	if h {
		if c == CharsetGSM7bit {
			o = int(ud[0]+1) * 8
			l -= o / 7
			o %= 7
			if o != 0 {
				o = 7 - o
				l--
			}
		} else {
			l -= int(ud[0] + 1)
		}
		u.UDH = UnmarshalUDHs(ud[0 : ud[0]+1])
		ud = ud[ud[0]+1:]
	}

	switch c {
	case CharsetGSM7bit:
		s := UnmarshalGSM7bitString(o, l, ud)
		u.Text = s.String()
	case Charset8bitData:
		u.Text = base64.StdEncoding.EncodeToString(ud)
	case CharsetUCS2:
		s := make([]uint16, l/2)
		for i := range s {
			s[i] = uint16(ud[2*i])<<8 | uint16(ud[2*i+1])
		}
		u.Text = string(utf16.Decode(s))
	}
	return nil
}

func (u UD) write(w *bytes.Buffer, d DCS) {
	c := CharsetGSM7bit
	if d != nil {
		c = d.Charset()
	}
	udh := MarshalUDHs(u.UDH)
	var ud []byte
	l := len(udh)

	switch c {
	case CharsetGSM7bit:
		o := l * 8
		l = o / 7
		o %= 7
		if o != 0 {
			o = 7 - o
			l++
		}
		s, _ := StringToGSM7bit(u.Text)
		ud = s.Marshal(o)
		l += s.septetLength()
	case Charset8bitData:
		var e error
		ud, e = base64.StdEncoding.DecodeString(u.Text)
		if e != nil {
			ud = []byte{}
		}
		l += len(ud)
	case CharsetUCS2:
		u := utf16.Encode([]rune(u.Text))
		ud = make([]byte, len(u)*2)
		for i, c := range u {
			ud[i*2] = byte((c >> 8) & 0xff)
			ud[i*2+1] = byte(c & 0xff)
		}
		l += len(ud)
	}

	w.WriteByte(byte(l))
	w.Write(udh)
	w.Write(ud)
}

func (u UD) isEmpty() bool {
	return len(u.Text) == 0 && len(u.UDH) == 0
}

// MakeSeparatedText generate splited data
func MakeSeparatedText(s string, c msgClass, id byte) (
	ud []UD, dcs GeneralDataCoding) {
	ud = []UD{}
	dcs = GeneralDataCoding{
		AutoDelete: false,
		Compressed: false,
		MsgClass:   c}

	dcs.MsgCharset = CharsetGSM7bit
	for _, r := range s {
		_, code := getCode(r)
		if code == 0xff {
			dcs.MsgCharset = CharsetUCS2
			break
		}
	}

	if dcs.MsgCharset == CharsetGSM7bit {
		r := []rune(s)
		maxlen := 160

		for len(r) > maxlen {
			ud = append(ud, UD{Text: string(r[:153])})
			r = r[153:]
			maxlen = 153
		}
		ud = append(ud, UD{Text: string(r)})
	} else {
		rs := make([]rune, 0, 70)
		maxlen := 140

		for _, r := range s {
			tmp := append(rs, r)
			if len(string(tmp)) > maxlen {
				ud = append(ud, UD{Text: string(rs)})
				rs = rs[:1]
				rs[0] = r
				maxlen = 134
			} else {
				rs = tmp
			}
		}
		ud = append(ud, UD{Text: string(rs)})
	}

	if len(ud) > 1 {
		for i := range ud {
			ud[i].UDH = append(ud[i].UDH, ConcatenatedSM{
				RefNum: id,
				MaxNum: byte(len(ud)),
				SeqNum: byte(i + 1)})
		}
	}

	return
}
