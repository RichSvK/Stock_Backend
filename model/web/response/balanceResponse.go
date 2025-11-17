package response

import "time"

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

type BalanceChangesResponse struct {
	HaveNext bool                `json:"have_next"`
	Data     []BalanceChangeData `json:"data"`
}

type BalanceChangeData struct {
	StockCode         string  `json:"stock_code"`
	PreviousOwnership uint64  `json:"previous_ownership"`
	CurrentOwnership  uint64  `json:"current_ownership"`
	ChangePercentage  float64 `json:"change_percentage"`
}
