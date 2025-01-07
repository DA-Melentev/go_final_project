package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"os"
	"time"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			} else {
				WriteError(w, http.StatusUnauthorized, fmt.Errorf("authentification required"))
				return
			}
			valid := validateJWT(jwt)

			if !valid {
				WriteError(w, http.StatusUnauthorized, fmt.Errorf("authentification required"))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	log.Printf("/api/signin POST")
	var req map[string]string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err)
		log.Printf("error while unmarshal payload: %v", err)
		return
	}
	password := req["password"]

	if len(password) == 0 {
		err := errors.New("password field is required")
		WriteError(w, http.StatusUnauthorized, err)
		log.Printf("error while handling password: %v", err)
		return
	}

	correctPassword := os.Getenv("TODO_PASSWORD")

	if password != correctPassword {
		err := errors.New("wrong password")
		WriteError(w, http.StatusUnauthorized, err)
		log.Println(err)
		return
	}

	token, err := generateToken()
	if err != nil {
		err = fmt.Errorf("error while generating token %v", err)
		WriteError(w, http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	WriteResponseJSON(w, http.StatusAccepted, map[string]interface{}{
		"token": token,
	})
}

func generateToken() (string, error) {
	var secretPhrase = getSecretPhrase()

	claims := jwt.MapClaims{
		"exp": time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretPhrase)
}

func validateJWT(tokenString string) bool {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getSecretPhrase(), nil
	})

	if err != nil {
		return false
	}

	var expUnix int64
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expUnixF, ok := claims["exp"].(float64)
		if !ok {
			return false
		}
		expUnix = int64(expUnixF)
	}

	exp := time.Unix(expUnix, 0)
	if time.Now().After(exp) {
		return false
	}

	return true
}

func getSecretPhrase() []byte {
	pass := os.Getenv("TODO_PASSWORD")

	hasher := sha256.New()
	hasher.Write([]byte(pass))
	return hasher.Sum(nil)
}
