package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (api *API) USSDProxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buffer, err := io.ReadAll(r.Body)
		_ = r.Body.Close()

		formValues, err := url.ParseQuery(string(buffer))
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		api.logg.Debug("proxy handler: form values", "form", formValues)

		version, err := api.cache.Get(r.Context(), formValues.Get("phoneNumber"))
		if err != nil {
			api.logg.Error("proxy handler: cache get error", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		api.logg.Debug("proxy handler: phone details", "phone", formValues.Get("phoneNumber"), "version", version)

		r.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))

		switch version {
		case 1:
			api.proxy.V1.ServeHTTP(w, r)
		case 2:
			api.proxy.V2.ServeHTTP(w, r)
		default:
			api.proxy.V1.ServeHTTP(w, r)
		}
	}
}
