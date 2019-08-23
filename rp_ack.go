package sms

// AckMO is MO RP-ACK RPDU
type AckMO struct {
	rpAnswer
}

// AckMT is MT RP-ACK RPDU
type AckMT struct {
	rpAnswer
}

// MarshalRP returns binary data
func (d AckMO) MarshalRP() []byte {
	return d.marshalAck(2)
}

// MarshalRP returns binary data
func (d AckMT) MarshalRP() []byte {
	return d.marshalAck(3)
}

// MarshalCP output byte data of this CPDU
func (d AckMO) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// MarshalCP output byte data of this CPDU
func (d AckMT) MarshalCP() []byte {
	return d.cpData.marshal(d.MarshalRP())
}

// UnmarshalAckMO decode Ack MO from bytes
func UnmarshalAckMO(b []byte) (a AckMO, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP reads binary data
func (d *AckMO) UnmarshalRP(b []byte) error {
	return d.unmarshalAck(b, 2)
}

// UnmarshalAckMT decode Ack MT from bytes
func UnmarshalAckMT(b []byte) (a AckMT, e error) {
	e = a.UnmarshalRP(b)
	return
}

// UnmarshalRP reads binary data
func (d *AckMT) UnmarshalRP(b []byte) error {
	return d.unmarshalAck(b, 3)
}

// UnmarshalCP get data of this CPDU
func (d *AckMO) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

// UnmarshalCP get data of this CPDU
func (d *AckMT) UnmarshalCP(b []byte) (e error) {
	if b, e = d.cpData.unmarshal(b); e == nil {
		e = d.UnmarshalRP(b)
	}
	return
}

func (d AckMO) String() string {
	return d.stringAck()
}

func (d AckMT) String() string {
	return d.stringAck()
}
