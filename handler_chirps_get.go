package main

import (
	"net/http"

	"github.com/DanilShapilov/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, req *http.Request) {
	authorId := req.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error
	var authorUUID uuid.UUID
	if authorId != "" {
		authorUUID, err = uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Incorrect format of author_id", err)
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthor(req.Context(), authorUUID)
	} else {
		chirps, err = cfg.db.GetChirps(req.Context())
	}
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

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, req *http.Request) {
	chirpIDString := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirps", err)
		return
	}

	jsonKeysChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, jsonKeysChirp)
}
