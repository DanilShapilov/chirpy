package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DanilShapilov/chirpy/internal/auth"
	"github.com/DanilShapilov/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, req *http.Request) {
	type reqData struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Couldn't find JWT",
			err,
		)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Couldn't validate JWT",
			err,
		)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := reqData{}
	err = decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "Email couldn't be empty", nil)
		return
	}
	if len(params.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "Password couldn't be empty", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
	}

	user, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user: %s", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
