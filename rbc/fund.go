package rbc

import (
	"crypto/md5"
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

func (f *Fund) PrintData() {
	f.Cache.printSummary()
}
