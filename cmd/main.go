/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2021/4/23 15:58
 */

package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
)

type Result struct {
	Number int
	HF     string
	HFF    float64
	HFCut  string
	File   string
}

const PATH = "./data"

func main() {

	files, err := ioutil.ReadDir(PATH)
	if err != nil {
		fmt.Println(err)
	}
	ResultChan := make(chan Result, len(files))
	wg := &sync.WaitGroup{}
	for _, file := range files {
		wg.Add(1)
		go Worker(wg, file.Name(), ResultChan)
	}
	var results Results
	verifyMap := make(map[string]int, 0)
	go monitorWorker(wg, ResultChan)
	for task := range ResultChan {
		//if _, ok := verifyMap[task.HFCut]; !ok {
			results = append(results, task)
			//verifyMap[task.HFCut] = 1
		//}
	}
	//fmt.Println(results)
	sort.Sort(results)
	//fmt.Println(results)
	WriteExcel(results)
	//计算结果超过25
	folderName := "ti" + time.Now().Format("20060102150405")
	CreateFolder(folderName)
	for _, val := range results {
		//var c float64
		var tempVal decimal.Decimal
		/**
		if key == 0 {
			c = 0
		} else {
			tempVal = decimal.NewFromFloat(val.HFF).Sub(decimal.NewFromFloat(results[0].HFF))
			c, _ = tempVal.Float64()
		}

		 */
		tempVal = decimal.NewFromFloat(val.HFF).Sub(decimal.NewFromFloat(results[0].HFF))
		d := tempVal.Mul(decimal.NewFromFloat(627.5))
		//fmt.Println(val.Number, val.HF, c, d)
		_, eok := verifyMap[val.HFCut]
		dok :=d.LessThanOrEqual(decimal.NewFromFloat(2.5))
		if  dok && !eok {
			fmt.Println(val.File)
			CopyFile(folderName+"/"+val.File, PATH+"/"+val.File)
		}

		verifyMap[val.HFCut] = 1
	}
}

func monitorWorker(wg *sync.WaitGroup, cs chan Result) {
	wg.Wait()
	close(cs)
}

func Worker(wg *sync.WaitGroup, fileName string, ResultChan chan Result) {
	defer wg.Done()
	fileNumber := GetFileNumber(fileName)
	var result Result
	hf := GetFileHF(fileName)

	if len(hf) > 0 {
		result.Number = cast.ToInt(fileNumber)
		result.HF = hf
		result.HFCut = fmt.Sprintf("%.5f", cast.ToFloat64(hf))
		result.HFF = cast.ToFloat64(hf)
		result.File = fileName
		ResultChan <- result
	}

}

func GetFileNumber(fileName string) string {
	//temp1 := strings.Split(fileName, ".")
	//temp2 := strings.Split(temp1[0], "-")
	str := `[-|_]0*([1-9][0-9]*)\.`
	Regexp := regexp.MustCompile(str)
	params := Regexp.FindStringSubmatch(fileName)
	//for _,param :=range params {
	//	fmt.Println(param)
	//}
	return params[1]
}

func GetFileHF(fileName string) string {
	fileName = "./data/" + fileName
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("open %s failed:%s", fileName, err)
	}
	//str := `HF=(-?\d+.\d+)\\`
	str := `HF=(-?\d+.\d+)\\`
	Regexp := regexp.MustCompile(str)
	//去除空白字符
	temp := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, string(file))
	params := Regexp.FindStringSubmatch(temp)
	if len(params) > 0 {
		return params[1]
	}
	return ""
}

func WriteExcel(Content []Result) {
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet1")
	// Set value of a cell.
	for key, val := range Content {
		f.SetCellValue("Sheet1", "A"+cast.ToString(key+1), val.Number)
		f.SetCellValue("Sheet1", "B"+cast.ToString(key+1), val.HF)
		if key == 0 {
			f.SetCellValue("Sheet1", "C"+cast.ToString(key+1), "0")
		} else {
			value := decimal.NewFromFloat(Content[key].HFF).Sub(decimal.NewFromFloat(Content[0].HFF))
			f.SetCellFormula("Sheet1", "C"+cast.ToString(key+1), value.String())
		}
		f.SetCellFormula("Sheet1", "D"+cast.ToString(key+1), "C"+cast.ToString(key+1)+"*627.5")
		f.SetCellValue("Sheet1", "E"+cast.ToString(key+1), val.HFCut)

	}

	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	fileName := time.Now().Format("20060102150405") + ".xlsx"
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}

}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateFolder(path string) {
	exist, err := PathExists(path)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}

	if exist {
		fmt.Printf("has dir![%v]\n", path)
	} else {
		fmt.Printf("no dir![%v]\n", path)
		// 创建文件夹
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

type Results []Result

//排序方法
func (r Results) Len() int { return len(r) }
func (r Results) Less(i, j int) bool {
	return r[i].HF > r[j].HF
}
func (r Results) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
