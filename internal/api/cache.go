package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/uptrace/bunrouter"
)

func (api *API) AddToCacheHandler(w http.ResponseWriter, req bunrouter.Request) error {
	var addCacheReq struct {
		PhoneNumber string `json:"phoneNumber"`
		Version     int    `json:"version"`
	}

	if err := json.NewDecoder(req.Body).Decode(&addCacheReq); err != nil {
		return err
	}

	if err := api.cache.Set(req.Context(), addCacheReq.PhoneNumber, addCacheReq.Version); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return ErrDuplicate
			}
		}
	}

	return nil
}

func (api *API) GetFromCacheHandler(w http.ResponseWriter, req bunrouter.Request) error {
	params := req.Params()
	phoneNumber := params.ByName("phoneNumber")

	version, err := api.cache.Get(req.Context(), phoneNumber)
	if err != nil {
		return err
	}

	if version < 1 {
		w.WriteHeader(http.StatusNotFound)
		_ = bunrouter.JSON(w, bunrouter.H{
			"message": "User not registered on canary",
		})
		return nil
	}

	w.WriteHeader(http.StatusOK)
	_ = bunrouter.JSON(w, bunrouter.H{
		"version": version,
	})
	return nil
}

func (api *API) UpdateCacheHandler(w http.ResponseWriter, req bunrouter.Request) error {
	var updateCacheReq struct {
		PhoneNumber string `json:"phoneNumber"`
		Version     int    `json:"version"`
	}

	if err := json.NewDecoder(req.Body).Decode(&updateCacheReq); err != nil {
		return err
	}

	if err := api.cache.Update(req.Context(), updateCacheReq.PhoneNumber, updateCacheReq.Version); err != nil {
		return err
	}

	return nil
}
