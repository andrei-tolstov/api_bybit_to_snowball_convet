package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"strconv"
	bybit "github.com/andrei-tolstov/bybit.go.api.yield.history"
	models "github.com/andrei-tolstov/bybit.go.api.yield.history/models"
	// "github.com/mitchellh/mapstructure"
)

type DepositRecords struct {
	Rows           []DepositRecord `json:"rows"`
	NextPageCursor string          `json:"nextPageCursor"`
}

type DepositRecord struct {
	Coin              string `json:"coin"`
	Chain             string `json:"chain"`
	Amount            string `json:"amount"`
	TxID              string `json:"txID"`
	Status            int    `json:"status"`
	ToAddress         string `json:"toAddress"`
	Tag               string `json:"tag"`
	DepositFee        string `json:"depositFee"`
	SuccessAt         string `json:"successAt"`
	Confirmations     string `json:"confirmations"`
	TxIndex           string `json:"txIndex"`
	BlockHash         string `json:"blockHash"`
	BatchReleaseLimit string `json:"batchReleaseLimit"`
	DepositType       string `json:"depositType"`
}

const (
    // Максимальный лимит записей на одну страницу API Bybit для депозита
    DEPOSIT_DEFAULT_API_LIMIT = 50
    // Интервал деления запросов на части в днях для депозита
    DEPOSIT_DAYS_CHUNK_INTERVAL = 30
)

// https://snowball-income.com/import-transactions
type snowBallData struct {
	Event           string
	Date            string
	Symbol          string
	Price           string
	Quantity        string
	Currency        string
	FeeTax          string
	Exchange        string
	NKD             string
	FeeCurrency     string
	DoNotAdjustCash string
	Note            string
}

var snowBallDataSlice = []snowBallData{}

// пока одно значение может в будущем будут еще
var notDefaultNameCurrency = map[string]string{
	"TON": "TON11419-USD",
}

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
				Event: "Buy",
				Date: ConvertUnixString(dataList[dataSliceIndex].SuccessAt), //  конвертировать в формат 2020-02-01
				Symbol: GetSymbol(dataList[dataSliceIndex].Coin), // конвертировать в COIN-USD
				Price: "0", // запросить цену на время депозита
				Quantity: dataList[dataSliceIndex].Amount, // кол-во
				Currency: "USD", // всегда сводим к USD
				FeeTax: GetWithdrawalFee(client, dataList[dataSliceIndex].Coin, dataList[dataSliceIndex].Chain), // рассчитать от комиссии за сеть??? использовать данные bybit (могут отличаться от использованых)
				Exchange: "Bybit",
				NKD: "0",
				FeeCurrency: "",
				DoNotAdjustCash: "True", // Не обновлять валюту
				Note: "deposit",
			}
			snowBallDataSlice = append(snowBallDataSlice, newSnowBallData)
		}
	}
}

func GetWithdrawalFee(client *bybit.Client, coin string, chain string) string {
	if value, ok := bufferWindowalFee[coin + "+" + chain]; ok {
		// fmt.Println("использую кэш")
        return value
    }
    params := map[string]interface{}{"coin": coin}
    serverResult, err := client.NewUtaBybitServiceWithParams(params).GetCoinInfo(context.Background())
    if err != nil {
		fmt.Errorf("ошибка GetWithdrawalFee: %w", err)
        return "0"
    }
    var apiResponse models.CoinInfoResult
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
			bufferWindowalFee[coin + "+" + chain] = apiResponse.Rows[0].Chains[index].WithdrawFee
            return bufferWindowalFee[coin + "+" + chain]
        }
    }

    //Если пара не найдена
	bufferWindowalFee[coin + "+" + chain] = "0"
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

