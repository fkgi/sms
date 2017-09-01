package sms

import "bytes"

type dcs interface {
	encode() byte
	String() string
	charset() charset
}

func decodeDCS(b byte) dcs {
	switch b & 0xc0 {
	case 0x00:
		return &GeneralDataCoding{
			AutoDelete: false,
			Compressed: b&0x20 == 0x20,
			MsgClass:   msgClass(b & 0x13),
			Charset:    charset(b & 0x0c)}
	case 0x40:
		return &GeneralDataCoding{
			AutoDelete: true,
			Compressed: b&0x20 == 0x20,
			MsgClass:   msgClass(b & 0x13),
			Charset:    charset(b & 0x0c)}
	}
	switch b & 0xf0 {
	case 0xc0:
		return &MessageWaiting{
			Behavior:    DiscardMessageGSM7bit,
			Active:      b&0x08 == 0x08,
			WaitingType: waitType(b & 0x03)}
	case 0xd0:
		return &MessageWaiting{
			Behavior:    StoreMessageGSM7bit,
			Active:      b&0x08 == 0x08,
			WaitingType: waitType(b & 0x03)}
	case 0xe0:
		return &MessageWaiting{
			Behavior:    StoreMessageUCS2,
			Active:      b&0x08 == 0x08,
			WaitingType: waitType(b & 0x03)}
	case 0xf0:
		return &DataCodingMessage{
			IsData:   b&0x04 == 0x04,
			MsgClass: msgClass((b & 0x03) | 0x10)}
	}
	return nil
}

func readDCS(r *bytes.Reader) (dcs, error) {
	p, e := r.ReadByte()
	if e != nil {
		return nil, e
	}
	d := decodeDCS(p)
	if d == nil {
		return nil, &InvalidDataError{
			Name:  "TP-DCS",
			Bytes: []byte{p}}
	}
	return d, nil
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

	// CharsetGSM7bit means GSM 7 bit default alphabet charset
	CharsetGSM7bit charset = 0x00
	// Charset8bitData means 8 bit data charset
	Charset8bitData charset = 0x04
	// CharsetUCS2 means UCS2 charset
	CharsetUCS2 charset = 0x08
)

// GeneralDataCoding is group of SMS Data Coding Scheme
type GeneralDataCoding struct {
	AutoDelete bool
	Compressed bool
	MsgClass   msgClass
	Charset    charset
}

func (c *GeneralDataCoding) encode() (b byte) {
	if c.AutoDelete {
		b = 0x40
	} else {
		b = 0x00
	}
	if c.Compressed {
		b |= 0x20
	}
	b |= byte(c.MsgClass)
	b |= byte(c.Charset & 0x0c)
	return
}

func (c *GeneralDataCoding) charset() charset {
	return c.Charset
}

func (c *GeneralDataCoding) String() string {
	var b bytes.Buffer
	b.WriteString("General Data Coding")
	if c.AutoDelete {
		b.WriteString("(Automatic Deletion)")
	}
	if c.Compressed {
		b.WriteString(", compressed")
	}
	switch c.MsgClass {
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
	switch c.Charset {
	case CharsetGSM7bit:
		b.WriteString(", GSM 7bit default alphabet")
	case Charset8bitData:
		b.WriteString(", 8 bit data")
	case CharsetUCS2:
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

func (c *MessageWaiting) encode() (b byte) {
	b = 0xc0
	b |= byte(c.Behavior & 0xc0)
	if c.Active {
		b |= 0x80
	}
	b |= byte(c.WaitingType & 0x03)
	return
}

func (c *MessageWaiting) charset() charset {
	if c.Behavior == StoreMessageUCS2 {
		return CharsetUCS2
	}
	return CharsetGSM7bit
}

func (c *MessageWaiting) String() string {
	var b bytes.Buffer
	b.WriteString("MessageWaiting")
	switch c.Behavior {
	case DiscardMessageGSM7bit:
		b.WriteString("(Discard Message with GSM 7bit default alphabet)")
	case StoreMessageGSM7bit:
		b.WriteString("(Store Message with GSM 7bit default alphabet)")
	case StoreMessageUCS2:
		b.WriteString("(Store Message with UCS2)")
	}
	if c.Active {
		b.WriteString(", active")
	} else {
		b.WriteString(", inactive")
	}
	switch c.WaitingType {
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

func (c *DataCodingMessage) encode() (b byte) {
	b = 0xf0
	if c.IsData {
		b |= 0x40
	}
	b |= byte(c.MsgClass & 0x03)
	return
}

func (c *DataCodingMessage) charset() charset {
	if c.IsData {
		return Charset8bitData
	}
	return CharsetGSM7bit
}

func (c *DataCodingMessage) String() string {
	var b bytes.Buffer
	b.WriteString("Data coding/message")
	if c.IsData {
		b.WriteString(", 8-bit data")
	} else {
		b.WriteString(", GSM 7 bit default alphabet")
	}
	switch c.MsgClass {
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
