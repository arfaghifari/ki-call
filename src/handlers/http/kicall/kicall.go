package kicall

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	usecase "github.com/arfaghifari/ki-call/src/usecase/kicall"
)

type Header struct {
	Error      string `json:"error_code"`
	StatusCode int    `json:"status_code"`
}

type KiCallRequest struct {
	Host    string
	Method  string
	Request map[string]interface{}
}

type MessageResponse struct {
	Header `json:"header"`
	Data   map[string]interface{} `json:"data"`
}

type SuccesMessage struct {
	Success bool `json:"success"`
}

type Handlers struct {
	usecase usecase.Usecase
}

func New(usecase usecase.Usecase) *Handlers {
	return &Handlers{usecase: usecase}
}

func GetHello(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello World")
	fmt.Fprintf(w, "HELLO Ki Ka Ku")
}

func (h *Handlers) GetListMethod(w http.ResponseWriter, r *http.Request) {
	var (
		statusCode = http.StatusBadRequest
		resp       MessageResponse
	)
	defer func() {
		w.Header().Set("Content-Type", "application/json")
		resp.StatusCode = statusCode
		responseWriter, err := json.Marshal(resp)
		if err != nil {
			log.Fatal("Failed build response")
		}
		w.Write(responseWriter)
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
		}
	}()

	res, err := h.usecase.GetListMethod()

	if err != nil {
		statusCode = http.StatusInternalServerError
		resp.Header.Error = err.Error()
		return
	}
	statusCode = http.StatusOK
	resp.Data = map[string]interface{}{
		"list_function": res,
	}

}

func (h *Handlers) GetRequestMethod(w http.ResponseWriter, r *http.Request) {
	var (
		statusCode = http.StatusBadRequest
		resp       MessageResponse
	)
	defer func() {
		w.Header().Set("Content-Type", "application/json")
		resp.StatusCode = statusCode
		responseWriter, err := json.Marshal(resp)
		if err != nil {
			log.Fatal("Failed build response")
		}
		w.Write(responseWriter)
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
		}
	}()

	methodName := r.URL.Query().Get("method")
	noEmptyStr := r.URL.Query().Get("no_empty")
	noEmpty, _ := strconv.ParseBool(noEmptyStr)

	res, err := h.usecase.GetRequestMethod(methodName, noEmpty)

	if err != nil {
		statusCode = http.StatusInternalServerError
		resp.Header.Error = err.Error()
		return
	}
	statusCode = http.StatusOK
	resp.Data = map[string]interface{}{
		"method":  methodName,
		"request": res,
	}
}

func (h *Handlers) KiCall(w http.ResponseWriter, r *http.Request) {
	var (
		statusCode = http.StatusBadRequest
		resp       MessageResponse
	)
	defer func() {
		w.Header().Set("Content-Type", "application/json")
		resp.StatusCode = statusCode
		responseWriter, err := json.Marshal(resp)
		if err != nil {
			log.Fatal("Failed build response")
		}
		w.Write(responseWriter)
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
		}
	}()

	req := KiCallRequest{}
	json.NewDecoder(r.Body).Decode(&req)
	res, err := h.usecase.KiCall(req.Host, req.Method, req.Request)

	if err != nil {
		statusCode = http.StatusInternalServerError
		resp.Header.Error = err.Error()
		return
	}
	statusCode = http.StatusOK
	resp.Data = map[string]interface{}{
		"response": res,
	}
}
