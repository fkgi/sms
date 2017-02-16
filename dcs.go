package sms

import (
	"bytes"
	"fmt"
	"unicode/utf16"
)

type dcs interface {
	encodeDCS() (b byte)
	unitSize() int
	String() string
	decodeData(b []byte) string
}

func decodeDCS(b byte) dcs {
	if b&0xc0 == 0x00 {
		return &GeneralDataCoding{
			AutoDelete: false,
			Compressed: b&0x20 == 0x20,
			MsgClass:   msgClass(b & 0x13),
			Charset:    charset(b & 0x0c)}
	}
	if b&0xc0 == 0x40 {
		return &GeneralDataCoding{
			AutoDelete: true,
			Compressed: b&0x20 == 0x20,
			MsgClass:   msgClass(b & 0x13),
			Charset:    charset(b & 0x0c)}
	}
	if b&0xf0 == 0xc0 {
		return &MessageWaiting{
			Behavior:    DiscardMessageGSM7bit,
			Active:      b&0x08 == 0x08,
			WaitingType: waitType(b & 0x03)}
	}
	if b&0xf0 == 0xd0 {
		return &MessageWaiting{
			Behavior:    StoreMessageGSM7bit,
			Active:      b&0x08 == 0x08,
			WaitingType: waitType(b & 0x03)}
	}
	if b&0xf0 == 0xe0 {
		return &MessageWaiting{
			Behavior:    StoreMessageUCS2,
			Active:      b&0x08 == 0x08,
			WaitingType: waitType(b & 0x03)}
	}
	if b&0xf0 == 0xf0 {
		return &DataCodingMessage{
			IsData:   b&0x04 == 0x04,
			MsgClass: msgClass((b & 0x03) | 0x10)}
	}
	return nil
}

type msgClass byte
type charset byte

const (
	// NoMessageClass means no message class
	NoMessageClass msgClass = 0x00
	// MessageClass0 means message class 0
	MessageClass0 msgClass = 0x10
	// MessageClass1 means message class 1
	MessageClass1 msgClass = 0x11
	// MessageClass2 means message class 2
	MessageClass2 msgClass = 0x12
	// MessageClass3 means message class 3
	MessageClass3 msgClass = 0x13

	// GSM7bitAlphabet means GSM 7 bit default alphabet charset
	GSM7bitAlphabet charset = 0x00
	// Data8bit means 8 bit data charset
	Data8bit charset = 0x04
	// UCS2 means UCS2 charset
	UCS2 charset = 0x08
)

// GeneralDataCoding is group of SMS Data Coding Scheme
type GeneralDataCoding struct {
	AutoDelete bool
	Compressed bool
	MsgClass   msgClass
	Charset    charset
}

func (s *GeneralDataCoding) encodeDCS() (b byte) {
	if s.AutoDelete {
		b = 0x40
	} else {
		b = 0x00
	}
	if s.Compressed {
		b = b | 0x20
	}
	b = b | byte(s.MsgClass)
	b = b | byte(s.Charset&0x0c)
	return
}

func (s *GeneralDataCoding) unitSize() int {
	if s.Charset == GSM7bitAlphabet {
		return 7
	}
	return 8
}

func (s *GeneralDataCoding) decodeData(b []byte) string {
	switch s.Charset {
	case GSM7bitAlphabet:
		return GSM7bitString(b).String()
	case Data8bit:
		return fmt.Sprintf("% x", b)
	case UCS2:
		return decodeUCS2(b)
	}
	return ""
}

func (s *GeneralDataCoding) String() string {
	var b bytes.Buffer
	b.WriteString("General Data Coding")
	if s.AutoDelete {
		b.WriteString("(Automatic Deletion)")
	}
	if s.Compressed {
		b.WriteString(", compressed")
	}
	switch s.MsgClass {
	case NoMessageClass:
		b.WriteString(", no class")
	case MessageClass0:
		b.WriteString(", class 0")
	case MessageClass1:
		b.WriteString(", class 1")
	case MessageClass2:
		b.WriteString(", class 2")
	case MessageClass3:
		b.WriteString(", class 3")
	}
	switch s.Charset {
	case GSM7bitAlphabet:
		b.WriteString(", GSM 7bit default alphabet")
	case Data8bit:
		b.WriteString(", 8 bit data")
	case UCS2:
		b.WriteString(", UCS2 (16bit)")
	}
	return b.String()
}

type waitType byte
type waitBehavior byte

const (
	// VoicemailMessageWaiting means waiting type
	VoicemailMessageWaiting waitType = 0x00
	// FaxMessageWaiting means waiting type
	FaxMessageWaiting waitType = 0x01
	// ElectronicMailMessageWaiting means waiting type
	ElectronicMailMessageWaiting waitType = 0x02
	// OtherMessageWaiting means waiting type
	OtherMessageWaiting waitType = 0x03

	// DiscardMessageGSM7bit means discard the contents
	DiscardMessageGSM7bit waitBehavior = 0x00
	// StoreMessageGSM7bit means store the contents GSM 7bit alphabet
	StoreMessageGSM7bit waitBehavior = 0x10
	// StoreMessageUCS2 means store the contents UCS2
	StoreMessageUCS2 waitBehavior = 0x20
)

// MessageWaiting is group of SMS Data Coding Scheme
type MessageWaiting struct {
	Behavior    waitBehavior
	Active      bool
	WaitingType waitType
}

func (s *MessageWaiting) encodeDCS() (b byte) {
	b = 0xc0
	b = b | byte(s.Behavior&0xc0)
	if s.Active {
		b = b | 0x80
	}
	b = b | byte(s.WaitingType&0x03)
	return
}

func (s *MessageWaiting) unitSize() int {
	if s.Behavior == StoreMessageUCS2 {
		return 8
	}

	return 7
}

func (s *MessageWaiting) decodeData(b []byte) string {
	if s.Behavior == StoreMessageUCS2 {
		return decodeUCS2(b)
	}
	return GSM7bitString(b).String()
}

func (s *MessageWaiting) String() string {
	var b bytes.Buffer
	b.WriteString("MessageWaiting")
	switch s.Behavior {
	case DiscardMessageGSM7bit:
		b.WriteString("(Discard Message with GSM 7bit default alphabet)")
	case StoreMessageGSM7bit:
		b.WriteString("(Store Message with GSM 7bit default alphabet)")
	case StoreMessageUCS2:
		b.WriteString("(Store Message with UCS2)")
	}
	if s.Active {
		b.WriteString(", active")
	} else {
		b.WriteString(", inactive")
	}
	switch s.WaitingType {
	case VoicemailMessageWaiting:
		b.WriteString(", Voicemail Message")
	case FaxMessageWaiting:
		b.WriteString(", Fax Message")
	case ElectronicMailMessageWaiting:
		b.WriteString(", Electronic Mail Message")
	case OtherMessageWaiting:
		b.WriteString(", Other Message")
	}
	return b.String()
}

// DataCodingMessage is group of SMS Data Coding Scheme
type DataCodingMessage struct {
	IsData   bool
	MsgClass msgClass
}

func (s *DataCodingMessage) encodeDCS() (b byte) {
	b = 0xf0
	if s.IsData {
		b = b | 0x40
	}
	b = b | byte(s.MsgClass&0x03)
	return
}

func (s *DataCodingMessage) unitSize() int {
	if s.IsData {
		return 8
	}
	return 7
}

func (s *DataCodingMessage) decodeData(b []byte) string {
	if s.IsData {
		return fmt.Sprintf("% x", b)
	}
	return GSM7bitString(b).String()
}

func (s *DataCodingMessage) String() string {
	var b bytes.Buffer
	b.WriteString("Data coding/message")
	if s.IsData {
		b.WriteString(", 8-bit data")
	} else {
		b.WriteString(", GSM 7 bit default alphabet")
	}
	switch s.MsgClass {
	case MessageClass0:
		b.WriteString(", class 0")
	case MessageClass1:
		b.WriteString(", class 1")
	case MessageClass2:
		b.WriteString(", class 2")
	case MessageClass3:
		b.WriteString(", class 3")
	}
	return b.String()
}

func encodeUCS2(s string) []byte {
	u := utf16.Encode([]rune(s))
	b := make([]byte, len(u)*2)
	for i, c := range u {
		b[i*2] = byte((c >> 8) & 0xff)
		b[i*2+1] = byte(c & 0xff)
	}
	return b
}

func decodeUCS2(b []byte) string {
	u := make([]uint16, len(b)/2)
	for i := range u {
		u[i] = uint16(b[2*i])<<8 | uint16(b[2*i+1])
	}
	return string(utf16.Decode(u))
}
