package sms

import "fmt"

// UnknownDataCodingError show invalid DCS
type UnknownDataCodingError struct {
	DCS byte
}

func (e UnknownDataCodingError) Error() string {
	return fmt.Sprintf("unknown DCS %x", e.DCS)
}

// UnknownGSM7bitRuneError show invalid DCS
type UnknownGSM7bitRuneError struct {
	R rune
}

func (e UnknownGSM7bitRuneError) Error() string {
	return fmt.Sprintf("unknown GSM 7bit rune %c", e.R)
}

// UnexpectedMessageTypeError show invalid SMS PDU type
type UnexpectedMessageTypeError struct {
	Expected, Actual byte
}

func (e UnexpectedMessageTypeError) Error() string {
	return fmt.Sprintf("unexpected message type %x is not %x", e.Actual, e.Expected)
}

// InvalidLengthError show invalid length for SMS TPDU
type InvalidLengthError struct {
}

func (e InvalidLengthError) Error() string {
	return fmt.Sprintf("invalid data length")
}

// UnexpectedInformationElementError show invalid SMS PDU type
type UnexpectedInformationElementError struct {
	Expected, Actual byte
}

func (e UnexpectedInformationElementError) Error() string {
	return fmt.Sprintf("unexpected IE %x is not %x", e.Actual, e.Expected)
}
