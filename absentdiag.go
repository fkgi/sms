package sms

import "fmt"

// AbsentDiag indicate Absent-User-Diagnostic-SM
type AbsentDiag int

const (
	// NoAbsentDiag is "no diag data"
	NoAbsentDiag AbsentDiag = iota
	// NoPagingRespMSC is "no paging response via the MSC"
	NoPagingRespMSC
	// IMSIDetached is "IMSI detached"
	IMSIDetached
	// RoamingRestrict is "roaming restriction"
	RoamingRestrict
	// DeregisteredNonGPRS is "deregistered in the HLR for non GPRS"
	DeregisteredNonGPRS
	// PurgedNonGPRS is "MS purged for non GPRS"
	PurgedNonGPRS
	// NoPagingRespSGSN is "no paging response via the SGSN"
	NoPagingRespSGSN
	// GPRSDetached is "GPRS detached"
	GPRSDetached
	// DeregisteredGPRS is "deregistered in the HLR for GPRS"
	DeregisteredGPRS
	// PurgedGPRS is "MS purged for GPRS"
	PurgedGPRS
	// UnidentifiedSubsMSC is "Unidentified subscriber via the MSC"
	UnidentifiedSubsMSC
	// UnidentifiedSubsSGSN is "Unidentified subscriber via the SGSN"
	UnidentifiedSubsSGSN
	// DeregisteredIMS is "deregistered in the HSS/HLR for IMS"
	DeregisteredIMS
	// NoRespIPSMGW is "no response via the IP-SM-GW"
	NoRespIPSMGW
	// TempUnavailable is "the MS is temporarily unavailable"
	TempUnavailable
)

func (a AbsentDiag) String() string {
	switch a {
	case NoAbsentDiag:
		return "no diag data"
	case NoPagingRespMSC:
		return "no paging response via the MSC"
	case IMSIDetached:
		return "IMSI detached"
	case RoamingRestrict:
		return "roaming restriction"
	case DeregisteredNonGPRS:
		return "deregistered in the HLR for non GPRS"
	case PurgedNonGPRS:
		return "MS purged for non GPRS"
	case NoPagingRespSGSN:
		return "no paging response via the SGSN"
	case GPRSDetached:
		return "GPRS detached"
	case DeregisteredGPRS:
		return "deregistered in the HLR for GPRS"
	case PurgedGPRS:
		return "MS purged for GPRS"
	case UnidentifiedSubsMSC:
		return "Unidentified subscriber via the MSC"
	case UnidentifiedSubsSGSN:
		return "Unidentified subscriber via the SGSN"
	case DeregisteredIMS:
		return "deregistered in the HSS/HLR for IMS"
	case NoRespIPSMGW:
		return "no response via the IP-SM-GW"
	case TempUnavailable:
		return "the MS is temporarily unavailable"
	}
	return fmt.Sprintf("unknown(0x%x)", byte(a-1))
}

// Byte make byte digit value
func (a AbsentDiag) Byte() byte {
	switch a {
	case NoPagingRespMSC:
		return 0
	case IMSIDetached:
		return 1
	case RoamingRestrict:
		return 2
	case DeregisteredNonGPRS:
		return 3
	case PurgedNonGPRS:
		return 4
	case NoPagingRespSGSN:
		return 5
	case GPRSDetached:
		return 6
	case DeregisteredGPRS:
		return 7
	case PurgedGPRS:
		return 8
	case UnidentifiedSubsMSC:
		return 9
	case UnidentifiedSubsSGSN:
		return 10
	case DeregisteredIMS:
		return 11
	case NoRespIPSMGW:
		return 12
	case TempUnavailable:
		return 13
	}
	return byte(a - 1)
}

// B2AbsDiag make AbsentDiag data from byte value
func B2AbsDiag(b byte) AbsentDiag {
	switch b {
	case 0:
		return NoPagingRespMSC
	case 1:
		return IMSIDetached
	case 2:
		return RoamingRestrict
	case 3:
		return DeregisteredNonGPRS
	case 4:
		return PurgedNonGPRS
	case 5:
		return NoPagingRespSGSN
	case 6:
		return GPRSDetached
	case 7:
		return DeregisteredGPRS
	case 8:
		return PurgedGPRS
	case 9:
		return UnidentifiedSubsMSC
	case 10:
		return UnidentifiedSubsSGSN
	case 11:
		return DeregisteredIMS
	case 12:
		return NoRespIPSMGW
	case 13:
		return TempUnavailable
	}
	return AbsentDiag(b + 1)
}
