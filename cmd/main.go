/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2021/4/23 15:58
 */

package main

import (
	"fmt"
	"github.com/spf13/cast"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"time"
)

type Result struct {
	Number int
	HF     string
	HFF    float64
	HFCut  float64
}

func main() {

	files, err := ioutil.ReadDir("./data")
	if err != nil {
		fmt.Println(err)
	}
	ResultChan := make(chan Result, len(files))
	for _, file := range files {
		go Worker(file.Name(), ResultChan)
	}
	var results Results
	for i := 0; i < len(files); i++ {
		results = append(results,<-ResultChan)
	}
	sort.Sort(results)
	//fmt.Println(results)
	WriteExcel(results)
}

func Worker(fileName string, ResultChan chan Result) {

	fileNumber := GetFileNumber(fileName)
	var result Result
	result.Number = cast.ToInt(fileNumber)
	hf := GetFileHF(fileName)
	result.HF =hf
	result.HFCut= cast.ToFloat64(hf[0:len(hf)-2])
	result.HFF = cast.ToFloat64(hf)
	ResultChan <- result

}

func GetFileNumber(fileName string) string {
	temp1 := strings.Split(fileName, ".")
	temp2 := strings.Split(temp1[0], "-")
	return temp2[1]
}

func GetFileHF(fileName string) string {
	fileName = "./data/" + fileName
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("open %s failed:%s", fileName, err)
	}
	str := `HF=(-?\d+.\d+)\\`
	Regexp := regexp.MustCompile(str)
	params := Regexp.FindStringSubmatch(string(file))
	if len(params) > 0 {
		return params[1]
	}
	return ""
}


func WriteExcel(Content []Result){
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet1")
	// Set value of a cell.
	for key,val := range Content{
		f.SetCellValue("Sheet1", "A"+cast.ToString(key), val.Number)
		f.SetCellValue("Sheet1", "B"+cast.ToString(key), val.HF)
	}

	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	fileName := time.Now().Format("20060102150405")+".xlsx"
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}

}

type Results []Result
//排序方法
func(r Results) Len() int{ return len(r)}
func(r Results) Less(i,j int) bool{
	return r[i].HF >r[j].HF
}
func(r Results) Swap(i,j int){r[i],r[j] = r[j],r[i]}
