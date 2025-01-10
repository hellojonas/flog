package users

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/hellojonas/flog/pkg/apps"
)

type userRouter struct {
	svc    *userService
	appSvc *apps.AppService
}

type HttpMessageResponse struct {
	Message string `json:"message"`
}

func NewRouter(svc *userService, appSvc *apps.AppService) *userRouter {
	ur := &userRouter{
		svc:    svc,
		appSvc: appSvc,
	}

	return ur
}

func (ur *userRouter) Route(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/{id}", ur.RetrieveById)
	mux.HandleFunc("POST /users", ur.Signup)
	mux.HandleFunc("GET /users/{id}/apps", ur.ListUserApps)
}

func (ur *userRouter) Signup(w http.ResponseWriter, r *http.Request) {
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

func (ur *userRouter) RetrieveById(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		sendJson(w, http.StatusBadRequest, HttpMessageResponse{
			Message: "invalid user id",
		})
	}

	user, err := ur.svc.FindById(uid)
	user.Password = ""

	// TODO: improve error to tell if user was not found
	if err != nil {
		msg := HttpMessageResponse{
			Message: "error retrieving user. " + err.Error(), // TODO: improve this, nonsense error might appear
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	sendJson(w, http.StatusOK, user)
}

func (ur *userRouter) ListUserApps(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		sendJson(w, http.StatusBadRequest, HttpMessageResponse{
			Message: "invalid user id",
		})
	}

	apps, err := ur.appSvc.ListUserApps(uid)

	// TODO: improve error to tell if user was not found
	if err != nil {
		msg := HttpMessageResponse{
			Message: "error retrieving user apps. " + err.Error(), // TODO: improve this, nonsense error might appear
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	sendJson(w, http.StatusOK, apps)
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
