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


# Пояснение к репозиторию

В файле "answer.md" находится текстовый ответ на вопрос.

В "/app/main.go" - код приложения, которое строит локальный order book. Да, мне стыдно за него - я написал его просто чтобы разобраться, как работает API Binance, и проверить своё понимание. Приложение выводит каждые 5 секунд топ 5 бидов и топ 5 асков ордер бука для "BTCUSDT", в не особо читаемом виде.

"screenshot_1.png" показывает моё умение нажимать F12 в Firefox

"screenshot_2_терминал.png" иллюстрирует три зоны в терминале, про которые спрашивается в задании. Красным обозначен стакан/order book, зеленым - график, синим - сделки/trades.