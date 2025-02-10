# Задание
На любой криптовалютной бирже на выбор с помощью подручных инструментов определить, из каких данных формируется:
- Order book
- Лента сделок
- График цены (snapshot + обновления)

Пояснение:
- Нужно определить, что это за данные, показать их, откуда они получаются и каким образом (с помощью каких запросов, через какие протоколы + описать общую механику)

- Нужно предоставить запросы, которые делает терминал, чтобы получить снэпшоты ордербука, сделок и графика цены

- Нужно предоставить запросы, которые шлёт терминал, чтобы получать обновления ордербука, сделок и графика цены

Нужно открыть терминал любой криптовалютной биржи, нажать F12 и с помощью вкладки Network описать какие данные приходят от сервера биржи по API и во что эти данные превращаются на экране + можно вооружиться документацией API биржи.


# Ответ (Binance.com)
Все три запрошенных типа данных получаются по одному принципу:

- Получаем снимок текущего состояния (snapshot) с помощью GET запроса на API endpoint
- Подписываемся на combined stream (endpoint wss://stream.binance.com:9443/stream)
- Фронтэнд (JavaScript) обновляет графики/стакан/информацию о сделках в реальном времени, используя данные, полученные из стрима
    
Единый subscribe message для websocket стрима (далее буду отдельно их приводить, для ясности)
```
{"method":"SUBSCRIBE","params":["!miniTicker@arr@3000ms","btcusdt@aggTrade","btcusdt@depth","btcusdt@kline_1d"],"id":1}
```

## Trades

### Snapshot request:
GET https://www.binance.com/api/v1/aggTrades?limit=80&symbol=BTCUSDT


### Updates:
Websocket connection to
    wss://stream.binance.com/stream
```
{"method":"SUBSCRIBE","params":["btcusdt@aggTrade"],"id":1}
```

### Описание из документации:
```
{
"e": "aggTrade",    // Event type
"E": 1672515782136, // Event time
"s": "BNBBTC",      // Symbol
"a": 12345,         // Aggregate trade ID
"p": "0.001",       // Price
"q": "100",         // Quantity
"f": 100,           // First trade ID
"l": 105,           // Last trade ID
"T": 1672515782136, // Trade time
"m": true,          // Is the buyer the market maker?
"M": true           // Ignore
}
```

### Пример:
```
{"stream":"btcusdt@aggTrade","data":{"e":"aggTrade","E":1739202462290,"s":"BTCUSDT","a":3427780726,"p":"97020.30000000","q":"0.00018000","f":4543224650,"l":4543224652,"T":1739202462289,"m":false,"M":true}}
```

## Стакан / Order book

### Snapshot request:
GET https://www.binance.com/api/v3/depth?symbol=BTCUSDT&limit=1000


### Updates:
Websocket connection to
    wss://stream.binance.com/stream

### Subscribe message:
```
{"method":"SUBSCRIBE","params":["btcusdt@depth"],"id":1}
```

### Описание из документации:
```
{
"e": "depthUpdate", // Event type
"E": 1672515782136, // Event time
"s": "BNBBTC",      // Symbol
"U": 157,           // First update ID in event
"u": 160,           // Final update ID in event
"b": [              // Bids to be updated
[
"0.0024",       // Price level to be updated
"10"            // Quantity
]
],
"a": [              // Asks to be updated
[
"0.0026",       // Price level to be updated
"100"           // Quantity
]
]
}
```

### Пример:
(Обрезал, так как очень много сообщений)
```
{"stream":"btcusdt@depth","data":{"e":"depthUpdate","E":1739202515014,"s":"BTCUSDT","U":60859031517,"u":60859034785,"b":[["97053.55000000","4.06247000"],["97053.29000000","0.00030000"],["97053.28000000","0.24154000"],["97053.27000000","0.56691000"],["97052.71000000","0.00037000"],["97052.48000000","0.03051000"],["97052.30000000","0.00024000"],["97052.00000000","0.03576000"],["97051.76000000","0.00012000"],["97051.50000000","0.04138000"],…..]}
```

## График

### Snapshot request:
GET https://www.binance.com/api/v3/uiKlines?limit=1000&symbol=BTCUSDT&interval=1d

### Updates:
Websocket connection to
    wss://stream.binance.com/stream


### subscribe message:
```
{"method":"SUBSCRIBE","params":["btcusdt@kline_1d"],"id":1}
```

### Описание из документации:
```
{
"e": "kline",         // Event type
"E": 1672515782136,   // Event time
"s": "BNBBTC",        // Symbol
"k": {
"t": 1672515780000, // Kline start time
"T": 1672515839999, // Kline close time
"s": "BNBBTC",      // Symbol
"i": "1m",          // Interval
"f": 100,           // First trade ID
"L": 200,           // Last trade ID
"o": "0.0010",      // Open price
"c": "0.0020",      // Close price
"h": "0.0025",      // High price
"l": "0.0015",      // Low price
"v": "1000",        // Base asset volume
"n": 100,           // Number of trades
"x": false,         // Is this kline closed?
"q": "1.0000",      // Quote asset volume
"V": "500",         // Taker buy base asset volume
"Q": "0.500",       // Taker buy quote asset volume
"B": "123456"       // Ignore
}
}
```

### Пример:

```
[ [ 1652832000000, "30444.93000000", "30709.99000000", "28654.47000000", "28715.32000000", "59749.15799000", 1652918399999, "1762843836.12693780", 1379212, "29501.76769000", "870623227.20705700", "0" ], ]
```

