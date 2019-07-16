package sms

import "fmt"

// InvalidDataError show invalid data for SMS TPDU
type InvalidDataError struct {
	Name  string
	Bytes []byte
}

func (e InvalidDataError) Error() string {
	return fmt.Sprintf("invalid data for %s: % x", e.Name, e.Bytes)
}
