package repo

import (
	"database/sql"

	"project/places/model"
)

type PlaceRepository struct {
	DB *sql.DB
}

func NewPlaceRepository(db *sql.DB) *PlaceRepository {
	data := make(map[string]model.Place)
	// data["1"] = model.Place{
	// 	ID:          "1",
	// 	Name:        "Эфелева башня",
	// 	Description: "Это знаменитое архитектурное сооружение, которое находится в центре Парижа, Франция. Эта башня является одной из самых узнаваемых и посещаемых достопримечательностей мира, а также символом как самого Парижа, так и Франции в целом.",
	// 	Rating:      4.5,
	// 	Cost:        "$$",
	// 	ImageURL:    "https://img-s-msn-com.akamaized.net/tenant/amp/entityid/AA1fmKP5.img",
	// }
	// data["2"] = model.Place{
	// 	ID:          "2",
	// 	Name:        "Эрмитаж",
	// 	Description: "Это один из самых знаменитых и крупнейших музеев мира, расположенный в Санкт-Петербурге, Россия. Этот музей является одной из наиболее значимых культурных достопримечательностей России и мировым центром искусства и культуры.",
	// 	Rating:      3.8,
	// 	Cost:        "$",
	// 	ImageURL:    "https://mykaleidoscope.ru/x/uploads/posts/2022-09/1663090921_7-mykaleidoscope-ru-p-zimnii-dvorets-sankt-peterburg-krasivo-7.jpg",
	// }
	// data["3"] = model.Place{
	// 	ID:          "3",
	// 	Name:        "МГТУ им. Баумана",
	// 	Description: "Является одним из ведущих технических университетов в России и весьма престижным учебным заведением.",
	// 	Rating:      5.0,
	// 	Cost:        "$$$",
	// 	ImageURL:    "https://sun6-23.userapi.com/XEbCUs5UIcV3L-JP87lxuKEWyRl9KgbNwaU91g/3ywb_ZTuGMs.jpg",
	// }
	// data["4"] = model.Place{
	// 	ID:          "4",
	// 	Name:        "Петра I памятник",
	// 	Description: "Памятник Петру I, также известный как Бронзовый всадник, - это памятник российскому императору Петру I, установленный в Санкт-Петербурге.",
	// 	Rating:      4.7,
	// 	Cost:        "$$",
	// 	ImageURL:    "https://img.tourister.ru/files/1/8/8/7/6/5/6/0/original.jpg",
	// }
	// data["5"] = model.Place{
	// 	ID:          "5",
	// 	Name:        "Статуя Свободы",
	// 	Description: "Статуя Свободы находится на острове Свободы в Нью-Йорке и является одним из символов Соединенных Штатов Америки.",
	// 	Rating:      4.6,
	// 	Cost:        "$$",
	// 	ImageURL:    "https://i.imgur.com/pOSbnHXh.jpg",
	// }
	// data["6"] = model.Place{
	// 	ID:          "6",
	// 	Name:        "Гренландия",
	// 	Description: "Гренландия - крупнейший остров в мире и административно-территориальное подразделение Королевства Дании.",
	// 	Rating:      4.2,
	// 	Cost:        "$$$",
	// 	ImageURL:    "https://10wallpaper.com/wallpaper/5120x2880/2103/Windows_10x_Microsoft_2021_Ocean_Glacier_5K_HD_Photo_5120x2880.jpg",
	// }
	// data["7"] = model.Place{
	// 	ID:          "7",
	// 	Name:        "Колизей",
	// 	Description: "Колизей - это амфитеатр в Риме, построенный в I веке н.э. и считающийся одним из величайших архитектурных и инженерных достижений древнего мира.",
	// 	Rating:      4.9,
	// 	Cost:        "$$",
	// 	ImageURL:    "https://sportishka.com/uploads/posts/2022-04/1650595488_15-sportishka-com-p-italiya-kolizei-krasivo-foto-15.jpg",
	// }
	// data["8"] = model.Place{
	// 	ID:          "8",
	// 	Name:        "Маяк Александрия",
	// 	Description: "Маяк Александрия был одним из семи чудес света и находился в древнем городе Александрия, в Египте.",
	// 	Rating:      4.4,
	// 	Cost:        "$$$",
	// 	ImageURL:    "https://polinka.top/uploads/posts/2023-06/thumbs/1685742571_polinka-top-p-kartinka-aleksandriya-yegipetskaya-instagr-44.jpg",
	// }
	// data["9"] = model.Place{
	// 	ID:          "9",
	// 	Name:        "Скайдайвинг в Нью Зеландии",
	// 	Description: "Нью Зеландия предлагает невероятные возможности для скайдайвинга с потрясающими видами на природную красоту страны.",
	// 	Rating:      4.8,
	// 	Cost:        "$$$",
	// 	ImageURL:    "https://www.dropzone.com/uploads/monthly_2020_01/Klavs.jpg.a67700687ec1204e6c6fedc92037bc8b.jpg",
	// }

	repo := &PlaceRepository{
		DB: db,
	}
	for _, place := range data {
		repo.AddPlace(&place)
	}
	return repo
}

func (r *PlaceRepository) AddPlace(place *model.Place) error {
	err := r.DB.QueryRow(
		`INSERT INTO place ("name", "description", "cost", "image_url")
        VALUES ($1, $2, $3, $4)`,
		place.Name,
		place.Description,
		place.Cost,
		place.ImageURL,
	).Scan()
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		//TODO: чужая ошибка, надо бы нормально обрабатывать
		return err
	}
	return nil
}

func (r *PlaceRepository) GetPlaces() (map[string]*model.Place, error) {
	places := make(map[string]*model.Place)
	rows, err := r.DB.Query("SELECT id, name, description, cost, image_url FROM place")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		place := &model.Place{}
		err = rows.Scan(&place.ID, &place.Name, &place.Description, &place.Cost, &place.ImageURL)
		if err != nil {
			return nil, err
		}
		places[place.ID] = place
	}
	return places, nil
}
