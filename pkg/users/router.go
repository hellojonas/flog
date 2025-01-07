package users

import (
	"encoding/json"
	"io"
	"net/http"
)

type usrRouter struct {
	svc *usrService
}

type HttpMessageResponse struct {
	Message string `json:"message"`
}

func NewRouter(svc *usrService) *usrRouter {
	ur := &usrRouter{
		svc: svc,
	}

	return ur
}

func (ur *usrRouter) Route(mux *http.ServeMux) {
	mux.HandleFunc("POST /signup", ur.Signup)
}

func (ur *usrRouter) Signup(w http.ResponseWriter, r *http.Request) {
	var usrInput UserCreateInput
	body, err := io.ReadAll(r.Body)

	if err != nil {
		msg := HttpMessageResponse{
			Message: "error parsing body",
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	err = json.Unmarshal(body, &usrInput)

	if err != nil {
		msg := HttpMessageResponse{
			Message: "error parsing body",
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	err = ur.svc.CreateUser(usrInput)

	if err != nil {
		msg := HttpMessageResponse{
			Message: "error creating user. " + err.Error(), // TODO: improve this, nonsense error might appear
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	sendJson(w, http.StatusCreated, "")
}

func sendJson(w http.ResponseWriter, status int, body interface{}) error {
	w.Header().Add("Content-type", "application/json")
	data, err := json.Marshal(body)

	if err != nil {
		return err
	}

	w.WriteHeader(status)
	w.Write(data)
	return nil
}
