package sms

// MakeCodedData generate splited data
func MakeCodedData(s string) (uds [][]byte, dcs GeneralDataCoding) {
	uds = [][]byte{}
	dcs = GeneralDataCoding{
		AutoDelete: false,
		Compressed: false,
		MsgClass:   NoMessageClass}

	dcs.Charset = GSM7bitAlphabet
	for _, r := range []rune(s) {
		if getCode(r) == 0x80 {
			dcs.Charset = UCS2
			break
		}
	}

	if dcs.Charset == GSM7bitAlphabet {
		r := []rune(s)
		maxlen := 160

		for len(r) > maxlen {
			ud, _ := dcs.Encode(string(r[:153]))
			uds = append(uds, ud)
			r = r[153:]
			maxlen = 153
		}
		ud, _ := dcs.Encode(string(r))
		uds = append(uds, ud)
	} else {
		rs := []rune{}
		maxlen := 140

		for _, r := range []rune(s) {
			tmp := append(rs, r)
			if len(string(tmp)) > maxlen {
				ud, _ := dcs.Encode(string(rs))
				uds = append(uds, ud)
				rs = []rune{r}
				maxlen = 135
			} else {
				rs = tmp
			}
		}
		ud, _ := dcs.Encode(string(rs))
		uds = append(uds, ud)
	}
	return
}
