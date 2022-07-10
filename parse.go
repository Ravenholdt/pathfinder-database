package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

}

func parseCopyright() {
	jsonMap := make(map[string]string)

	csvFile, _ := os.Open("d20-spell-copyright.csv")
	csvReader := csv.NewReader(csvFile)
	data, _ := csvReader.ReadAll()
	for idx := range data {
		if idx == 0 {
			continue
		}

		var copyright string

		if strings.Contains(data[idx][5], "©") {
			copyright = strings.Split(data[idx][5], "©")[0]
		} else if strings.Contains(data[idx][5], "Copyright") {
			copyright = strings.Split(data[idx][5], "Copyright")[0]
		} else {
			if data[idx][5] != "null" {
				copyright = data[idx][5]
			}
		}

		if strings.Contains(copyright, "15") {
			fmt.Println(data[idx][4])
		}

		copyright = strings.Trim(copyright, ". ")

		jsonMap[data[idx][4]] = copyright

		writeFile, _ := json.MarshalIndent(jsonMap, "", " ")

		ioutil.WriteFile("spells-copyright.json", writeFile, 0644)
	}
}
