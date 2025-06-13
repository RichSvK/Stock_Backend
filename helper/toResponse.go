package helper

import (
	"backend/model/entity"
	"backend/model/web/output"
	"backend/model/web/response"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func ToWebResponse(status int, message string, data interface{}) interface{} {
	return output.WebResponse{
		Code:    status,
		Message: message,
		Data:    data,
	}
}

func ToFailedResponse(status int, message string) interface{} {
	return output.FailedResponse{
		Code:    status,
		Message: message,
	}
}

func ToBalanceResponse(stock entity.Stock) response.BalanceResponse {
	return response.BalanceResponse{
		Date:         stock.Date,
		ListedShares: stock.ListedShares,
		LocalIS:      stock.LocalIS,
		LocalCP:      stock.LocalCP,
		LocalPF:      stock.LocalPF,
		LocalIB:      stock.LocalIB,
		LocalID:      stock.LocalID,
		LocalMF:      stock.LocalMF,
		LocalSC:      stock.LocalSC,
		LocalFD:      stock.LocalFD,
		LocalOT:      stock.LocalOT,
		ForeignIS:    stock.ForeignIS,
		ForeignCP:    stock.ForeignCP,
		ForeignPF:    stock.ForeignPF,
		ForeignIB:    stock.ForeignIB,
		ForeignID:    stock.ForeignID,
		ForeignMF:    stock.ForeignMF,
		ForeignSC:    stock.ForeignSC,
		ForeignFD:    stock.ForeignFD,
		ForeignOT:    stock.ForeignOT,
	}
}

func ToBalanceResponses(listBalance []entity.Stock) []response.BalanceResponse {
	var balanceResponses []response.BalanceResponse = nil
	for _, balance := range listBalance {
		balanceResponses = append(balanceResponses, ToBalanceResponse(balance))
	}
	return balanceResponses
}

func ToIpoResponse(ipo entity.Ipo) response.IpoResponse {
	responseIpo := response.IpoResponse{
		StockCode:       ipo.StockCode,
		Price:           ipo.Price,
		IPO_Shares:      ipo.IPO_Shares,
		ListedShares:    ipo.ListedShares,
		Equity:          ipo.Equity,
		Nominal:         ipo.Nominal,
		MCB:             ipo.MCB,
		IsAffiliated:    ipo.IsAffiliated,
		IsAcceleration:  ipo.IsAcceleration,
		IsNew:           ipo.IsNew,
		LockUp:          ipo.LockUp,
		SubscribedStock: ipo.SubscribedStock,
		AllUnderwriter:  "{",
		Amount:          ipo.Amount,
	}

	if ipo.Warrant == 0 {
		responseIpo.Warrant = "0"
	} else {
		warrant := big.NewRat(int64(ipo.IPO_Shares), int64(ipo.Warrant))
		responseIpo.Warrant = fmt.Sprintf("%s : %s", warrant.Num().String(), warrant.Denom().String())
	}

	underwriter := strings.Split(ipo.AllUnderwriter, ",")
	uwShares := strings.Split(ipo.AllShares, ",")
	size := len(underwriter) - 1
	var percentage float64 = 0

	for i := 0; i < size; i++ {
		share, _ := strconv.ParseFloat(uwShares[i], 64)
		percentage = share / float64(ipo.IPO_Shares) * 100
		responseIpo.AllUnderwriter += fmt.Sprintf("%s : %.2f%%, ", underwriter[i], percentage)
	}
	share, _ := strconv.ParseFloat(uwShares[size], 64)
	percentage = share / float64(ipo.IPO_Shares) * 100
	responseIpo.AllUnderwriter += fmt.Sprintf("%s : %.2f%%}", underwriter[size], percentage)

	return responseIpo
}

func ToIpoResponses(listIpo []entity.Ipo) []response.IpoResponse {
	var IpoResponses []response.IpoResponse = nil
	for _, ipo := range listIpo {
		IpoResponses = append(IpoResponses, ToIpoResponse(ipo))
	}
	return IpoResponses
}

func ToBrokerResponse(broker entity.Broker) response.BrokerResponse {
	return response.BrokerResponse{
		Code: broker.Broker_Code,
		Name: broker.Broker_Code + " - " + broker.Broker_Name,
	}
}

func ToBrokerResponses(listBroker []entity.Broker) []response.BrokerResponse {
	var brokerResponses []response.BrokerResponse = nil
	for _, broker := range listBroker {
		brokerResponses = append(brokerResponses, ToBrokerResponse(broker))
	}
	return brokerResponses
}
