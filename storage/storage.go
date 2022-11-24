package storage

import (
	usergrade "HW_WB"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var store sync.Map

func SetStore(user usergrade.UserGrade) {
	store.Store(user.UserId, user)
}

func GetStore(userId string) (usergrade.UserGrade, bool) {
	ug, inStore := store.Load(userId)
	if inStore {
		return ug.(usergrade.UserGrade), inStore
	}

	return usergrade.UserGrade{}, inStore
}

func Backup(store sync.Map) {
	createTime := time.Now().Format("2006-01-02-15-04-05.000")
	fileName := "backup_" + createTime + ".csv"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	var str = make([]string, 0, 0)
	store.Range(func(key, value interface{}) bool {
		k := key.(string)
		v, _ := json.Marshal(value)
		str = append(str, k+"|"+string(v))
		return true
	})
	e := writer.Write(str)
	if e != nil {
		fmt.Println(e)
	}
	writer.Flush()
}

func GetBackup(fileName string) (sync.Map, time.Time) {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		fmt.Println(err)

	}
	name := strings.Trim(file.Name(), "backup_")
	name = strings.Trim(name, ".csv")
	backupTime, err := time.Parse("2006-01-02-15-04-05.000", name)
	if err != nil {
		fmt.Println(err)
	}

	reader := csv.NewReader(file)
	rows, err := reader.Read()
	if err != nil {
		fmt.Println(err)
	}
	var store sync.Map

	for _, v := range rows {
		str := strings.Split(v, "|")
		var key string
		var userGrade usergrade.UGfromRequest
		for i := 0; i < 2; i++ {
			if i == 0 {
				key = str[0]
			} else { // ОШИБКА!
				err := json.Unmarshal([]byte(str[1]), &userGrade) //json: cannot unmarshal object into Go struct field UGfromRequest.postpaid_limit of type int

				fmt.Println(err)
				fmt.Println(userGrade)
			}
		}
		store.Store(key, userGrade)
	}
	fmt.Println(store)

	return store, backupTime
}
