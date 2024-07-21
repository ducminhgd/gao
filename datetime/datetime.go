package datetime

import "time"

// https://stackoverflow.com/a/51915792
// In go, the format layout for values of for math are:
// - 1: Month
// - 2: Day of Month
// - 3 or 15: Hours
// - 4: minutes
// - 5: seconds
// - 6: year
// - 7: Timezone
// - 0 or 9: for partial second
const (
	RFC3339      = time.RFC3339
	RFC3339Micro = "2006-01-02T15:04:05.999999Z07:00"
	RFC3339Nano  = time.RFC3339Nano

	// Omit timezone, assume that this is UTC
	RFC3339_UTC_NoZone      = "2006-01-02T15:04:05Z"
	RFC3339Micro_UTC_NoZone = "2006-01-02T15:04:05.999999Z"
	RFC3339Nano_UTC_NoZone  = "2006-01-02T15:04:05.999999999Z"

	ISO8601Date          = "2006-01-02"
	ISO8601YMD           = "20060102"
	ISO8601Time          = "15:04:05"
	ISO8601Simple        = "2006-01-02T15:04:05"
	ISO8601Timezone      = "2006-01-02T15:04:05-0700"
	ISO8601UTC           = "2006-01-02T15:04:05Z"
	ISO8601TimeMilli     = "15:04:05.999"
	ISO8601SimpleMilli   = "2006-01-02T15:04:05.999"
	ISO8601TimezoneMilli = "2006-01-02T15:04:05.999-0700"
	ISO8601UTCMilli      = "2006-01-02T15:04:05.999Z"
	ISO8601TimeMicro     = "15:04:05.999999"
	ISO8601SimpleMicro   = "2006-01-02T15:04:05.999999"
	ISO8601TimezoneMicro = "2006-01-02T15:04:05.999999-0700"
	ISO8601UTCMicro      = "2006-01-02T15:04:05.999999Z"
	ISO8601UTCMicroNoZ   = "2006-01-02T15:04:05.999999"
)
