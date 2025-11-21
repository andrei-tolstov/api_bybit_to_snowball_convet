# api_bybit_to_snowball_convet

# https://bybit-exchange.github.io/docs/v5/intro

for start and end time
https://currentmillis.com/

go mod init github.com/andrei-tolstov/api_bybit_to_snowball_convet
go mod tidy

// go mod download

export BYBIT_API_KEY=
export BYBIT_API_SECRET=


логика работы
deposit создает сделку buy "валюта"-usd
withdraw создает сделку sell "валюта"-usd

начисление earn создает прочие начисления в паре "валюта"-usd

покупка(обмен) валют создает две сделки
sell "валюта_за_которую_купили"-usd(расчитываем_курс на момент сделки)
buy "валюта_покупки"-usd(берем из сделки sell)

пара USDT-USD - рассчитываю 1:1