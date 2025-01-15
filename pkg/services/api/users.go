package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/hellojonas/flog/pkg/services"
)

type userRouter struct {
	svc    *services.UserService
	appSvc *services.AppService
}

func NewUserRouter(svc *services.UserService, appSvc *services.AppService) *userRouter {
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
	var usrInput services.UserCreateInput
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
