package helper

import (
	"backend/model/entity"
	"log"
	"os"
	"strconv"
)

func RemoveFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Failed to delete %s\n", filePath)
		return err
	}
	log.Printf("File %s telah dihapus\n", filePath)
	return nil
}

func MakeCSV(code string, listStock []entity.Stock) error {
	file, err := os.OpenFile("./Resource/"+code+".csv", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("Date,Code,Local IS,Local CP,Local PF,Local IB,Local ID,Local MF,Local SC,Local FD,Local OT,Foreign IS,Foreign CP,Foreign PF,Foreign IB,Foreign ID,Foreign MF,Foreign SC,Foreign FD,Foreign OT\n")

	for _, stock := range listStock {
		WriteStockData(file, &stock)
	}
	return nil
}

func WriteStockData(file *os.File, stock *entity.Stock) {
	formattedDate := stock.Date.Format("02-01-2006")
	file.WriteString(formattedDate + ",")
	file.WriteString(stock.Code + ",")
	file.WriteString(strconv.FormatUint(stock.LocalIS, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.LocalCP, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.LocalPF, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.LocalIB, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.LocalID, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.LocalMF, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.LocalSC, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.LocalFD, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.LocalOT, 10) + ",")

	file.WriteString(strconv.FormatUint(stock.ForeignIS, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.ForeignCP, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.ForeignPF, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.ForeignIB, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.ForeignID, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.ForeignMF, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.ForeignSC, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.ForeignFD, 10) + ",")
	file.WriteString(strconv.FormatUint(stock.ForeignOT, 10) + "\n")
}
