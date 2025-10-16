package utils

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Order struct {
	TranactionType    uint8
	TransactionMethod uint8
	OrderType         uint8
	Ticker            string
	Quantity          uint32
	Price             uint32
	OrderDate         uint64
	GoodUntil         uint64
	TraderId          string
	ClientOrderId     string
}

var URI_LEN = 84
var ORDER_BYTE_LEN = 63

func OrderUriParser(uri string) (*Order, error) {
	var cleanUri string

	if uri[0] == '/' {
		cleanUri = strings.Replace(uri, "/", "", 1)
	}

	if len(cleanUri) != URI_LEN {
		return nil, fmt.Errorf("ParseOrderUri: Wrong URI length recieved - expected: %d, recieved: %d", URI_LEN, len(cleanUri))
	}

	orderBytes, orderBytesErr := base64.URLEncoding.DecodeString(cleanUri)
	if orderBytesErr != nil {
		return nil, orderBytesErr
	}

	if len(orderBytes) != ORDER_BYTE_LEN {
		return nil, fmt.Errorf("ParseOrderUri: Wrong byte length recieved - expected: %d, recieved: %d", ORDER_BYTE_LEN, len(orderBytes))
	}

	tranactionTypeByte := orderBytes[0]
	transactionMethodByte := orderBytes[1]
	orderTypeByte := orderBytes[2]
	tickerBytes := orderBytes[3:7]
	quantityBytes := orderBytes[7:11]
	priceBytes := orderBytes[11:15]
	orderDateBytes := orderBytes[15:23]
	goodUntilBytes := orderBytes[23:31]
	traderIdBytes := orderBytes[31:47]
	clientOrderIdBytes := orderBytes[47:63]

	order := Order{}

	if tranactionTypeByte < 1 || tranactionTypeByte > 2 {
		return nil, fmt.Errorf("ParseOrderUri: Wrong transaction type recieved: %v", tranactionTypeByte)
	} else {
		order.TranactionType = uint8(tranactionTypeByte)
	}

	if transactionMethodByte < 1 || transactionMethodByte > 3 {
		return nil, fmt.Errorf("ParseOrderUri: Wrong transaction method recieved: %v", transactionMethodByte)
	} else {
		order.TransactionMethod = uint8(transactionMethodByte)
	}

	if orderTypeByte < 1 || orderTypeByte > 2 {
		return nil, fmt.Errorf("ParseOrderUri: Wrong order type recieved: %v", orderTypeByte)
	} else {
		order.OrderType = uint8(orderTypeByte)
	}

	for _, byte := range tickerBytes {
		if byte > 0 && (byte < 65 || byte > 90) {
			return nil, fmt.Errorf("ParseOrderUri: Wrong ticker recieved: %s", string(tickerBytes))
		}
	}

	traderUuid, traderUuidErr := uuid.FromBytes(traderIdBytes)
	if traderUuidErr != nil {
		return nil, fmt.Errorf("ParseOrderUri: %w", traderUuidErr)
	}

	orderUuid, orderUuidErr := uuid.FromBytes(clientOrderIdBytes)
	if orderUuidErr != nil {
		return nil, fmt.Errorf("ParseOrderUri: %w", orderUuidErr)
	}

	order.Ticker = string(tickerBytes)
	order.Quantity = uint32(binary.BigEndian.Uint32(quantityBytes))
	order.Price = uint32(binary.BigEndian.Uint32(priceBytes))
	order.OrderDate = uint64(binary.BigEndian.Uint64(orderDateBytes))
	order.GoodUntil = uint64(binary.BigEndian.Uint64(goodUntilBytes))
	order.TraderId = traderUuid.String()
	order.ClientOrderId = orderUuid.String()

	return &order, nil
}
