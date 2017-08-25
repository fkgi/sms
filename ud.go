package sms

// MakeCodedData generate splited data
func MakeCodedData(s string) (uds [][]byte, dcs GeneralDataCoding) {
	uds = [][]byte{}
	dcs = GeneralDataCoding{
		AutoDelete: false,
		Compressed: false,
		MsgClass:   NoMessageClass}

	if _, e := GetGSM7bitString(s); e == nil {
		dcs.Charset = GSM7bitAlphabet
	} else {
		dcs.Charset = UCS2
	}
	if ud, e := dcs.Encode(s); e == nil {
		for len(ud) > 140 {
			uds = append(uds, ud[:134])
			ud = ud[134:]
		}
		uds = append(uds, ud)
	}
	return
}
