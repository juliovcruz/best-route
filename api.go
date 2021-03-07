package main

import (
	"best-route/database"
	"best-route/models"
	"best-route/router"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type server struct {
	database.DatabaseClient
	router.RouterClient
	Addr string
}

func RunAPI(db *database.Database, djk *router.Router, addr string) error {
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", "wrong method, try again with POST"))
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", "body cannot be empty"))
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", fmt.Sprintf("%+v", err.Error())))
		return
	}

	var route *models.Route

	if err := json.Unmarshal(reqBody, &route); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", fmt.Sprintf("%+v", err.Error())))
		return
	}

	route = route.Trim()

	if resErr := validateAPIInsertRequest(route); resErr != nil {
		w.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprintf(w, "%+v", string(resErr.ToJSON()))
		return
	}

	res, err := s.DatabaseClient.InsertOneRoute(route)
	if err != nil {
		if err.Error() == "already exists" {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", fmt.Sprintf("%+v", err.Error())))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", fmt.Sprintf("%+v", err.Error())))
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", fmt.Sprintf("%+v", err.Error())))
		return
	}
	fmt.Fprintf(w, "%+v", string(data))
}

func (s *server) handleBestRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", "wrong method, try again with GET"))
		return
	}

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", "body cannot be empty"))
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", fmt.Sprintf("%+v", err.Error())))
		return
	}

	var route *models.Route

	if err := json.Unmarshal(reqBody, &route); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", fmt.Sprintf("%+v", err.Error())))
		return
	}

	route = route.Trim()

	if resErr := validateAPIBestRouteRequest(route); resErr != nil {
		w.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprintf(w, "%+v", string(resErr.ToJSON()))
		return
	}

	routes, err := s.DatabaseClient.GetAllRoutes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", fmt.Sprintf("%+v", err.Error())))
		return
	}

	res, resErr := s.RouterClient.BestRoute(route.Start, route.Target, routes)
	if resErr != nil {
		w.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprintf(w, "%+v", string(resErr.ToJSON()))
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v", models.NewResponseErrorToJSON("", fmt.Sprintf("%+v", err.Error())))
		return
	}
	fmt.Fprintf(w, "%+v", string(data))
}
