package mortgage

import "errors"

var (
	ErrInvalidSum    = errors.New("неверная сумма")
	ErrInvalidRate   = errors.New("неверная ставка")
	ErrInvalidPeriod = errors.New("неверный срок")
	ErrInvalidDate   = errors.New("неверная дата начала")
	ErrInvalidAmount = errors.New("неверная сумма платежа")
)
