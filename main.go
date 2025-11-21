package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	bybit "github.com/andrei-tolstov/bybit.go.api.yield.history"
	models "github.com/andrei-tolstov/bybit.go.api.yield.history/models"
	// "github.com/mitchellh/mapstructure"
)




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

	// client := bybit.NewBybitHttpClient(BYBIT_API_KEY, BYBIT_API_SECRET, bybit.WithBaseURL(bybit.MAINNET))
	// максимальный срок хранения 2года
	now := time.Now()
	endDate := now
	// endDate := now.UnixMilli()
    startDate := now.AddDate(-2, 0, 0)
	// deposit
	// https://bybit-exchange.github.io/docs/v5/asset/deposit/deposit-record
	// 
	// GetDepositRecords(client, startDate, endDate)
	GetTimeSlice(startDate, endDate, 30)



	// торговый аккаунт
	// GetTransaction(BYBIT_API_KEY, BYBIT_API_SECRET)
	// актуальный баланс
	// GetAccountWallet(BYBIT_API_KEY, BYBIT_API_SECRET)


	// earn order history
	// GetEarnRedeemOrder(BYBIT_API_KEY, BYBIT_API_SECRET)
	// get earn out history
	// GetYieldHistory(BYBIT_API_KEY, BYBIT_API_SECRET)

	
}

// получаю записи о депозитах 
func GetDepositRecords(client *bybit.Client, startDate time.Time, endDate time.Time) {

        params := map[string]interface{}{"limit": 50, "startTime": currentStart.UnixMilli(), "endTime": currentEnd.UnixMilli()}
		serverResult, err := client.NewUtaBybitServiceWithParams(params).GetDepositRecords(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(bybit.PrettyPrint(serverResult.Result))
    
}

func GetTimeSlice(startDate time.Time, endDate time.Time, days int) []map[string]int64 {
    currentStart := startDate
	timeChunk := []map[string]int64{}
    for currentStart.Before(endDate) {
        currentEnd := currentStart.AddDate(0, 0, days)
        if currentEnd.After(endDate) {
            currentEnd = endDate
        }
		currentIterator := make(map[string]int64)
		currentIterator["start"] = currentStart.UnixMilli()
		currentIterator["end"] = currentEnd.UnixMilli()
		timeChunk = append(timeChunk, currentIterator)
		currentStart = currentEnd
    }
	return timeChunk
}




func GetTransaction(BYBIT_API_KEY, BYBIT_API_SECRET string) {
	client := bybit.NewBybitHttpClient(BYBIT_API_KEY, BYBIT_API_SECRET, bybit.WithBaseURL(bybit.MAINNET))
	params := map[string]interface{}{"accountType": "UNIFIED", "startTime": 1762290000000, "endTime": 1762549200000}
	accountResult, err := client.NewUtaBybitServiceWithParams(params).GetTransactionLog(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	// конвертирую interface в struct
	var httpAnsver models.TransactionLogInfo
	tempBytes, _ := json.Marshal(accountResult.Result)
	err = json.Unmarshal(tempBytes, &httpAnsver)
	if err != nil {
		panic(err)
	}
	fmt.Println(httpAnsver.List[0].Type)
}

func GetAccountWallet(BYBIT_API_KEY, BYBIT_API_SECRET string) {
	client := bybit.NewBybitHttpClient(BYBIT_API_KEY, BYBIT_API_SECRET, bybit.WithBaseURL(bybit.MAINNET))
	params := map[string]interface{}{"accountType": "UNIFIED"}
	accountResult, err := client.NewUtaBybitServiceWithParams(params).GetAccountWallet(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(accountResult))
}



func GetEarnRedeemOrder(BYBIT_API_KEY, BYBIT_API_SECRET string) {
	client := bybit.NewBybitHttpClient(BYBIT_API_KEY, BYBIT_API_SECRET, bybit.WithBaseURL(bybit.MAINNET))
	params := map[string]interface{}{"category": "FlexibleSaving", "startTime": 1762290000000, "endTime": 1762549200000}
	serverResult, err := client.NewUtaBybitServiceWithParams(params).GetEarnRedeemOrder(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(serverResult))
}

func GetYieldHistory(BYBIT_API_KEY, BYBIT_API_SECRET string) {
	client := bybit.NewBybitHttpClient(BYBIT_API_KEY, BYBIT_API_SECRET, bybit.WithBaseURL(bybit.MAINNET))
	params := map[string]interface{}{"category": "FlexibleSaving"}
	serverResult, err := client.NewUtaBybitServiceWithParams(params).GetYieldHistory(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(serverResult))
}

