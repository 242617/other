package mortgage

import (
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

// intermediatePrecision — количество знаков после запятой для промежуточных
// вычислений аннуитета перед финальным Round(2). Значения 30 с запасом
// достаточно для любых рублёвых сумм и при этом далеко от предела decimal.Decimal.
const intermediatePrecision int32 = 30

// ScheduleRow — одна строка графика погашения.
type ScheduleRow struct {
	Number    int             // порядковый номер месяца от даты начала (с 1)
	Date      time.Time       // дата платежа
	Payment   decimal.Decimal // общий платёж за месяц
	Principal decimal.Decimal // часть платежа, идущая в погашение тела
	// (может быть отрицательной, если платёж не
	// покрыл начисленные проценты)
	Interest decimal.Decimal // проценты, начисленные за месяц
	Balance  decimal.Decimal // остаток тела после этой строки
}

// Schedule возвращает защитную копию текущего графика. Изменение возвращённого
// слайса не влияет на объект Mortgage.
func (m *Mortgage) Schedule() []ScheduleRow {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]ScheduleRow, len(m.schedule))
	copy(out, m.schedule)
	return out
}

// rebuild пересчитывает m.schedule на основе m.payments по правилам
// аннуитетного кредита.
//
// Формулы:
//
//	r = годовая_ставка / 12            — месячная ставка
//	I_k = B_{k-1} · r                  — проценты месяца k
//	D_k = платёж − I_k                 — погашение тела в месяце k
//	B_k = B_{k-1} − D_k                — остаток тела после месяца k
//	A = P · r(1+r)ⁿ / ((1+r)ⁿ − 1)     — аннуитетный платёж
//
// Модель платежей:
//
// Кредит итерируется по месяцам (не по слайсу платежей). Для каждого месяца:
//   - Если есть явные платежи пользователя — используется их сумма.
//   - Если явных платежей нет — используется currentA (текущий аннуитет).
//
// currentA — текущий минимальный платёж. Изначально равен аннуитету для полной
// суммы и исходного срока. Меняется только при стратегии ReducePayment.
//
// Стратегии (срабатывают, когда платёж превысил currentA):
//
//   - ReduceTerm: пересчитываем оставшийся срок через ceilMonths — наименьшее
//     число месяцев, за которое currentA погасит остаток. Срок сокращается.
//   - ReducePayment: пересчитываем currentA на новый срок — он уменьшается.
//     Срок сохраняется. Все будущие месяцы без явных платежей автоматически
//     используют новый (меньший) currentA.
func (m *Mortgage) rebuild() error {
	if err := m.validate(); err != nil {
		return fmt.Errorf("validate mortgage: %w", err)
	}
	if len(m.payments) == 0 {
		m.schedule = nil
		return nil
	}

	/// 1. Сортировка и слияние платежей по дате

	sorted := make([]SinglePayment, len(m.payments))
	copy(sorted, m.payments)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Date.Before(sorted[j].Date)
	})

	type slot struct {
		amount   decimal.Decimal
		strategy PaymentStrategy
	}
	slotMap := make(map[time.Time]*slot, len(sorted))
	dateKeys := make([]time.Time, 0, len(sorted))
	for _, sp := range sorted {
		s, ok := slotMap[sp.Date]
		if !ok {
			s = &slot{}
			slotMap[sp.Date] = s
			dateKeys = append(dateKeys, sp.Date)
		}
		s.amount = s.amount.Add(sp.Amount)
		// ReducePayment «побеждает»: если хотя бы один платёж в дате несёт
		// эту стратегию, вся дата обрабатывается как ReducePayment.
		if sp.Strategy == ReducePayment {
			s.strategy = ReducePayment
		}
	}
	sort.Slice(dateKeys, func(i, j int) bool { return dateKeys[i].Before(dateKeys[j]) })

	/// 2. Инициализация

	monthlyRate := m.rate.DivRound(decimal.NewFromInt(12), intermediatePrecision)
	balance := m.sum
	currentA := m.initialPayment
	remainMonths := m.period
	oneKopeck := decimal.New(1, -2)

	schedule := make([]ScheduleRow, 0, m.period)
	keyIdx := 0

	/// 3. Цикл по месяцам

	for rowNum := 1; remainMonths > 0 && balance.GreaterThan(oneKopeck); rowNum++ {
		payDate := m.begin.AddDate(0, rowNum-1, 0)

		// Проценты месяца: I_k = B_{k-1} · r.
		interest := balance.Mul(monthlyRate).Round(2)

		// Собираем явные платежи за этот месяц.
		totalAmount := decimal.Zero
		strategy := ReduceTerm
		hasPayment := false
		for keyIdx < len(dateKeys) && !dateKeys[keyIdx].After(payDate) {
			s := slotMap[dateKeys[keyIdx]]
			totalAmount = totalAmount.Add(s.amount)
			strategy = s.strategy
			hasPayment = true
			keyIdx++
		}
		// Если явных платежей нет — платим текущий аннуитет (auto-fill).
		if !hasPayment {
			totalAmount = currentA
		}

		// Погашение тела: D_k = платёж − I_k.
		// Допускаем отрицательное тело (платёж не покрыл проценты) —
		// тогда остаток растёт (капитализация процентов).
		principal := totalAmount.Sub(interest)

		// Не переплачиваем: ограничиваем тело остатком баланса.
		if principal.GreaterThan(balance) {
			principal = balance
		}

		actualPayment := interest.Add(principal)
		balance = balance.Sub(principal)

		schedule = append(schedule, ScheduleRow{
			Number:    rowNum,
			Date:      payDate,
			Payment:   actualPayment,
			Principal: principal,
			Interest:  interest,
			Balance:   balance,
		})

		// Кредит закрыт — выходим.
		if !balance.GreaterThan(oneKopeck) {
			break
		}

		remainMonths--
		if remainMonths <= 0 {
			break
		}

		/// 4. Обработка стратегий
		// Срабатывает только когда платёж превысил текущий аннуитет
		// (досрочное погашение).
		if !hasPayment || !totalAmount.GreaterThan(currentA) {
			continue
		}

		switch strategy {
		case ReduceTerm:
			// Пересчитываем оставшийся срок: сколько месяцев потребуется,
			// чтобы погасить остаток текущим аннуитетом.
			if newN := m.ceilMonths(balance, currentA); newN < remainMonths {
				remainMonths = newN
			}
		case ReducePayment:
			// Пересчитываем аннуитет на оставшийся срок.
			// Все будущие месяцы без явных платежей автоматически
			// получат новый (меньший) currentA.
			currentA = m.annuityPayment(balance, remainMonths)
		}
	}

	m.schedule = schedule
	return nil
}

// annuityPayment вычисляет аннуитетный платёж по формуле
//
//	A = P · r(1+r)ⁿ / ((1+r)ⁿ − 1)
//
// где r — месячная ставка, n — срок в месяцах. При r = 0 сводится к равному
// делению суммы на срок.
func (m *Mortgage) annuityPayment(sum decimal.Decimal, months int) decimal.Decimal {
	if months <= 0 {
		months = 1
	}
	if m.rate.IsZero() {
		return sum.DivRound(decimal.NewFromInt(int64(months)), 2)
	}

	one := decimal.NewFromInt(1)
	r := m.rate.DivRound(decimal.NewFromInt(12), intermediatePrecision)
	onePlusR := one.Add(r)
	onePlusRN := decimalPow(onePlusR, months)

	numerator := r.Mul(onePlusRN)
	denominator := onePlusRN.Sub(one)
	factor := numerator.DivRound(denominator, intermediatePrecision)

	return sum.Mul(factor).Round(2)
}

// ceilMonths вычисляет минимальное количество месяцев, за которое платёж
// payment погасит остаток balance при текущей годовой ставке.
//
// Использует бинарный поиск по формуле аннуитета (без float64): находит
// наименьшее n, при котором annuityPayment(balance, n) ≤ payment.
//
// Если платёж не покрывает даже проценты первого месяца, возвращает m.period
// (погашение невозможно при данном платеже).
func (m *Mortgage) ceilMonths(balance, payment decimal.Decimal) int {
	monthlyRate := m.rate.DivRound(decimal.NewFromInt(12), intermediatePrecision)

	// Если платёж не покрывает проценты, погашение невозможно.
	firstInterest := balance.Mul(monthlyRate).Round(2)
	if payment.LessThanOrEqual(firstInterest) {
		return m.period
	}

	// Бинарный поиск: наименьшее n, при котором annuityPayment(balance, n) ≤ payment.
	// annuityPayment мононно убывает по n, поэтому бинарный поиск корректен.
	lo, hi := 1, m.period
	for lo < hi {
		mid := (lo + hi) / 2
		if m.annuityPayment(balance, mid).LessThanOrEqual(payment) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

// decimalPow возводит base в степень exp бинарным методом. Используется в
// annuityPayment; exp ограничен сверху m.period (≤ 600).
func decimalPow(base decimal.Decimal, exp int) decimal.Decimal {
	result := decimal.NewFromInt(1)
	for exp > 0 {
		if exp&1 == 1 {
			result = result.Mul(base)
		}
		base = base.Mul(base)
		exp >>= 1
	}
	return result
}
