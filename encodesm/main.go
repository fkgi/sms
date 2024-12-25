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
	pdutype := flag.String("t", "",
		"PDU type for encoding `submit|submitreport|deliver|deliverreport|command|statusreport`")
	revert := flag.Bool("r", false, "decode SMS PDU from bindary to JSON")
	flag.Parse()

	in, e := io.ReadAll(os.Stdin)
	if e != nil {
		fmt.Fprintf(os.Stderr, "failed to read data from stdin: %s", e)
		os.Exit(1)
	}

	var res []byte
	if *revert {
		res, e = decode(in, *pdutype)
	} else {
		res, e = encode(in, *pdutype)
	}

	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	os.Stdout.Write(res)
}

func encode(jsondata []byte, pdutype string) (r []byte, e error) {
	var pdu sms.TPDU

	switch pdutype {
	case "submit":
		d := new(sms.Submit)
		e = json.Unmarshal(jsondata, d)
		pdu = d
	case "deliver":
		d := new(sms.Deliver)
		e = json.Unmarshal(jsondata, d)
		pdu = d
	case "submitreport":
		d := new(sms.SubmitReport)
		e = json.Unmarshal(jsondata, d)
		pdu = d
	case "deliverreport":
		d := new(sms.DeliverReport)
		e = json.Unmarshal(jsondata, d)
		pdu = d
	case "command":
		d := new(sms.Command)
		e = json.Unmarshal(jsondata, d)
		pdu = d
	case "statusreport":
		d := new(sms.StatusReport)
		e = json.Unmarshal(jsondata, d)
		pdu = d
	default:
		e = fmt.Errorf("invalid PDU type: %s", pdutype)
	}

	if e == nil {
		r = pdu.MarshalTP()
	}
	return
}

func decode(bindata []byte, pdutype string) (r []byte, e error) {
	var pdu sms.TPDU
	switch pdutype {
	case "submit", "deliverreport", "command":
		pdu, e = sms.UnmarshalTPMO(bindata)
	case "deliver", "submitreport", "statusreport":
		pdu, e = sms.UnmarshalTPMT(bindata)
	default:
		e = fmt.Errorf("invalid PDU type: %s", pdutype)
	}
	if e != nil {
		return
	}

	switch pdutype {
	case "submit":
		if _, ok := pdu.(sms.Submit); !ok {
			e = fmt.Errorf("PDU type missmatch")
		}
	case "deliverreport":
		if _, ok := pdu.(sms.DeliverReport); !ok {
			e = fmt.Errorf("PDU type missmatch")
		}
	case "command":
		if _, ok := pdu.(sms.Command); !ok {
			e = fmt.Errorf("PDU type missmatch")
		}
	case "deliver":
		if _, ok := pdu.(sms.Deliver); !ok {
			e = fmt.Errorf("PDU type missmatch")
		}
	case "submitreport":
		if _, ok := pdu.(sms.SubmitReport); !ok {
			e = fmt.Errorf("PDU type missmatch")
		}
	case "statusreport":
		if _, ok := pdu.(sms.StatusReport); !ok {
			e = fmt.Errorf("PDU type missmatch")
		}
	}
	if e == nil {
		r, e = json.Marshal(pdu)
	}
	return
}
