package main

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"log"
	"net/http"
)

// Batch is a batch of items.
type Batch []Item

// Item is some abstract item.
type Item struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Information string `json:"information"`
}

func main() {
	GetLimit()
	SendBatch()
}

// GetLimit получает лимиты сервисы
func GetLimit() {
	query := "http://localhost:8080/limit"
	response, err := http.Get(query)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}

// NewBatch создает новый набор элементов для отправки
func NewBatch(size int) Batch {
	batch := make(Batch, size)

	for i := 0; i < size; i++ {
		item := Item{
			Id:          i + 1,
			Name:        fmt.Sprint("Item", i+1),
			Information: fmt.Sprint("Information about Item", i+1),
		}
		batch[i] = item
	}

	return batch
}

// SendBatch отправляет данные на сервис
func SendBatch() {
	size := 100
	batch := NewBatch(size)
	jsonData, err := json.Marshal(batch)
	if err != nil {
		log.Fatal("Ошибкка преобразования данных")
	}
	query := "http://localhost:8080/send"
	buffer := bytes.NewBuffer(jsonData)
	response, err := http.Post(query, "application/json", buffer)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
