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
	"strings"
)

type Result struct {
	Number int
	HF     float64
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
	for i := 0; i < len(files); i++ {
		temp := <-ResultChan
		fmt.Println(temp)
	}
}

func Worker(fileName string, ResultChan chan Result) {

	fileNumber := GetFileNumber(fileName)
	var result Result
	result.Number = cast.ToInt(fileNumber)
	hf := GetFileHF(fileName)
	result.HF = cast.ToFloat64(hf)
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
	/**
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(fileName+"Read file error!", err)
				return ""
			}
		}
		params :=Regexp.FindStringSubmatch(line)
		if len(params) >0{
			return params[1]
		}
	}

	*/
	return ""
}
