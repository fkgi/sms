package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fkgi/sms"
)

func main() {
	pdutype := flag.String("t", "",
		"PDU type for encoding `submit|submitreport|deliver|deliverreport|command|statusreport`")
	layer := flag.String("l", "tp", "PDU layer for encoding `tp|rp|cp")
	revert := flag.Bool("r", false, "decode SMS PDU from bindary to JSON")
	flag.Parse()

	in, e := io.ReadAll(os.Stdin)
	if e != nil {
		fmt.Fprintf(os.Stderr, "failed to read data from stdin: %s", e)
		os.Exit(1)
	}

	var res []byte
	if *revert {
		res, e = decode(in, *pdutype, *layer)
	} else {
		res, e = encode(in, *pdutype, *layer)
	}

	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	os.Stdout.Write(res)
}

func encode(jsondata []byte, pdutype, layer string) (r []byte, e error) {
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

	if e != nil {
		return
	}
	switch layer {
	case "tp":
		r = pdu.MarshalTP()
	case "rp":
		r = pdu.MarshalRP()
	case "cp":
		r = pdu.MarshalCP()
	default:
		e = fmt.Errorf("invalid PDU layer: %s", layer)
	}
	return
}

func decode(bindata []byte, pdutype, layer string) (r []byte, e error) {
	var pdu sms.CPDU
	switch pdutype {
	case "submit", "deliverreport", "command":
		switch layer {
		case "tp":
			pdu, e = sms.UnmarshalTPMO(bindata)
		case "rp":
			pdu, e = sms.UnmarshalRPMO(bindata)
		case "cp":
			pdu, e = sms.UnmarshalCPMO(bindata)
		default:
			e = fmt.Errorf("invalid PDU layer: %s", layer)
		}
	case "deliver", "submitreport", "statusreport":
		switch layer {
		case "tp":
			pdu, e = sms.UnmarshalTPMT(bindata)
		case "rp":
			pdu, e = sms.UnmarshalRPMT(bindata)
		case "cp":
			pdu, e = sms.UnmarshalCPMT(bindata)
		default:
			e = fmt.Errorf("invalid PDU layer: %s", layer)
		}
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
	if e != nil {
		return
	}
	if r, e = json.Marshal(pdu); e != nil {
		return
	}
	tmp1 := map[string]json.RawMessage{}
	tmp2 := map[string]json.RawMessage{}
	if e = json.Unmarshal(r, &tmp1); e != nil {
		return
	}
	for k, v := range tmp1 {
		switch layer {
		case "tp":
			if strings.HasPrefix(k, "tp-") {
				tmp2[k] = v
			}
		case "rp":
			if !strings.HasPrefix(k, "cp-") {
				tmp2[k] = v
			}
		default:
			tmp2[k] = v
		}
	}
	r, e = json.Marshal(tmp2)

	return
}
