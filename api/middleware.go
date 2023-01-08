package api

import (
	"a21hc3NpZ25tZW50/model"
	"context"
	"encoding/json"
	"net/http"
)

func (api *API) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// sessionToken := "" // TODO: replace this
		var respErr model.ErrorResponse
		// memperoleh session token dari requests cookies
		c, err := r.Cookie("session_token")
		if err != nil {
			w.WriteHeader(401)
			respErr.Error = "http: named cookie not present"

			// encode json
			jsonInBytes, err := json.Marshal(respErr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(jsonInBytes)
			return
		}
		sessionToken := c.Value

		sessionFound, err := api.sessionsRepo.CheckExpireToken(sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "username", sessionFound.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) Get(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: answer here
		var respErr model.ErrorResponse
		if r.Method != "GET" {
			w.WriteHeader(405)
			respErr.Error = "Method is not allowed!"
			// w.Write([]byte("Method is not allowed"))

			// w.Header().Set("Content-Type", "application/json")

			// encode json
			jsonInBytes, err := json.Marshal(respErr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(jsonInBytes)

			return
		}
		next.ServeHTTP(w, r)
	})
}

func (api *API) Post(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: answer here

		var respErr model.ErrorResponse
		if r.Method != "POST" {

			respErr.Error = "Method is not allowed!"

			// marshal json
			jsonInBytes, err := json.Marshal(respErr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//write response
			w.WriteHeader(405)
			w.Write(jsonInBytes)
			return
		}

		// if err := r.ParseMultipartForm(1024); err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		next.ServeHTTP(w, r)
	})
}
