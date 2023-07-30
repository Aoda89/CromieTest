package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")

const (
	urlClient = "localhost:8080"
)

// Service defines external service that can process batches of items.
type Service interface {
	GetLimits() (uint64, time.Duration)
	Process(batch Batch) error
}

type Limit struct {
	Time     time.Duration `json:"time"`
	Quantity uint64        `json:"quantity"`
}

// GetLimits Функция получения лимитов сервиса
func (l *Limit) GetLimits() (uint64, time.Duration) {
	timing := 10 * time.Second // 10 секунд
	var quantity uint64 = 100
	return quantity, timing
}

/*
1)Проверям колличество данных на лимит установленный сервисом

	2)Запускаем таймер который отчсчитывает время которое установленно сервисом для обработки данных. Если время выходит
	  и данные не успевают обработаться сервис блокируется на 10 минут
	3) Данные деляться на блоки по 10 элементов. Каждый блок обрабатывается в отдельной горутине
*/
func (l *Limit) Process(batch Batch) error {
	quantity, timing := l.GetLimits()
	size := len(batch)
	if uint64(size) > quantity {
		return ErrBlocked
	}
	timer := time.NewTimer(timing)
	for i := 0; i <= size; i += 10 {
		go func(startIndex int) {
			for j := startIndex; j < startIndex+10 && j != size; j++ {
				// Логикка работы с полученными данными
			}
		}(i)
	}
	select {
	case <-timer.C:
		time.Sleep(10 * time.Minute)
		return ErrBlocked
	default:
		return nil
	}
	return nil
}

// SendLimits Функция получения лимитов сервиса
func SendLimits(с *gin.Context) {
	limit := Limit{}
	quantity, time := limit.GetLimits()
	limit.Time = time
	limit.Quantity = quantity
	с.JSON(http.StatusOK, limit)
}

// ProcessData Читаем полученные данные и отправляем на выполнение
func ProcessData(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	limit := Limit{}
	if err != nil {
		c.JSON(http.StatusBadRequest, "Ошибка чтения тела запроса")
	}
	var batch Batch
	err = json.Unmarshal(body, &batch)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Ошибка преобразования JSON данных")
	}
	err = limit.Process(batch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, error.Error)
	} else {
		c.JSON(http.StatusOK, "Успешное выполнение")
	}
}

// Batch is a batch of items.
type Batch []Item

// Item is some abstract item.
type Item struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Information string `json:"information"`
}

func main() {
	router := gin.Default()
	router.GET("/limit", SendLimits)
	router.POST("/send", ProcessData)
	router.Run(urlClient)
}
