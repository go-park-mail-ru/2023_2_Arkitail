package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Place struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	Cost        string  `json:"cost"`
	ImageURL    string  `json:"imageUrl"`
}

var places = map[string]Place{
	"1": {
		ID:          "1",
		Name:        "Эфелева башня",
		Description: "Это знаменитое архитектурное сооружение, которое находится в центре Парижа, Франция. Эта башня является одной из самых узнаваемых и посещаемых достопримечательностей мира, а также символом как самого Парижа, так и Франции в целом. Она была построена для Всемирной выставки 1889 года, которая отмечала столетие Великой французской революции.",
		Rating:      4.5,
		Cost:        "$$",
		ImageURL:    "https://example.com/image1.jpg",
	},
	"2": {
		ID:          "2",
		Name:        "Эрмитаж",
		Description: "Это один из самых знаменитых и крупнейших музеев мира, расположенный в Санкт-Петербурге, Россия. Этот музей является одной из наиболее значимых культурных достопримечательностей России и мировым центром искусства и культуры.",
		Rating:      3.8,
		Cost:        "$",
		ImageURL:    "https://example.com/image2.jpg",
	},
	"3": {
		ID:          "3",
		Name:        "МГТУ им. Баумана",
		Description: "Является одним из ведущих технических университетов в России и весьма престижным учебным заведением.",
		Rating:      5.0,
		Cost:        "$$$",
		ImageURL:    "https://example.com/image3.jpg",
	},
	"4": {
        ID:          "4",
        Name:        "Петра I памятник",
        Description: "Памятник Петру I, также известный как Бронзовый всадник, - это памятник российскому императору Петру I, установленный в Санкт-Петербурге.",
        Rating:      4.7,
        Cost:        "$$",
        ImageURL:    "https://example.com/image4.jpg",
    },
    "5": {
        ID:          "5",
        Name:        "Статуя Свободы",
        Description: "Статуя Свободы находится на острове Свободы в Нью-Йорке и является одним из символов Соединенных Штатов Америки.",
        Rating:      4.6,
        Cost:        "$$",
        ImageURL:    "https://example.com/image5.jpg",
    },
    "6": {
        ID:          "6",
        Name:        "Гренландия",
        Description: "Гренландия - крупнейший остров в мире и административно-территориальное подразделение Королевства Дании.",
        Rating:      4.2,
        Cost:        "$$$",
        ImageURL:    "https://example.com/image6.jpg",
    },
    "7": {
        ID:          "7",
        Name:        "Колизей",
        Description: "Колизей - это амфитеатр в Риме, построенный в I веке н.э. и считающийся одним из величайших архитектурных и инженерных достижений древнего мира.",
        Rating:      4.9,
        Cost:        "$$",
        ImageURL:    "https://example.com/image7.jpg",
    },
    "8": {
        ID:          "8",
        Name:        "Маяк Александрия",
        Description: "Маяк Александрия был одним из семи чудес света и находился в древнем городе Александрия, в Египте.",
        Rating:      4.4,
        Cost:        "$$$",
        ImageURL:    "https://example.com/image8.jpg",
    },
    "9": {
        ID:          "9",
        Name:        "Скайдайвинг в Нью Зеландии",
        Description: "Нью Зеландия предлагает невероятные возможности для скайдайвинга с потрясающими видами на природную красоту страны.",
        Rating:      4.8,
        Cost:        "$$$",
        ImageURL:    "https://example.com/image9.jpg",
    },
}


func CreatePlace(w http.ResponseWriter, r *http.Request) {
	var newPlace Place
	err := json.NewDecoder(r.Body).Decode(&newPlace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newPlace.ID = generateUniqueID()

	places[newPlace.ID] = newPlace

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": newPlace.ID})
}

func GetPlaces(w http.ResponseWriter, r *http.Request) {
	placeList := []Place{}
	for _, place := range places {
		placeList = append(placeList, place)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(placeList)
}

func generateUniqueID() string {
	return fmt.Sprintf("id%d", len(places)+1)
}
