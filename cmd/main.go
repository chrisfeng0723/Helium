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
)

type Result struct {
	Number int
	HF     string
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
	//找出最后两位不同的内容
	resultMap := make(map[string]int,0)
	for _,val := range results{
		tempHF:=val.HF
		temp :=tempHF[0:len(tempHF)-2]
		resultMap[temp] = val.Number
	}
	//转成results数组
	filterResult :=make(Results,0)
	for key,val :=range resultMap{
		filterResult = append(filterResult,Result{
			Number: val,
			HF:     key,
		})
	}
	sort.Sort(results)
	fmt.Println(results)
}

func Worker(fileName string, ResultChan chan Result) {

	fileNumber := GetFileNumber(fileName)
	var result Result
	result.Number = cast.ToInt(fileNumber)
	hf := GetFileHF(fileName)
	result.HF =hf
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

type Results []Result
//排序方法
func(r Results) Len() int{ return len(r)}
func(r Results) Less(i,j int) bool{
	return r[i].HF >r[j].HF
}
func(r Results) Swap(i,j int){r[i],r[j] = r[j],r[i]}
