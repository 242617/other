// Package mortgage рассчитывает графики платежей по аннуитетному кредиту
// (ипотеке) с поддержкой досрочных погашений.
//
// Объект [Mortgage] создаётся функцией [New], после чего наполняется платежами
// через [Mortgage.Add]. После каждого Add график перестраивается
// детерминированно из списка платежей, поэтому один и тот же набор платежей
// всегда даёт один и тот же график.
//
// Платежи реализуют интерфейс [Payment]. В пакете есть две готовые реализации:
// [SinglePayment] (разовое досрочное погашение на конкретную дату) и
// [MonthlyPayment] (регулярный платёж, разворачиваемый в слайс SinglePayment
// от Begin до End).
package mortgage

import (
	"fmt"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// Mortgage описывает аннуитетный кредит с фиксированной годовой ставкой и
// сроком в месяцах.
type Mortgage struct {
	sum    decimal.Decimal
	rate   decimal.Decimal
	begin  time.Time
	period int

	mu sync.RWMutex

	initialPayment decimal.Decimal
	payments       []SinglePayment
	schedule       []ScheduleRow
}

// New создаёт провалидированный объект Mortgage.
//
//	sum    — тело кредита, в рублях;
//	rate   — годовая ставка в виде доли (0.06 = 6%);
//	begin  — дата начала графика (используется компонент дня);
//	period — срок в месяцах.
func New(sum, rate decimal.Decimal, begin time.Time, period int) (*Mortgage, error) {
	m := &Mortgage{
		sum:    sum,
		rate:   rate,
		begin:  normalizeDate(begin),
		period: period,
	}
	if err := m.validate(); err != nil {
		return nil, err
	}
	m.initialPayment = m.annuityPayment(m.sum, m.period)
	return m, nil
}

func (m *Mortgage) Validate() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.validate()
}

func (m *Mortgage) validate() error {
	switch {
	case m.sum.IsNegative() || m.sum.IsZero():
		return fmt.Errorf("%w: должна быть положительной, получено %s", ErrInvalidSum, FormatMoney(m.sum))
	case m.rate.IsNegative():
		return fmt.Errorf("%w: не может быть отрицательной, получено %s", ErrInvalidRate, m.rate.String())
	case m.period <= 0:
		return fmt.Errorf("%w: должен быть положительным, получено %d", ErrInvalidPeriod, m.period)
	case m.period > 600:
		return fmt.Errorf("%w: не может превышать 600 месяцев (50 лет), получено %d", ErrInvalidPeriod, m.period)
	case m.begin.IsZero():
		return fmt.Errorf("%w: должна быть задана", ErrInvalidDate)
	}
	return nil
}

func (m *Mortgage) Sum() decimal.Decimal {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sum
}

func (m *Mortgage) Rate() decimal.Decimal {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rate
}

func (m *Mortgage) Begin() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.begin
}

func (m *Mortgage) Period() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.period
}

// InitialPayment возвращает аннуитетный платёж для полной суммы и исходного
// срока. Значение не учитывает досрочные погашения.
func (m *Mortgage) InitialPayment() decimal.Decimal {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.initialPayment
}

// Add регистрирует [Payment] и перестраивает график.
//
// Платеж разворачивается в слайс [SinglePayment] методом [Payment.Payments],
// нормализуется и добавляется. После этого график пересчитывается с нуля.
// Метод безопасно вызывать из нескольких горутин.
func (m *Mortgage) Add(p Payment) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, sp := range p.Payments() {
		sp.Date = normalizeDate(sp.Date)
		if !sp.Amount.IsPositive() {
			return fmt.Errorf("%w: должна быть положительной, получено %s", ErrInvalidAmount, FormatMoney(sp.Amount))
		}
		m.payments = append(m.payments, sp)
	}
	return m.rebuild()
}
