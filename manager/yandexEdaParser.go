package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"YandexEdaParser/db"
	"YandexEdaParser/structs"
)

var err error

func (y YandexManager) RunParser(ctx context.Context, jobs *int, ch chan string) (err error) {
	log.Println("Начал многопоточный парсинг, потоков ", 8)

	getRestHttp := fmt.Sprintf("https://eda.yandex.ru/api/v2/catalog?latitude=%f&longitude=%f&rating=4.8", y.Latitude, y.Longitude)

	var allRestorationListResponse *http.Response
	if allRestorationListResponse, err = http.Get(getRestHttp); err != nil {
		return errors.Wrap(err, "Error Get from Yandex")
	}
	var bodyAllRestorationListResponse []byte
	if bodyAllRestorationListResponse, err = io.ReadAll(allRestorationListResponse.Body); err != nil {
		return errors.Wrap(err, "Error Read Body allRestorationListResponse")
	}

	var Response structs.ResponseRestaurant
	err = json.Unmarshal(bodyAllRestorationListResponse, &Response)
	restInput := make(chan structs.Restaurant, *jobs)
	done := make(chan bool)

	//запускаем джобы ранера (ожидает на вход рестораны и ищет в них филу)
	for i := 0; i < *jobs; i++ {
		go runner(restInput, done, y.Repository, i)
	}

	//передаем рестораны в джобы
	for _, place := range Response.Payload.FoundPlaces {
		if place.Restaurant.Rating >= y.Rating {
			restInput <- place.Restaurant
		}
	}
	var zeroRest structs.Restaurant
	for i := 0; i < *jobs; i++ {
		restInput <- zeroRest
		<-done
	}
	ch <- "Done Parse"
	return nil
}

// Поток-runner.
func runner(restInput chan structs.Restaurant, done chan bool, DB db.Repository, i int) {
	// Бесконечный цикл.
	log.Println("Запустил обработчик в потоке", i+1)
	for {
		rest := <-restInput
		// Если YandexId пустое, нас просят завершиться.
		if rest.YandexId == 0 {
			break
		}
		//log.Println("Обрабатываю ", rest.Name, ",в потоке ", i+1)
		getMenuHttp := fmt.Sprintf("https://eda.yandex.ru/api/v2/menu/retrieve/%s?autoTranslate=false", rest.Slug)
		allRestorationListResponse, _ := http.Get(getMenuHttp)
		bodyMenuListResponse, _ := io.ReadAll(allRestorationListResponse.Body)
		var Response structs.ResponseMenu
		json.Unmarshal(bodyMenuListResponse, &Response)
		//парсим меню ресторана, ищем филу, есди находим сохраняем ресторан и позицию в базу
		for _, category := range Response.Payload.Categories {
			for _, item := range category.Items {
				if strings.Contains(item.Name, "Филадельфия с лососем") {
					item.Measure.ValueInt, _ = strconv.Atoi(item.Measure.Value)
					//расчитываю внутрений рейтинг позиции и меняю рейтинг ресторана если эта позиция лучше той которая в базе сейчас.
					item.InternalRating = (item.Price)/float64(item.Measure.ValueInt) + rest.Rating
					rest.InternalRating = item.InternalRating
					//получаю текущий рейтинг
					var currentRating float64
					if currentRating, err = DB.GetRestaurantInternalRating(rest.ID); err != nil {
						log.Println(err)
					}
					if currentRating < rest.InternalRating {
						if item.RestaurantId, err = DB.AddOrUpdateRestaurant(rest); err != nil {
							log.Println(err)
						}
					}
					DB.AddOrUpdateItem(item)
				}
			}
		}
	}
	// Посылаем сообщение, что поток завершился.
	log.Println("Завершил обаботку в потоке", i+1)
	done <- true
}

func (y *YandexManager) GetAllRestaurants() (Restaurants []structs.Restaurant, err error) {
	log.Println("Запросили все рестораны - отдаю")
	if Restaurants, err = y.Repository.GetAllRestaurants(); err != nil {
		return nil, err
	}
	return Restaurants, nil
}

func (y *YandexManager) GetRestaurantPrices(RestId int) (Prices []int, err error) {
	log.Println("Запросили цены ресторана - отдаю")
	if Prices, err = y.Repository.GetRestaurantPrices(RestId); err != nil {
		return nil, err
	}
	return Prices, nil
}
