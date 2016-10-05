package rbc

type SummaryWrapper struct {
	Success bool
	Data    string `json:"chartData"`
}

type Summary struct {
	//"BearMarketEvents": null,
	CumulativeContrib                 float64
	CumulativeWithdrawal              float64
	MaxDurationCummTotalDistrib       float64
	MaxDurationMonth                  int64
	MaxDurationTotalCFInitial         float64
	MaxDurationTotalCFWContrib        float64
	MaxDurationYear                   int64
	MinDurationMonth                  int64
	MinDurationYear                   int64
	PortfolioZeroMonth                int64
	PortfolioZeroYear                 int64
	SinceInceptionTotalReturn         float64
	SinceInceptionTotalReturnLeapYear float64
	TotalReturnInitialInvestment      float64
	initialInvestment                 int64
	totalAmtInvested                  int64
}
