package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	fileCount = 200
	fileSize  = int64(100000000)
	fileDir   = "/Users/zhangyinhao/work/softDev/data/"
	csvDir    = "/Users/zhangyinhao/work/softDev/data/csv/"
)

/**
* 初始化短码数据导入脚本  200亿约500G的CSV文件
 */
func main() {
	rmAllFile()
	buildData()
	shuffleFile()
	writeSerialNumber()
	println("数据处理完成")
}

func rmAllFile() {
	for i := 0; i < fileCount; i++ {
		if err := os.Remove(getFileName(i)); err != nil {
			log.Println("未找到文件")
		}
	}
}

func buildData() {
	total := int64(0)
	for index := 0; index < fileCount; index++ {
		var buf strings.Builder
		for c := int64(0); c < fileSize; c++ {
			buf.WriteString(strconv.FormatInt(total, 10) + "\n")
			total++
		}
		file, _ := os.OpenFile(getFileName(index), os.O_WRONLY|os.O_CREATE, 0666)
		_, err := file.WriteString(buf.String())
		if err != nil {
			panic(err)
		}
	_:
		file.Close()
	}
}

func getFileName(index int) string {
	return fileDir + strconv.Itoa(index) + ".txt"
}

func csvFileName(index int) string {
	return csvDir + strconv.Itoa(index) + ".csv"
}

func shuffleFile() {
	arr := createFileArr()
	for j := 0; j < 2; j++ {
		shuffle(arr)
		for i := 0; i < fileCount; i = i + 2 {
			merge(arr[i], arr[i+1])
		}
	}
}

func merge(index_a int, index_b int) {
	arr1 := readFile(index_a)
	arr2 := readFile(index_b)
	arr1 = append(arr1, arr2...)
	shuffleInt64(arr1)
	tmp1 := arr1[0:fileSize]
	tmp2 := arr1[fileSize:]
	writeFile(index_a, tmp1)
	writeFile(index_b, tmp2)
}

func writeSerialNumber() {
	serialNumber := int64(0)
	for i := 0; i < fileCount; i++ {
		arr := readFile(i)
	_:
		os.Remove(csvFileName(i))
		var buf strings.Builder
		for _, val := range arr {
			buf.WriteString(strconv.FormatInt(val, 10) + "," + strconv.FormatInt(serialNumber, 10) + "\n")
			serialNumber++
		}
		file, _ := os.OpenFile(csvFileName(i), os.O_WRONLY|os.O_CREATE, 0666)
		file.WriteString("code,serial_number" + "\n")
		file.WriteString(buf.String())
		file.Close()
	}
}

func readFile(index int) []int64 {
	arr := make([]int64, fileSize, fileSize*2)
	file, _ := os.OpenFile(getFileName(index), os.O_RDONLY, 0666)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	var i = int64(0)
	for scanner.Scan() {
		val, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		arr[i] = val
		i++
	}
	return arr
}

func writeFile(index int, arr []int64) {
	os.Remove(getFileName(index))
	var buf strings.Builder
	for _, val := range arr {
		buf.WriteString(strconv.FormatInt(val, 10) + "\n")
	}
	file, _ := os.OpenFile(getFileName(index), os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	file.WriteString(buf.String())
}

func shuffleInt64(arr []int64) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}

func shuffle(arr []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}

func createFileArr() []int {
	var files []int
	for i := 0; i < fileCount; i++ {
		files = append(files, i)
	}
	return files
}
