package utils

import (
	"encoding/csv"
	"go.uber.org/zap"
	"math/big"
	"os"
	"strconv"
	"time"
)

type GasData struct {
	Average float64 `json:"average"`
}

/*
The data for daily gas average is retrieved from https://etherscan.io/chart/gasprice
*/
func ReadMarch2025AverageGasPriceFromCSV() *big.Int {
	// Open the CSV file
	file, err := os.Open("utils/gasData.csv")
	if err != nil {
		zap.S().Errorln("Error:", err)
		return nil
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		zap.S().Errorln("Error:", err)
		return nil
	}

	now := time.Now()
	// Calculate the threshold for 30 days ago
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	avgPrice := big.NewInt(0)
	numberOfDays := 0
	// Print the CSV data
	for _, row := range records {

		timeStamp := row[1]
		unixTime, err := strconv.ParseInt(timeStamp, 10, 64)
		if err != nil {
			continue
		}

		t := time.Unix(unixTime, 0)
		if t.After(thirtyDaysAgo) || t.Equal(thirtyDaysAgo) {
			value, _ := new(big.Int).SetString(row[2], 10)
			numberOfDays = numberOfDays + 1
			avgPrice.Add(avgPrice, value)
			//zap.S().Info("date : ", t, "\t gas price: ", value, "\t total price: ", avgPrice)
		}

	}

	avgPrice.Div(avgPrice, big.NewInt(int64(numberOfDays)))

	return avgPrice
}
