package sms

import (
	"bytes"
	"fmt"
)

type udh interface {
	encode() []byte
	decode([]byte)
	fmt.Stringer
	Key() byte
	Value() []byte
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
			u = &GenericIEI{K: k}
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
	K byte
	V []byte
}

// Key of this IEI
func (h *GenericIEI) Key() byte {
	if h == nil {
		return 0x00
	}
	return h.K
}

// Value of this IEI
func (h *GenericIEI) Value() []byte {
	if h == nil || h.V == nil {
		return []byte{}
	}
	return h.V
}

func (h *GenericIEI) encode() []byte {
	if h == nil {
		return []byte{}
	}

	r := make([]byte, len(h.V)+2)
	r[0] = h.K
	r[1] = byte(len(h.V))
	for i := range h.V {
		r[i+2] = h.V[i]
	}
	return r
}

func (h *GenericIEI) decode(b []byte) {
	if h != nil && b != nil {
		h.V = b
	}
}

func (h *GenericIEI) String() string {
	if h == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Generic(%x): % x", h.K, h.V)
}

// ConcatenatedSM is User Data Header
type ConcatenatedSM struct {
	RefNum byte
	MaxNum byte
	SeqNum byte
}

// Key of this IEI
func (h *ConcatenatedSM) Key() byte {
	return 0x00
}

// Value of this IEI
func (h *ConcatenatedSM) Value() []byte {
	if h == nil {
		return []byte{}
	}
	return []byte{h.RefNum, h.MaxNum, h.SeqNum}
}

func (h *ConcatenatedSM) encode() []byte {
	if h == nil {
		return []byte{}
	}
	return []byte{0x00, 0x03, h.RefNum, h.MaxNum, h.SeqNum}
}

func (h *ConcatenatedSM) decode(b []byte) {
	if h != nil && b != nil && len(b) > 3 {
		h.RefNum = b[0]
		h.MaxNum = b[1]
		h.SeqNum = b[2]
	}
}

func (h *ConcatenatedSM) String() string {
	if h == nil {
		return "<nil>"
	}
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

// Key of this IEI
func (h *ConcatenatedSM16bit) Key() byte {
	return 0x08
}

// Value of this IEI
func (h *ConcatenatedSM16bit) Value() []byte {
	if h == nil {
		return []byte{}
	}
	return []byte{byte(h.RefNum >> 8), byte(h.RefNum & 0x00ff), h.MaxNum, h.SeqNum}
}

func (h *ConcatenatedSM16bit) encode() []byte {
	if h == nil {
		return []byte{}
	}
	return []byte{0x08, 0x04,
		byte(h.RefNum >> 8), byte(h.RefNum & 0x00ff), h.MaxNum, h.SeqNum}
}

func (h *ConcatenatedSM16bit) decode(b []byte) {
	if h != nil && b != nil && len(b) > 4 {
		h.RefNum = (uint16(b[0]) << 8) | uint16(b[1])
		h.MaxNum = b[2]
		h.SeqNum = b[3]
	}
}

func (h *ConcatenatedSM16bit) String() string {
	if h == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"Concatenated SM (16bit ref number): Ref=%d, Max=%d, Seq=%d",
		h.RefNum, h.MaxNum, h.SeqNum)
}
