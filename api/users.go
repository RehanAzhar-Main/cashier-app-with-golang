package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"net/http"
	"path"
	"text/template"
	"time"

	"github.com/google/uuid"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	// Read username and password request with FormValue.
	creds := model.Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	} // TODO: replace this

	// Handle request if creds is empty send response code 400, and message "Username or Password empty"
	// TODO: answer here

	// tidak memberikan request body
	var respErr model.ErrorResponse

	if len(creds.Username) <= 0 && len(creds.Username) <= 0 {

		respErr.Error = "Username or Password empty"

		jsonInBytes, err := json.Marshal(respErr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(400)
		w.Write(jsonInBytes)
		return
	}

	// username dan password kosong
	if creds.Username == "" && creds.Password == "" {

		respErr.Error = "Username or Password empty"

		jsonInBytes, err := json.Marshal(respErr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(400)
		w.Write(jsonInBytes)
		return
	}

	// semua Oke

	err := api.usersRepo.AddUser(creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	filepath := path.Join("views", "status.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	var data = map[string]string{"name": creds.Username, "message": "register success!"}

	w.WriteHeader(200)
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	// Read usernmae and password request with FormValue.
	creds := model.Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	} // TODO: replace this

	// Handle request if creds is empty send response code 400, and message "Username or Password empty"
	// TODO: answer here

	// tidak memberikan request body
	var respErr model.ErrorResponse

	if len(creds.Username) <= 0 && len(creds.Username) <= 0 {

		respErr.Error = "Username or Password empty"

		jsonInBytes, err := json.Marshal(respErr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(400)
		w.Write(jsonInBytes)
		return
	}

	// username dan password kosong
	if creds.Username == "" && creds.Password == "" {

		respErr.Error = "Username or Password empty"

		jsonInBytes, err := json.Marshal(respErr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(400)
		w.Write(jsonInBytes)
		return
	}

	//  check keberadaan username dan password
	err := api.usersRepo.LoginValid(creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized) // 401
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Wrong User or Password!"})
		return
	} else {
		// Membuat random session token menggunakan "github.com/google/uuid" library untuk generate UUIDs
		sessionToken := uuid.NewString()
		expiresAt := time.Now().Add(5 * time.Hour) // 5 jam

		// set cookie sisi client
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Path:    "./data/sessions.json",
			Expires: expiresAt, // expired 5 jam
		})

		// set cookie dalam struct
		session := model.Session{
			Token:    sessionToken,
			Username: creds.Username,
			Expiry:   expiresAt,
		}

		// menambahkan session
		err = api.sessionsRepo.AddSessions(session)
		if err != nil {
			return
		}

		w.WriteHeader(200)
	}

	// Generate Cookie with Name "session_token", Path "/", Value "uuid generated with github.com/google/uuid", Expires time to 5 Hour.
	// TODO: answer here

	api.dashboardView(w, r)
}

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	// //Read session_token and get Value:
	// sessionToken := "" // TODO: replace this

	// find cookie session token
	var respErr model.ErrorResponse
	c, err := r.Cookie("session_token")
	if err != nil {
		// jika session tidak ada
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

	w.WriteHeader(200)

	// get session token
	sessionToken := c.Value

	api.sessionsRepo.DeleteSessions(sessionToken)

	//Set Cookie name session_token value to empty and set expires time to Now:
	// TODO: answer here
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})

	filepath := path.Join("views", "login.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}
