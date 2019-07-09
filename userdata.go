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
	Text string `json:"text"`
	UDH  []udh  `json:"hdr"`
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

type judh struct {
	Key byte   `json:"key"`
	Val []byte `json:"value"`
}

// MarshalJSON provide custom marshaller
func (u UD) MarshalJSON() ([]byte, error) {
	ud := struct {
		Text string `json:"text"`
		UDH  []judh `json:"hdr"`
	}{
		Text: u.Text,
		UDH:  make([]judh, len(u.UDH))}
	for i, h := range u.UDH {
		ud.UDH[i] = judh{Key: h.Key(), Val: h.Value()}
	}
	return json.Marshal(ud)
}

// UnmarshalJSON provide custom marshaller
func (u *UD) UnmarshalJSON(b []byte) (e error) {
	type alias UD
	al := struct {
		UDH []judh `json:"hdr"`
		*alias
	}{
		alias: (*alias)(u)}
	if e = json.Unmarshal(b, &al); e != nil {
		return
	}
	u.UDH = make([]udh, len(al.UDH))
	for _, h := range al.UDH {
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
	if h == nil || u == nil {
		return
	}
	if u.UDH == nil {
		u.UDH = []udh{h}
	} else {
		u.UDH = append(u.UDH, h)
	}
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
		c = d.charset()
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
		u.UDH = decodeUDH(ud[0 : ud[0]+1])
		ud = ud[ud[0]+1:]
	}

	switch c {
	case CharsetGSM7bit:
		s := GSM7bitString(make([]rune, 0, l))
		s.decode(o, ud)
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

func (u *UD) write(w *bytes.Buffer, d DCS) {
	c := CharsetGSM7bit
	if d != nil {
		c = d.charset()
	}
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