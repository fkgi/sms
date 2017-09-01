package sms

import (
	"bytes"
	"fmt"
)

type udh interface {
	encode() []byte
	decode([]byte)
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
		default:
			u = &GenericIEI{Key: k}
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

// GenericIEI is User Data Header
type GenericIEI struct {
	Key   byte
	Value []byte
}

func (h *GenericIEI) encode() []byte {
	r := make([]byte, len(h.Value)+2)
	r[0] = h.Key
	r[1] = byte(len(h.Value))
	for i := range h.Value {
		r[i+2] = h.Value[i]
	}
	return r
}

func (h *GenericIEI) decode(b []byte) {
	h.Value = b
}

func (h *GenericIEI) String() string {
	return fmt.Sprintf("Generic(%x): % x", h.Key, h.Value)
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
		"Concatenated SM: Ref=%d, Max=%d, Seq=%d",
		h.RefNum, h.MaxNum, h.SeqNum)
}

// ConcatenatedSM16bit is User Data Header
type ConcatenatedSM16bit struct {
	RefNum uint16
	MaxNum byte
	SeqNum byte
}

func (h *ConcatenatedSM16bit) encode() []byte {
	return []byte{0x08, 0x04,
		byte(h.RefNum >> 8), byte(h.RefNum & 0x00ff), h.MaxNum, h.SeqNum}
}

func (h *ConcatenatedSM16bit) decode(b []byte) {
	h.RefNum = (uint16(b[0]) << 8) | uint16(b[1])
	h.MaxNum = b[2]
	h.SeqNum = b[3]
}

func (h *ConcatenatedSM16bit) String() string {
	return fmt.Sprintf(
		"Concatenated SM (16bit ref number): Ref=%d, Max=%d, Seq=%d",
		h.RefNum, h.MaxNum, h.SeqNum)
}
