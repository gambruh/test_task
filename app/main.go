package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	symbol       = "BTCUSDT"
	snapshotURL  = "https://api.binance.com/api/v3/depth?symbol=" + symbol + "&limit=1000"
	websocketURL = "wss://stream.binance.com/stream"
)

type WebSocketResponse struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

type OrderBookSnapshot struct {
	LastUpdateID int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type DepthUpdate struct {
	E    string     `json:"e"` // Event type (e.g., "depthUpdate")
	Es   int        `json:"E"` // Event time in milliseconds
	S    string     `json:"s"` // Symbol (e.g., "BTCUSDT")
	U    int        `json:"U"` // First update ID
	U2   int        `json:"u"` // Last update ID
	Bids [][]string `json:"b"` // Bids (Price, Quantity)
	Asks [][]string `json:"a"` // Asks (Price, Quantity)
}

var (
	currentOrderBook OrderBookSnapshot
)

func fetchOrderBookSnapshot() OrderBookSnapshot {
	resp, err := http.Get(snapshotURL)
	if err != nil {
		log.Fatalf("Error fetching order book snapshot: %v", err)
	}
	defer resp.Body.Close()

	var snapshot OrderBookSnapshot
	if err := json.NewDecoder(resp.Body).Decode(&snapshot); err != nil {
		log.Fatalf("Error decoding snapshot: %v", err)
	}

	return snapshot
}

func applyDepthUpdate(update DepthUpdate, isFirst *bool, firstMsgU *int) {

	// initialize the order book if not done yet:
	if *isFirst {
		initializeOrderBook()
		*isFirst = false
	}

	if *firstMsgU > currentOrderBook.LastUpdateID {
		initializeOrderBook()
	}

	if update.U2 < currentOrderBook.LastUpdateID {
		log.Println("skipping the message, as u < lastupdateID")
		return
	}

	for _, bid := range update.Bids {
		updateOrder(&currentOrderBook.Bids, bid[0], bid[1], false)
	}
	for _, ask := range update.Asks {
		updateOrder(&currentOrderBook.Asks, ask[0], ask[1], true)
	}

	currentOrderBook.LastUpdateID = update.U2
}

func updateOrder(orders *[][]string, priceStr, sizeStr string, isAsk bool) {
	price, _ := strconv.ParseFloat(priceStr, 64)
	size, _ := strconv.ParseFloat(sizeStr, 64)

	for i, order := range *orders {
		// fmt.Println(order)
		existingPrice, _ := strconv.ParseFloat(order[0], 64)
		if existingPrice == price {
			if size == 0 {
				*orders = append((*orders)[:i], (*orders)[i+1:]...)
			} else {
				(*orders)[i][1] = sizeStr
			}
			return
		}
	}

	if size > 0 {
		*orders = append(*orders, []string{priceStr, sizeStr})
		sort.Slice(*orders, func(i, j int) bool {
			priceI, _ := strconv.ParseFloat((*orders)[i][0], 64)
			priceJ, _ := strconv.ParseFloat((*orders)[j][0], 64)
			if isAsk {
				return priceI < priceJ
			} else {
				return priceI > priceJ
			}
		})
	}

	// fmt.Println(orders)
}

func handleDepthUpdate(data json.RawMessage, isFirst *bool, firstMessageU *int) {
	var depthData DepthUpdate
	if err := json.Unmarshal(data, &depthData); err != nil {
		log.Printf("Error unmarshaling depth data: %v\n", err)
		return
	}
	if *isFirst {
		*firstMessageU = depthData.U
	}

	// updateBuffer = append(updateBuffer, depthData)

	applyDepthUpdate(depthData, isFirst, firstMessageU)
}

func initializeOrderBook() {

	for {
		log.Println("Fetching fresh order book snapshot...")
		snapshot := fetchOrderBookSnapshot()

		// if len(updateBuffer) > 0 && snapshot.LastUpdateID < updateBuffer[0].U {
		// 	continue
		// }

		// filteredBuffer := []DepthUpdate{}
		// for _, update := range updateBuffer {
		// 	if update.U2 > snapshot.LastUpdateID {
		// 		filteredBuffer = append(filteredBuffer, update)
		// 	}
		// }

		currentOrderBook = snapshot
		// updateBuffer = filteredBuffer
		log.Println("Order book initialized with snapshot")
		break
	}
}

func subscribeToOBStream(conn *websocket.Conn) {

	subscribeParam := strings.ToLower(symbol) + "@depth"
	// Build a subscription message
	subscribeMessage := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{
			subscribeParam,
		},
		"id": 1,
	}
	subscribeMessageJSON, err := json.Marshal(subscribeMessage)
	if err != nil {
		log.Fatal("Error marshaling subscribe message: ", err)
	}

	// Send the subscription message
	err = conn.WriteMessage(websocket.TextMessage, subscribeMessageJSON)
	if err != nil {
		log.Fatal("Error subscribing: ", err)
	}
}

func streamOrderBookUpdates() {

	conn, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		log.Fatalf("WebSocket connection error: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to Binance WebSocket")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			printOrderBook()
		}
	}()
	subscribeToOBStream(conn)
	initializeOrderBook()

	isFirst := true
	firstMessageU := 0
	for {

		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			return
		}

		// Parse the incoming message
		var wsResponse WebSocketResponse
		if err := json.Unmarshal(message, &wsResponse); err != nil {
			log.Println("Error unmarshaling WebSocket message:", err)
			continue
		}

		handleDepthUpdate(wsResponse.Data, &isFirst, &firstMessageU)

	}

}

func printOrderBook() {

	fmt.Print("\033[H\033[2J") // Clear terminal output
	fmt.Println("Updated Order Book:")
	fmt.Println("Bids:", currentOrderBook.Bids[:5])
	fmt.Println("Asks:", currentOrderBook.Asks[:5])
}

func main() {
	streamOrderBookUpdates()
}
