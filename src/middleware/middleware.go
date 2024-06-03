package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func MiddlewarePanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := parseError(recover()); err != nil {
				text := fmt.Sprintf("Panic: %s", err.Error())
				log.Println(text)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func parseError(msg interface{}) error {
	if msg != nil {
		switch err := msg.(type) {
		case string:
			return errors.New(err)
		case error:
			return err
		default:
			return errors.New("unknown error")
		}
	}
	return nil
}
