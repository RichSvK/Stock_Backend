package response

import (
	"backend/model/entity"
	query_model "backend/model/query"
	"time"
)

type UploadBalanceResponse struct {
	Message string `json:"message"`
}

type GetBalanceResponse struct {
	Message string            `json:"message"`
	Data    []BalanceResponse `json:"data"`
}

type BalanceResponse struct {
	Date         time.Time `json:"date"`
	ListedShares uint64    `json:"listed_shares"`
	LocalIS      uint64    `json:"local_is"`
	LocalCP      uint64    `json:"local_cp"`
	LocalPF      uint64    `json:"local_pf"`
	LocalIB      uint64    `json:"local_ib"`
	LocalID      uint64    `json:"local_id"`
	LocalMF      uint64    `json:"local_mf"`
	LocalSC      uint64    `json:"local_sc"`
	LocalFD      uint64    `json:"local_fd"`
	LocalOT      uint64    `json:"local_ot"`
	ForeignIS    uint64    `json:"foreign_is"`
	ForeignCP    uint64    `json:"foreign_cp"`
	ForeignPF    uint64    `json:"foreign_pf"`
	ForeignIB    uint64    `json:"foreign_ib"`
	ForeignID    uint64    `json:"foreign_id"`
	ForeignMF    uint64    `json:"foreign_mf"`
	ForeignSC    uint64    `json:"foreign_sc"`
	ForeignFD    uint64    `json:"foreign_fd"`
	ForeignOT    uint64    `json:"foreign_ot"`
}

type BalanceChangeResponse struct {
	HaveNext bool                        `json:"have_next"`
	Data     []query_model.BalanceChange `json:"data"`
}

func ToBalanceResponse(stock *entity.Stock) BalanceResponse {
	return BalanceResponse{
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

func ToBalanceResponses(listBalance []entity.Stock) []BalanceResponse {
	balanceResponses := make([]BalanceResponse, 0, len(listBalance))
	for _, balance := range listBalance {
		balanceResponses = append(balanceResponses, ToBalanceResponse(&balance))
	}
	return balanceResponses
}

type GetScriptlessChangeResponse struct {
	Message string               `json:"message"`
	Data    []ScriptlessResponse `json:"data"`
}

type ScriptlessResponse struct {
	Code               string  `json:"code"`
	FirstShare         uint64  `json:"first_share"`
	SecondShare        uint64  `json:"second_share"`
	FirstListedShares  uint64  `json:"first_listed_share"`
	SecondListedShares uint64  `json:"second_listed_share"`
	Change             int64   `json:"change"`
	ChangePercentage   float64 `json:"change_percentage"`
}
