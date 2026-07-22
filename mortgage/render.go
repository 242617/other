package mortgage

import (
	"bytes"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shopspring/decimal"
)

// String отрисовывает график в виде ASCII-таблицы со строкой итогов.
func (m *Mortgage) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.schedule) == 0 {
		return "График не рассчитан"
	}

	var (
		buf           bytes.Buffer
		totalPaid     decimal.Decimal
		totalInterest decimal.Decimal
	)

	t := table.NewWriter()
	t.SetOutputMirror(&buf)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"#", "Дата", "Платёж", "Тело кредита", "Проценты", "Остаток"})

	for _, row := range m.schedule {
		t.AppendRow(table.Row{
			row.Number,
			row.Date.Format(TimeFormat),
			FormatMoney(row.Payment),
			FormatMoney(row.Principal),
			FormatMoney(row.Interest),
			FormatMoney(row.Balance),
		})
		totalPaid = totalPaid.Add(row.Payment)
		totalInterest = totalInterest.Add(row.Interest)
	}

	t.AppendSeparator()
	t.AppendFooter(table.Row{
		"", "", FormatMoney(totalPaid),
		FormatMoney(m.sum),
		FormatMoney(totalInterest), "",
	})

	t.Render()
	return strings.TrimRight(buf.String(), "\n")
}
