package main

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func compressUUID(input string) ([]byte, error) {
	uuid, err := uuid.Parse(input)
	if err != nil {
		return nil, err
	}
	return uuid[:], nil
}

func uint64ToBytes(input uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, input)

	return buf
}

func main() {
	timeNow := time.Now()

	tranactionType := 1                                            // 1 byte
	transactionMethod := 1                                         // 1 byte
	orderType := 1                                                 // 1 byte
	ticker := "XYZQ"                                               // 4 bytes
	quantity := 4294967295                                         // ?? bytes
	price := 770                                                   // ?? bytes
	var orderDate uint64 = uint64(timeNow.AddDate(1, 0, 0).Unix()) // ?? bytes
	var goodUntil uint64 = 18446744073709551615                    // ?? bytes
	traderId := uuid.NewString()                                   // 36 bytes
	clientOrderId := uuid.NewString()                              // 36 bytes

	tickerByte := []byte(ticker)                        // 4 bytes
	quantityByte := make([]byte, 4)                     // 4 bytes
	priceByte := make([]byte, 4)                        // 4 bytes
	orderDateByte := uint64ToBytes(orderDate)           // 8 bytes
	goodUntilByte := uint64ToBytes(goodUntil)           // 8 bytes
	traderIdByte, _ := compressUUID(traderId)           // 16 bytes
	clientOrderIdByte, _ := compressUUID(clientOrderId) // 16 bytes

	binary.BigEndian.PutUint32(quantityByte, uint32(quantity))
	binary.BigEndian.PutUint32(priceByte, uint32(price))

	orderPayload := []byte{}

	orderPayload = append(orderPayload, byte(tranactionType))
	orderPayload = append(orderPayload, byte(transactionMethod))
	orderPayload = append(orderPayload, byte(orderType))
	orderPayload = append(orderPayload, tickerByte...)
	orderPayload = append(orderPayload, quantityByte...)
	orderPayload = append(orderPayload, priceByte...)
	orderPayload = append(orderPayload, orderDateByte...)
	orderPayload = append(orderPayload, goodUntilByte...)
	orderPayload = append(orderPayload, traderIdByte...)
	orderPayload = append(orderPayload, clientOrderIdByte...)
	orderPayloadUrlEncoded := base64.URLEncoding.EncodeToString(orderPayload)

	req, err := http.NewRequest("GET", "http://localhost:8080/"+orderPayloadUrlEncoded, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "market-broker")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	if res.StatusCode > 399 {
		var errorResponse map[string]string

		jsonErr := json.Unmarshal(body, &errorResponse)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		errorMsg := errorResponse["error"]
		if errorMsg != "" {
			log.Println(errorMsg)
		}
	}

	defer res.Body.Close()

	log.Printf("Response: %s", res.Status)

}
