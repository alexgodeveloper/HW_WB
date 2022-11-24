package usergrade

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
)

type UGfromRequest struct {
	UserId        string  `json:"user_id" validate:"required"`
	PostpaidLimit JSONInt `json:"postpaid_limit"`
	Spp           JSONInt `json:"spp"`
	ShippingFee   JSONInt `json:"shipping_fee"`
	ReturnFee     JSONInt `json:"return_fee"`
}

type UserGrade struct {
	UserId        string `json:"user_id" validate:"required"`
	PostpaidLimit int    `json:"postpaid_limit"`
	Spp           int    `json:"spp"`
	ShippingFee   int    `json:"shipping_fee"`
	ReturnFee     int    `json:"return_fee"`
}

type JSONInt struct {
	Value int
	Valid bool
	Set   bool
}

func (i *JSONInt) UnmarshalJSON(data []byte) error {
	// Если метод был вызван, устанавливаем true
	i.Set = true

	if string(data) == "null" {
		// Если значение пустое, устанавливаем в валидацию false и завершаем обработку поля.
		i.Valid = false
		return nil
	}

	// Если в поле есть значение, десериализируем его.
	var temp int
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Valid = true
	return nil
}

func MatchUserGrade(userGradeFromStore UserGrade, userGrade UGfromRequest) UserGrade {
	if userGrade.PostpaidLimit.Valid && userGrade.PostpaidLimit.Set {
		userGradeFromStore.PostpaidLimit = userGrade.PostpaidLimit.Value
	}
	if userGrade.Spp.Valid && userGrade.Spp.Set {
		userGradeFromStore.Spp = userGrade.Spp.Value
	}
	if userGrade.ShippingFee.Valid && userGrade.ShippingFee.Set {
		userGradeFromStore.ShippingFee = userGrade.ShippingFee.Value
	}
	if userGrade.ReturnFee.Valid && userGrade.ReturnFee.Set {
		userGradeFromStore.ReturnFee = userGrade.ReturnFee.Value
	}
	return userGradeFromStore
}

func ValidateStruct(u UGfromRequest) error {
	var validate *validator.Validate
	validate = validator.New()
	err := validate.Struct(u)
	if err != nil {
		return err
	}
	return nil
}

func JsonToUserGradeFR(body []byte) UGfromRequest {
	var user UGfromRequest
	if err := json.Unmarshal([]byte(body), &user); err != nil {
		log.Fatalf("Произошла ошибка при десериализации JSON. Err: %s", err)
	}
	return user
}

func JsonToUserGrade(body []byte) UserGrade {
	var user UserGrade
	if err := json.Unmarshal([]byte(body), &user); err != nil {
		log.Fatalf("Произошла ошибка при десериализации JSON. Err: %s", err)
	}
	return user
}

func GetUserGradeToJson(userGrade UserGrade) []byte {
	ugJson, err := json.Marshal(userGrade)
	if err != nil {
		log.Fatalf("Произошла ошибка при сериализации JSON. Err: %s", err)
	}
	return ugJson

}
