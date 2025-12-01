package main

import (
	"context"
	"encoding/json"
	"fmt"
	bybit "github.com/andrei-tolstov/bybit.go.api.yield.history"
	"strconv"
	"time"
	"strings"
	"log"
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
	symbol := strings.ToUpper(coin + pairCoin)
	if symbol == "USDTUSDT" {
		return "1"
	}

	// Получаем временные метки 
	startTs, endTs, err := GetMinuteBoundsFromString(timeUnixMilli)
	if err != nil {
		fmt.Printf("ошибка времени: %w", err)
		return "0"
	}

	params := map[string]interface{}{
		"symbol":   symbol,
		"interval": "1", 
		"start":    startTs,
		"end":      endTs,
		"limit":    1,
	}

	serverResult, err := client.NewUtaBybitServiceWithParams(params).GetMarketKline(context.Background())
	if err != nil {
		log.Printf("ошибка API Bybit: %w", err)
		return "0"
	}

	jsonBytes, err := json.Marshal(serverResult.Result)
	if err != nil {
		log.Printf("ошибка маршалинга результата: %w", err)
		return "0"
	}

candle, err := ConvertRawKLine(jsonBytes)
	if err != nil {
		log.Printf("Внимание: не удалось получить свечу для %s: %v", symbol, err)
		return "0" // Или return "0", err если хотите прервать выполнение
	}

	return candle.OpenPrice
}

func GetMinuteBoundsFromString(timeUnixMilliStr string) (startMilli int64, endMilli int64, err error) {
	inputMilli, err := strconv.ParseInt(timeUnixMilliStr, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось преобразовать строку '%s': %w", timeUnixMilliStr, err)
	}

	t := time.UnixMilli(inputMilli)
	startOfMinute := t.Truncate(time.Minute)
	
	// Конец минуты - это начало следующей минуты
	endOfMinuteExclusive := startOfMinute.Add(time.Minute)

	return startOfMinute.UnixMilli(), endOfMinuteExclusive.UnixMilli(), nil
}

func ConvertRawKLine(jsonBytes []byte) (*MarketlineCandle, error) {
	var rawResponse KLineResponseRaw

	if err := json.Unmarshal(jsonBytes, &rawResponse); err != nil {
		return nil, fmt.Errorf("ошибка распаковки JSON: %w", err)
	}

	// Проверяем, есть ли данные.
	if len(rawResponse.List) == 0 {
		return nil, fmt.Errorf("получен пустой список свечей")
	}

	// Проверяем формат внутренней свечи (должно быть минимум 5 полей)
	firstCandle := rawResponse.List[0]
	if len(firstCandle) < 5 {
		return nil, fmt.Errorf("некорректный формат данных свечи")
	}

	return &MarketlineCandle{
		StartTime: firstCandle[0],
		OpenPrice: firstCandle[1],
		HighPrice: firstCandle[4],
	}, nil
}