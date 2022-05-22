package shop

import (
	"BeachAPI/api"
	"BeachAPI/app"
	"BeachAPI/user"
	"time"
)

type Shop struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	UserID      int       `json:"user_id" db:"user_id"`
	Status      int       `json:"status"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Balance     float64   `json:"balance"`
	TimeUpdated time.Time `json:"time_updated" db:"time_updated"`
	Image       string    `json:"image"`
}

type Item struct {
	ID          int       `json:"id"`
	ShopID      int       `json:"shop_id" db:"shop_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Count       int       `json:"count"`
	Price       float64   `json:"price"`
	Created     time.Time `json:"created"`
	Status      int       `json:"status"`
	TimeUpdated time.Time `json:"time_updated" db:"time_updated"`
	Image       string    `json:"image"`
}

func Exist(shopID string) bool {
	var i bool
	err := app.DB.Get(&i, "select exists(select id from shops where id=$1)", shopID)
	api.CheckErrInfo(err, "Exist")
	return i
}

type DataShop struct {
	Items []Item    `json:"items"`
	Shop  Shop      `json:"shop"`
	User  user.User `json:"user"`
}

func GetData(shopID string) DataShop {
	var shop Shop
	err := app.DB.Get(&shop, "SELECT * FROM shops where id=$1", shopID)
	api.CheckErrInfo(err, "ShopID")
	items := GetItems(shopID)
	return DataShop{Items: items, Shop: shop}
}

func GetItem(itemID int) Item {
	var i Item
	err := app.DB.Get(&i, "SELECT * FROM items where id=$1", itemID)
	api.CheckErrInfo(err, "GetItem")
	return i
}

func GetItems(shopID string) []Item {
	rows, err := app.DB.Queryx("SELECT * FROM items where shop_id=$1", shopID)
	api.CheckErrInfo(err, "Users")

	i := []Item{}

	for rows.Next() {
		var item Item
		err = rows.StructScan(&item)
		api.CheckErrInfo(err, "Items")
		i = append(i, item)
	}
	_ = rows.Close()
	return i
}
