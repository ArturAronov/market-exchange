// go test ./test -fuzz=FuzzOrderUriParser -fuzztime=120s

package test

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"market-exchange/utils"

	"github.com/google/uuid"
)

func compressUUID(input string) ([]byte, error) {
	if input == "" {
		return nil, fmt.Errorf("empty UUID")
	}

	if len(input) != 36 {
		return nil, fmt.Errorf("invalid UUID length")
	}

	validUUIDErr := uuid.Validate(input)
	if validUUIDErr != nil {
		return nil, validUUIDErr
	}

	parsedUUID, parsedUUIDErr := uuid.Parse(input)
	if parsedUUIDErr != nil {
		return nil, parsedUUIDErr
	}

	return parsedUUID[:], nil
}

func uint64ToBytes(input uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(input))

	return buf
}

func generateInput(input utils.Order) (string, error) {
	tickerByte := []byte(input.Ticker)
	quantityByte := make([]byte, 4)
	priceByte := make([]byte, 4)
	orderDateByte := uint64ToBytes(input.OrderDate)
	goodUntilByte := uint64ToBytes(input.GoodUntil)
	traderIdByte, traderIdByteErr := compressUUID(input.TraderId)
	clientOrderIdByte, clientOrderIdByteErr := compressUUID(input.ClientOrderId)

	if traderIdByteErr != nil {
		return "", traderIdByteErr
	}

	if clientOrderIdByteErr != nil {
		return "", clientOrderIdByteErr
	}

	binary.BigEndian.PutUint32(quantityByte, uint32(input.Quantity))
	binary.BigEndian.PutUint32(priceByte, uint32(input.Price))

	orderPayload := []byte{}

	orderPayload = append(orderPayload, byte(input.TranactionType))
	orderPayload = append(orderPayload, byte(input.TransactionMethod))
	orderPayload = append(orderPayload, byte(input.OrderType))
	orderPayload = append(orderPayload, tickerByte...)
	orderPayload = append(orderPayload, quantityByte...)
	orderPayload = append(orderPayload, priceByte...)
	orderPayload = append(orderPayload, orderDateByte...)
	orderPayload = append(orderPayload, goodUntilByte...)
	orderPayload = append(orderPayload, traderIdByte...)
	orderPayload = append(orderPayload, clientOrderIdByte...)

	output := base64.URLEncoding.EncodeToString(orderPayload)

	return output, nil
}

func FuzzOrderUriParser(f *testing.F) {
	testcases := []utils.Order{
		{
			TranactionType:    1,
			TransactionMethod: 1,
			OrderType:         1,
			Ticker:            "XYZQ",
			Quantity:          4294967295,
			Price:             770,
			OrderDate:         uint64(time.Now().AddDate(1, 0, 0).Unix()),
			GoodUntil:         18446744073709551615,
			TraderId:          uuid.NewString(),
			ClientOrderId:     uuid.NewString(),
		},
		{
			TranactionType:    8,
			TransactionMethod: 8,
			OrderType:         8,
			Ticker:            "#$%^",
			Quantity:          4294967295,
			Price:             770000000,
			OrderDate:         uint64(time.Now().AddDate(1, 0, 0).Unix()),
			GoodUntil:         1,
			TraderId:          "46090691-dc41-4280-960f-a",
			ClientOrderId:     "2527ee99-4f5c-4adf-8a10-4",
		},
		{
			TranactionType:    0,
			TransactionMethod: 0,
			OrderType:         0,
			Ticker:            "",
			Quantity:          0,
			Price:             0,
			OrderDate:         0,
			GoodUntil:         0,
			TraderId:          "",
			ClientOrderId:     "",
		},
	}

	for _, tc := range testcases {
		encodedOrder, encodedOrderErr := generateInput(tc)
		if encodedOrderErr != nil {
			continue
		}

		f.Add(encodedOrder)
	}

	f.Add("invalid base64!")
	f.Add("====")
	f.Add(string([]byte{0xFF, 0xFF}))

	f.Fuzz(func(t *testing.T, encodedOrder string) {
		parsedOrder, parsedOrderErr := utils.OrderUriParser(encodedOrder)
		if parsedOrderErr != nil {
			t.Skip()
		}

		reEncodedOrder, reEncodedOrderErr := generateInput(*parsedOrder)
		if reEncodedOrderErr != nil {
			t.Errorf("Failed to re-encode parsed order %v", reEncodedOrderErr)
		}

		if reEncodedOrder != encodedOrder {
			t.Errorf("Re-encoded order does not match encoded order. Got re-encoded: %v vs encoded: %v", reEncodedOrder, encodedOrder)
		}
	})
}
