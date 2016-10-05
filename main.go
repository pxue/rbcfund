// this script scrapes rbc portfolio and attempts to
// find the best fitting solution for you. it wont replace
// your advisor, but it will help you sift through the hundres of solutions they
// offer

package main

import (
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/pxue/rbcfund/rbc"
)

var (
	app        = kingpin.New("rbcfund", "A CLI RBC fund manager")
	fundSeries = app.Flag("series", "RBC fund series code").Default("A").String()
	flushCache = portfolio.Flag("flush", "flush and rewrite the fund performance cache").Bool()

	portfolio      = app.Command("portfolio", "Generate a new portfolio")
	portfolioFunds = portfolio.Arg("funds", "List of funds").Required().Strings()

	summarize          = app.Command("summarize", "Generate a summarized report of all RBC funds")
	summarizeSortField = summarize.Arg("sort", "Field to aggregate and sort the funds on").Required().String()
)

func init() {
	if err := rbc.GetFundList(*fundSeries); err != nil {
		panic(err)
	}
	rbc.Setup(*flushCache)
}

func main() {
	app.Parse(os.Args[1:])
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case portfolio.FullCommand():
		portfolio, err := rbc.NewPortfolio(*portfolioFunds)
		if err != nil {
			panic(err)
		}
		portfolio.PrintSummary()
	case summarize.FullCommand():
		rbc.FundCache.Query(*summarizeSortField)
	}

}
