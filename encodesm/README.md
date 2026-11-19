# encodesm
`encodesm` is a command-line tool written in Go for encoding and decoding SMS PDUs (Protocol Data Units).  
It supports conversion between JSON and binary PDU formats, and handles various SMS PDU types such as submit, deliver, and reports.
Binary PDU is represented as a string in hexadecimal format.

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
- `-t` : Specify PDU type (mandatory)  
  `submit` | `submitreport` | `deliver` | `deliverreport` | `command` | `statusreport`
- `-l` : Specify layer of PDU (default: `tp`)
  `tp` | `rp` | `cp`
- `-r` : Decode mode (PDU → JSON). If omitted, encodes (JSON → PDU)
- `-b` : Output/Input of binary will be raw binary data
(not hexadecimal format string)

## JSON Example
JSON for submit.
```json
{
  "tp-rd": false,
  "tp-srr": false,
  "tp-rp": false,
  "tp-mr": 64,
  "tp-da": {
    "ton": 1,
    "npi": 1,
    "addr": "819012345678"
  },
  "tp-pid": 0,
  "tp-dcs": 8,
  "tp-ud": {
    "text": "あいうえお"
  }
}
```
binary for submit.
```
01000c9118092143658700081800540065007300740020006d006500730073006100670065
```

## License
MIT License

## Notes
- PDU specifications and JSON structure follow [github.com/fkgi/sms](https://github.com/fkgi/sms).
- All input and output is via standard input/output. Use redirection for file operations.