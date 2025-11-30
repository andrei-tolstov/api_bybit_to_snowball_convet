package main

import (
	"context"
	"encoding/json"
	"fmt"
	bybit "github.com/andrei-tolstov/bybit.go.api.yield.history"
	"strconv"
	"time"
	"strings"
)

// пока одно значение может в будущем будут еще
var notDefaultNameCurrency = map[string]string{
	"TON": "TON11419-USD",
}

func GetWithdrawalFee(client *bybit.Client, coin string, chain string) string {
	if value, ok := bufferWindowalFee[coin+"+"+chain]; ok {
		// fmt.Println("использую кэш")
		return value
	}
	params := map[string]interface{}{"coin": coin}
	serverResult, err := client.NewUtaBybitServiceWithParams(params).GetCoinInfo(context.Background())
	if err != nil {
		fmt.Errorf("ошибка GetWithdrawalFee: %w", err)
		return "0"
	}
	var apiResponse CoinInfoResult
	tempBytes, err := json.Marshal(serverResult.Result)
	if err != nil {
		fmt.Errorf("ошибка GetWithdrawalFee: %w", err)
		return "0"
	}

	if err := json.Unmarshal(tempBytes, &apiResponse); err != nil {
		fmt.Errorf("ошибка GetWithdrawalFee: %w", err)
		return "0"
	}

	for index := range apiResponse.Rows[0].Chains {
		// Комиссия фиксирована для каждой пары (Монета + Сеть)
		if apiResponse.Rows[0].Chains[index].Chain == chain {
			bufferWindowalFee[coin+"+"+chain] = apiResponse.Rows[0].Chains[index].WithdrawFee
			return bufferWindowalFee[coin+"+"+chain]
		}
	}

	//Если пара не найдена
	bufferWindowalFee[coin+"+"+chain] = "0"
	return "0"
}

func GetSymbol(symbol string) string {
	if value, ok := notDefaultNameCurrency[symbol]; ok {
		return value
	}
	return symbol + "-USD"
}

func ConvertUnixString(unixString string) string {
	// Преобразование строки в целое число
	unixMilli, err := strconv.ParseInt(unixString, 10, 64)
	if err != nil {
		fmt.Errorf("ошибка при парсинге строки UnixTime: %w", err)
		return ""
	}
	// Преобразование миллисекунд в тип time.Time
	t := time.UnixMilli(unixMilli)
	formattedDate := t.Format("2006-01-02")
	return formattedDate
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

func GetMarketKline(client *bybit.Client, coin string, pairCoin string, timeUnixMilli string) string {
	symbol := strings.ToUpper(coin + pairCoin) // like BTCUSDT, uppercase only
	if symbol == "USDTUSDT" {
		return "1"
	}
	startUnixSec, endUnixSec, err := GetMinuteBoundsFromString(timeUnixMilli)
	if err != nil {
		fmt.Errorf("ошибка получения цены: %w", err)
		return "0"
	}

	params := map[string]interface{}{"symbol": symbol, "interval": 1, "start": startUnixSec, "end": endUnixSec, "limit": "1"}
	serverResult, err := client.NewUtaBybitServiceWithParams(params).GetMarketKline(context.Background())
	if err != nil {
		fmt.Errorf("ошибка получения цены: %w", err)
		return "0"
	}
	tempBytes, err := json.Marshal(serverResult.Result)
	if err != nil {
		fmt.Errorf("ошибка получения цены: %w", err)
		return "0"
	}

	httpRespon, err := ConvertRawKLine(string(tempBytes))
	if err != nil {
		fmt.Errorf("ошибка получения цены: %w", err)
		return "0"
	}

	// fmt.Println(httpRespon[0].OpenPrice)
	return httpRespon[0].OpenPrice
}

func GetMinuteBoundsFromString(timeUnixMilliStr string) (startUnixSec int64, endUnixSec int64, err error) {
    timeUnixMilli, err := strconv.ParseInt(timeUnixMilliStr, 10, 64)
    if err != nil {
        return 0, 0, fmt.Errorf("не удалось преобразовать строку '%s' в число: %w", timeUnixMilliStr, err)
    }
	t := time.UnixMilli(timeUnixMilli)
	startOfMinute := t.Truncate(time.Minute)
	endOfMinuteExclusive := startOfMinute.Add(time.Minute)
	startUnixSec = startOfMinute.UnixMilli()
	endUnixSec = endOfMinuteExclusive.UnixMilli()
	return startUnixSec, endUnixSec, nil
}

func ConvertRawKLine(jsonString string) ([]MarketlineCandle, error) {
	var rawResponse KLineResponseRaw
    // Распаковка сырого ответа
	if err := json.Unmarshal([]byte(jsonString), &rawResponse); err != nil {
		return nil, fmt.Errorf("ошибка распаковки сырого JSON: %w", err)
	}
	var dataList []MarketlineCandle
      
	dataList = append(dataList, MarketlineCandle{
		StartTime:      rawResponse.List[0][0],
		OpenPrice:      rawResponse.List[0][1],
		HighPrice:      rawResponse.List[0][4],
	})
	
	return dataList, nil
}