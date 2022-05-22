package token

import (
	"BeachAPI/api"
	"BeachAPI/app"
	"BeachAPI/request"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

func GenerateJWT(id int, number string) (string, error) {
	var mySigningKey = []byte(app.CFG.SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = api.ToString(id)
	claims["number"] = number
	claims["exp"] = time.Now().Add(time.Hour * 3000).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Authorization"] == nil {
			json.NewEncoder(w).Encode(request.Response{Error: &request.Error{
				Message: "No Token Found",
				Code:    17,
			}})
			return
		}

		var mySigningKey = []byte(app.CFG.SecretKey)

		s := strings.Split(r.Header["Authorization"][0], " ")

		token, err := jwt.Parse(s[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			//todo delete token
			json.NewEncoder(w).Encode(request.Response{Error: &request.Error{
				Message: "Your Token has been expired",
				Code:    16,
			}})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			r.Header.Set("Id", claims["id"].(string))
			r.Header.Set("Number", claims["number"].(string))
			handler.ServeHTTP(w, r)
			return
		}

		json.NewEncoder(w).Encode(request.Response{Error: &request.Error{
			Message: "Not Authorized",
			Code:    15,
		}})
	}
}
