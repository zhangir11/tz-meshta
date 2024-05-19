package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Структура для представления произвольных данных
type Item map[string]interface{}

// Функция для генерации случайных строк
func randomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Функция для генерации массива структур с произвольными полями и значениями
func generateItems(count int) []Item {
	items := make([]Item, count)
	for i := 0; i < count; i++ {
		item := make(Item)
		numFields := rand.Intn(5) + 1 // Случайное количество полей от 1 до 5
		for j := 0; j < numFields; j++ {
			fieldName := randomString(rand.Intn(5) + 1)    // Случайное имя поля длиной от 1 до 5 символов
			fieldValue := rand.Intn(100000000) - 200000000 // Случайное значение от -100000000 до 100000000
			item[fieldName] = fieldValue
		}
		items[i] = item
	}
	return items
}

func main() {
	// Запрос имени файла у пользователя
	var fileName string
	fmt.Print("Введите имя файла (без расширения): ")
	fmt.Scan(&fileName)
	fileName += ".json" // Добавляем расширение .json
	var ElementSize string
	var ElementSizeNum int
	for {
		fmt.Print("Введите количество элементов: ")
		fmt.Scan(&ElementSize)
		if num, err := strconv.Atoi(ElementSize); err == nil {
			ElementSizeNum = num
			break
		} else {
			fmt.Println("Неправильное число. Ошибка:",err.Error())
		}
	}
	// Генерация массива структур
	items := generateItems(ElementSizeNum) // Например, генерируем 10 элементов

	// Преобразование массива структур в JSON
	jsonData, err := json.MarshalIndent(items, "", "    ")
	if err != nil {
		log.Fatalf("Ошибка при сериализации в JSON: %s", err)
	}

	// Запись JSON в файл
	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		log.Fatalf("Ошибка при записи файла: %s", err)
	}

	fmt.Printf("JSON данные успешно записаны в файл %s\n", fileName)
}
