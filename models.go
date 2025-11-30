package main

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

// CoinInfoResult represents the structure for coin info results.
type CoinInfoResult struct {
	Rows []CoinInfoRow `json:"rows"`
}

// CoinInfoRow represents the structure for each row of coin information.
type CoinInfoRow struct {
	Name         string          `json:"name"`
	Coin         string          `json:"coin"`
	RemainAmount string          `json:"remainAmount"`
	Chains       []CoinChainInfo `json:"chains"`
}

// CoinChainInfo represents the structure for each chain's information for a coin.
type CoinChainInfo struct {
	Chain                 string `json:"chain"`
	ChainType             string `json:"chainType"`
	Confirmation          string `json:"confirmation"`
	WithdrawFee           string `json:"withdrawFee"`
	DepositMin            string `json:"depositMin"`
	WithdrawMin           string `json:"withdrawMin"`
	MinAccuracy           string `json:"minAccuracy"`
	ChainDeposit          string `json:"chainDeposit"`
	ChainWithdraw         string `json:"chainWithdraw"`
	WithdrawPercentageFee string `json:"withdrawPercentageFee"`
}

type KLineResponseRaw struct {
	List     [][]string `json:"list"`
	Symbol   string     `json:"symbol"`
}

type MarketlineCandle struct {
	StartTime  string `json:"startTime"`
	OpenPrice  string `json:"openPrice"`
	HighPrice  string `json:"highPrice"`
}