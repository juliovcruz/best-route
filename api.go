package main

import (
	"best-route/database"
	"best-route/dijkstra"
	"best-route/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type server struct {
	database.DatabaseClient
	dijkstra.DijkstraClient
}

func RunAPI(db *database.Database, djk *dijkstra.Dijkstra) error {
	startHandles(&server{
		DatabaseClient: db.Client,
		DijkstraClient: djk.Client,
	})
	return nil
}

func startHandles(s *server) {
	http.HandleFunc("/add", s.handleInsert)
	http.HandleFunc("/best", s.handleBestRoute)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func (s *server) handleBestRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		fmt.Fprintf(w, "%+v", errors.New("wrong method: try again with GET"))
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	var route *models.Route

	if err := json.Unmarshal(reqBody, &route); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%+v", err)
		return
	}

	// TODO: VALIDATE REQUEST

	routes, err := s.DatabaseClient.GetAllRoutes()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%+v", err)
		return
	}

	res, err := s.DijkstraClient.BestRoute(route.Start, route.Target, routes)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%+v", err)
		return
		// TODO: not found in graph
	}

	data, _ := json.Marshal(res)
	fmt.Fprintf(w, "%+v", string(data))
}

func (s *server) handleInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		fmt.Fprintf(w, "%+v", errors.New("wrong method: try again with POST"))
		return
	}

	// TODO: VALIDATE REQUEST

	reqBody, _ := ioutil.ReadAll(r.Body)

	var route *models.Route

	if err := json.Unmarshal(reqBody, &route); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%+v", err)
		return
	}

	res, err := s.DatabaseClient.InsertOneRoute(route)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%+v", err)
		return
	}

	data, _ := json.Marshal(res)
	fmt.Fprintf(w, "%+v", string(data))
}
