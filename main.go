package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Item map[string]int

func main() {

	var fileName string
	var fileSample *os.File
	// Get Name of sample and open it
	for {
		fmt.Print("Введите имя файла (без расширения): ")
		fmt.Scan(&fileName)
		if f, err := os.Open(fmt.Sprintf("./samples/%v.json", fileName)); err != nil {
			fmt.Println("Ошибка: ", err.Error())
		} else {
			fileSample = f
			break
		}
	}
	defer fileSample.Close()

	var routine string
	var routineNum int
	// Get num of goroutines
	for {
		fmt.Print("Введите количество гороутин(1-4): ")
		fmt.Scan(&routine)
		if num, err := strconv.Atoi(routine); err == nil && num > 0 && num < 5 {
			routineNum = num
			break
		} else if err != nil {
			fmt.Println("Неправильное число. Ошибка:", err.Error())
		} else {
			fmt.Println("Неправильное число. Ошибка: перевысил либо не дотянул лимит 1-4")
		}
	}
	// Разделение на блоков для гороутин
	InfoFile, err := fileSample.Stat()
	if err != nil {
		fmt.Println("Ошибка при чтений файла: ", err.Error())
		return
	}
	// 8 is minimum for valid array of struct
	if InfoFile.Size() < 8 {
		fmt.Println("Ошибка формата")
		return
	}
	blockForRoutine := int(InfoFile.Size() / int64(routineNum))
	// 6 is minimum for valid block
	if blockForRoutine < 6 {
		blockForRoutine = int(InfoFile.Size())
	}

	FileReader := bufio.NewReader(fileSample)
	Ostatok := []byte{}
	// getting to start of array and escaping [
	FileReader.ReadBytes('[')
	chanOfSum := make(chan int)
	chanOfErr := make(chan error)
	numRealRoutine := 0
	for i := 1; i <= routineNum; i++ {
		if i == routineNum {
			blockForRoutine += routineNum
		}
		// getting left bytes
		BytesToGoRoutine := make([]byte, 0, int(blockForRoutine)+len(Ostatok))
		BytesToGoRoutine = append(BytesToGoRoutine, '[')
		BytesToGoRoutine = append(BytesToGoRoutine, Ostatok...)

		BytesToRead := make([]byte, int(blockForRoutine))
		numReadBytes, err := FileReader.Read(BytesToRead)
		if err != nil {
			fmt.Println("reading error:", err.Error())
		}
		// Ended
		if numReadBytes == 0 {
			break
		}
		BytesToRead = BytesToRead[:numReadBytes]
		isEnd := bytes.LastIndex(BytesToRead, []byte{']'})
		if isEnd == -1 {
			lastStructEnd := bytes.LastIndex(BytesToRead, []byte{'}', ','})
			if lastStructEnd != -1 {
				if len(BytesToRead) > lastStructEnd+2 {
					Ostatok = BytesToRead[lastStructEnd+2:]
				} else {
					Ostatok = []byte{}
				}
				BytesToGoRoutine = append(BytesToGoRoutine, BytesToRead[:lastStructEnd+1]...)
			}
		} else {
			lastStructEnd := bytes.LastIndex(BytesToRead, []byte{'}'})
			BytesToGoRoutine = append(BytesToGoRoutine, BytesToRead[:lastStructEnd+1]...)
		}
		BytesToGoRoutine = append(BytesToGoRoutine, ']')
		numRealRoutine++
		go func(bytes []byte) {
			sum, err := getSumFromJson(bytes)
			chanOfSum <- sum
			chanOfErr <- err
		}(BytesToGoRoutine)

	}
	res:=0
	for i:=1;i<=numRealRoutine;i++{
		res+=<-chanOfSum
		err=<-chanOfErr
		if err!=nil{
			fmt.Println(err)
		}
	}
	fmt.Println("Сумма:",res)
}

func getSumFromJson(jsonStr []byte) (int, error) {
	res := 0
	// Decoding (unmarshaling) the item struct
	var items []Item
	err := json.Unmarshal(jsonStr, &items)
	if err != nil {
		fmt.Println(json.Valid(jsonStr))
		return res, err
	}
	for _, item := range items {
		for _, num := range item {
			res += num
		}
	}
	return res, nil
}
