package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/fkgi/sms"
)

func main() {
	jsondata, e := io.ReadAll(os.Stdin)
	if e != nil {
		fmt.Fprintf(os.Stderr, "failed to read data from stdin: %s", e)
		os.Exit(1)
	}
	tpdus := map[string]json.RawMessage{}
	e = json.Unmarshal(jsondata, &tpdus)
	if e != nil {
		fmt.Fprintf(os.Stderr, "not correct json data: %s", e)
		os.Exit(1)
	}
	var pdu sms.TPDU
	for k, v := range tpdus {
		switch k {
		case "submit":
			d := new(sms.Submit)
			e = json.Unmarshal(v, d)
			pdu = d
		case "deliver":
			d := new(sms.Deliver)
			e = json.Unmarshal(v, d)
			pdu = d
		case "submitreport":
			d := new(sms.SubmitReport)
			e = json.Unmarshal(v, d)
			pdu = d
		case "deliverreport":
			d := new(sms.DeliverReport)
			e = json.Unmarshal(v, d)
			pdu = d
		case "command":
			d := new(sms.Command)
			e = json.Unmarshal(v, d)
			pdu = d
		case "statusreport":
			d := new(sms.StatusReport)
			e = json.Unmarshal(v, d)
			pdu = d
		}
		break
	}
	if e != nil {
		fmt.Fprintf(os.Stderr, "not correct json data: %s", e)
		os.Exit(1)
	}
	os.Stdout.Write(pdu.MarshalTP())
}
