package mortgage

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Summary содержит ключевые результаты расчёта кредита.
// Все денежные поля — в рублях.
type Summary struct {
	MonthlyPayment   decimal.Decimal // первый ежемесячный платёж из графика
	TotalPaid        decimal.Decimal // суммарно выплачено за срок кредита
	Overpayment      decimal.Decimal // суммарно выплачено процентов (переплата)
	LastDate         time.Time       // дата последнего планового платежа
	RemainingBalance decimal.Decimal // остаток тела после последней строки
}

// String возвращает однострочное форматированное представление сводки.
func (s Summary) String() string {
	if s.MonthlyPayment.IsZero() && s.TotalPaid.IsZero() {
		return "Нет данных"
	}
	return fmt.Sprintf(
		"ежемес.: %15s RUB | всего: %15s RUB | переплата: %15s RUB | конец: %12s | остаток: %15s RUB",
		FormatMoney(s.MonthlyPayment),
		FormatMoney(s.TotalPaid),
		FormatMoney(s.Overpayment),
		s.LastDate.Format(TimeFormat),
		FormatMoney(s.RemainingBalance),
	)
}

// Summary возвращает агрегированные результаты графика платежей.
//
// Если график пуст, возвращается нулевое значение. Если платежи не покрывают
// кредит, RemainingBalance будет положительным, а LastDate — датой последней
// строки графика.
func (m *Mortgage) Summary() Summary {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.schedule) == 0 {
		return Summary{}
	}

	var (
		totalPaid     = decimal.Zero
		totalInterest = decimal.Zero
	)
	for _, row := range m.schedule {
		totalPaid = totalPaid.Add(row.Payment)
		totalInterest = totalInterest.Add(row.Interest)
	}

	last := m.schedule[len(m.schedule)-1]
	return Summary{
		MonthlyPayment:   m.schedule[0].Payment,
		TotalPaid:        totalPaid,
		Overpayment:      totalInterest,
		LastDate:         last.Date,
		RemainingBalance: last.Balance,
	}
}
