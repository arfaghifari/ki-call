package src

import (
	"net/http"
	"time"

	myClient "github.com/arfaghifari/ki-call/src/client"
	kcHandlers "github.com/arfaghifari/ki-call/src/handlers/http/kicall"
	"github.com/arfaghifari/ki-call/src/middleware"
	server "github.com/arfaghifari/ki-call/src/server"
	kcUsecase "github.com/arfaghifari/ki-call/src/usecase/kicall"
	"github.com/gorilla/mux"
)

func Main() {

	// Init serve HTTP
	router := mux.NewRouter()
	kicallHandlers := kcHandlers.New(kcUsecase.NewUsecase())
	// routes http
	router.HandleFunc("/hello", kcHandlers.GetHello).Methods(http.MethodGet)
	router.HandleFunc("/ls-svc", kicallHandlers.GetListService).Methods(http.MethodGet)
	router.HandleFunc("/ls-func", kicallHandlers.GetListMethod).Methods(http.MethodGet)
	router.HandleFunc("/requests", kicallHandlers.GetRequestMethod).Methods(http.MethodGet)

	router.HandleFunc("/ki-call", kicallHandlers.KiCall).Methods(http.MethodPost)

	routerMw := middleware.MiddlewarePanic(router)

	myClient.ClientKitex.RegisterAllClient("")

	serverConfig := server.Config{
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  5 * time.Second,
		Port:         9700,
	}
	server.Serve(serverConfig, &routerMw)
}
