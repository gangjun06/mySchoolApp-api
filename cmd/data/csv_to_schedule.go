package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type output struct {
	Dow         int    `json:"dow"`
	Period      int    `json:"period"`
	Grade       int    `json:"grade"`
	Class       int    `json:"class"`
	Subject     string `json:"subject"`
	Teacher     string `json:"teacher"`
	Description string `json:"description"`
	ClassRoom   string `json:"classRoom"`
}

var splitDays = []int{7, 7, 6, 7, 7}
var classSplit = []int{2, 2, 3}

func main() {
	var class [][]int
	for i, d := range classSplit {
		for j := 0; j < d; j++ {
			class = append(class, []int{(i + 1), (j + 1)})
		}
	}
	file, err := os.Open("./cmd/data/schedule.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	var result []output

	realIndex := 0
	for i := range records {
		if i != 0 && i%3 != 0 {
			continue
		}
		fmt.Println(i, "-----", class[realIndex][0], "/", class[realIndex][1])
		index := 0
		for iDOW, d := range splitDays {
			fmt.Print("요일: ", iDOW)
			for j := 0; j < d; j++ {
				fmt.Print(records[i][index])
				result = append(result, output{
					Dow:         iDOW + 1,
					Period:      j + 1,
					Grade:       class[realIndex][0],
					Class:       class[realIndex][1],
					Subject:     records[i][index],
					Teacher:     records[i+1][index],
					ClassRoom:   records[i+2][index],
					Description: "",
				})
				index += 1
			}
			fmt.Println("")
		}
		realIndex += 1
	}
	resultFile, err := json.Marshal(&result)
	if err != nil {
		log.Fatalln(err)
	}
	ioutil.WriteFile("output.json", resultFile, 0644)
}
