# encodesm
`encodesm` is a command-line tool written in Go for encoding and decoding SMS PDUs (Protocol Data Units).  
It supports conversion between JSON and binary PDU formats, and handles various SMS PDU types such as submit, deliver, and reports.

## Installation
```sh
go build -o encodesm ./encodesm
```

## Usage
### Encode (JSON → PDU binary)
```sh
cat submit.json | ./encodesm -t submit > submit.pdu
```

### Decode (PDU binary → JSON)
```sh
cat submit.pdu | ./encodesm -t submit -r > submit.json
```

#### Options
- `-t` : Specify PDU type  
  `submit` | `submitreport` | `deliver` | `deliverreport` | `command` | `statusreport`
- `-r` : Decode mode (PDU → JSON). If omitted, encodes (JSON → PDU)

## JSON Example
```json
{
  "tp-mti": 1,
  "tp-rd": false,
  "tp-vpf": 0,
  "tp-rp": false,
  "tp-mr": 0,
  "tp-da": {
    "ton": 1,
    "npi": 1,
    "addr": "819012345678"
  },
  "tp-pid": 0,
  "tp-dcs": {
    "msgcharset": 0
  },
  "tp-vp": "2025-09-26T12:00:00+09:00",
  "tp-ud": {
    "text": "Test message"
  }
}
```

## License
MIT License

## Notes
- PDU specifications and JSON structure follow [github.com/fkgi/sms](https://github.com/fkgi/sms).
- All input and output is via standard input/output. Use redirection for file operations.