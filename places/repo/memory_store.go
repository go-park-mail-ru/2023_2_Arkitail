package repo

import (
    "sync"
    "project/places/model"
)

type MemoryStore struct {
    places map[string]model.Place
    mu     sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
    places := map[string]model.Place{
        "1": {
			ID:          "1",
			Name:        "Эфелева башня",
			Description: "Это знаменитое архитектурное сооружение, которое находится в центре Парижа, Франция. Эта башня является одной из самых узнаваемых и посещаемых достопримечательностей мира, а также символом как самого Парижа, так и Франции в целом.",
			Rating:      4.5,
			Cost:        "$$",
			ImageURL:    "https://img-s-msn-com.akamaized.net/tenant/amp/entityid/AA1fmKP5.img",
		},
		"2": {
			ID:          "2",
			Name:        "Эрмитаж",
			Description: "Это один из самых знаменитых и крупнейших музеев мира, расположенный в Санкт-Петербурге, Россия. Этот музей является одной из наиболее значимых культурных достопримечательностей России и мировым центром искусства и культуры.",
			Rating:      3.8,
			Cost:        "$",
			ImageURL:    "https://mykaleidoscope.ru/x/uploads/posts/2022-09/1663090921_7-mykaleidoscope-ru-p-zimnii-dvorets-sankt-peterburg-krasivo-7.jpg",
		},
		"3": {
			ID:          "3",
			Name:        "МГТУ им. Баумана",
			Description: "Является одним из ведущих технических университетов в России и весьма престижным учебным заведением.",
			Rating:      5.0,
			Cost:        "$$$",
			ImageURL:    "https://sun6-23.userapi.com/XEbCUs5UIcV3L-JP87lxuKEWyRl9KgbNwaU91g/3ywb_ZTuGMs.jpg",
		},
		"4": {
			ID:          "4",
			Name:        "Петра I памятник",
			Description: "Памятник Петру I, также известный как Бронзовый всадник, - это памятник российскому императору Петру I, установленный в Санкт-Петербурге.",
			Rating:      4.7,
			Cost:        "$$",
			ImageURL:    "https://img.tourister.ru/files/1/8/8/7/6/5/6/0/original.jpg",
		},
		"5": {
			ID:          "5",
			Name:        "Статуя Свободы",
			Description: "Статуя Свободы находится на острове Свободы в Нью-Йорке и является одним из символов Соединенных Штатов Америки.",
			Rating:      4.6,
			Cost:        "$$",
			ImageURL:    "https://i.imgur.com/pOSbnHXh.jpg",
		},
		"6": {
			ID:          "6",
			Name:        "Гренландия",
			Description: "Гренландия - крупнейший остров в мире и административно-территориальное подразделение Королевства Дании.",
			Rating:      4.2,
			Cost:        "$$$",
			ImageURL:    "https://10wallpaper.com/wallpaper/5120x2880/2103/Windows_10x_Microsoft_2021_Ocean_Glacier_5K_HD_Photo_5120x2880.jpg",
		},
		"7": {
			ID:          "7",
			Name:        "Колизей",
			Description: "Колизей - это амфитеатр в Риме, построенный в I веке н.э. и считающийся одним из величайших архитектурных и инженерных достижений древнего мира.",
			Rating:      4.9,
			Cost:        "$$",
			ImageURL:    "https://sportishka.com/uploads/posts/2022-04/1650595488_15-sportishka-com-p-italiya-kolizei-krasivo-foto-15.jpg",
		},
		"8": {
			ID:          "8",
			Name:        "Маяк Александрия",
			Description: "Маяк Александрия был одним из семи чудес света и находился в древнем городе Александрия, в Египте.",
			Rating:      4.4,
			Cost:        "$$$",
			ImageURL:    "https://polinka.top/uploads/posts/2023-06/thumbs/1685742571_polinka-top-p-kartinka-aleksandriya-yegipetskaya-instagr-44.jpg",
		},
		"9": {
			ID:          "9",
			Name:        "Скайдайвинг в Нью Зеландии",
			Description: "Нью Зеландия предлагает невероятные возможности для скайдайвинга с потрясающими видами на природную красоту страны.",
			Rating:      4.8,
			Cost:        "$$$",
			ImageURL:    "https://www.dropzone.com/uploads/monthly_2020_01/Klavs.jpg.a67700687ec1204e6c6fedc92037bc8b.jpg",
		},
    }

    return &MemoryStore{
        places: places,
    }
}

func (ms *MemoryStore) GetAllPlaces() ([]model.Place, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var result []model.Place
	for _, place := range ms.places {
		result = append(result, place)
	}
	return result, nil
}

func (ms *MemoryStore) AddPlace(place model.Place) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.places[place.ID] = place
	return nil
}
