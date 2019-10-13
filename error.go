package sms

import (
	"errors"
	"fmt"
)

// UnknownDataCodingError show invalid DCS
type UnknownDataCodingError struct {
	DCS byte
}

func (e UnknownDataCodingError) Error() string {
	return fmt.Sprintf("unknown DCS % x", e.DCS)
}

// UnknownGSM7bitRuneError show invalid rune for GSM 7bit string
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
	return fmt.Sprintf("unexpected message type %x is not %x",
		e.Actual, e.Expected)
}

// UnknownMessageTypeError show invalid SMS PDU type
type UnknownMessageTypeError struct {
	Actual byte
}

func (e UnknownMessageTypeError) Error() string {
	return fmt.Sprintf("unknown message type %x", e.Actual)
}

// UnexpectedInformationElementError show invalid SMS IE
type UnexpectedInformationElementError struct {
	Expected, Actual byte
}

func (e UnexpectedInformationElementError) Error() string {
	return fmt.Sprintf("unexpected IE %x is not %x", e.Actual, e.Expected)
}

var (
	// ErrInvalidLength show invalid length for SMS PDU data
	ErrInvalidLength = errors.New("invalid data length")

	// ErrExtraData show extra data for SMS PDU
	ErrExtraData = errors.New("extra data")
)
