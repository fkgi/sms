package sms

type dcs interface {
	encodeDCS() (b byte)
	unitSize() int
}

func decodeDCS(b byte) dcs {
	if b&0xc == 0x00 {
		return &GeneralDataCoding{
			AutoDelete: false,
			Compressed: b&0x20 == 0x20,
			MsgClass:   msgClass(b & 0x13),
			Charset:    charset(b & 0x0c)}
	}
	if b&0xc == 0x40 {
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
	NoMessageClass = msgClass(0x00)
	// MessageClass0 means message class 0
	MessageClass0 = msgClass(0x10)
	// MessageClass1 means message class 1
	MessageClass1 = msgClass(0x11)
	// MessageClass2 means message class 2
	MessageClass2 = msgClass(0x12)
	// MessageClass3 means message class 3
	MessageClass3 = msgClass(0x13)

	// GSM7bitAlphabet means GSM 7 bit default alphabet charset
	GSM7bitAlphabet = charset(0x00)
	// Data8bit means 8 bit data charset
	Data8bit = charset(0x04)
	// UCS2 means UCS2 charset
	UCS2 = charset(0x08)
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

type waitType byte
type waitBehavior byte

const (
	// VoicemailMessageWaiting means waiting type
	VoicemailMessageWaiting = waitType(0x00)
	// FaxMessageWaiting means waiting type
	FaxMessageWaiting = waitType(0x01)
	// ElectronicMailMessageWaiting means waiting type
	ElectronicMailMessageWaiting = waitType(0x02)
	// OtherMessageWaiting means waiting type
	OtherMessageWaiting = waitType(0x03)

	// DiscardMessageGSM7bit means discard the contents
	DiscardMessageGSM7bit = waitBehavior(0x00)
	// StoreMessageGSM7bit means store the contents GSM 7bit alphabet
	StoreMessageGSM7bit = waitBehavior(0x10)
	// StoreMessageUCS2 means store the contents UCS2
	StoreMessageUCS2 = waitBehavior(0x20)
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
