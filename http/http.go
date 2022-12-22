package http

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"YandexEdaParser/structs"
)

type PublicHandler interface {
	RunParser(ctx context.Context, jobs *int, ch chan string) (err error)
	GetAllRestaurants() (Restaurants []structs.Restaurant, err error)
	GetRestaurantPrices(RestId int) (Prices []int, err error)
}

type application struct {
	auth struct {
		username string
		password string
	}
}

func RunPublic(
	ctx context.Context,
	h PublicHandler,
) error {
	r := mux.NewRouter()
	var jobs *int = flag.Int("jobs", 8, "number of concurrent jobs")
	app := new(application)

	app.auth.username = "test" //os.Getenv("AUTH_USERNAME") //это надо долько если мы считываем логин и пароль как переменные окржения
	app.auth.password = "test" //os.Getenv("AUTH_PASSWORD") //это надо долько если мы считываем логин и пароль как переменные окржения

	//проверяем чтобы переменные в окружение были не нулевыми
	//if app.auth.username == "" {
	//	log.Fatal("basic auth username must be provided")
	//}
	//
	//if app.auth.password == "" {
	//	log.Fatal("basic auth password must be provided")
	//}
	var isParsing bool = false
	//запускаем парсинг ресторанов
	r.HandleFunc("/parse", locker(app.basicAuth(ParserHandler(h, jobs, &isParsing)), &isParsing)).Methods("GET")
	r.HandleFunc("/restaurant", app.basicAuth(RestaurantHandler(h))).Methods("GET")
	r.HandleFunc("/restaurant/{id}", app.basicAuth(OneRestaurantHandler(h))).Methods("GET")

	log.Println("Listening...")
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		//WriteTimeout: 3 * time.Second,
		//ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
	return nil
}

func ParserHandler(h PublicHandler, jobs *int, isParsing *bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body.Read()
		ch := make(chan string, 1)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		go h.RunParser(ctx, jobs, ch)
		select {
		case result := <-ch:
			if result == "Done Parse Timeout" {
				*isParsing = false
				log.Printf(result)
				w.WriteHeader(http.StatusRequestTimeout)
			}
			if result == "Done Parse" {
				log.Printf(result)
				*isParsing = false
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}

func RestaurantHandler(h PublicHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp []structs.Restaurant
		var err error
		if resp, err = h.GetAllRestaurants(); err != nil {
			log.Println(err)
		}
		encoder := json.NewEncoder(w)
		if err = encoder.Encode(&resp); err != nil {
			log.Println("Не смог кодировать в json", err)
		}
		log.Println("Отдал рестораны!")
	}
}

func OneRestaurantHandler(h PublicHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ids := vars["id"]
		var err error
		var id int
		if id, err = strconv.Atoi(ids); err != nil {
			log.Println(err)
		}
		var Prices []int
		if Prices, err = h.GetRestaurantPrices(id); err != nil {
			log.Println(err)
		}
		encoder := json.NewEncoder(w)
		if err = encoder.Encode(&Prices); err != nil {
			log.Println("Не смог кодировать в json", err)
		}
	}
}

//базовая авторизация
func (app *application) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			//хешируем родные и пришедшие значения, чтобы исключить атаку временем
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(app.auth.username))
			expectedPasswordHash := sha256.Sum256([]byte(app.auth.password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

//Locker проверяет выполнятся ли сейчас парсинг
func locker(next http.HandlerFunc, isParsing *bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !*isParsing {
			log.Println("Нет блокировки, запустил парсинг!")
			*isParsing = true
			next.ServeHTTP(w, r)
		} else {
			log.Println("Парсинг уже запущен! Обработан повторный реквест - 208 отдал.")
			w.WriteHeader(http.StatusAlreadyReported)
		}
	})
}
