package rbc

import (
	"fmt"
	"math"

	"github.com/fatih/color"
)

type DistributionWrapper struct {
	Items Distributions `json:"distributions"`
}

type Distribution struct {
	Month                          int64
	Year                           int64
	TotalCashflowInitial           float64
	TotalCashflowWithContributions float64
	TotalDistribution              float64
}

type Distributions []*Distribution

var (
	G = color.New(color.FgGreen).Add(color.Bold)
	Y = color.New(color.FgYellow)
	R = color.New(color.FgRed).Add(color.Italic)
)

func (dst Distributions) printSummary() {
	var (
		totalDist = 0.0
		yearInit  = 0.0
		yearDist  = 0.0

		initAmount = 0.0
		endAmount  = 0.0
	)

	var totalCapital float64
	for i, d := range dst {
		if i == 0 {
			yearInit = d.TotalCashflowInitial
			totalCapital = d.TotalCashflowInitial
			initAmount = yearInit
			fmt.Printf("Started with %f On %d/%d.\n", d.TotalCashflowInitial, d.Month, d.Year)
			fmt.Printf("******\t %d \t******\n", d.Year)
			continue
		}

		if d.Month == 1 {
			yearDist = 0.0
			yearInit = d.TotalCashflowInitial
			fmt.Printf("******\t %d \t******\n", d.Year)
		}

		// trend indicator
		prevCap := dst[i-1].TotalCashflowInitial
		if d.TotalCashflowInitial > prevCap {
			G.Printf("▲")
		} else {
			R.Printf("▼")
		}

		if d.TotalCashflowInitial < float64(totalCapital) {
			R.Printf("%d: %f, %f\n", d.Month, d.TotalCashflowInitial, d.TotalDistribution)
		} else if d.TotalCashflowInitial < prevCap {
			Y.Printf("%d: %f, %f\n", d.Month, d.TotalCashflowInitial, d.TotalDistribution)
		} else {
			G.Printf("%d: %f, %f\n", d.Month, d.TotalCashflowInitial, d.TotalDistribution)
		}
		totalDist += d.TotalDistribution
		yearDist += d.TotalDistribution

		if d.Month == 12 || i == len(dst)-1 {
			capRet := (d.TotalCashflowInitial + yearDist) / yearInit
			capRetPct := (capRet - 1.0) * 100.0
			fmt.Printf("Year: %d Return: %.2f%%\n\n", d.Year, capRetPct)

			endAmount = d.TotalCashflowInitial
		}
	}

	fmt.Printf("Sum Distribution: %f\n", totalDist)
	ret := math.Cbrt(((endAmount + totalDist) / initAmount))
	G.Printf("Return: %f%%\n", (ret-1.0)*100.0)
}
