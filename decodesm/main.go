package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/fkgi/sms"
)

func main() {
	flag.Parse()
	ismo := flag.Arg(0) == "mo"

	bindata, e := io.ReadAll(os.Stdin)
	if e != nil {
		fmt.Fprintf(os.Stderr, "failed to read data from stdin: %s", e)
		os.Exit(1)
	}

	var pdu sms.TPDU
	if ismo {
		pdu, e = sms.UnmarshalTPMO(bindata)
	} else {
		pdu, e = sms.UnmarshalTPMT(bindata)
	}
	if e != nil {
		mostr := "MT"
		if ismo {
			mostr = "MO"
		}
		fmt.Fprintf(os.Stderr, "not correct TPDU data for %s: %s", mostr, e)
		os.Exit(1)
	}

	jsondata, e := json.Marshal(pdu)
	if e != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal TPDU data to JSON: %s", e)
		os.Exit(1)
	}
	out := map[string]json.RawMessage{}

	switch pdu.(type) {
	case sms.Submit:
		out["submit"] = json.RawMessage(jsondata)
	case sms.Deliver:
		out["deliver"] = json.RawMessage(jsondata)
	case sms.SubmitReport:
		out["submitreport"] = json.RawMessage(jsondata)
	case sms.DeliverReport:
		out["deliverreport"] = json.RawMessage(jsondata)
	case sms.Command:
		out["command"] = json.RawMessage(jsondata)
	case sms.StatusReport:
		out["statusreport"] = json.RawMessage(jsondata)
	}
	jsondata, e = json.Marshal(out)
	if e != nil {
		panic(e)
	}
	os.Stdout.Write(jsondata)
}
