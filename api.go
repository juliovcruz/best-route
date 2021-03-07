package main

import (
	"best-route/database"
	"best-route/models"
	route_calculator "best-route/route_calculator"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type server struct {
	DatabaseClient database.Client
	RouterClient   route_calculator.Client
	Addr           string
}

func RunAPI(db *database.Database, djk *route_calculator.Router, addr string) error {
	fmt.Printf("start api in localhost%v\n", addr)

	if err := startHandles(&server{
		DatabaseClient: db.Client,
		RouterClient:   djk.Client,
		Addr:           addr,
	}); err != nil {
		return err
	}
	return nil
}

func startHandles(s *server) error {
	http.HandleFunc("/add", s.handleInsert)
	http.HandleFunc("/best", s.handleBestRoute)
	return http.ListenAndServe(s.Addr, nil)
}

func (s *server) handleInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, &models.ResponseError{Message: "wrong method, try again with POST"})
		return
	}

	if r.Body == nil {
		httpError(w, http.StatusBadRequest, &models.ResponseError{Message: "body cannot be empty"})
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpError(w, http.StatusInternalServerError, &models.ResponseError{Message: err.Error()})
		return
	}

	var route *models.Route

	if err := json.Unmarshal(reqBody, &route); err != nil {
		httpError(w, http.StatusInternalServerError, &models.ResponseError{Message: err.Error()})
		return
	}

	route = route.Trim()

	if resErr := validateAPIInsertRequest(route); resErr != nil {
		httpError(w, http.StatusPreconditionFailed, resErr)
		return
	}

	res, err := s.DatabaseClient.InsertOneRoute(route)
	if err != nil {
		if err.Error() == database.AlreadyExistsMessageError {
			httpError(w, http.StatusPreconditionFailed, &models.ResponseError{Message: err.Error()})
			return
		}

		httpError(w, http.StatusInternalServerError, &models.ResponseError{Message: err.Error()})
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		httpError(w, http.StatusInternalServerError, &models.ResponseError{Message: err.Error()})
		return
	}
	fmt.Fprintf(w, "%+v", string(data))
}

func (s *server) handleBestRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpError(w, http.StatusMethodNotAllowed, &models.ResponseError{Message: "wrong method, try again with GET"})
		return
	}

	if r.Body == nil {
		httpError(w, http.StatusBadRequest, &models.ResponseError{Message: "body cannot be empty"})
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpError(w, http.StatusInternalServerError, &models.ResponseError{Message: err.Error()})
		return
	}

	var route *models.Route

	if err := json.Unmarshal(reqBody, &route); err != nil {
		httpError(w, http.StatusInternalServerError, &models.ResponseError{Message: err.Error()})
		return
	}

	route = route.Trim()

	if resErr := validateAPIBestRouteRequest(route); resErr != nil {
		httpError(w, http.StatusPreconditionFailed, resErr)
		return
	}

	routes, err := s.DatabaseClient.GetAllRoutes()
	if err != nil {
		httpError(w, http.StatusInternalServerError, &models.ResponseError{Message: err.Error()})
		return
	}

	res, resErr := s.RouterClient.BestRoute(route.Start, route.Target, routes)
	if resErr != nil {
		httpError(w, http.StatusPreconditionFailed, resErr)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		httpError(w, http.StatusInternalServerError, &models.ResponseError{Message: err.Error()})
		return
	}
	fmt.Fprintf(w, "%+v", string(data))
}

func httpError(w http.ResponseWriter, statusCode int, resErr *models.ResponseError) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%+v", string(resErr.ToJSON()))
}
