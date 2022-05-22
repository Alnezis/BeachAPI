package card

import (
	"BeachAPI/api"
	"BeachAPI/app"
	"BeachAPI/shop"
)

type BuyItem struct {
	Id       int     `json:"id"`
	BuyCount int     `json:"buyCount"`
	Price    float64 `json:"price"`
}

func NewBuy(i BuyItem, userID int, shopID string, hashCard string) int {
	var id int
	_, err := app.DB.Exec(`INSERT INTO buys (shop_id, item_id, user_id, price, count, hash_card)
VALUES ($1,$2,$3,$4,$5,$6)`, shopID, i.Id, userID, i.Price, i.BuyCount, hashCard)
	api.CheckErrInfo(err, "NewBuy")
	return id
}

type card struct {
	HashCard   string      `json:"hash_card" db:"hash_card"`
	UserID     string      `json:"user_id" db:"user_id"`
	CountItems int         `json:"count_items" db:"count_items"`
	SumPrice   float64     `json:"sum_price" db:"sum_price"`
	Items      []shop.Item `json:"items,omitempty" db:"item,omitempty"`
}

func GetCards(idUser int) []card {
	//    SELECT hash_card, count(*) as count, sum(price) as price FROM buys GROUP BY hash_card;
	//    SELECT hash_card, user_id, count(*) as count, sum(price) as price FROM buys GROUP BY hash_card, user_id having user_id=4;
	rows, err := app.DB.Queryx(" SELECT hash_card, user_id, count(*) as count_items, sum(price) as sum_price FROM buys GROUP BY hash_card, user_id having user_id=$1;", idUser)
	api.CheckErrInfo(err, "Users")

	i := []card{}

	for rows.Next() {
		var item card
		err = rows.StructScan(&item)
		i = append(i, item)
		api.CheckErrInfo(err, "Items")
	}
	_ = rows.Close()
	return i
}

func CheckCard(idUser int, idShop string, hashCard string) []shop.Item {
	//    SELECT hash_card, count(*) as count, sum(price) as price FROM buys GROUP BY hash_card;
	//    SELECT hash_card, user_id, count(*) as count, sum(price) as price FROM buys GROUP BY hash_card, user_id having user_id=4;
	rows, err := app.DB.Queryx("SELECT * from buys where user_id=$1 and hash_card=$2 and shop_id=$3;", idUser, hashCard, idShop)
	api.CheckErrInfo(err, "GetCard")

	items := []shop.Item{}
	for rows.Next() {
		var i shop.Item
		err = rows.StructScan(&i)
		items = append(items, i)
		api.CheckErrInfo(err, "Items 2")
	}

	_ = rows.Close()
	return items
}
