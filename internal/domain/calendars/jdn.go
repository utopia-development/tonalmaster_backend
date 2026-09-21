package calendars

import "time"

// GregorianToJDN converts a Gregorian calendar date at midnight to its
// Julian Day Number (integer day convention).
func GregorianToJDN(date time.Time) int64 {
	y := int64(date.Year())
	m := int64(date.Month())
	d := int64(date.Day())

	a := (14 - m) / 12
	y2 := y + 4800 - a
	m2 := m + 12*a - 3

	return d + (153*m2+2)/5 + 365*y2 + y2/4 - y2/100 + y2/400 - 32045
}
