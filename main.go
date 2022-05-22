package main

import (
	"BeachAPI/api"
	"BeachAPI/app"
	"BeachAPI/controllers"
	"BeachAPI/token"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {
	_cors := cors.New(cors.Options{
		//AllowedOrigins: []string{"http://yooga.ru:3006",
		//	"http://localhost:3006",
		//	"https://localhost:3006",
		//	"http://maykop-concert.ru",
		//	"https://maykop-concert.ru",
		//	"http://localhost:3000",
		//	"http://bilet.yoogline.ru:3000",
		//	"https://bilet.yoogline.ru:3000",
		//	"http://nalmes.online:3006",
		//	"http://kassa01.ru:3006",
		//	"https://kassa01.ru:3006",
		//	"https://nalmes.online:3006",
		//	"http://nalmes.online",
		//	"https://nalmes.online",
		//	"http://ticket.maykop-concert.ru:3006",
		//	"http://ticket.maykop-concert.ru:3006",
		//	"http://46.175.120.206:3006",
		//	"https://46.175.120.206:3006",
		//	"https://ticket.kassa01.ru:3006",
		//	"http://ticket.kassa01.ru:3006",
		//	"http://ticket.kassa01.ru",
		//	"https://ticket.kassa01.ru",
		//},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	router := mux.NewRouter()

	//{"number_phone": "79384846350"}

	//image
	router.HandleFunc("/uploadImage", controllers.UploadFile).Methods("POST")
	router.PathPrefix("/images").Handler(http.StripPrefix("/images", http.FileServer(http.Dir("images/")))).Methods("GET")

	router.HandleFunc("/auth/code", controllers.Code).Methods("POST")
	//{"user_id": 1, "code": "1989"}
	router.HandleFunc("/auth/checkCode", controllers.CheckCode).Methods("POST")

	router.HandleFunc("/user/get", token.IsAuthorized(controllers.GetUser)).Methods("GET")
	//router.HandleFunc("/users", controllers.Users).Methods("GET")
	router.HandleFunc("/items/{shop_id}", controllers.Items).Methods("GET")

	router.HandleFunc("/shop/{shop_id}", token.IsAuthorized(controllers.ShopData)).Methods("GET")

	router.HandleFunc("/shops/get", token.IsAuthorized(controllers.Shops)).Methods("GET")
	router.HandleFunc("/shops/{shop_id}/buy", token.IsAuthorized(controllers.Buy)).Methods("POST")

	//router.HandleFunc("/shops/{shop_id}/cards", token.IsAuthorized(controllers.Cards)).Methods("POST")

	//{"value": float}
	router.HandleFunc("/user/addBalance", token.IsAuthorized(controllers.AddBalance)).Methods("POST")

	router.HandleFunc("/shop/get/{shop_id}", token.IsAuthorized(controllers.ShopID)).Methods("GET")

	//router.HandleFunc("/shop/{shop_id}/cards", controllers.GetCardsShop).Methods("GET")
	//shops/1/check?hash_card=aYVAG
	router.HandleFunc("/shop/{shop_id}/check/{hash_card}", token.IsAuthorized(controllers.Check)).Methods("GET")

	router.HandleFunc("/user/cards", token.IsAuthorized(controllers.GetCardsUser)).Methods("GET")

	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NTMwODA0MDAsImlkIjoxLCJudW1iZXIiOiI3OTk5NDI3MDQ0NiJ9.j5GOlNoZWBxEuxg7T99Xku3Agz8HTxGWP6ZkNjEHesI

	cert := "/etc/letsencrypt/live/alnezis.riznex.ru/fullchain.pem"
	key := "/etc/letsencrypt/live/alnezis.riznex.ru/privkey.pem"
	if _, err := os.Stat(cert); err != nil {
		if os.IsNotExist(err) {
			log.Println("no ssl")
			handler := _cors.Handler(router)
			err := http.ListenAndServe(fmt.Sprintf(":%d", app.CFG.Port), handler)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
	log.Println("yes ssl")
	handler := _cors.Handler(router)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", app.CFG.Port), cert, key, handler)
	if err != nil {
		api.CheckErrInfo(err, "ListenAndServeTLS")
		//	return
	}
}
