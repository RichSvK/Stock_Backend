package service

import (
	"backend/helper"
	"backend/model/entity"
	domain_error "backend/model/error"
	"backend/model/web/response"
	"backend/repository"
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
	var path = filepath.Join("Resource", fileName)
	file, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return nil, domain_error.ErrInternalServer
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println("Error closing file:", err)
		}
	}()

	reader := bufio.NewReader(file)

	_, _, err = reader.ReadLine()
	if err != nil {
		return nil, domain_error.ErrInternalServer
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
		if stockData[2] == "CORPORATE BOND" {
			break
		}

		if len(stockData[1]) != 4 {
			continue
		}

		stock.Date, err = time.Parse(dateFormatter, string(stockData[0]))
		if err != nil {
			return nil, domain_error.ErrInternalServer
		}

		stock.Date, err = time.Parse("02-01-2006", stock.Date.Format("02-01-2006"))
		if err != nil {
			return nil, domain_error.ErrInternalServer
		}

		stock.Code = string(stockData[1])
		stock.ListedShares, _ = strconv.ParseUint(string(stockData[3]), 10, 64)
		stock.Price, _ = strconv.ParseUint(string(stockData[4]), 10, 64)
		stock.LocalIS, _ = strconv.ParseUint(string(stockData[5]), 10, 64)
		stock.LocalCP, _ = strconv.ParseUint(string(stockData[6]), 10, 64)
		stock.LocalPF, _ = strconv.ParseUint(string(stockData[7]), 10, 64)
		stock.LocalIB, _ = strconv.ParseUint(string(stockData[8]), 10, 64)
		stock.LocalID, _ = strconv.ParseUint(string(stockData[9]), 10, 64)
		stock.LocalMF, _ = strconv.ParseUint(string(stockData[10]), 10, 64)
		stock.LocalSC, _ = strconv.ParseUint(string(stockData[11]), 10, 64)
		stock.LocalFD, _ = strconv.ParseUint(string(stockData[12]), 10, 64)
		stock.LocalOT, _ = strconv.ParseUint(string(stockData[13]), 10, 64)

		stock.ForeignIS, _ = strconv.ParseUint(string(stockData[15]), 10, 64)
		stock.ForeignCP, _ = strconv.ParseUint(string(stockData[16]), 10, 64)
		stock.ForeignPF, _ = strconv.ParseUint(string(stockData[17]), 10, 64)
		stock.ForeignIB, _ = strconv.ParseUint(string(stockData[18]), 10, 64)
		stock.ForeignID, _ = strconv.ParseUint(string(stockData[19]), 10, 64)
		stock.ForeignMF, _ = strconv.ParseUint(string(stockData[20]), 10, 64)
		stock.ForeignSC, _ = strconv.ParseUint(string(stockData[21]), 10, 64)
		stock.ForeignFD, _ = strconv.ParseUint(string(stockData[22]), 10, 64)
		stock.ForeignOT, _ = strconv.ParseUint(string(stockData[23]), 10, 64)

		listStock = append(listStock, stock)
	}

	err = service.BalanceRepository.Create(ctx, listStock)
	if err != nil {
		return nil, domain_error.ErrInternalServer
	}

	response := &response.UploadBalanceResponse{
		Message: "Data uploaded and data inserted successfully",
	}

	return response, err
}

func (service *BalanceServiceImpl) ExportCode(ctx context.Context, code string) error {
	listStock, err := service.BalanceRepository.GetBalanceStock(ctx, code)

	if err != nil {
		return err
	}

	if len(listStock) == 0 {
		return domain_error.ErrBalanceNotFound
	}

	err = helper.MakeCSV(code, listStock)
	if err != nil {
		return domain_error.ErrFailedWriteCSV
	}

	return err
}

func (service *BalanceServiceImpl) GetBalanceData(ctx context.Context, code string) (*response.GetBalanceResponse, error) {
	listBalance, err := service.BalanceRepository.GetBalanceStock(ctx, code)
	if err != nil {
		return nil, err
	}

	if len(listBalance) == 0 {
		return nil, domain_error.ErrBalanceNotFound
	}

	response := &response.GetBalanceResponse{
		Message: fmt.Sprintf("Balance data for stock code %s retrieved successfully", code),
		Data:    response.ToBalanceResponses(listBalance),
	}

	return response, err
}

func (service *BalanceServiceImpl) GetScriptlessChange(ctx context.Context, startTime time.Time, endTime time.Time) (*response.GetScriptlessChangeResponse, error) {
	if startTime.After(endTime) {
		return nil, domain_error.ErrInvalidDateRange
	}

	if endTime.Before(startTime.AddDate(0, 1, 0)) {
		return nil, domain_error.ErrInvalidDateRange
	}

	// Check if dates are not in the future
	now := time.Now()
	if endTime.After(now) || startTime.After(endTime) {
		return nil, domain_error.ErrInvalidDateRange
	}

	listStock, err := service.BalanceRepository.GetScriptlessChange(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}

	count := len(listStock)
	if count == 0 {
		return nil, domain_error.ErrBalanceNotFound
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
	if page < 1 {
		return nil, domain_error.ErrPaginationInvalid
	}

	var AllowedColumns = map[string]bool{
		"local_is": true, "local_cp": true, "local_pf": true,
		"local_ib": true, "local_id": true, "local_mf": true,
		"local_sc": true, "local_fd": true, "local_ot": true,
		"foreign_is": true, "foreign_cp": true, "foreign_pf": true,
		"foreign_ib": true, "foreign_id": true, "foreign_mf": true,
		"foreign_sc": true, "foreign_fd": true, "foreign_ot": true,
	}

	if !AllowedColumns[shareholderType] {
		return nil, domain_error.ErrInvalidShareholderType
	}

	if change != "Increase" && change != "Decrease" {
		return nil, domain_error.ErrInvalidChangeRequest
	}

	now := time.Now()

	base := time.Date(
		now.Year(),
		now.Month(),
		1,
		0, 0, 0, 0,
		now.Location(),
	)

	prevYM := base.AddDate(0, -1, 0).Format("2006-01")
	prev2YM := base.AddDate(0, -2, 0).Format("2006-01")

	listBalanceChange, err := service.BalanceRepository.GetBalanceChangeData(ctx, shareholderType, change, page, prev2YM, prevYM)

	if err != nil {
		return nil, err
	}

	if len(listBalanceChange) == 0 {
		return nil, domain_error.ErrBalanceNotFound
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
