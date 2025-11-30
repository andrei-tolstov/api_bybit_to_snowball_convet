package main

import (
	"fmt"
	bybit "github.com/andrei-tolstov/bybit.go.api.yield.history"
	"os"
	"time"
)

var snowBallDataSlice = []snowBallData{}

var bufferWindowalFee = make(map[string]string)

func main() {
	BYBIT_API_KEY := os.Getenv("BYBIT_API_KEY")
	if BYBIT_API_KEY == "" {
		fmt.Println("Ошибка получения BYBIT_API_KEY")
		return
	}
	BYBIT_API_SECRET := os.Getenv("BYBIT_API_SECRET")
	if BYBIT_API_SECRET == "" {
		fmt.Println("Ошибка получения BYBIT_API_SECRET")
		return
	}

	client := bybit.NewBybitHttpClient(BYBIT_API_KEY, BYBIT_API_SECRET, bybit.WithBaseURL(bybit.MAINNET))
	// максимальный срок хранения 2года
	now := time.Now()
	endDate := now
	// endDate := now.UnixMilli()
	startDate := now.AddDate(-2, 0, 0)
	// deposit
	// https://bybit-exchange.github.io/docs/v5/asset/deposit/deposit-record
	//
	GetDepositRecords(client, startDate, endDate)
	for dataSliceIndex := range snowBallDataSlice {
		fmt.Println(snowBallDataSlice[dataSliceIndex])
	}

	// торговый аккаунт
	// GetTransaction(BYBIT_API_KEY, BYBIT_API_SECRET)
	// актуальный баланс
	// GetAccountWallet(BYBIT_API_KEY, BYBIT_API_SECRET)

	// earn order history
	// GetEarnRedeemOrder(BYBIT_API_KEY, BYBIT_API_SECRET)
	// get earn out history
	// GetYieldHistory(BYBIT_API_KEY, BYBIT_API_SECRET)

}
