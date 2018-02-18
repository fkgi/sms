package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf16"
)

// UD is TP-UD
type UD struct {
	Text string
	UDH  []udh
}

type jud struct {
	Text []byte `json:"text"`
	UDH  []judh `json:"hdr"`
}
type judh struct {
	Key byte   `json:"key"`
	Val []byte `json:"value"`
}

func (u UD) String(d DCS) string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "%sTP-UD:", Indent)

	for _, h := range u.UDH {
		fmt.Fprintf(w, "\n%s%s%s", Indent, Indent, h)
	}
	if len(u.Text) != 0 {
		switch d.charset() {
		case CharsetGSM7bit, CharsetUCS2:
			fmt.Fprintf(w, "\n%s%s%s", Indent, Indent, u.Text)
		case Charset8bitData:
			fmt.Fprintf(w, "\n%s%s% x",
				Indent, Indent, []byte(u.Text))
		}
	}
	return w.String()
}

// MarshalJSON provide custom marshaller
func (u UD) MarshalJSON() ([]byte, error) {
	ud := jud{
		Text: []byte(u.Text),
		UDH:  make([]judh, len(u.UDH))}
	for i, h := range u.UDH {
		ud.UDH[i] = judh{Key: h.Key(), Val: h.Value()}
	}
	return json.Marshal(ud)
}

// UnmarshalJSON provide custom marshaller
func (u *UD) UnmarshalJSON(b []byte) (e error) {
	ud := jud{}
	if e = json.Unmarshal(b, &ud); e != nil {
		return
	}
	u.Text = string(ud.Text)
	u.UDH = make([]udh, 0)
	for _, h := range ud.UDH {
		switch h.Key {
		case 0x00:
			t := &ConcatenatedSM{}
			t.decode(h.Val)
			u.UDH = append(u.UDH, t)
		default:
			t := &GenericIEI{K: h.Key}
			t.decode(h.Val)
			u.UDH = append(u.UDH, t)
		}
	}
	return nil
}

// AddUDH add additional User-Data header
func (u *UD) AddUDH(h udh) {
	if u.UDH == nil {
		u.UDH = []udh{h}
	} else {
		u.UDH = append(u.UDH, h)
	}
}

func (u *UD) read(r *bytes.Reader, d DCS, h bool) error {
	p, e := r.ReadByte()
	if e != nil {
		return e
	}

	c := d.charset()
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
		u.UDH = decodeUDH(ud[0 : ud[0]+1])
		ud = ud[ud[0]+1:]
	}

	switch c {
	case CharsetGSM7bit:
		s := GSM7bitString(make([]rune, l))
		s.decode(o, ud)
		u.Text = s.String()
	case Charset8bitData:
		u.Text = string(ud)
	case CharsetUCS2:
		s := make([]uint16, l/2)
		for i := range s {
			s[i] = uint16(ud[2*i])<<8 | uint16(ud[2*i+1])
		}
		u.Text = string(utf16.Decode(s))
	}
	return nil
}

func (u *UD) write(w *bytes.Buffer, d DCS) {
	c := d.charset()
	udh := encodeUDH(u.UDH)
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
		ud = s.encode(o)
		l += s.Length()
	case Charset8bitData:
		ud = []byte(u.Text)
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

// MakeSeparatedText generate splited data
func MakeSeparatedText(s string, c msgClass, id byte) (
	ud []UD, dcs GeneralDataCoding) {
	ud = []UD{}
	dcs = GeneralDataCoding{
		AutoDelete: false,
		Compressed: false,
		MsgClass:   c}

	dcs.Charset = CharsetGSM7bit
	for _, r := range s {
		_, code := getCode(r)
		if code == 0xff {
			dcs.Charset = CharsetUCS2
			break
		}
	}

	if dcs.Charset == CharsetGSM7bit {
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
			ud[i].AddUDH(&ConcatenatedSM{
				RefNum: id,
				MaxNum: byte(len(ud)),
				SeqNum: byte(i + 1)})
		}
	}

	return
}
