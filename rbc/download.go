package rbc

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var (
	DlURL = "http://funds.rbcgam.com/pdf/fund-pages/monthly/%s_e.pdf"
)

func DownloadPDF(series string) error {
	cwrapper := "<html><body>%s</body></html>"
	cbody := ""

	catLookup := map[string][]*Fund{}
	for _, fund := range FundCache {
		if strings.Contains(fund.FundName, "US$") {
			// skip USD funds
			continue
		}
		catLookup[fund.AssetClass] = append(catLookup[fund.AssetClass], fund)

		fileName := fmt.Sprintf("data/portfolios/%s.pdf", fund.FundCode)
		if _, err := os.Stat(fileName); !os.IsNotExist(err) {
			// already exist, skip
			log.Printf("skipping %s", fund.FundName)
			continue
		}

		output, err := os.Create(fileName)
		if err != nil {
			log.Printf(errors.Wrapf(err, "create fund %s", fund.FundName).Error())
			continue
		}

		resp, err := http.Get(fmt.Sprintf(DlURL, strings.ToLower(fund.FundCode)))
		if err != nil {
			log.Printf(errors.Wrapf(err, "dl fund %s", fund.FundName).Error())
			continue
		}
		if _, err := io.Copy(output, resp.Body); err != nil {
			log.Printf(errors.Wrapf(err, "write fund %s", fund.FundName).Error())
			continue
		}

		resp.Body.Close()
		output.Close()
		log.Printf("downloaded %s", fund.FundName)
	}

	for cls, funds := range catLookup {
		cbody += fmt.Sprintf("<h1>%s</h1>", cls)
		cbody += "<ul>"
		for _, fund := range funds {
			cbody += "<li>"
			cbody += fmt.Sprintf("<a href='./portfolios/%s.pdf' target='_blank'>%s - %s</a>", fund.FundCode, fund.FundCode, fund.FundName)
			cbody += "</li>"
		}
		cbody += "</ul>"
	}
	cwrapper = fmt.Sprintf(cwrapper, cbody)

	ptf, err := os.Create("data/portfolio.html")
	if err != nil {
		return err
	}
	if _, err := ptf.WriteString(cwrapper); err != nil {
		return err
	}
	return ptf.Close()
}
