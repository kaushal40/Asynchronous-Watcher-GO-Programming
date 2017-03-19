//directory file changes watcher
//when you add any file to the source folder it will trigger the routine
package main

import (
	"fmt"
	//get access to file system
	"os"
	//convert time
	"encoding/csv"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const watchedPath = "./inv"

func main() {
	//infinite loop for watching
	for {
		d, _ := os.Open(watchedPath)
		//negative value return all files if positive number used it will limit the number of files
		files, _ := d.Readdir(-1)
		//we are using two CPU here
		//		runtime.GOMAXPROCS(2)
		for _, fi := range files {
			filePath := watchedPath + "/" + fi.Name()
			f, _ := os.Open(filePath)
			data, _ := ioutil.ReadAll(f)
			// we dont need to put defer call because we don't need to wait
			f.Close()
			// we can remove the file once used
			os.Remove(filePath)
			go func(data string) {
				reader := csv.NewReader(strings.NewReader(data))
				records, _ := reader.ReadAll()
				for _, r := range records {
					invoice := new(Invoice)
					invoice.Number = r[0]
					//string to float
					invoice.Amount, _ = strconv.ParseFloat(r[1], 64)
					//string to integer
					invoice.PurchaseOrderNumber, _ = strconv.Atoi(r[2])
					unixTime, _ := strconv.ParseInt(r[3], 10, 64)
					invoice.InvoiceDate = time.Unix(unixTime, 0)

					fmt.Printf("Received Invoice '%v' for $%.2f and submitted for processing\n", invoice.Number, invoice.Amount)
				}
			}(string(data))
		}
		d.Close()
		time.Sleep(100 * time.Millisecond)
	}
}

type Invoice struct {
	Number              string
	Amount              float64
	PurchaseOrderNumber int
	InvoiceDate         time.Time
}
