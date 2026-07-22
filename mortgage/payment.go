package mortgage

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Payment — всё, что можно развернуть в слайс [SinglePayment] методом Payments.
// В пакете есть две готовые реализации: сам [SinglePayment] и [MonthlyPayment]
// (который разворачивается в поток).
type Payment interface {
	Payments() []SinglePayment
}

// PaymentStrategy описывает, как досрочный платёж влияет на график.
type PaymentStrategy int

const (
	// ReduceTerm сохраняет размер регулярного ежемесячного платежа и
	// закрывает кредит раньше срока. Это значение по умолчанию (ноль).
	ReduceTerm PaymentStrategy = iota
	// ReducePayment сохраняет исходный срок и уменьшает регулярный
	// ежемесячный платёж после внесения досрочного.
	ReducePayment
)

// String возвращает стабильный идентификатор стратегии ("reduce_term" /
// "reduce_payment"), подходящий для логов и сериализации. Для неизвестных
// значений возвращается "early_repayment(N)".
func (s PaymentStrategy) String() string {
	switch s {
	case ReduceTerm:
		return "reduce_term"
	case ReducePayment:
		return "reduce_payment"
	default:
		return fmt.Sprintf("early_repayment(%d)", int(s))
	}
}

// SinglePayment — единичный платёж с конкретной датой. Это либо самостоятельное
// досрочное погашение, либо одна строка из [MonthlyPayment.Payments].
//
//   - Date     — дата платежа (время нормализуется до полуночи UTC);
//   - Amount   — сумма в рублях;
//   - Strategy — см. [PaymentStrategy]. ReducePayment имеет смысл прежде всего
//     для платежей, порождённых [MonthlyPayment]: она помечает их
//     как «регулярный» поток, размер которого может пересчитываться
//     при уменьшении баланса из-за досрочного (ReduceTerm) платежа.
type SinglePayment struct {
	Date     time.Time
	Amount   decimal.Decimal
	Strategy PaymentStrategy
}

// Payments реализует [Payment], возвращая одноэлементный слайс.
func (sp SinglePayment) Payments() []SinglePayment {
	return []SinglePayment{sp}
}

// MonthlyPayment — регулярный платёж в каждый месяц со дня Begin (включительно)
// до End (исключительно). Метод Payments разворачивает его в эквивалентный
// плоский слайс [SinglePayment], каждый из которых несёт те же Amount и Strategy.
//
// Пример: поток на 12 месяцев с 2026-06-01 по 2027-06-01 даёт 12 SinglePayment
// с датами с 2026-06-01 по 2027-05-01.
type MonthlyPayment struct {
	Begin, End time.Time
	Amount     decimal.Decimal
	Strategy   PaymentStrategy
}

// Payments реализует [Payment], разворачивая рекуррентный поток в слайс.
func (mp MonthlyPayment) Payments() []SinglePayment {
	payments := make([]SinglePayment, 0, 12)
	current := normalizeDate(mp.Begin)
	end := normalizeDate(mp.End)
	for current.Before(end) {
		payments = append(payments, SinglePayment{
			Date:     current,
			Amount:   mp.Amount,
			Strategy: mp.Strategy,
		})
		current = current.AddDate(0, 1, 0)
	}
	return payments
}
