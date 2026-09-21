package calendars

import "time"

const (
	TonalpohualliCASOID = "tonalpohualli_caso"
	// CASOAnchorJDN is the Julian Day Number for 13 Aug 1521 (Julian),
	// identified as 1-Coatl and the beginning of trecena 1 in the
	// documented Alfonso Caso correlation.
	CASOAnchorJDN int64 = 2276828
)

var tonalpohualliSigns = [...]string{
	"Cipactli", "Ehecatl", "Calli", "Cuetzpalin", "Coatl",
	"Miquiztli", "Mazatl", "Tochtli", "Atl", "Itzcuintli",
	"Ozomahtli", "Malinalli", "Acatl", "Ocelotl", "Cuauhtli",
	"Cozcacuauhtli", "Ollin", "Tecpatl", "Quiahuitl", "Xochitl",
}

type TonalpohualliCASO struct {
	BaseJDN int64
}

func NewTonalpohualliCASO() TonalpohualliCASO {
	return TonalpohualliCASO{BaseJDN: CASOAnchorJDN}
}

func (s TonalpohualliCASO) ID() string {
	return TonalpohualliCASOID
}

func (s TonalpohualliCASO) Convert(date time.Time) (Result, error) {
	jdn := GregorianToJDN(date)
	// offset in [0, 259] relative to the Caso anchor (1-Coatl / trecena 1).
	offset := int64(mod(jdn-s.BaseJDN, 260))

	return Result{
		System:    s.ID(),
		Date:      date,
		JDN:       jdn,
		Trecena:   int(offset)/13 + 1,
		Sign:      tonalpohualliSigns[mod(offset+4, 20)],
		DayNumber: mod(offset, 13) + 1,
	}, nil
}
