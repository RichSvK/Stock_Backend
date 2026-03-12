package entity

type Ipo struct {
	StockCode       string `json:"stock_code"`
	Price           uint32 `json:"price"`
	IPO_Shares      uint64 `json:"ipo_shares"`
	ListedShares    uint64 `json:"listed_shares"`
	Equity          int64  `json:"equity"`
	Warrant         uint64 `json:"warrant"`
	Nominal         uint32 `json:"nominal"`
	MCB             uint64 `json:"mcb"`
	IsAffiliated    bool   `json:"is_affiliated"`
	IsAcceleration  bool   `json:"is_acceleration"`
	IsNew           bool   `json:"is_new"`
	LockUp          int8   `json:"lock_up"`
	SubscribedStock uint64 `json:"subscribed_stock"`
	AllUnderwriter  string `json:"all_underwriter"`
	AllShares       string `json:"all_shares"`
	Amount          uint64 `json:"amount"`
}
