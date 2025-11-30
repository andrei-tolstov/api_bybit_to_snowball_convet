package main

import (
	"context"
	"encoding/json"
	"fmt"
	bybit "github.com/andrei-tolstov/bybit.go.api.yield.history"
	models "github.com/andrei-tolstov/bybit.go.api.yield.history/models"
)

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
