package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "../service"
	"github.com/julienschmidt/httprouter"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/observer"
)

type MerchantController struct {
	MerchantService *MerchantService
}

//GetAll Get All merchant Observer.
func (m *MerchantController) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	merchants, err := m.MerchantService.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, merchants)
}

//GetAllOb Get All merchant Observer.
func (m *MerchantController) GetAllOb(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	onNext := handlers.NextFunc(func(merchants interface{}) {
		//fmt.Printf("Processing: %v\n", merchants)
		respondWithJSON(w, http.StatusOK, merchants)
	})
	onError := handlers.ErrFunc(func(err error) {
		fmt.Printf("Encountered error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "ERROR")
	})
	watcher := observer.New(onNext, onError)
	m.MerchantService.GetAllObserve().Subscribe(watcher)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
