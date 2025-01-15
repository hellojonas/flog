package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/hellojonas/flog/pkg/services"
)

type appRouter struct {
	appSvc *services.AppService
	logSvc *services.LogService
}

func NewAppRouter(appSvc *services.AppService, logSvc *services.LogService) *appRouter {
	ar := &appRouter{
		appSvc: appSvc,
		logSvc: logSvc,
	}

	return ar
}

func (ar *appRouter) Route(mux *http.ServeMux) {
	mux.HandleFunc("POST /apps", ar.CreateApp)
	mux.HandleFunc("GET /apps/{id}", ar.RetrieveById)
	mux.HandleFunc("POST /apps/{id}/members", ar.SetMembers)
	mux.HandleFunc("GET /apps/{id}/members", ar.ListAppMembers)
	mux.HandleFunc("GET /apps/{id}/logs", ar.ListAppLogs)
	mux.HandleFunc("GET /apps/{id}/logs/{name}", ar.DowloadLog)
}

func (ar *appRouter) DowloadLog(w http.ResponseWriter, r *http.Request) {
	aid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		sendJson(w, http.StatusBadRequest, HttpMessageResponse{
			Message: "invalid appId id",
		})
	}

	logName := r.PathValue("name")
	app, err := ar.appSvc.FindById(aid)

	if err != nil {
		sendJson(w, http.StatusBadRequest, HttpMessageResponse{
			Message: "error retrieving app info",
		})
	}

	logPath := ar.logSvc.ResolveLogFile(app.Name, logName)

	if err != nil {
		sendJson(w, http.StatusBadRequest, HttpMessageResponse{
			Message: "error retrieving app logs",
		})
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(app.Name+"_"+logName))
	w.Header().Set("Content-Type", "application/octect-stream")
	http.ServeFile(w, r, logPath)
}

func (ur *appRouter) CreateApp(w http.ResponseWriter, r *http.Request) {
	var appInput services.AppCreateInput
	body, err := io.ReadAll(r.Body)

	if err != nil {
		msg := HttpMessageResponse{
			Message: "error parsing body",
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	err = json.Unmarshal(body, &appInput)

	if err != nil {
		msg := HttpMessageResponse{
			Message: "error parsing body",
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	app, err := ur.appSvc.CreateApp(appInput)

	if err != nil {
		msg := HttpMessageResponse{
			Message: "error creating app. " + err.Error(), // TODO: improve this, nonsense error might appear
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	sendJson(w, http.StatusCreated, app)
}

func (ur *appRouter) RetrieveById(w http.ResponseWriter, r *http.Request) {
	aid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		sendJson(w, http.StatusBadRequest, HttpMessageResponse{
			Message: "invalid appId id",
		})
	}

	user, err := ur.appSvc.FindById(aid)

	// TODO: improve error to tell if user was not found
	if err != nil {
		msg := HttpMessageResponse{
			Message: "error retrieving app. " + err.Error(), // TODO: improve this, nonsense error might appear
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	sendJson(w, http.StatusOK, user)
}

func (ur *appRouter) SetMembers(w http.ResponseWriter, r *http.Request) {
	aid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		sendJson(w, http.StatusBadRequest, HttpMessageResponse{
			Message: "invalid appId id",
		})
	}

	var appMemberInput services.AppMemberInput
	body, err := io.ReadAll(r.Body)

	if err != nil {
		msg := HttpMessageResponse{
			Message: "error parsing body",
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	err = json.Unmarshal(body, &appMemberInput)

	if err != nil {
		msg := HttpMessageResponse{
			Message: "error parsing body",
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	err = ur.appSvc.SetMembers(aid, appMemberInput.Members)

	if err != nil {
		msg := HttpMessageResponse{
			Message: "error setting app members. " + err.Error(), // TODO: improve this, nonsense error might appear
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	sendResponse(w, http.StatusOK)
}

func (ur *appRouter) ListAppMembers(w http.ResponseWriter, r *http.Request) {
	aid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		sendJson(w, http.StatusBadRequest, HttpMessageResponse{
			Message: "invalid appId id",
		})
	}

	apps, err := ur.appSvc.ListAppMembers(aid)

	// TODO: improve error to tell if user was not found
	if err != nil {
		msg := HttpMessageResponse{
			Message: "error retrieving app members. " + err.Error(), // TODO: improve this, nonsense error might appear
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	sendJson(w, http.StatusOK, apps)
}

func (ur *appRouter) ListAppLogs(w http.ResponseWriter, r *http.Request) {
	aid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		sendJson(w, http.StatusBadRequest, HttpMessageResponse{
			Message: "invalid appId id",
		})
	}

	logs, err := ur.logSvc.ListAppLogs(aid)

	// TODO: improve error to tell if user was not found
	if err != nil {
		msg := HttpMessageResponse{
			Message: "error retrieving app logs. " + err.Error(), // TODO: improve this, nonsense error might appear
		}
		sendJson(w, http.StatusBadRequest, msg)
		return
	}

	sendJson(w, http.StatusOK, logs)
}
