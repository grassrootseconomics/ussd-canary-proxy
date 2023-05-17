package api

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (api *API) AddToCacheHandler(w http.ResponseWriter, req bunrouter.Request) error {
	return nil
}

func (api *API) GetFromCacheHandler(w http.ResponseWriter, req bunrouter.Request) error {
	return nil
}
