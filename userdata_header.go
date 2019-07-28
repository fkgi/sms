package sms

import (
	"bytes"
	"fmt"
	"reflect"
)

// UDH is user data header
type UDH interface {
	Marshal() []byte
	// unmarshal([]byte)
	fmt.Stringer
	Key() byte
	Value() []byte
	Equal(UDH) bool
}

// UnmarshalUDHs make UDHs
func UnmarshalUDHs(b []byte) (h []UDH) {
	if len(b) == 0 {
		return
	}
	buf := bytes.NewBuffer(b[1:])
	for buf.Len() != 0 {
		k, _ := buf.ReadByte()
		l, _ := buf.ReadByte()
		v := make([]byte, l)
		buf.Read(v)
		switch k {
		case 0x00:
			u := UnmarshalConcatenatedSM(v)
			h = append(h, u)
		case 0x08:
			u := UnmarshalConcatenatedSM16bit(v)
			h = append(h, u)
		default:
			u := UnmarshalGeneric(v)
			u.K = k
			h = append(h, u)
		}
	}
	return
}

// MarshalUDHs ganerate binary data of this UDHs
func MarshalUDHs(h []UDH) []byte {
	if len(h) == 0 {
		return []byte{}
	}

	var b bytes.Buffer
	b.WriteByte(0x00)
	for _, u := range h {
		b.Write(u.Marshal())
	}
	r := b.Bytes()
	r[0] = byte(len(r) - 1)
	return r
}

// GenericIEI is User Data Header
type GenericIEI struct {
	K byte   `json:"key"`
	V []byte `json:"value"`
}

// Equal reports a and b are same
func (h GenericIEI) Equal(b UDH) bool {
	a, ok := b.(GenericIEI)
	if !ok {
		return false
	}
	if a.K != h.K {
		return false
	}
	return reflect.DeepEqual(h.V, a.V)
}

// Key of this IEI
func (h GenericIEI) Key() byte {
	return h.K
}

// Value of this IEI
func (h GenericIEI) Value() []byte {
	if h.V == nil {
		return []byte{}
	}
	return h.V
}

// Marshal generate binary data of this UDH
func (h GenericIEI) Marshal() []byte {
	r := make([]byte, len(h.V)+2)
	r[0] = h.K
	r[1] = byte(len(h.V))
	for i := range h.V {
		r[i+2] = h.V[i]
	}
	return r
}

// UnmarshalGeneric make Generic UDH
func UnmarshalGeneric(b []byte) (h GenericIEI) {
	if b != nil {
		h.V = b
	}
	return
}

func (h GenericIEI) String() string {
	return fmt.Sprintf("Generic(0x%x): % x", h.K, h.V)
}

// ConcatenatedSM is User Data Header
type ConcatenatedSM struct {
	RefNum byte
	MaxNum byte
	SeqNum byte
}

// Equal reports a and b are same
func (h ConcatenatedSM) Equal(b UDH) bool {
	a, ok := b.(ConcatenatedSM)
	if !ok {
		return false
	}
	if a.RefNum != h.RefNum {
		return false
	}
	if a.MaxNum != h.MaxNum {
		return false
	}
	if a.SeqNum != h.SeqNum {
		return false
	}
	return true
}

// Key of this IEI
func (h ConcatenatedSM) Key() byte {
	return 0x00
}

// Value of this IEI
func (h ConcatenatedSM) Value() []byte {
	return []byte{h.RefNum, h.MaxNum, h.SeqNum}
}

// Marshal generate binary data of this UDH
func (h ConcatenatedSM) Marshal() []byte {
	return []byte{0x00, 0x03, h.RefNum, h.MaxNum, h.SeqNum}
}

// UnmarshalConcatenatedSM make ConcatenatedSM UDH
func UnmarshalConcatenatedSM(b []byte) (h ConcatenatedSM) {
	if b != nil && len(b) >= 3 {
		h.RefNum = b[0]
		h.MaxNum = b[1]
		h.SeqNum = b[2]
	}
	return
}

func (h ConcatenatedSM) String() string {
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

// Equal reports a and b are same
func (h ConcatenatedSM16bit) Equal(b UDH) bool {
	a, ok := b.(ConcatenatedSM16bit)
	if !ok {
		return false
	}
	if a.RefNum != h.RefNum {
		return false
	}
	if a.MaxNum != h.MaxNum {
		return false
	}
	if a.SeqNum != h.SeqNum {
		return false
	}
	return true
}

// Key of this IEI
func (h ConcatenatedSM16bit) Key() byte {
	return 0x08
}

// Value of this IEI
func (h ConcatenatedSM16bit) Value() []byte {
	return []byte{byte(h.RefNum >> 8), byte(h.RefNum & 0x00ff), h.MaxNum, h.SeqNum}
}

// Marshal generate binary data of this UDH
func (h ConcatenatedSM16bit) Marshal() []byte {
	return []byte{0x08, 0x04,
		byte(h.RefNum >> 8), byte(h.RefNum & 0x00ff), h.MaxNum, h.SeqNum}
}

// UnmarshalConcatenatedSM16bit make ConcatenatedSM16bit UDH
func UnmarshalConcatenatedSM16bit(b []byte) (h ConcatenatedSM16bit) {
	if b != nil && len(b) >= 4 {
		h.RefNum = (uint16(b[0]) << 8) | uint16(b[1])
		h.MaxNum = b[2]
		h.SeqNum = b[3]
	}
	return
}

func (h ConcatenatedSM16bit) String() string {
	return fmt.Sprintf(
		"Concatenated SM (16bit ref number): Ref=%d, Max=%d, Seq=%d",
		h.RefNum, h.MaxNum, h.SeqNum)
}
