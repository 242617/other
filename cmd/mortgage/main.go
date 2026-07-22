package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/shopspring/decimal"

	"github.com/242617/other/mortgage"
)

func init() { log.SetFlags(log.Lshortfile) }
func main() {
	sum := mortgage.Money(12_000_000, 0)
	rate := decimal.New(6, -2) // 0.06 = 6%
	begin := time.Date(2026, 7, 25, 0, 0, 0, 0, time.UTC)

	{
		m, err := mortgage.New(sum, rate, begin, 18*12)
		die(err)
		end := begin.AddDate(0, m.Period()-1, 0)
		die(m.Add(mortgage.MonthlyPayment{Begin: begin, End: end,
			Amount:   m.InitialPayment(),
			Strategy: mortgage.ReducePayment,
		}))
		fmt.Println(m.Summary(), "\t18 лет, минимальтный платёж")
	}
	{
		m, err := mortgage.New(sum, rate, begin, 30*12)
		die(err)
		end := begin.AddDate(0, m.Period()-1, 0)
		die(m.Add(mortgage.MonthlyPayment{Begin: begin, End: end,
			Amount:   mortgage.Money(120_000, 0),
			Strategy: mortgage.ReducePayment,
		}))
		fmt.Println(m.Summary(), "\t30 лет, 120 000 рублей")
		fmt.Println(m)
	}
}

func die(args ...any) {
	if len(args) == 0 {
		return
	}
	if err, ok := args[len(args)-1].(error); ok && err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: %s", file, line, err.Error())
		os.Exit(1)
	}
}
