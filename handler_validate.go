package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type reqData struct {
	Body string `json:"body"`
}

type returnVal struct {
	Valid bool `json:"valid"`
}

func handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := reqData{}
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVal{
		Valid: true,
	})
}
