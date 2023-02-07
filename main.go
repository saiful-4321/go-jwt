package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var SECRET = []byte("super-secret-auth-key")
var API_KEY = "a1b2c3d4e5f6g7h8i9j10k11l12m13"

func main() {
	http.Handle("/", ValidateJWT(Home))
	http.HandleFunc("/generate-token", GetJWT)
	http.ListenAndServe(":3000", nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wellcome to home section")
}

// Generating token for authentication
func CreateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenStr, err := token.SignedString(SECRET)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return tokenStr, nil
}

// validating token from request
// acts like middlewire
func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)

				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Sorry! You are not authorized"))
				}
				return SECRET, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Not authorized:" + err.Error()))
			}

			if token.Valid {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Token not found"))
		}
	})
}

// Getting JWT token via request
func GetJWT(w http.ResponseWriter, r *http.Request) {
	if r.Header["Access-Key"] != nil {
		if r.Header["Access-Key"][0] == API_KEY {
			token, err := CreateJWT()
			if err != nil {
				return
			}
			fmt.Fprintf(w, token)
		}
	}
}
