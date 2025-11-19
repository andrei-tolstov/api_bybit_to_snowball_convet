package main

import (
	"context"
	"fmt"
	"os"
	bybit "github.com/andrei-tolstov/bybit.go.api.yield.history"
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
	// торговый аккаунт
	// GetTransaction(BYBIT_API_KEY, BYBIT_API_SECRET)
	// актуальный баланс
	// GetAccountWallet(BYBIT_API_KEY, BYBIT_API_SECRET)
	// deposit
	// GetDepositRecords(BYBIT_API_KEY, BYBIT_API_SECRET)
	// earn order history
	// GetEarnRedeemOrder(BYBIT_API_KEY, BYBIT_API_SECRET)
	// get earn out history
	GetYieldHistory(BYBIT_API_KEY, BYBIT_API_SECRET)

	
}

func GetTransaction(BYBIT_API_KEY, BYBIT_API_SECRET string) {
	client := bybit.NewBybitHttpClient(BYBIT_API_KEY, BYBIT_API_SECRET, bybit.WithBaseURL(bybit.MAINNET))
	params := map[string]interface{}{"accountType": "UNIFIED", "startTime": 1762290000000, "endTime": 1762549200000}
	accountResult, err := client.NewUtaBybitServiceWithParams(params).GetTransactionLog(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(accountResult))
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

func GetDepositRecords(BYBIT_API_KEY, BYBIT_API_SECRET string) {
	client := bybit.NewBybitHttpClient(BYBIT_API_KEY, BYBIT_API_SECRET, bybit.WithBaseURL(bybit.MAINNET))
	params := map[string]interface{}{}
	serverResult, err := client.NewUtaBybitServiceWithParams(params).GetDepositRecords(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(serverResult))
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