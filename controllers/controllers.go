package controllers

import (
	"BeachAPI/api"
	"BeachAPI/app"
	"BeachAPI/card"
	"BeachAPI/shop"
	"BeachAPI/sms"
	"BeachAPI/user"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Result interface{} `json:"result"`
	Error  *Error      `json:"error"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type addBalance struct {
	Value float64 `json:"value"`
}

func AddBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var i addBalance
	err := json.NewDecoder(r.Body).Decode(&i)

	idStr := r.Header["Id"][0]

	idUser, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "ID error",
				Code:    10,
			},
		})
		return
	}

	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: err.Error(),
				Code:    0,
			},
		})
		return
	}

	user.AddBalance(idUser, i.Value)

	json.NewEncoder(w).Encode(&Response{
		Result: map[string]string{"status": "OK"},
	})
}

type buy struct {
	BuyItems []card.BuyItem `json:"buyItems"`
}

func Buy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var i buy
	err := json.NewDecoder(r.Body).Decode(&i)
	api.Print(i)
	params := mux.Vars(r)
	if params["shop_id"] == "" {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "shop_id отсутствует",
				Code:    20,
			},
		})
		return
	}

	idStr := r.Header["Id"][0]

	idUser, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "ID error",
				Code:    10,
			},
		})
		return
	}

	shopID := params["shop_id"]
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: err.Error(),
				Code:    0,
			},
		})
		return
	}

	hashCard := api.RandString(5)
	var price float64
	for _, item := range i.BuyItems {
		price = price + item.Price
	}

	if price > user.GetBalance(idUser) {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "Баланса не хватает.",
				Code:    25,
			},
		})
		return
	}

	for _, item := range i.BuyItems {
		card.NewBuy(item, idUser, shopID, hashCard)
	}
	user.AddBalance(idUser, -price)

	json.NewEncoder(w).Encode(&Response{
		Result: map[string]string{"status": "OK", "hash_card": hashCard},
	})
}

type code struct {
	NumberPhone string `json:"number_phone"`
}

func Code(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var i code
	err := json.NewDecoder(r.Body).Decode(&i)

	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: err.Error(),
				Code:    0,
			},
		})
		return
	}

	if i.NumberPhone == "" {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "Номер отсутствует.",
				Code:    7,
			},
		})
		return
	}

	if len(i.NumberPhone) != 11 {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "Number no valid",
				Code:    1,
			},
		})
		return
	}
	num, err := strconv.Atoi(i.NumberPhone)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "Number no valid",
				Code:    2,
			},
		})
		return
	}

	var id int
	if !user.Exists(num) {
		id = user.NewUser(num)
	} else {
		id = user.GetUserOnNum(num).ID
	}

	res := sms.SendSMS(i.NumberPhone)

	if res == nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "Неизвестная ошибка получения смс.",
				Code:    3,
			},
		})
		return
	}
	if res.Status == "ERROR" {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: res.StatusText,
				Code:    4,
			},
		})
		return
	}
	if res.Status != "OK" {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "Сервер не вернул статус OK.",
				Code:    5,
			},
		})
		return
	}

	user.NewCode(id, res.Code)

	json.NewEncoder(w).Encode(&Response{
		Result: map[string]string{"status": "OK", "user_id": api.ToString(id)},
	})
}

type checkCode struct {
	UserID int `json:"user_id"`
	Code   string
}

func CheckCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var i checkCode
	err := json.NewDecoder(r.Body).Decode(&i)

	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: err.Error(),
				Code:    0,
			},
		})
		return
	}

	if i.UserID == 0 {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "UserID отсутствует.",
				Code:    7,
			},
		})
		return
	}

	if user.CheckCode(i.UserID, i.Code) {
		user.Confirmed(i.UserID)
		i := user.GetUser(i.UserID)
		token := i.CreateToken()
		json.NewEncoder(w).Encode(&Response{
			Result: map[string]string{"token": token, "status": "OK"},
		})
		return
	} else {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "Код неверный!",
				Code:    8,
			},
		})
		return
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	idStr := r.Header["Id"][0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "ID error",
				Code:    10,
			},
		})
		return
	}

	json.NewEncoder(w).Encode(&Response{
		Result: user.GetUser(id),
	})
}

func GetCardsUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	idStr := r.Header["Id"][0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "ID error",
				Code:    10,
			},
		})
		return
	}

	json.NewEncoder(w).Encode(&Response{
		Result: card.GetCards(id),
	})
}

func Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	//fmt.Println(r.Header)
	//idStr := r.Header["Id"][0]
	//
	//id, err := strconv.Atoi(idStr)
	//if err != nil {
	//	json.NewEncoder(w).Encode(&Response{
	//		Error: &Error{
	//			Message: "ID error",
	//			Code:    10,
	//		},
	//	})
	//	return
	//}

	rows, err := app.DB.Queryx("SELECT * FROM users LIMIT 20;")
	api.CheckErrInfo(err, "Users")

	var users []user.User

	for rows.Next() {
		var user user.User
		err = rows.StructScan(&user)
		users = append(users, user)
		api.CheckErrInfo(err, "Buys")
	}
	_ = rows.Close()

	json.NewEncoder(w).Encode(&Response{
		Result: users,
	})
}

func Shops(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	rows, err := app.DB.Queryx("SELECT * FROM shops")
	api.CheckErrInfo(err, "Users")

	var ii []shop.Shop

	for rows.Next() {
		var i shop.Shop
		err = rows.StructScan(&i)
		ii = append(ii, i)
		api.CheckErrInfo(err, "Shops")
	}
	_ = rows.Close()

	json.NewEncoder(w).Encode(&Response{
		Result: ii,
	})
}

func ShopID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := mux.Vars(r)
	if params["shop_id"] == "" {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "shop_id отсутствует",
				Code:    20,
			},
		})
		return
	}
	shopID := params["shop_id"]
	var i bool
	err := app.DB.Get(&i, "select exists(select id from shops where id=$1)", shopID)
	api.CheckErrInfo(err, "CheckCode")

	if !i {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: fmt.Sprintf("Магазин %s отсутсвует в системе", shopID),
				Code:    20,
			},
		})
		return
	}

	var s shop.Shop
	err = app.DB.Get(&s, "SELECT * FROM shops where id=$1", shopID)
	api.CheckErrInfo(err, "ShopID")

	json.NewEncoder(w).Encode(&Response{
		Result: s,
	})
}

func Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	idStr := r.Header["Id"][0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "ID error",
				Code:    10,
			},
		})
		return
	}
	params := mux.Vars(r)
	if params["shop_id"] == "" {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "shop_id отсутствует",
				Code:    20,
			},
		})
		return
	}
	shopID := params["shop_id"]
	if params["hash_card"] == "" {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "hash_card отсутствует",
				Code:    20,
			},
		})
		return
	}
	hashCard := params["hash_card"]
	var i bool
	err = app.DB.Get(&i, "select exists(select hash_card from buys where hash_card=$1)", hashCard)
	api.CheckErrInfo(err, "CheckHash")

	if !i {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: fmt.Sprintf("Магазин %s отсутсвует в системе", hashCard),
				Code:    20,
			},
		})
		return
	}

	//items

	s := card.CheckCard(id, shopID, hashCard)

	json.NewEncoder(w).Encode(&Response{
		Result: s,
	})
}

func ShopData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := mux.Vars(r)
	if params["shop_id"] == "" {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "shop_id отсутствует",
				Code:    20,
			},
		})
		return
	}
	if !shop.Exist(params["shop_id"]) {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: fmt.Sprintf("Магазин %s отсутсвует в системе", params["shop_id"]),
				Code:    20,
			},
		})
		return
	}

	idStr := r.Header["Id"][0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "ID error",
				Code:    10,
			},
		})
		return
	}

	i := shop.GetData(params["shop_id"])
	i.User = user.GetUser(id)

	json.NewEncoder(w).Encode(&Response{
		Result: i,
	})
}

func Items(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	params := mux.Vars(r)
	if params["shop_id"] == "" {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "shop_id отсутствует",
				Code:    20,
			},
		})
		return
	}
	shopID := params["shop_id"]

	if !shop.Exist(shopID) {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: fmt.Sprintf("Магазин %s отсутсвует в системе", shopID),
				Code:    20,
			},
		})
		return
	}

	json.NewEncoder(w).Encode(&Response{
		Result: shop.GetItems(shopID),
	})
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	//  Ensure our file does not exceed 5MB
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)

	file, handler, err := r.FormFile("image")

	// Capture any errors that may arise
	if err != nil {
		fmt.Fprintf(w, "Error getting the file")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded file name: %+v\n", handler.Filename)
	fmt.Printf("Uploaded file size %+v\n", handler.Size)
	fmt.Printf("File mime type %+v\n", handler.Header)

	// Get the file content type and access the file extension
	fileType := strings.Split(handler.Header.Get("Content-Type"), "/")[1]

	// Create the temporary file name
	fileName := fmt.Sprintf("upload-*.%s", fileType)
	// Create a temporary file with a dir folder
	tempFile, err := ioutil.TempFile("images", fileName)

	if err != nil {
		fmt.Println(err)
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully uploaded file")
}
