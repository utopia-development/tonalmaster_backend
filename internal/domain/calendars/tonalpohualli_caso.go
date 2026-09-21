package calendars

import (
	"time"
)

const TonalpohualliCASOID = "tonalpohualli_caso"

type TonalpohualliCASO struct {
	BaseJDN int64
}

func NewTonalpohualliCASO(baseJDN int64) TonalpohualliCASO {
	return TonalpohualliCASO{BaseJDN: baseJDN}
}

func (s TonalpohualliCASO) ID() string {
	return TonalpohualliCASOID
}

func (s TonalpohualliCASO) Convert(date time.Time) (Result, error) {
	jdn := GregorianToJDN(date)
	offset := mod(jdn-s.BaseJDN, 260)

	// The correlation-specific mapping is intentionally isolated here.
	// These coefficients must be confirmed against the authoritative
	// correlation before adding calendar-result fixtures.
	return Result{
		System: s.ID(),
		Date: date,
		JDN: jdn,
		Trecena: mod(offset, 13) + 1,
		Sign: "",
		DayNumber: mod(offset, 20) + 1,
		NightLord: "",
	}, nil
}

func mod(value, divisor int64) int {
	r := value % divisor
	if r < 0 {
		return int(r + divisor)
	}
	return int(r)
}
