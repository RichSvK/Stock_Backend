package helper

import (
	"backend/model/entity"
	"encoding/csv"
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
	file, err := os.Create("./Resource/" + code + ".csv")
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println("Error closing file:", err)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{
		"Date", "Code", "Local IS", "Local CP", "Local PF", "Local IB", "Local ID", "Local MF",
		"Local SC", "Local FD", "Local OT", "Foreign IS", "Foreign CP", "Foreign PF",
		"Foreign IB", "Foreign ID", "Foreign MF", "Foreign SC", "Foreign FD", "Foreign OT",
	}); err != nil {
		return err
	}

	// Write each stock
	for _, stock := range listStock {
		if err := WriteStockData(writer, &stock); err != nil {
			return err
		}
	}

	return nil
}

func WriteStockData(w *csv.Writer, stock *entity.Stock) error {
	record := []string{
		stock.Date.Format("02-01-2006"),
		stock.Code,
		strconv.FormatUint(stock.LocalIS, 10),
		strconv.FormatUint(stock.LocalCP, 10),
		strconv.FormatUint(stock.LocalPF, 10),
		strconv.FormatUint(stock.LocalIB, 10),
		strconv.FormatUint(stock.LocalID, 10),
		strconv.FormatUint(stock.LocalMF, 10),
		strconv.FormatUint(stock.LocalSC, 10),
		strconv.FormatUint(stock.LocalFD, 10),
		strconv.FormatUint(stock.LocalOT, 10),
		strconv.FormatUint(stock.ForeignIS, 10),
		strconv.FormatUint(stock.ForeignCP, 10),
		strconv.FormatUint(stock.ForeignPF, 10),
		strconv.FormatUint(stock.ForeignIB, 10),
		strconv.FormatUint(stock.ForeignID, 10),
		strconv.FormatUint(stock.ForeignMF, 10),
		strconv.FormatUint(stock.ForeignSC, 10),
		strconv.FormatUint(stock.ForeignFD, 10),
		strconv.FormatUint(stock.ForeignOT, 10),
	}

	return w.Write(record)
}
