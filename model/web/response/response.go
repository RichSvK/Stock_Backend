package response

import (
	"time"
)

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

type BrokerResponse struct {
	Code string `json:"broker_code"`
	Name string `json:"name"`
}
