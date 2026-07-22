package mortgage

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// TimeFormat — формат YYYY-MM-DD, используемый в пакете для отображения
// значений [time.Time].
const TimeFormat = "2006-01-02"

// Money формирует сумму в рублях из рублёвой и копеечной частей.
//
//	Money(12_345, 67)  → 12 345.67 RUB
//	Money(0, 50)       → 0.50 RUB
//
// Оба аргумента должны быть неотрицательными, а kopecks — в диапазоне [0, 99];
// иначе функция паникует. Отрицательные суммы не имеют смысла для предметной
// области (кредиты, тела, платежи), с которой работает пакет.
func Money(rubles int64, kopecks int64) decimal.Decimal {
	if rubles < 0 {
		panic(fmt.Sprintf("mortgage.Money: rubles must be non-negative, got %d", rubles))
	}
	if kopecks < 0 || kopecks > 99 {
		panic(fmt.Sprintf("mortgage.Money: kopecks must be in [0, 99], got %d", kopecks))
	}
	return decimal.New(rubles*100+kopecks, -2)
}

// FormatMoney форматирует сумму в рублях с двумя знаками после запятой и
// разделителем разрядов (пробел без NBSP), например "1 234 567.89" или
// "-1 234 567.89".
//
// Примечание: разделитель дробной части — "." (английский стиль), а разделитель
// тысяч — " " (русский стиль). Такой гибрид намеренно используется для
// человекочитаемого вывода в CLI.
func FormatMoney(v decimal.Decimal) string {
	s := v.StringFixed(2)

	sign := ""
	if strings.HasPrefix(s, "-") {
		sign = "-"
		s = s[1:]
	}

	parts := strings.SplitN(s, ".", 2)
	rubles := parts[0]
	kopecks := parts[1]

	var buf strings.Builder
	for i, c := range rubles {
		if i > 0 && (len(rubles)-i)%3 == 0 {
			buf.WriteByte(' ')
		}
		buf.WriteRune(c)
	}

	return sign + buf.String() + "." + kopecks
}

// normalizeDate усекает время до точности дня, сохраняя календарную дату
// в исходной локации.
func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
