package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func WithJWT(handlerFunc http.HandlerFunc, s Store) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tokenStr := getTokenFromRequest(r)
		token, err := validateJWT(tokenStr)

		if err != nil {
			permissionDenied(w, "failed to authenticate token")
			return
		}

		if !token.Valid {
			permissionDenied(w, "failed to authenticate token")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(string)

		_, err = s.GetUserByID(userID)
		if err != nil {
			permissionDenied(w, "failed to get user")
			return
		}

		handlerFunc(w, r)
	}
}

func permissionDenied(w http.ResponseWriter, m string) {
	log.Println(m)
	WriteJSON(w, http.StatusUnauthorized, ErrorResponse{
		Error: fmt.Errorf("permission denied").Error(),
	})
}

func getTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token != "" {
		return token
	}

	token = r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	return ""
}

func validateJWT(tokenStr string) (*jwt.Token, error) {
	secret := Envs.JWTSecret

	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func CreateJWT(secret []byte, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiredAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
