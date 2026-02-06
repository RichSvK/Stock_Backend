package entity

import "time"

type Stock struct {
	Date         time.Time `json:"date" gorm:"type:DATE;primaryKey"`
	Price        uint64    `json:"price" gorm:"type:BIGINT UNSIGNED"`
	Code         string    `json:"code" gorm:"type:CHAR(4);primaryKey"`
	ListedShares uint64    `json:"listed_shares" gorm:"type: BIGINT UNSIGNED"`
	LocalIS      uint64    `json:"local_is" gorm:"type:BIGINT UNSIGNED"`
	LocalCP      uint64    `json:"local_cp" gorm:"type:BIGINT UNSIGNED"`
	LocalPF      uint64    `json:"local_pf" gorm:"type:BIGINT UNSIGNED"`
	LocalIB      uint64    `json:"local_ib" gorm:"type:BIGINT UNSIGNED"`
	LocalID      uint64    `json:"local_id" gorm:"type:BIGINT UNSIGNED"`
	LocalMF      uint64    `json:"local_mf" gorm:"type:BIGINT UNSIGNED"`
	LocalSC      uint64    `json:"local_sc" gorm:"type:BIGINT UNSIGNED"`
	LocalFD      uint64    `json:"local_fd" gorm:"type:BIGINT UNSIGNED"`
	LocalOT      uint64    `json:"local_ot" gorm:"type:BIGINT UNSIGNED"`
	ForeignIS    uint64    `json:"foreign_is" gorm:"type:BIGINT UNSIGNED"`
	ForeignCP    uint64    `json:"foreign_cp" gorm:"type:BIGINT UNSIGNED"`
	ForeignPF    uint64    `json:"foreign_pf" gorm:"type:BIGINT UNSIGNED"`
	ForeignIB    uint64    `json:"foreign_ib" gorm:"type:BIGINT UNSIGNED"`
	ForeignID    uint64    `json:"foreign_id" gorm:"type:BIGINT UNSIGNED"`
	ForeignMF    uint64    `json:"foreign_mf" gorm:"type:BIGINT UNSIGNED"`
	ForeignSC    uint64    `json:"foreign_sc" gorm:"type:BIGINT UNSIGNED"`
	ForeignFD    uint64    `json:"foreign_fd" gorm:"type:BIGINT UNSIGNED"`
	ForeignOT    uint64    `json:"foreign_ot" gorm:"type:BIGINT UNSIGNED"`
}

type StockIPO struct {
	StockCode         string `gorm:"type:char(4);primaryKey"`
	Price             uint32 `gorm:"type:BIGINT UNSIGNED;not null"`
	IPO_Shares        uint64 `gorm:"type:BIGINT UNSIGNED;not null"`
	ListedShares      uint64 `gorm:"type:BIGINT UNSIGNED;not null"`
	Equity            int64  `gorm:"not null"`
	Warrant           uint64 `gorm:"type:BIGINT UNSIGNED;not null"`
	Nominal           uint32 `gorm:"type:INT UNSIGNED;not null"`
	MCB               uint64 `gorm:"type:BIGINT UNSIGNED;not null"`
	IsAffiliated      bool   `gorm:"not null"`
	IsAcceleration    bool   `gorm:"not null"`
	IsNew             bool   `gorm:"not null"`
	IsFullCommitment  bool   `gorm:"not null"`
	IsNotInvolvedCase bool   `gorm:"not null"`
	LockUp            int8   `gorm:"not null"`
	SubscribedStock   uint64 `gorm:"type:BIGINT UNSIGNED;not null"`

	// Relationship
	Detail []IPO_Detail `gorm:"foreignKey:stock_code;references:stock_code"`
}

type IPO_Detail struct {
	StockCode string `gorm:"type:char(4);primaryKey"`
	UW_Code   string `gorm:"type:char(2);primaryKey"`
	UwShares  uint64 `gorm:"type:BIGINT UNSIGNED"`
}

// Make table name from default "stocks" to "stock"
func (stock *Stock) TableName() string {
	return "stock"
}

func (stock *StockIPO) TableName() string {
	return "stock_ipo"
}

// Make table name from default "ipo_details" to "ipo_detail"
func (broker *IPO_Detail) TableName() string {
	return "ipo_detail"
}
