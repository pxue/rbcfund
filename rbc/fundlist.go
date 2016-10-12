package rbc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type FundList struct {
	Success bool
	Data    []*Fund
}

type Funds map[string]*Fund

func (fd Funds) Query(filter string) {
	type summary struct {
		StartCap     float64
		EndCap       float64
		Distribution float64
	}
	for _, f := range fd {
		fmt.Printf(f.FundName)
		s := summary{}
		for _, d := range f.Cache {
			if d.Year == 2016 {
				// do summary
				if d.Month == 1 {
					s.StartCap = d.TotalCashflowInitial
				}
				s.EndCap = d.TotalCashflowInitial
				s.Distribution += d.TotalDistribution
			}
		}
		capRet := (s.EndCap + s.Distribution) / s.StartCap
		capRetPct := (capRet - 1.0) * 100.0
		fmt.Printf("\t%.2f%%\n", capRetPct)
	}
}

func GetFundList(fundSeries string) error {
	// check if cached
	var (
		reader  io.Reader
		doWrite bool
	)

	cached, err := os.Open("./data/fundList.json")
	if err != nil {
		log.Print(err)

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

	if reader == nil {
		// read from cache
		reader = cached
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

		// load distribution cache if available
		cache, err := os.Open(fmt.Sprintf("./data/cache/%s.json", f.Hash()))
		if err != nil {
			continue
		}

		resBody, err := ioutil.ReadAll(cache)
		if err != nil {
			continue
		}
		var wrapper *DistributionWrapper
		if err := json.Unmarshal(resBody, &wrapper); err != nil {
			continue
		}
		f.Cache = wrapper.Items
	}

	return nil
}
