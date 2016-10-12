package rbc

import (
	"crypto/md5"
	"errors"
	"fmt"
)

type Fund struct {
	MorningstarCategoryName string
	FundCode                string
	FundShareClassID        string
	Series                  string
	WSODIssue               string
	InceptionDate           int64
	FundNameEN              string
	FundNameFR              string
	InvestmentSolutionEN    string
	InvestmentSolutionFR    string
	FundName                string
	InvestmentSolution      string
	AssetClass              string

	// Cached distribution data
	Cache Distributions
	// Cached summary data
	Summary Summary
}

func (f *Fund) Hash() string {
	h := md5.New()
	h.Write([]byte(f.FundCode))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetFund(fcode string) (*Fund, error) {
	fund, ok := FundCache[fcode]
	if !ok {
		return nil, errors.New(fmt.Sprintf("unknown fund code %s", fcode))
	}

	// cache miss. load from rbc
	if fund.Cache == nil {
		p := NewPortfolio([]string{fcode})
		if err := p.GetDistribution(); err != nil {
			return nil, err
		}
		fund.Cache = p.Distributions
	}

	return fund, nil
}

func (f *Fund) PrintData() {
	f.Cache.printSummary()
}
