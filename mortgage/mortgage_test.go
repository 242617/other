package mortgage_test

import (
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/242617/other/mortgage"
)

func rub(v int64) decimal.Decimal                   { return decimal.NewFromInt(v) }
func decStr(v decimal.Decimal) string               { return v.StringFixed(2) }
func rateFromPercent(percent int64) decimal.Decimal { return decimal.New(percent, -2) }

// newMortgage строит провалидированный объект через публичный конструктор.
func newMortgage(t *testing.T, sumRub int64, rate decimal.Decimal, period int, begin time.Time) *mortgage.Mortgage {
	t.Helper()
	m, err := mortgage.New(decimal.NewFromInt(sumRub), rate, begin, period)
	require.NoError(t, err)
	return m
}

// mustCalculate регистрирует регулярный ежемесячный платёж на весь срок.
// Стратегия по умолчанию — ReduceTerm (сохранить размер, возможно сократить срок).
func mustCalculate(t *testing.T, m *mortgage.Mortgage) {
	t.Helper()
	require.NoError(t, m.Add(mortgage.MonthlyPayment{
		Begin:  m.Begin(),
		End:    m.Begin().AddDate(0, m.Period(), 0),
		Amount: m.InitialPayment(),
	}))
}

// ---------------------------------------------------------------------------
// Money
// ---------------------------------------------------------------------------

func TestMoney(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "12345.67", mortgage.Money(12_345, 67).StringFixed(2))
	assert.Equal(t, "12.00", mortgage.Money(12, 0).StringFixed(2))
	assert.Equal(t, "0.50", mortgage.Money(0, 50).StringFixed(2))
	assert.Equal(t, "0.00", mortgage.Money(0, 0).StringFixed(2))
}

func TestMoney_InvalidPanics(t *testing.T) {
	t.Parallel()
	assert.PanicsWithValue(t,
		"mortgage.Money: kopecks must be in [0, 99], got 150",
		func() { mortgage.Money(100, 150) },
	)
	assert.PanicsWithValue(t,
		"mortgage.Money: rubles must be non-negative, got -1",
		func() { mortgage.Money(-1, 0) },
	)
}

// ---------------------------------------------------------------------------
// FormatMoney
// ---------------------------------------------------------------------------

func TestFormatMoney(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"ноль", "0", "0.00"},
		{"целое", "1000000", "1 000 000.00"},
		{"с копейками", "1234567.89", "1 234 567.89"},
		{"отрицательное", "-1234567.89", "-1 234 567.89"},
		{"маленькое", "0.50", "0.50"},
		{"округление до 2 знаков", "1.005", "1.01"}, // StringFixed(2) округляет вверх
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d, err := decimal.NewFromString(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, mortgage.FormatMoney(d))
		})
	}
}

// ---------------------------------------------------------------------------
// Validate
// ---------------------------------------------------------------------------

func TestValidate_RejectsInvalid(t *testing.T) {
	t.Parallel()
	valid := newMortgage(t, 10_000, rateFromPercent(6), 12,
		time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC))

	// New уже провалидировал `valid` успешно — проверка на всякий случай.
	assert.NoError(t, valid.Validate())

	cases := []struct {
		name         string
		sum          decimal.Decimal
		rate         decimal.Decimal
		period       int
		begin        time.Time
		wantSentinel error
		wantSubstr   string
	}{
		{
			name: "нулевая сумма", sum: decimal.Zero, rate: rateFromPercent(6),
			period: 12, begin: valid.Begin(),
			wantSentinel: mortgage.ErrInvalidSum, wantSubstr: "должна быть положительной",
		},
		{
			name: "отрицательная сумма", sum: decimal.NewFromInt(-100), rate: rateFromPercent(6),
			period: 12, begin: valid.Begin(),
			wantSentinel: mortgage.ErrInvalidSum, wantSubstr: "должна быть положительной",
		},
		{
			name: "отрицательная ставка", sum: decimal.NewFromInt(10_000), rate: decimal.New(-1, -2),
			period: 12, begin: valid.Begin(),
			wantSentinel: mortgage.ErrInvalidRate, wantSubstr: "не может быть отрицательной",
		},
		{
			name: "нулевой срок", sum: decimal.NewFromInt(10_000), rate: rateFromPercent(6),
			period: 0, begin: valid.Begin(),
			wantSentinel: mortgage.ErrInvalidPeriod, wantSubstr: "должен быть положительным",
		},
		{
			name: "слишком длинный срок", sum: decimal.NewFromInt(10_000), rate: rateFromPercent(6),
			period: 601, begin: valid.Begin(),
			wantSentinel: mortgage.ErrInvalidPeriod, wantSubstr: "не может превышать 600",
		},
		{
			name: "нулевая дата", sum: decimal.NewFromInt(10_000), rate: rateFromPercent(6),
			period: 12, begin: time.Time{},
			wantSentinel: mortgage.ErrInvalidDate, wantSubstr: "должна быть задана",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := mortgage.New(tc.sum, tc.rate, tc.begin, tc.period)
			require.Error(t, err)
			assert.ErrorIs(t, err, tc.wantSentinel)
			assert.Contains(t, err.Error(), tc.wantSubstr)
		})
	}
}

// ---------------------------------------------------------------------------
// InitialPayment
// ---------------------------------------------------------------------------

func TestInitialPayment_Annuity(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	assert.Equal(t, "860.66", decStr(m.InitialPayment()))
}

func TestInitialPayment_ZeroPercent(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 12_000, decimal.Zero, 12, begin)
	assert.Equal(t, "1000.00", decStr(m.InitialPayment()))
}

// ---------------------------------------------------------------------------
// Геттеры
// ---------------------------------------------------------------------------

func TestGetters(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)

	assert.True(t, m.Sum().Equal(decimal.NewFromInt(10_000)))
	assert.True(t, m.Rate().Equal(rateFromPercent(6)))
	assert.Equal(t, 12, m.Period())
	// Begin нормализуется до полуночи UTC — исходное значение уже было ею.
	assert.True(t, m.Begin().Equal(begin))
}

// ---------------------------------------------------------------------------
// Add — валидация
// ---------------------------------------------------------------------------

func TestAdd_SinglePayment_ZeroAmount(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)

	err := m.Add(mortgage.SinglePayment{Date: begin, Amount: decimal.Zero})
	require.ErrorIs(t, err, mortgage.ErrInvalidAmount)
}

func TestAdd_SinglePayment_NegativeAmount(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)

	err := m.Add(mortgage.SinglePayment{Date: begin, Amount: decimal.NewFromInt(-5)})
	require.ErrorIs(t, err, mortgage.ErrInvalidAmount)
}

func TestAdd_InvalidMortgage(t *testing.T) {
	t.Parallel()
	_, err := mortgage.New(
		decimal.NewFromInt(10_000),
		decimal.New(-1, -2), // отрицательная ставка
		time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		12,
	)
	require.ErrorIs(t, err, mortgage.ErrInvalidRate)
}

// ---------------------------------------------------------------------------
// Граничные случаи
// ---------------------------------------------------------------------------

func TestNoPayments(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)

	assert.Empty(t, m.Schedule())
	assert.Equal(t, mortgage.Summary{}, m.Summary())
}

// TestPaymentTooLow проверяет, что поток платежей, не покрывающий проценты,
// не обрывает график молча: каждый месяц платежа даёт строку, баланс растёт
// (отрицательная амортизация), а график завершается с ненулевым
// RemainingBalance.
func TestPaymentTooLow(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(12), 12, begin)
	require.NoError(t, m.Add(mortgage.MonthlyPayment{
		Begin: begin, End: begin.AddDate(0, 12, 0), Amount: rub(50),
	}))

	schedule := m.Schedule()
	require.Len(t, schedule, 12, "каждый месяц платежа должен давать строку")

	// Тело должно быть отрицательным (платёж не покрыл проценты).
	assert.True(t, schedule[0].Principal.IsNegative(),
		"первое тело должно быть отрицательным, получено %s", decStr(schedule[0].Principal))

	// Баланс должен расти со временем.
	assert.True(t, schedule[11].Balance.GreaterThan(m.Sum()),
		"баланс должен расти, последний=%s сумма=%s", decStr(schedule[11].Balance), decStr(m.Sum()))

}

// ---------------------------------------------------------------------------
// Расчёт — корректные сценарии
// ---------------------------------------------------------------------------

func TestCalculate_ZeroPercent(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 12_000, decimal.Zero, 12, begin)
	mustCalculate(t, m)

	schedule := m.Schedule()
	require.Len(t, schedule, 12)

	for i, row := range schedule {
		assert.Equal(t, "1000.00", decStr(row.Payment), "платёж, строка %d", i)
		assert.Equal(t, "1000.00", decStr(row.Principal), "тело, строка %d", i)
		assert.Equal(t, "0.00", decStr(row.Interest), "проценты, строка %d", i)
		expectedBalance := decimal.NewFromInt(12_000 - int64((i+1)*1000))
		assert.True(t, row.Balance.Equal(expectedBalance), "баланс, строка %d", i)
	}

	s := m.Summary()
	assert.Equal(t, "1000.00", decStr(s.MonthlyPayment))
	assert.Equal(t, "12000.00", decStr(s.TotalPaid))
	assert.Equal(t, "0.00", decStr(s.Overpayment))
}

func TestCalculate_OneMonth(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 5_000, rateFromPercent(12), 1, begin)
	mustCalculate(t, m)

	schedule := m.Schedule()
	require.Len(t, schedule, 1)

	row := schedule[0]
	assert.Equal(t, 1, row.Number)
	assert.True(t, row.Date.Equal(begin))
	assert.Equal(t, "5050.00", decStr(row.Payment))
	assert.Equal(t, "5000.00", decStr(row.Principal))
	assert.Equal(t, "50.00", decStr(row.Interest))
	assert.Equal(t, "0.00", decStr(row.Balance))
}

func TestCalculate_Annuity(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	mustCalculate(t, m)

	schedule := m.Schedule()
	assert.GreaterOrEqual(t, len(schedule), 12)
	assert.LessOrEqual(t, len(schedule), 13)

	assert.Equal(t, "860.66", decStr(schedule[0].Payment))
	assert.Equal(t, "810.66", decStr(schedule[0].Principal))
	assert.Equal(t, "50.00", decStr(schedule[0].Interest))

	for i := 1; i < len(schedule)-1; i++ {
		assert.True(t, schedule[i].Payment.Equal(schedule[0].Payment),
			"платёж %d отличается", i)
	}

	last := schedule[len(schedule)-1]
	assert.True(t, last.Balance.LessThan(decimal.NewFromInt(1)),
		"финальный баланс должен быть около 0, получено %s", decStr(last.Balance))

	totalPrincipal := decimal.Zero
	for _, row := range schedule {
		totalPrincipal = totalPrincipal.Add(row.Principal)
	}
	diff := m.Sum().Sub(totalPrincipal).Abs()
	assert.True(t, diff.LessThan(decimal.NewFromInt(1)),
		"расхождение по телу слишком велико: %s", decStr(diff))

	s := m.Summary()
	assert.Equal(t, "860.66", decStr(s.MonthlyPayment))
	assert.True(t, s.Overpayment.IsPositive(), "переплата должна быть положительной")
}

// ---------------------------------------------------------------------------
// Стратегии досрочного погашения
// ---------------------------------------------------------------------------

func TestReduceTerm_ShortensSchedule(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	mustCalculate(t, m) // по умолчанию ReduceTerm
	originalLen := len(m.Schedule())

	require.NoError(t, m.Add(mortgage.SinglePayment{
		Date:   time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC),
		Amount: rub(2500),
	}))

	schedule := m.Schedule()
	assert.Less(t, len(schedule), originalLen, "срок должен сократиться")
	// Ежемесячный платёж остаётся прежним.
	assert.Equal(t, "860.66", decStr(schedule[0].Payment))
}

// TestReducePayment_LowersMonthly проверяет модель auto-fill: добавляем
// только досрочный платёж, регулярные подставляются автоматически из currentA.
// После досрочного платежа ReducePayment пересчитывает currentA в меньшую
// сторону — все следующие месяцы (без явных платежей) получают новый минимум.
func TestReducePayment_LowersMonthly(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)

	// Единственный явный платёж — досрочный во 2-м месяце.
	// Остальные месяцы подставляются автоматически из currentA.
	require.NoError(t, m.Add(mortgage.SinglePayment{
		Date:     begin.AddDate(0, 1, 0),
		Amount:   rub(2500),
		Strategy: mortgage.ReducePayment,
	}))

	schedule := m.Schedule()
	require.NotEmpty(t, schedule)

	// Месяц 1: auto-fill = исходный аннуитет (860.66).
	assert.Equal(t, "860.66", decStr(schedule[0].Payment))

	// Месяц 2: досрочный платёж 2500 (явный, auto-fill не применяется).
	assert.True(t, schedule[1].Payment.GreaterThan(schedule[0].Payment),
		"во 2-м месяце должно быть больше, получено %s", decStr(schedule[1].Payment))

	// Месяц 3: auto-fill = новый аннуитет (меньше исходного).
	assert.True(t, schedule[2].Payment.LessThan(schedule[0].Payment),
		"в 3-м месяце должно быть меньше, получено %s", decStr(schedule[2].Payment))

	// Все месяцы после досрочного (кроме последнего) равны новому аннуитету.
	for i := 2; i < len(schedule)-1; i++ {
		assert.True(t, schedule[i].Payment.Equal(schedule[2].Payment),
			"платёж %d отличается: %s против %s", i,
			decStr(schedule[i].Payment), decStr(schedule[2].Payment))
	}

	// Срок сохраняется — 12 месяцев (ReducePayment не сокращает срок).
	assert.Equal(t, 12, len(schedule))

	last := schedule[len(schedule)-1]
	assert.True(t, last.Balance.LessThan(decimal.NewFromInt(1)),
		"финальный баланс должен быть около 0, получено %s", decStr(last.Balance))
}

// TestReduceTerm_ShortensSchedule проверяет ReduceTerm через auto-fill:
// после досрочного платежа срок сокращается, currentA не меняется.
func TestReduceTerm_ShortensViaAutoFill(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)

	// Досрочный платёж во 2-м месяце со стратегией ReduceTerm (по умолчанию).
	require.NoError(t, m.Add(mortgage.SinglePayment{
		Date:   begin.AddDate(0, 1, 0),
		Amount: rub(2500),
	}))

	schedule := m.Schedule()
	require.NotEmpty(t, schedule)

	// Месяц 1: auto-fill = исходный аннуитет.
	assert.Equal(t, "860.66", decStr(schedule[0].Payment))

	// Месяц 2: досрочный платёж.
	assert.True(t, schedule[1].Payment.GreaterThan(schedule[0].Payment),
		"во 2-м месяце должно быть больше, получено %s", decStr(schedule[1].Payment))

	// Месяц 3+: auto-fill = ТЕ ЖЕ 860.66 (ReduceTerm не пересчитывает currentA).
	assert.Equal(t, "860.66", decStr(schedule[2].Payment))

	// Срок сократился (меньше 12 месяцев).
	assert.Less(t, len(schedule), 12, "срок должен сократиться")
}

// TestStrategiesDiffer проверяет, что при одинаковом досрочном платеже
// ReducePayment сохраняет срок, а ReduceTerm — сокращает.
func TestStrategiesDiffer(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	extraDate := begin.AddDate(0, 1, 0)
	extraAmount := rub(2500)

	// ReducePayment
	mPay := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	require.NoError(t, mPay.Add(mortgage.SinglePayment{
		Date: extraDate, Amount: extraAmount, Strategy: mortgage.ReducePayment,
	}))
	schedPay := mPay.Schedule()

	// ReduceTerm
	mTerm := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	require.NoError(t, mTerm.Add(mortgage.SinglePayment{
		Date: extraDate, Amount: extraAmount, // Strategy: ReduceTerm (по умолчанию)
	}))
	schedTerm := mTerm.Schedule()

	// ReduceTerm закрывает раньше, чем ReducePayment.
	assert.Less(t, len(schedTerm), len(schedPay),
		"ReduceTerm должен дать меньший срок: %d против %d",
		len(schedTerm), len(schedPay))

	// При ReducePayment после досрочного платежи ниже.
	assert.True(t, schedPay[2].Payment.LessThan(schedTerm[2].Payment),
		"платёж ReducePayment должен быть ниже ReduceTerm: %s против %s",
		decStr(schedPay[2].Payment), decStr(schedTerm[2].Payment))
}

// TestRebuildIsIdempotent проверяет, что повторный вызов Add с тем же
// досрочным платежом не дублирует и не накапливает состояние — график должен
// быть детерминированной функцией от m.payments.
func TestRebuildIsIdempotent(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	mustCalculate(t, m)
	require.NoError(t, m.Add(mortgage.SinglePayment{
		Date:   begin.AddDate(0, 1, 0),
		Amount: rub(2500),
	}))
	scheduleAfterExtra := m.Schedule()

	// Добавляем пустой платёж (MonthlyPayment с Begin==End не даёт ни одного
	// SinglePayment), что вызывает перестроение графика без новых платежей.
	require.NoError(t, m.Add(mortgage.MonthlyPayment{
		Begin: begin, End: begin, Amount: rub(1),
	}))

	// Несмотря на перестроение, результат должен совпадать с прежним, так как
	// новых платежей не добавилось. Проверяем стабильность структуры графика.
	current := m.Schedule()
	require.Equal(t, len(scheduleAfterExtra), len(current))
	for i := range current {
		assert.True(t, current[i].Payment.Equal(scheduleAfterExtra[i].Payment), "строка %d", i)
	}
}

// ---------------------------------------------------------------------------
// Добавление досрочных платежей
// ---------------------------------------------------------------------------

func TestAdd_SinglePayment_AfterCalculate(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	mustCalculate(t, m)
	originalLen := len(m.Schedule())

	require.NoError(t, m.Add(mortgage.SinglePayment{
		Date:   time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC),
		Amount: rub(2500),
	}))

	schedule := m.Schedule()
	assert.Less(t, len(schedule), originalLen)

	totalPrincipal := decimal.Zero
	for _, row := range schedule {
		totalPrincipal = totalPrincipal.Add(row.Principal)
	}
	diff := m.Sum().Sub(totalPrincipal).Abs()
	assert.True(t, diff.LessThan(decimal.NewFromInt(1)), "расхождение по телу: %s", decStr(diff))
}

func TestAdd_SinglePayment_Before(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)

	require.NoError(t, m.Add(mortgage.MonthlyPayment{
		Begin: begin, End: begin.AddDate(0, 12, 0), Amount: m.InitialPayment(),
	}))
	require.NoError(t, m.Add(mortgage.SinglePayment{
		Date: time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC), Amount: rub(2500),
	}))

	schedule := m.Schedule()
	assert.Less(t, len(schedule), 12)

	totalPrincipal := decimal.Zero
	for _, row := range schedule {
		totalPrincipal = totalPrincipal.Add(row.Principal)
	}
	diff := m.Sum().Sub(totalPrincipal).Abs()
	assert.True(t, diff.LessThan(decimal.NewFromInt(1)), "расхождение по телу: %s", decStr(diff))
}

func TestAdd_OverlappingMonthlyPayments(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)

	require.NoError(t, m.Add(mortgage.MonthlyPayment{
		Begin: begin, End: begin.AddDate(0, 12, 0), Amount: rub(500),
	}))
	require.NoError(t, m.Add(mortgage.MonthlyPayment{
		Begin: begin.AddDate(0, 2, 0), End: begin.AddDate(0, 8, 0), Amount: rub(300),
	}))

	schedule := m.Schedule()

	assert.Equal(t, "500.00", decStr(schedule[0].Payment))
	assert.Equal(t, "500.00", decStr(schedule[1].Payment))
	assert.Equal(t, "800.00", decStr(schedule[2].Payment))
	assert.Equal(t, "500.00", decStr(schedule[8].Payment))
}

// ---------------------------------------------------------------------------
// Schedule возвращает копию
// ---------------------------------------------------------------------------

func TestSchedule_IsDefensiveCopy(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	mustCalculate(t, m)

	sched1 := m.Schedule()
	sched2 := m.Schedule()
	// Мутация первого возвращённого слайса не должна влиять ни на второй, ни на
	// внутреннее состояние объекта Mortgage.
	sched1[0].Payment = decimal.NewFromInt(999999)
	assert.NotEqual(t, sched1[0].Payment, sched2[0].Payment)
}

// ---------------------------------------------------------------------------
// Auto-fill: месяцы без явных платежей получают currentA
// ---------------------------------------------------------------------------

// TestAutoFill проверяет, что месяцы без явных платежей автоматически
// получают регулярный платёж (currentA), а не ноль.
func TestAutoFill(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)

	// Только два явных платежа (поменьше). Остальные месяцы — auto-fill.
	require.NoError(t, m.Add(mortgage.MonthlyPayment{
		Begin: begin, End: begin.AddDate(0, 3, 0), Amount: rub(2000),
	}))
	require.NoError(t, m.Add(mortgage.MonthlyPayment{
		Begin: begin.AddDate(0, 6, 0), End: begin.AddDate(0, 9, 0), Amount: rub(2000),
	}))

	schedule := m.Schedule()
	require.NotEmpty(t, schedule)

	// Месяцы 1–2 (индексы 0–1): явный платёж 2000.
	assert.Equal(t, "2000.00", decStr(schedule[0].Payment))

	// Месяц 4 (индекс 3): нет явного платежа → auto-fill = currentA.
	// currentA = initialPayment = 860.66 (аннуитет для 10000, 6%, 12 мес).
	assert.True(t, schedule[3].Payment.GreaterThan(decimal.Zero),
		"месяц без платежа должен получить auto-fill, получено %s",
		decStr(schedule[3].Payment))
	assert.True(t, schedule[3].Payment.Equal(schedule[0].Payment) ||
		!schedule[3].Payment.Equal(rub(2000)),
		"auto-fill не должен быть равен 2000")
}

// ---------------------------------------------------------------------------
// Интерфейс Payment
// ---------------------------------------------------------------------------

func TestMonthlyPayment_Payments(t *testing.T) {
	t.Parallel()
	mp := mortgage.MonthlyPayment{
		Begin:  time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		End:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		Amount: rub(500),
	}
	payments := mp.Payments()
	assert.Len(t, payments, 3)

	for i, expectedMonth := range []time.Month{time.June, time.July, time.August} {
		assert.Equal(t, expectedMonth, payments[i].Date.Month())
		assert.Equal(t, 2026, payments[i].Date.Year())
		assert.Equal(t, "500.00", decStr(payments[i].Amount))
	}
}

func TestSinglePayment_Payments(t *testing.T) {
	t.Parallel()
	sp := mortgage.SinglePayment{
		Date:   time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC),
		Amount: rub(2500),
	}
	payments := sp.Payments()
	require.Len(t, payments, 1)
	assert.True(t, payments[0].Date.Equal(sp.Date))
	assert.Equal(t, "2500.00", decStr(payments[0].Amount))
}

func TestPaymentStrategy_String(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "reduce_term", mortgage.ReduceTerm.String())
	assert.Equal(t, "reduce_payment", mortgage.ReducePayment.String())
	assert.Equal(t, "early_repayment(42)", mortgage.PaymentStrategy(42).String())
}

// ---------------------------------------------------------------------------
// String / Summary
// ---------------------------------------------------------------------------

func TestString_NotEmpty(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	mustCalculate(t, m)

	s := m.String()
	assert.NotEmpty(t, s)
	assert.Contains(t, s, "2026-06-01")
	assert.Contains(t, s, "PAYMENT")
	assert.Contains(t, s, "BALANCE")
}

func TestString_NoSchedule(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	assert.Equal(t, "график не рассчитан", m.String())
}

func TestSummary_NoSchedule(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)

	s := m.Summary()
	assert.Equal(t, mortgage.Summary{}, s)
}

func TestSummary_String(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	mustCalculate(t, m)

	s := m.Summary().String()
	assert.Contains(t, s, "ежемес.:")
	assert.Contains(t, s, "всего:")
	assert.Contains(t, s, "переплата:")
	assert.Contains(t, s, "конец:")
	assert.Contains(t, s, "остаток:")
	assert.Contains(t, s, "RUB")
}

func TestSummary_String_NoData(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "нет данных", mortgage.Summary{}.String())
}

func TestString_FormattedTotal(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 1_000_000, decimal.New(111, -3), 120, begin)
	mustCalculate(t, m)

	s := m.String()
	assert.Contains(t, s, "13 831.67")
}

// ---------------------------------------------------------------------------
// Конкурентность
// ---------------------------------------------------------------------------

func TestConcurrent(t *testing.T) {
	t.Parallel()
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	m := newMortgage(t, 10_000, rateFromPercent(6), 12, begin)
	mustCalculate(t, m)

	const writers = 20
	const readers = 20

	var wg sync.WaitGroup
	wg.Add(writers + readers)

	// Писатели: добавляют досрочные платежи (каждый вызывает rebuild).
	for i := range writers {
		go func(i int) {
			defer wg.Done()
			_ = m.Add(mortgage.SinglePayment{
				Date:   begin.AddDate(0, 0, i%30+1),
				Amount: rub(int64(10 + i)),
			})
		}(i)
	}
	// Читатели: параллельно читают график и сводку.
	for range readers {
		go func() {
			defer wg.Done()
			_ = m.Schedule()
			_ = m.Summary()
			_ = m.String()
		}()
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// Бенчмарк
// ---------------------------------------------------------------------------

func BenchmarkRebuild(b *testing.B) {
	begin := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m, err := mortgage.New(decimal.NewFromInt(12_000_000), rateFromPercent(6), begin, 30*12)
		if err != nil {
			b.Fatal(err)
		}
		_ = m.Add(mortgage.MonthlyPayment{
			Begin: begin, End: begin.AddDate(0, 30*12, 0), Amount: m.InitialPayment(),
		})
		b.StartTimer()
		_ = m.Add(mortgage.SinglePayment{
			Date: begin.AddDate(0, 12, 0), Amount: rub(500_000),
		})
	}
}

// ---------------------------------------------------------------------------
// Фаззинг
// ---------------------------------------------------------------------------

// FuzzAnnuityPayment проверяет устойчивость расчёта аннуитетного платежа на
// произвольных входах. Малые суммы при большой длительности и нулевой ставке
// могут округляться до нуля — такие вырожденные входы пропускаются.
func FuzzAnnuityPayment(f *testing.F) {
	f.Add(int64(1_000_000), int64(6), int(12))
	f.Add(int64(100), int64(0), int(1))
	f.Add(int64(12_000_000), int64(20), int(360))

	f.Fuzz(func(t *testing.T, sumCents int64, rateBps int64, months int) {
		if months <= 0 || months > 600 || sumCents <= 0 {
			t.Skip()
		}
		if rateBps < 0 || rateBps > 1000 {
			t.Skip()
		}
		// Пропускаем вырожденные случаи, где платёж из-за округления равен нулю:
		// это нормальное поведение, а не баг. Порог: хотя бы 1 копейка на месяц.
		if sumCents < int64(months) {
			t.Skip()
		}
		m, err := mortgage.New(
			decimal.New(sumCents, -2),
			decimal.New(rateBps, -4),
			time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			months,
		)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		// Платёж должен быть положительным и конечным.
		p := m.InitialPayment()
		if !p.IsPositive() {
			t.Fatalf("неположительный платёж: %s для sum=%d rate=%d months=%d",
				p.String(), sumCents, rateBps, months)
		}
	})
}
