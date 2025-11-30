package main

import (
	"context"
	"encoding/json"
	"fmt"
	bybit "github.com/andrei-tolstov/bybit.go.api.yield.history"
	"time"
)

const (
	// Максимальный лимит записей на одну страницу API Bybit для депозита
	DEPOSIT_DEFAULT_API_LIMIT = 50
	// Интервал деления запросов на части в днях для депозита
	DEPOSIT_DAYS_CHUNK_INTERVAL = 30
)

// получаю записи о депозитах
func GetDepositRecords(client *bybit.Client, startDate time.Time, endDate time.Time) {
	timeChunk := GetTimeSlice(startDate, endDate, DEPOSIT_DAYS_CHUNK_INTERVAL)
	for timeChunkIndex := range timeChunk {
		currentStart := timeChunk[timeChunkIndex]["start"]
		currentEnd := timeChunk[timeChunkIndex]["end"]

		// Логирование в читаемом формате
		// startStr := time.UnixMilli(currentStart).Format(time.RFC3339)
		// endStr := time.UnixMilli(currentEnd).Format(time.RFC3339)
		// fmt.Printf("Запрос депозитов с %s по %s...\n", startStr, endStr)

		dataList := GetSliceDepositRecords(client, currentStart, currentEnd, "")
		for dataSliceIndex := range dataList {
			// fmt.Println(dataList[dataSliceIndex])
			// all deposit = buy

			newSnowBallData := snowBallData{
				Event:           "Buy",
				Date:            ConvertUnixString(dataList[dataSliceIndex].SuccessAt),                                   // конвертировать в формат 2020-02-01
				Symbol:          GetSymbol(dataList[dataSliceIndex].Coin),                                                // конвертировать в COIN-USD
				Price:           GetMarketKline(client, dataList[dataSliceIndex].Coin, "USDT", dataList[dataSliceIndex].SuccessAt),        // запросить цену на время депозита

				// Price:           GetMarketKline(client, dataList[dataSliceIndex].Coin, "USDT", dataList[dataSliceIndex].SuccessAt),        // запросить цену на время депозита
				Quantity:        dataList[dataSliceIndex].Amount,                                                         // кол-во
				Currency:        "USD",                                                                                   // всегда сводим к USD
				FeeTax:          GetWithdrawalFee(client, dataList[dataSliceIndex].Coin, dataList[dataSliceIndex].Chain), // расчет от комиссии за сеть по данным bybit
				Exchange:        "Bybit",
				NKD:             "0",
				FeeCurrency:     "",
				DoNotAdjustCash: "True", // Не обновлять валюту
				Note:            "deposit",
			}
			snowBallDataSlice = append(snowBallDataSlice, newSnowBallData)
		}
	}
}

func GetSliceDepositRecords(client *bybit.Client, startTime int64, endTime int64, cursor string) []DepositRecord {
	params := map[string]interface{}{"limit": DEPOSIT_DEFAULT_API_LIMIT, "startTime": startTime, "endTime": endTime, "cursor": cursor}
	if cursor == "" {
		delete(params, "cursor")
	}
	serverResult, err := client.NewUtaBybitServiceWithParams(params).GetDepositRecords(context.Background())
	if err != nil {
		fmt.Println(err)
		return []DepositRecord{}
	}
	// fmt.Println(bybit.PrettyPrint(serverResult.Result))
	var httpAnsver DepositRecords
	tempBytes, _ := json.Marshal(serverResult.Result)
	err = json.Unmarshal(tempBytes, &httpAnsver)
	if err != nil {
		fmt.Printf("Ошибка декодирования JSON: %v\n", err)
		return []DepositRecord{}
	}
	// fmt.Println(httpAnsver.List[0].Type)
	rows := httpAnsver.Rows
	if httpAnsver.NextPageCursor != "" {
		iterRows := GetSliceDepositRecords(client, startTime, endTime, httpAnsver.NextPageCursor)
		if len(iterRows) != 0 {
			rows = append(rows, iterRows...)
		}
	}
	return rows
}
