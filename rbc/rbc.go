package rbc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type ChartOption struct {
	StartMonth int64  `json:"initMonth"`
	StartYear  int64  `json:"initYear"`
	EndMonth   int64  `json:"endMonth"`
	EndYear    int64  `json:"endYear"`
	Frequency  string `json:"frequency"` // "Monthly"
	Reinvest   string `json:"reinvest"`  // "0" or "1"
	//"SinceInceptionTotalReturn": 5.745664293439767,
	//"TotalReturnInitialInvestment": 5.59364968111582,
	//"TotalReturnInitialInvestmentWithExtraLogic": 5.59364968111582,
	//"bearMarket": "0",
	//"byAddress": null,
	//"byCompany": null,
	//"byEmail": null,
	//"byName": null,
	//"byPhone": null,
	//"byTitle": "",
	//"comments": "N/A",
	//"duration": 80,
	//"height": 294,
	//"width": 1152
	//"prepFor": " ",
}

type ChartItem struct {
	Balance      int64  `json:"balance"`
	Category     string `json:"category"`
	Contribution int64  `json:"contrib"`
	Family       string `json:"family"`
	GUID         string `json:"guid"`
	Series       string `json:"series"`
	Symbol       string `json:"symbol"`
	WD           int64  `json:"wd"`
	WSODIssue    string `json:"wsodIssue"`
}

type RBC struct {
	Opt *ChartOption

	// opt: flush the existing cached fund data
	flushCache bool
}

var (
	ApiURL      = "https://services.rbcgam.com/portfolio-tools/public/investment-performance"
	ContentType = "application/x-www-form-urlencoded; charset=UTF-8"
	// Cached fundlist
	FundCache = Funds{}

	app *RBC
)

func NewRequest(apiPath string, p *Portfolio) ([]byte, error) {
	// check cache
	cachePath := fmt.Sprintf("./data/cache/%s.json", p.nameHash)
	if !app.flushCache { // don't use cache
		if f, err := os.Open(cachePath); err == nil {
			return ioutil.ReadAll(f)
		}
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", ApiURL, apiPath), nil)

	// chart options
	optStr, _ := json.Marshal(app.Opt)
	req.AddCookie(&http.Cookie{Name: "ips.chart", Value: url.QueryEscape(string(optStr))})

	items := make([]*ChartItem, len(p.Funds))
	for i, f := range p.Funds {
		items[i] = &ChartItem{
			Balance:      25000,
			WSODIssue:    f.WSODIssue,
			Symbol:       f.FundCode,
			Category:     f.MorningstarCategoryName,
			Family:       "RBC",
			Series:       "A",
			Contribution: 0,
			WD:           0,
		}
	}
	// chart items
	itemsStr, _ := json.Marshal(items)
	req.AddCookie(&http.Cookie{Name: "ips.portfolio", Value: url.QueryEscape(string(itemsStr))})

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// write to cache
	f, err := os.Create(cachePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	f.Write(resBody)

	return resBody, nil
}

func Setup(flushCache bool) {
	app = &RBC{flushCache: flushCache}
	// TODO: configurable start/end times
	app.Opt = &ChartOption{
		StartMonth: 1,
		StartYear:  2013,
		EndMonth:   10,
		EndYear:    2016,
		Frequency:  "Monthly",
		Reinvest:   "0",
	}
}
