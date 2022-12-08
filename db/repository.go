package db

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"YandexEdaParser/structs"
)

const (
	insertRestaurant         = `INSERT INTO YandexEdaParser.restaurants (yandexId, created_at, updated_at, name, slag, rating, minimalOrderPrice, internalrating) VALUES ($1,now(), now(),$2, $3, $4, $5, $6) RETURNING id`
	selectRestById           = `SELECT * FROM YandexEdaParser.restaurants WHERE id = $1`
	selectRestByYandexId     = `SELECT * FROM YandexEdaParser.restaurants WHERE yandexid = $1`
	updateRestaurant         = `UPDATE YandexEdaParser.restaurants SET name = $1, slag = $2, rating = $3, minimalOrderPrice = $4, internalrating = $5, updated_at = now() WHERE id = $6`
	selectInternalRatingById = `SELECT restaurants.internalRating FROM YandexEdaParser.restaurants WHERE  restaurants.id = $1`
	selectAllRestaurant      = `SELECT *  FROM YandexEdaParser.restaurants ORDER BY restaurants.InternalRating desc;`
	SelectRestaurantPrices   = `SELECT menuitem.price FROM YandexEdaParser.menuitem WHERE restaurantid = $1;`

	insertItem           = `INSERT INTO YandexEdaParser.menuitem (yandexId, restaurantid,  created_at, updated_at, name, description, price, value, internalRating) VALUES ($1, $2, now(), now(), $3, $4, $5, $6, $7)`
	selectItemById       = `SELECT * FROM YandexEdaParser.menuitem WHERE id = $1`
	selectItemByYandexId = `SELECT * FROM YandexEdaParser.menuitem WHERE yandexid = $1`
	updateItem           = `UPDATE YandexEdaParser.menuitem SET name = $1, description = $2, price = $3, value = $4, internalRating = $5, restaurantid = $6, updated_at = now() WHERE yandexId = $7`
)

type Repository struct {
	DB *sql.DB
}

func (r *Repository) Init(db structs.DataBase) (err error) {
	pgSql := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable binary_parameters=yes",
		db.Host, db.Port, db.Username, db.Password, db.Database)

	//создаем конект к используемым базам
	if r.DB, err = sql.Open("postgres", pgSql); err != nil {
		return errors.Wrap(err, "error create conn to pgSqlPartner")
	}
	return nil
}

// AddOrUpdateRestaurant функция ищет ресторан в базе, если находить проверяет изменились ли параметры - обновляет,  если нет в базе  - добавляет. На выход ID ресторана в базе
func (r *Repository) AddOrUpdateRestaurant(restaurant structs.Restaurant) (Id int, err error) {
	var restaurantInBase structs.Restaurant
	if restaurantInBase, err = r.GetRestaurantByYandexId(restaurant.YandexId); err != nil {
		if err == sql.ErrNoRows {
			if Id, err = r.AddRestaurant(restaurant); err != nil {
				return Id, err
			}
		} else {
			return 0, err
		}
	}
	if !r.CheckEqualRestaurant(restaurantInBase, restaurant) {
		if r.UpdateRestaurant(restaurant, restaurantInBase.ID); err != nil {
			return 0, err
		}
	}
	return restaurantInBase.ID, nil
}

// GetRestaurantById получаем ресторан по его внутреннему id
func (r *Repository) GetRestaurantById(id int) (rest structs.Restaurant, err error) {
	if err := r.DB.QueryRow(selectRestById, id).Scan(&rest.ID, &rest.YandexId, &rest.CreateAt, &rest.UpdateAt, &rest.Name, &rest.Slug, &rest.Rating, &rest.MinimalOrderPrice, &rest.InternalRating); err != nil {
		if err != sql.ErrNoRows {
			return rest, err
		}
	}
	return rest, nil
}

// GetRestaurantByYandexId получаем ресторан из базы по его яндекс id
func (r *Repository) GetRestaurantByYandexId(id int) (rest structs.Restaurant, err error) {
	if err = r.DB.QueryRow(selectRestByYandexId, id).Scan(&rest.ID, &rest.YandexId, &rest.CreateAt, &rest.UpdateAt, &rest.Name, &rest.Slug, &rest.Rating, &rest.MinimalOrderPrice, &rest.InternalRating); err != nil {
		return rest, err
	}
	return rest, nil
}

// AddRestaurant добавляем ресторан в базу
func (r *Repository) AddRestaurant(rest structs.Restaurant) (ID int, err error) {
	if err := r.DB.QueryRow(insertRestaurant, rest.YandexId, rest.Name, rest.Slug, rest.Rating, rest.MinimalOrderPrice, &rest.InternalRating).Scan(&ID); err != nil {
		return 0, err
	}
	return ID, nil
}

// GetRestaurantInternalRating получаем текущий внутрений рейтинг
func (r *Repository) GetRestaurantInternalRating(id int) (internalRating float64, err error) {
	if err = r.DB.QueryRow(selectInternalRatingById, id).Scan(&internalRating); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return internalRating, nil
}

// UpdateRestaurant обновляем данные ресторана в базе
func (r *Repository) UpdateRestaurant(rest structs.Restaurant, ID int) (err error) {
	if _, err := r.DB.Exec(updateRestaurant, rest.Name, rest.Slug, rest.Rating, rest.MinimalOrderPrice, &rest.InternalRating, ID); err != nil {
		return err
	}
	return nil
}

// GetAllRestaurants получаем все рестораны по внутреннему рейтингу
func (r *Repository) GetAllRestaurants() (Restaurants []structs.Restaurant, err error) {
	var allItems *sql.Rows
	if allItems, err = r.DB.Query(selectAllRestaurant); err != nil {
		return nil, err
	}
	for allItems.Next() {
		var rest structs.Restaurant
		if err := allItems.Scan(&rest.ID, &rest.YandexId, &rest.CreateAt, &rest.UpdateAt, &rest.Name, &rest.Slug, &rest.Rating, &rest.MinimalOrderPrice, &rest.InternalRating); err != nil {
			return nil, err
		}
		Restaurants = append(Restaurants, rest)
	}
	return Restaurants, nil
}

// GetRestaurantPrices получаю цены на филадельфии в конкретном ресторане
func (r *Repository) GetRestaurantPrices(Id int) (Prices []int, err error) {
	var allItems *sql.Rows
	if allItems, err = r.DB.Query(SelectRestaurantPrices, Id); err != nil {
		return nil, err
	}
	for allItems.Next() {
		var priese int
		if err := allItems.Scan(&priese); err != nil {
			return nil, err
		}
		Prices = append(Prices, priese)
	}
	return Prices, nil
}

// CheckEqualRestaurant сравниваем параметра ресторанов
func (r *Repository) CheckEqualRestaurant(rest1 structs.Restaurant, rest2 structs.Restaurant) bool {
	if rest1.Slug == rest2.Slug && rest1.Name == rest2.Name && rest1.Rating == rest2.Rating && rest1.InternalRating == rest2.InternalRating {
		return true
	} else {
		return false
	}
}

/// ПОЗИЦИИ МЕНЮ

func (r *Repository) AddOrUpdateItem(item structs.Items) (err error) {
	var itemInBase structs.Items
	if itemInBase, err = r.GetItemsByYandexId(item.YandexId); err != nil {
		if err == sql.ErrNoRows {
			if err = r.AddItems(item); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if !r.CheckEqualItems(itemInBase, item) {
		if r.UpdateItems(item, itemInBase.YandexId); err != nil {
			return err
		}
	}
	return nil
}

// AddItems добавляем позицию в базу
func (r *Repository) AddItems(item structs.Items) (err error) {
	if _, err := r.DB.Exec(insertItem, item.YandexId, item.RestaurantId, item.Name, item.Description, item.Price, item.Measure.ValueInt, item.InternalRating); err != nil {
		return err
	}
	return nil
}

// UpdateItems обновляем данные в базе
func (r *Repository) UpdateItems(item structs.Items, yandexID int) (err error) {
	if _, err := r.DB.Exec(updateItem, item.Name, item.Description, item.Price, item.Measure.ValueInt, item.InternalRating, item.RestaurantId, yandexID); err != nil {
		return err
	}
	return nil
}

// GetItemsByYandexId получаем позицию из базы по его яндекс id
func (r *Repository) GetItemsByYandexId(id int) (item structs.Items, err error) {
	if err = r.DB.QueryRow(selectItemByYandexId, id).Scan(&item.ID, &item.YandexId, &item.RestaurantId, &item.CreateAt, &item.UpdateAt, &item.Name, &item.Description, &item.Price, &item.Measure.ValueInt, &item.InternalRating); err != nil {
		return item, err
	}
	return item, nil
}

// GetItemsById  получаем позицию по его внутреннему id
func (r *Repository) GetItemsById(id int) (item structs.Items, err error) {
	if err := r.DB.QueryRow(selectItemById, id).Scan(&item.ID, &item.YandexId, &item.RestaurantId, &item.CreateAt, &item.UpdateAt, &item.Name, &item.Description, &item.Price, &item.Measure.ValueInt, &item.InternalRating); err != nil {
		if err != sql.ErrNoRows {
			return item, err
		}
	}
	return item, err
}

// CheckEqualItems сравниваем параметра позицию
func (r *Repository) CheckEqualItems(item1 structs.Items, item2 structs.Items) bool {
	if item1 == item2 {
		return true
	} else {
		return false
	}
}
