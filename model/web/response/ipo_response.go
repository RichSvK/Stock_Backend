package response

import (
	"backend/model/entity"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type GetIpoResponse struct {
	Message string        `json:"message"`
	Data    []IpoResponse `json:"data"`
}

type IpoResponse struct {
	StockCode       string `json:"stock_code"`
	Price           uint32 `json:"price"`
	IPO_Shares      uint64 `json:"ipo_shares"`
	ListedShares    uint64 `json:"listed_shares"`
	Equity          int64  `json:"equity"`
	Warrant         string `json:"warrant"`
	Nominal         uint32 `json:"nominal"`
	MCB             uint64 `json:"mcb"`
	IsAffiliated    bool   `json:"is_affiliated"`
	IsAcceleration  bool   `json:"is_acceleration"`
	IsNew           bool   `json:"is_new"`
	LockUp          int8   `json:"lock_up"`
	SubscribedStock uint64 `json:"subscribed_stock"`
	AllUnderwriter  string `json:"all_underwriter"`
	Amount          uint64 `json:"amount"`
}

func ToIpoResponse(ipo *entity.Ipo) IpoResponse {
	responseIpo := IpoResponse{
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

	for i := 0; i < size; i++ {
		share, _ := strconv.ParseFloat(uwShares[i], 64)
		percentage := share / float64(ipo.IPO_Shares) * 100
		responseIpo.AllUnderwriter += fmt.Sprintf("%s : %.2f%%, ", underwriter[i], percentage)
	}
	share, _ := strconv.ParseFloat(uwShares[size], 64)
	percentage := share / float64(ipo.IPO_Shares) * 100
	responseIpo.AllUnderwriter += fmt.Sprintf("%s : %.2f%%}", underwriter[size], percentage)

	return responseIpo
}

func ToIpoResponses(listIpo []entity.Ipo) []IpoResponse {
	ipoResponses := make([]IpoResponse, 0, len(listIpo))

	for _, ipo := range listIpo {
		ipoResponses = append(ipoResponses, ToIpoResponse(&ipo))
	}
	return ipoResponses
}
