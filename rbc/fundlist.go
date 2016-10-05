package rbc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type FundList struct {
	Success bool
	Data    []*Fund
}

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
}

type Funds map[string]*Fund

var (
	// Cached fund list for easy lookup
	FundCache = Funds{}
)

func (fd Funds) Query(filter string) {

	//var result struct {
	//FundCode string
	//FundName string
	//Value    float64
	//}
	//rList := []result{}
	//for _, f := range fd {
	//p, err := NewPortfolio([]string{f.FundCode})
	//if err != nil {
	//panic(err)
	//}
	//}

}

func GetFundList(fundSeries string) error {
	// check if cached
	var (
		reader  io.Reader
		doWrite bool
	)

	if f, err := os.Open("./data/fundList.json"); err == nil {
		reader = f
	} else {
		payload := url.Values{}
		payload.Set("fundGroup", "RBC Funds")
		payload.Set("fundSeries", fundSeries)
		res, err := http.Post(fmt.Sprintf("%s/fundList", ApiURL), ContentType, bytes.NewBuffer([]byte(payload.Encode())))
		if err != nil {
			return err
		}
		defer res.Body.Close()
		reader = res.Body
		doWrite = true
	}

	rawFundList, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	if doWrite {
		// save it to the data dir
		f, err := os.Create("./data/fundList.json")
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := f.Write(rawFundList); err != nil {
			return err
		}
	}

	var fundList *FundList
	if err := json.Unmarshal(rawFundList, &fundList); err != nil {
		return err
	}

	// cache the results
	for _, f := range fundList.Data {
		FundCache[f.FundCode] = f
	}

	return nil
}
