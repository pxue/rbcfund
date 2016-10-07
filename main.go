// this script scrapes rbc portfolio and attempts to
// find the best fitting solution for you. it wont replace
// your advisor, but it will help you sift through the hundres of solutions they
// offer

package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/pxue/rbcfund/rbc"
)

var (
	app        = kingpin.New("rbcfund", "A CLI RBC fund manager")
	flushCache = app.Flag("flush", "flush and rewrite the fund performance cache").Bool()
	fundSeries = app.Flag("series", "RBC fund series code").Default("A").String()

	fund      = app.Command("fund", "Look up a single fund")
	portfolio = app.Command("portfolio", "Generate a new portfolio")
	summarize = app.Command("summarize", "Generate a summarized report of all RBC funds")

	fundCode = fund.Arg("code", "Fund code to lookup").Required().String()

	portfolioFunds     = portfolio.Arg("funds", "List of funds").Required().Strings()
	summarizeSortField = summarize.Arg("sort", "Field to aggregate and sort the funds on").Required().String()
)

func main() {
	cmd, err := app.Parse(os.Args[1:])
	if err := rbc.GetFundList(*fundSeries); err != nil {
		panic(err)
	}
	rbc.Setup(*flushCache)

	switch kingpin.MustParse(cmd, err) {
	case fund.FullCommand():
		f, ok := rbc.FundCache[*fundCode]
		if !ok {
			fmt.Errorf("Unknown fund code %s", *fundCode)
			return
		}
		fmt.Println(f.FundName)
		f.PrintData()
	case portfolio.FullCommand():
		fmt.Println(*portfolioFunds)
		portfolio, err := rbc.NewPortfolio(*portfolioFunds)
		if err != nil {
			panic(err)
		}
		portfolio.PrintSummary()
	case summarize.FullCommand():
		rbc.FundCache.Query(*summarizeSortField)
	}

}
