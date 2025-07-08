package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps from db", err)
		return
	}

	var res = make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		jsonKeysChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		}
		res[i] = jsonKeysChirp
	}

	respondWithJSON(w, http.StatusOK, res)
}
