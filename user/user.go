package user

import (
	"BeachAPI/api"
	"BeachAPI/app"
	"BeachAPI/token"
	"crypto/md5"
	"encoding/hex"
	"time"
)

func IsConfirmed(phoneNumber int) bool {
	var i bool
	err := app.DB.Get(&i, "select exists(select id from users where confirmed=true and phone_number=$1)", phoneNumber)
	api.CheckErrInfo(err, "IsConfirmed")
	return i
}

func Exists(phoneNumber int) bool {
	var i bool
	err := app.DB.Get(&i, "select exists(select id from users where phone_number=$1)", phoneNumber)
	api.CheckErrInfo(err, "Exists")
	return i
}

//func IsTokenValid(serverKey, nick, password string) bool {
//	var i bool
//	err := app.DB.Get(&i, "select exists(select id from accounts where serverkey=$1 and password=$2)", serverKey, password)
//	api.CheckErrInfo(err, "IsTokenValid")
//	return i
//}

func Confirmed(id int) {
	_, err := app.DB.Exec("update users set confirmed=true where id=$1", id)
	api.CheckErrInfo(err, "Confirmed")
}

func MD5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func (i User) CreateToken() string {
	var t, err = token.GenerateJWT(i.ID, i.PhoneNumber)
	api.CheckErrInfo(err, "GenerateJWT")
	_, err = app.DB.Exec(`INSERT INTO tokens (user_id, token, created) VALUES ($1, $2, $3)`, i.ID, t, time.Now())
	api.CheckErrInfo(err, "CreateToken")
	return t
}

func IsToken(id int, token string) bool {
	var i bool
	err := app.DB.Get(&i, "select exists(select id from tokens where id=$1 and token=$2)", id, token)
	api.CheckErrInfo(err, "IsToken")
	return i
}

func GetUserOnNum(phoneNumber int) User {
	var i User
	err := app.DB.Get(&i, `select * from users where phone_number=$1`, phoneNumber)
	api.CheckErrInfo(err, "GetUserOnNum")
	return i
}

type User struct {
	ID           int     `json:"id,omitempty"`
	PhoneNumber  string  `json:"phone_number" db:"phone_number"`
	Balance      float64 `json:"balance" db:"balance"`
	Confirmed    bool    `json:"confirmed"`
	RoleID       int     `json:"role_id" db:"role_id"`
	BalanceBonus float64 `json:"balance_bonus" db:"balance_bonus"`
}

func GetUser(id int) User {
	var i User
	err := app.DB.Get(&i, `select * from users where id=$1`, id)
	api.CheckErrInfo(err, "GetUser")
	return i
}

func NewUser(phoneNumber int) int {
	var id int
	err := app.DB.Get(&id, `INSERT INTO users (phone_number) VALUES ($1) returning id`, phoneNumber)
	api.CheckErrInfo(err, "NewUser")
	return id
}

func NewCode(id, code int) bool {
	_, err := app.DB.Exec(`INSERT INTO codes (user_id, code, created) VALUES ($1, $2, $3)`,
		id, code, time.Now())
	api.CheckErrInfo(err, "NewCode")
	return true
}

func CheckCode(id int, code string) bool {
	var i bool
	err := app.DB.Get(&i, "select exists(select id from codes where user_id=$1 and code=$2)", id, code)
	api.CheckErrInfo(err, "CheckCode")
	return i
}

func GetBalance(id int) float64 {
	row := app.DB.QueryRow("SELECT balance FROM users WHERE id = $1", id)
	var i float64
	err := row.Scan(&i)
	if err != nil {
		return 0
	}
	return i
}

func AddBalance(id int, count float64) bool {
	res, err := app.DB.Exec("update users set balance=balance+$1 where id=$2", count, id)
	api.CheckErrInfo(err, "addBalance")
	result, _ := res.RowsAffected()
	if result == 0 {
		return false
	}
	return true
}
