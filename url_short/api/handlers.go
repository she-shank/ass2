package api

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/its-kos/assignment1/types"
	"github.com/its-kos/assignment1/urls"
)

const defaultTtl = 86400 // 24 hours in seconds

// Retrieve the original URL associated with the given short identifier
func (s *ApiServer) handleGetURLByID(w http.ResponseWriter, r *http.Request) {
	s.requestCount++
	urlID := chi.URLParam(r, "id")
	url, err := s.db.GetURL(urlID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := types.GetURLResponse{
		ID: url.URL,
	}

	s.succesful++
	w.WriteHeader(http.StatusMovedPermanently)
	json.NewEncoder(w).Encode(response)
}

// Update the URL associated with the given short identifier
func (s *ApiServer) handleUpdateURLByID(w http.ResponseWriter, r *http.Request) {
	s.requestCount++
	var req types.UpdateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": %s"}`, err.Error())))
		return
	}

	urlID := chi.URLParam(r, "id")
	newUrl, err := s.db.GetURL(urlID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		//w.Write([]byte(`{"error": "id not found"}`))
		return
	}

	if !urls.ValidateURL(req.URL) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "not a valid URL"}`))
		return
	}

	newUrl.URL = req.URL

	err = s.db.UpdateURL(newUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		w.Write([]byte(fmt.Sprintf(`{"error": %s"}`, err.Error())))
		return
	}

	s.succesful++
	w.WriteHeader(http.StatusOK)
}

// Delete the mapping for the given short identifier
func (s *ApiServer) handleDeleteURLByID(w http.ResponseWriter, r *http.Request) {
	s.requestCount++
	urlID := chi.URLParam(r, "id")
	err := s.db.DeleteURL(urlID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		//w.Write([]byte(`{"error": "internal server error"}`))
		return
	}
	s.succesful++
	w.WriteHeader(http.StatusNoContent)
}

// Retrieve a list of all short identifiers (:id) stored in the service.
func (s *ApiServer) handleGetAllIDs(w http.ResponseWriter, r *http.Request) {
	s.requestCount++
	urls, err := s.db.GetAllURLs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		w.Write([]byte(fmt.Sprintf(`{"error": %s"}`, err.Error())))
		return
	}

	ids := make([]string, len(urls))
	for i, url := range urls {
		ids[i] = url.ID
	}

	response := types.GetAllURLSResponse{
		IDs: ids,
	}

	s.succesful++
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Create a new short identifier (:id) for a given URL.
func (s *ApiServer) handleCreateURLAlias(w http.ResponseWriter, r *http.Request) {
	s.requestCount++
	var req types.CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": %s"}`, err.Error())))
		return
	}

	if !urls.ValidateURL(req.URL) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": not a valid URL}`))
		return
	}

	// if ttl is specified the use if otherwise use default value
	var ttl int
	if req.Ttl == nil {
		slog.Info(fmt.Sprint("did not recieve time to live, using default value : ", defaultTtl, " seconds"))
		ttl = defaultTtl
	} else {
		slog.Info(fmt.Sprint("time to live recieved : ", &req.Ttl, " seconds"))
		ttl = *req.Ttl

	}

	//id, err := urls.Shorten(urls.StripDomain(req.URL))
	id, err := urls.Shorten(req.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": %s"}`, err.Error())))
		return
	}

	response := types.CreateURLResponse{
		ID: id,
	}

	found, _ := s.db.GetURL(id)

	if found != nil {
		s.succesful++
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)

	} else {
		url := types.URL{
			URL:        req.URL,
			ID:         id,
			CreatedAt:  time.Now(),
			Hits:       0,
			TimeToLive: ttl,
		}

		err = s.db.CreateURL(url)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(`{"error": %s"}`, err.Error())))
			return
		}

		s.succesful++
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// 404 always
func (s *ApiServer) handleDeleteAllURLs(w http.ResponseWriter, r *http.Request) {
	s.requestCount++
	err := s.db.DeleteAllURLs()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": %s"}`, err.Error())))
		return
	}
	s.succesful++
	w.WriteHeader(http.StatusNotFound)
}

func (s *ApiServer) handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	s.requestCount++
	urls, err := s.db.GetAllURLs()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": %s"}`, err.Error())))
		return
	}

	s.succesful++
	response := types.GetMetricsResponse{
		TotalURLs: fmt.Sprintf("%d", len(urls)),
		TotalRequests: fmt.Sprintf("%d", s.requestCount),
		RequestRate: fmt.Sprintf("%f", float64(s.requestCount)/time.Since(s.startTime).Seconds()),
		SuccessfulRequests: fmt.Sprintf("%d", s.succesful),
		SuccessRate: fmt.Sprintf("%f", float64(s.succesful)/float64(s.requestCount)),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *ApiServer) handleGetMetricsByID(w http.ResponseWriter, r *http.Request) {
	s.requestCount++
	urlID := chi.URLParam(r, "id")
	url, err := s.db.GetURL(urlID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := types.GetMetricsByIDResponse{
		ID:  url.ID,
		URL: url.URL,
		Hits: url.Hits,
	}

	s.succesful++
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
