package rbc

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
)

type Portfolio struct {
	Funds         []*Fund
	Distributions Distributions

	nameHash string
}

func NewPortfolio(fundNames []string) (*Portfolio, error) {
	portfolio := &Portfolio{}
	portfolio.Funds = make([]*Fund, len(fundNames))
	h := md5.New()
	for i, fcode := range fundNames {
		if f, ok := FundCache[fcode]; ok {
			fmt.Printf("%s\n", f.FundName)
			portfolio.Funds[i] = f
			h.Write([]byte(f.FundCode))
		}
	}
	portfolio.nameHash = fmt.Sprintf("%x", h.Sum(nil))

	// get distribution
	portfolio.GetDistribution()
	return portfolio, nil
}

// GetDistribution pulls the portfolio distribution data from cached data or RBC
func (p *Portfolio) GetDistribution() error {
	resBody, err := NewRequest("distributions", p)
	if err != nil {
		return err
	}

	var wrapper *DistributionWrapper
	if err := json.Unmarshal(resBody, &wrapper); err != nil {
		return err
	}

	p.Distributions = wrapper.Items
	return nil
}

func (p *Portfolio) PrintSummary() {
	// Distribution summary
	p.Distributions.printSummary()
}

func (p *Portfolio) GetPortfolioPerformance() error {
	// portfolio.GetChartData(options, items)
	return p.GetDistribution()
}
