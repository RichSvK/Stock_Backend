package service

import (
	"backend/internal/entity"
	"backend/internal/helper"
	"backend/internal/model/domainerr"
	"backend/internal/model/query_filter"
	"backend/internal/model/response"
	"backend/internal/repository"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BalanceService interface {
	Create(ctx context.Context, fileName string) (*response.UploadBalanceResponse, error)
	ExportCode(ctx context.Context, code string) error
	GetBalanceData(ctx context.Context, code string) (*response.GetBalanceResponse, error)
	GetScriptlessChange(ctx context.Context, startTime time.Time, endTime time.Time) (*response.GetScriptlessChangeResponse, error)
	GetBalanceChangeData(ctx context.Context, shareholderType string, change string, page int) (*response.BalanceChangeResponse, error)
}

type BalanceServiceImpl struct {
	BalanceRepository repository.BalanceRepository
}

func NewBalanceService(repositoryBalance repository.BalanceRepository) BalanceService {
	return &BalanceServiceImpl{
		BalanceRepository: repositoryBalance,
	}
}

func (service *BalanceServiceImpl) Create(ctx context.Context, fileName string) (*response.UploadBalanceResponse, error) {
	var path = filepath.Join("resource", fileName)
	file, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		log.Printf("Error opening file: %v\n", err)
		return nil, domainerr.ErrInternalServer
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println("Error closing file:", err)
		}
	}()

	reader := bufio.NewReader(file)

	headerLine, _, err := reader.ReadLine()
	if err != nil {
		return nil, domainerr.ErrInternalServer
	}

	headers := strings.Split(string(headerLine), "|")
	headerMap := make(map[string]int)
	for i, value := range headers {
		headerMap[strings.TrimSpace(value)] = i
	}

	get := func(row []string, key string) string {
		if idx, ok := headerMap[key]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	var rowsData []byte
	var listStock []entity.Stock = nil
	var stock = entity.Stock{}
	dateFormatter := "02-Jan-2006"

	for {
		rowsData, _, err = reader.ReadLine()
		if err == io.EOF {
			break
		}

		stockData := strings.Split(string(rowsData), "|")
		if get(stockData, "Type") != "EQUITY" {
			break
		}

		if len(get(stockData, "Code")) != 4 {
			continue
		}

		stock.Date, err = time.Parse(dateFormatter, get(stockData, "Date"))
		if err != nil {
			return nil, domainerr.ErrInternalServer
		}

		stock.Date, err = time.Parse("02-01-2006", stock.Date.Format("02-01-2006"))
		if err != nil {
			return nil, domainerr.ErrInternalServer
		}

		stock.Code = get(stockData, "Code")
		stock.ListedShares, _ = strconv.ParseUint(string(stockData[3]), 10, 64)
		stock.Price, _ = strconv.ParseUint(get(stockData, "Price"), 10, 64)

		stock.LocalIS, _ = strconv.ParseUint(get(stockData, "Local IS"), 10, 64)
		stock.LocalCP, _ = strconv.ParseUint(get(stockData, "Local CP"), 10, 64)
		stock.LocalPF, _ = strconv.ParseUint(get(stockData, "Local PF"), 10, 64)
		stock.LocalIB, _ = strconv.ParseUint(get(stockData, "Local IB"), 10, 64)
		stock.LocalID, _ = strconv.ParseUint(get(stockData, "Local ID"), 10, 64)
		stock.LocalMF, _ = strconv.ParseUint(get(stockData, "Local MF"), 10, 64)
		stock.LocalSC, _ = strconv.ParseUint(get(stockData, "Local SC"), 10, 64)
		stock.LocalFD, _ = strconv.ParseUint(get(stockData, "Local FD"), 10, 64)
		stock.LocalOT, _ = strconv.ParseUint(get(stockData, "Local OT"), 10, 64)

		stock.ForeignIS, _ = strconv.ParseUint(get(stockData, "Foreign IS"), 10, 64)
		stock.ForeignCP, _ = strconv.ParseUint(get(stockData, "Foreign CP"), 10, 64)
		stock.ForeignPF, _ = strconv.ParseUint(get(stockData, "Foreign PF"), 10, 64)
		stock.ForeignIB, _ = strconv.ParseUint(get(stockData, "Foreign IB"), 10, 64)
		stock.ForeignID, _ = strconv.ParseUint(get(stockData, "Foreign ID"), 10, 64)
		stock.ForeignMF, _ = strconv.ParseUint(get(stockData, "Foreign MF"), 10, 64)
		stock.ForeignSC, _ = strconv.ParseUint(get(stockData, "Foreign SC"), 10, 64)
		stock.ForeignFD, _ = strconv.ParseUint(get(stockData, "Foreign FD"), 10, 64)
		stock.ForeignOT, _ = strconv.ParseUint(get(stockData, "Foreign OT"), 10, 64)

		listStock = append(listStock, stock)
	}

	err = service.BalanceRepository.Create(ctx, listStock)
	if err != nil {
		return nil, err
	}

	response := &response.UploadBalanceResponse{
		Message: "Data uploaded and inserted successfully",
	}

	return response, err
}

func (service *BalanceServiceImpl) ExportCode(ctx context.Context, code string) error {
	listStock, err := service.BalanceRepository.GetBalanceStock(ctx, code)

	if err != nil {
		return err
	}

	if len(listStock) == 0 {
		return domainerr.ErrBalanceNotFound
	}

	err = helper.MakeCSV(code, listStock)
	if err != nil {
		return domainerr.ErrFailedWriteCSV
	}

	return err
}

func (service *BalanceServiceImpl) GetBalanceData(ctx context.Context, code string) (*response.GetBalanceResponse, error) {
	listBalance, err := service.BalanceRepository.GetBalanceStock(ctx, code)
	if err != nil {
		return nil, err
	}

	if len(listBalance) == 0 {
		return nil, domainerr.ErrBalanceNotFound
	}

	response := &response.GetBalanceResponse{
		Message: fmt.Sprintf("Balance data for stock code %s retrieved successfully", code),
		Data:    response.ToBalanceResponses(listBalance),
	}

	return response, err
}

func (service *BalanceServiceImpl) GetScriptlessChange(ctx context.Context, startTime time.Time, endTime time.Time) (*response.GetScriptlessChangeResponse, error) {
	// Check if dates are invalid
	now := time.Now()
	if endTime.After(now) || startTime.After(endTime) || endTime.Before(startTime.AddDate(0, 1, 0)){
		return nil, domainerr.ErrInvalidDateRange
	}

	startTime = time.Date(startTime.Year(), startTime.Month(), 1, 0, 0, 0, 0, startTime.Location())
	endTime = time.Date(endTime.Year(), endTime.Month(), 1, 0, 0, 0, 0, endTime.Location())

	dateRange := query_filter.DateRangeQuery{
		StartTime:     startTime.Format("2006-01-02"),
		StartTimeLast: startTime.AddDate(0, 1, 0).Format("2006-01-02"),
		EndTime:       endTime.Format("2006-01-02"),
		EndTimeLast:   endTime.AddDate(0, 1, 0).Format("2006-01-02"),
	}
	listStock, err := service.BalanceRepository.GetScriptlessChange(ctx, dateRange)
	
	if err != nil {
		return nil, err
	}

	count := len(listStock)
	if count == 0 {
		return nil, domainerr.ErrBalanceNotFound
	}

	var listResponseChange []response.ScriptlessResponse
	stock := response.ScriptlessResponse{}
	for i := 0; i < count; i++ {
		stock = response.ScriptlessResponse{}

		if i < count-1 && listStock[i].Code == listStock[i+1].Code {
			stock.Code = listStock[i].Code
			stock.FirstShare = TotalShares(&listStock[i])
			stock.SecondShare = TotalShares(&listStock[i+1])
			stock.FirstListedShares = listStock[i].ListedShares
			stock.SecondListedShares = listStock[i+1].ListedShares
			stock.Change = int64(stock.SecondShare) - int64(stock.FirstShare)
			stock.ChangePercentage = float64(stock.Change) / float64(stock.FirstShare) * 100

			if stock.Change != 0 {
				listResponseChange = append(listResponseChange, stock)
			}
			i++
		} else {
			stock.Code = listStock[i].Code

			if listStock[i].Date.Month() != startTime.Month() {
				// IPO Stock
				stock.FirstShare = 0
				stock.FirstListedShares = 0
				stock.SecondShare = TotalShares(&listStock[i])
				stock.SecondListedShares = listStock[i].ListedShares
				stock.Change = int64(stock.SecondShare)
				stock.ChangePercentage = 100
			} else {
				// Delisted Stock
				stock.FirstShare = TotalShares(&listStock[i])
				stock.FirstListedShares = listStock[i].ListedShares
				stock.SecondShare = 0
				stock.SecondListedShares = 0
				stock.Change = -int64(stock.FirstShare)
				stock.ChangePercentage = -100
			}
			listResponseChange = append(listResponseChange, stock)
		}
	}

	sort.Slice(listResponseChange, func(i, j int) bool {
		return listResponseChange[i].ChangePercentage > listResponseChange[j].ChangePercentage
	})

	response := &response.GetScriptlessChangeResponse{
		Message: "Scriptless change data retrieved successfully",
		Data:    listResponseChange,
	}

	return response, err
}

func TotalShares(s *entity.Stock) uint64 {
	return s.LocalIS + s.LocalCP + s.LocalPF + s.LocalIB + s.LocalID + s.LocalMF +
		s.LocalSC + s.LocalFD + s.LocalOT + s.ForeignIS + s.ForeignCP + s.ForeignPF +
		s.ForeignIB + s.ForeignID + s.ForeignMF + s.ForeignSC + s.ForeignFD + s.ForeignOT
}

func (service *BalanceServiceImpl) GetBalanceChangeData(ctx context.Context, shareholderType string, change string, page int) (*response.BalanceChangeResponse, error) {
	var AllowedColumns = map[string]bool{
		"local_is": true, "local_cp": true, "local_pf": true,
		"local_ib": true, "local_id": true, "local_mf": true,
		"local_sc": true, "local_fd": true, "local_ot": true,
		"foreign_is": true, "foreign_cp": true, "foreign_pf": true,
		"foreign_ib": true, "foreign_id": true, "foreign_mf": true,
		"foreign_sc": true, "foreign_fd": true, "foreign_ot": true,
	}

	if !AllowedColumns[shareholderType] {
		return nil, domainerr.ErrInvalidShareholderType
	}

	now := time.Now()

	base := time.Date(
		now.Year(),
		now.Month(),
		1,
		0, 0, 0, 0,
		now.Location(),
	)

	prevYM := base.AddDate(0, -1, 0)
	prev2YM := base.AddDate(0, -2, 0)
	dateRangeQuery := query_filter.DateRangeQuery{
		StartTime:     prev2YM.Format("2006-01-02"),
		StartTimeLast: prev2YM.AddDate(0, 1, 0).Format("2006-01-02"),
		EndTime:       prevYM.Format("2006-01-02"),
		EndTimeLast:   prevYM.AddDate(0, 1, 0).Format("2006-01-02"),
	}

	listBalanceChange, err := service.BalanceRepository.GetBalanceChangeData(ctx, shareholderType, change, page, dateRangeQuery)

	if err != nil {
		return nil, err
	}

	if len(listBalanceChange) == 0 {
		return nil, domainerr.ErrBalanceNotFound
	}

	listBalanceChangeLen := len(listBalanceChange)

	if listBalanceChangeLen == 11 {
		listBalanceChange = listBalanceChange[:10]
	}

	response := response.BalanceChangeResponse{
		HaveNext: listBalanceChangeLen == 11,
		Data:     listBalanceChange,
	}

	return &response, nil
}
