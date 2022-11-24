package main

import (
	usergrade "HW_WB"
	"HW_WB/nats"
	"HW_WB/storage"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Нельзя хранить в коде, в энв ОС?
var (
	username = "test"
	password = "12345"
)

func BasicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			fmt.Println("Ошибка разбора данных")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if user != username {
			fmt.Printf("Некорректное имя пользователя: %s\n", user)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if pass != password {
			fmt.Printf("Некорректный пароль: %s\n", pass)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func SetHeaderWithJsonMsg(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{" + msg + "}"))

}

func SetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		SetHeaderWithJsonMsg(w, "error: Неправильный метод передачи данных, нужен "+http.MethodPost)
	} else {

		body, _ := ioutil.ReadAll(r.Body)
		userGradeFR := usergrade.JsonToUserGradeFR(body)
		ve := usergrade.ValidateStruct(userGradeFR)
		if ve != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Fatalf("Данные не провалидированы. Err: %s", ve)
		} else {
			var userGrade usergrade.UserGrade
			userGradeFromStore, inStore := storage.GetStore(userGradeFR.UserId)
			if inStore {
				userGrade = usergrade.MatchUserGrade(userGradeFromStore, userGradeFR)
			} else {
				userGrade = usergrade.UserGrade{
					UserId:        userGradeFR.UserId,
					PostpaidLimit: userGradeFR.PostpaidLimit.Value,
					Spp:           userGradeFR.Spp.Value,
					ShippingFee:   userGradeFR.ShippingFee.Value,
					ReturnFee:     userGradeFR.ReturnFee.Value,
				}
			}

			storage.SetStore(userGrade)
			updateUg, _ := json.Marshal(userGrade)
			w.Header().Set("Content-Type", "application/json")
			w.Write(updateUg)
			nats.Publish("UserGrades.*", AppName, userGrade)
		}
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		SetHeaderWithJsonMsg(w, "error: Неправильный метод передачи данных, нужен "+http.MethodGet)
	} else {
		userId := r.URL.Query().Get("user_id")
		if userId != "" {
			if ug, inStore := storage.GetStore(userId); inStore {
				userGrade := usergrade.GetUserGradeToJson(ug)
				w.Header().Set("Content-Type", "application/json")
				w.Write(userGrade)
			} else {
				w.WriteHeader(http.StatusNotFound)
				SetHeaderWithJsonMsg(w, "{error: Не обнаружен пользователь с таким user_id}")
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{error: Не обнаружен user_id}"))
		}
	}
}

var AppName string
var SetPort string
var GetPort string

func main() {
	flag.StringVar(&AppName, "app_name", "app", "приложение")
	flag.StringVar(&SetPort, "set_port", "8080", "порт для установки данных")
	flag.StringVar(&GetPort, "get_port", "8081", "порт для отправки данных")
	flag.Parse()

	fmt.Println(AppName, SetPort, GetPort)

	nats.NatsJsConnect()

	Mux1 := http.NewServeMux()
	Mux1.HandleFunc("/set", BasicAuthMiddleware(SetHandler))

	Mux2 := http.NewServeMux()
	Mux2.HandleFunc("/get", GetHandler)

	go func() {
		nats.Subscribe("UserGrades.*", AppName) // надо добавить время

	}()

	fmt.Println("Запускаем сервер на порту 8080...")
	go func() {
		log.Fatalln(http.ListenAndServe(":"+SetPort, Mux1))

	}()
	fmt.Println("Запускаем сервер на порту 8081...")

	log.Fatalln(http.ListenAndServe(":"+GetPort, Mux2))

}
