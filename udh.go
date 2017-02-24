package sms

import (
	"bytes"
	"fmt"
)

type udh interface {
	encode() []byte
	decode(b []byte)
	String() string
}

func decodeUDH(b []byte) (h []udh) {
	if len(b) == 0 {
		return
	}
	buf := bytes.NewBuffer(b[1:])
	for buf.Len() != 0 {
		var u udh
		k, _ := buf.ReadByte()
		switch k {
		case 0x00:
			u = &ConcatenatedSM{}
		}
		l, _ := buf.ReadByte()
		v := make([]byte, l)
		buf.Read(v)
		u.decode(v)
		h = append(h, u)
	}
	return
}

func encodeUDH(h []udh) []byte {
	if len(h) == 0 {
		return []byte{}
	}

	var b bytes.Buffer
	b.WriteByte(0x00)
	for _, u := range h {
		b.Write(u.encode())
	}
	r := b.Bytes()
	r[0] = byte(len(r) - 1)
	return r
}

// ConcatenatedSM is User Data Header
type ConcatenatedSM struct {
	RefNum byte
	MaxNum byte
	SeqNum byte
}

func (h *ConcatenatedSM) encode() []byte {
	return []byte{0x00, 0x03, h.RefNum, h.MaxNum, h.SeqNum}
}

func (h *ConcatenatedSM) decode(b []byte) {
	h.RefNum = b[0]
	h.MaxNum = b[1]
	h.SeqNum = b[2]
}

func (h *ConcatenatedSM) String() string {
	return fmt.Sprintf(
		"Concatenated SM: ReferenceNum=%d, MaximumNum=%d, SequenceNum=%d",
		h.RefNum, h.MaxNum, h.SeqNum)
}
