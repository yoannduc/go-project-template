package example

import (
	"bytes"
	"encoding/json/v2"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/yoannduc/go-project-template/internal/dtos"
	"github.com/yoannduc/go-project-template/internal/ports"
	"github.com/yoannduc/go-project-template/pkg/logger"
)

const (
	idPathValue = "id"
	prefixAPI   = "/api/v1/example"
)

var (
	routeGetAll  = prefixAPI
	routeGetByID = prefixAPI + "/{" + idPathValue + "}"
	routeCreate  = prefixAPI
	routeUpdate  = prefixAPI + "/{" + idPathValue + "}"
	routeDelete  = prefixAPI + "/{" + idPathValue + "}"
)

// exampleHandler is the http handler for Examples.
// It is designed to be used with std lib http.ServeMux.
type exampleHandler struct {
	srv   ports.ExampleService
	loggr *slog.Logger
}

// New is the function that returns the handler that will rely on
// inputed ports.ExampleService & slog.Logger.
func New(srv ports.ExampleService, loggr *slog.Logger) ports.Handler {
	return exampleHandler{
		srv:   srv,
		loggr: loggr,
	}
}

// LoadRoutes is the function that will load routes to the
// std lib http.ServeMux.
func (hdl exampleHandler) LoadRoutes(mux *http.ServeMux) {
	mux.Handle(http.MethodGet+" "+routeGetAll, http.HandlerFunc(hdl.getAllWithSearch))
	mux.Handle(http.MethodGet+" "+routeGetByID, http.HandlerFunc(hdl.getByID))
	mux.Handle(http.MethodPost+" "+routeCreate, http.HandlerFunc(hdl.post))
	mux.Handle(http.MethodPut+" "+routeUpdate, http.HandlerFunc(hdl.update))
	mux.Handle(http.MethodDelete+" "+routeDelete, http.HandlerFunc(hdl.delete))
}

// getAllWithSearch is the handler to find a list of Example.
// If any query string "search" is passed, the list will be filtered
// using this search input, else all elements are returned in the list,
// json encoded. It uses ports.ExampleService GetByLabelContaining or
// GetAll depending on if search param was used.
func (hdl exampleHandler) getAllWithSearch(w http.ResponseWriter, r *http.Request) {
	var out []dtos.Example
	var err error
	if search := r.URL.Query().Get("search"); search != "" {
		out, err = hdl.srv.GetByLabelContaining(r.Context(), search)
	} else {
		out, err = hdl.srv.GetAll(r.Context())
	}
	if err != nil {
		hdl.loggr.LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	buf := bytes.Buffer{}
	if err = json.MarshalWrite(&buf, out); err != nil {
		hdl.loggr.LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

// getByID is the handler to find a single Example based on its id
// found in uri path json encoded if found.
// It uses ports.ExampleService getByID.
func (hdl exampleHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue(idPathValue))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, http.StatusText(http.StatusBadRequest))
		return
	}

	out, err := hdl.srv.GetByID(r.Context(), id)
	if err != nil {
		hdl.loggr.LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	buf := bytes.Buffer{}
	if err = json.MarshalWrite(&buf, out); err != nil {
		hdl.loggr.LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

// post is the handler to create a new Example.
// It returns newly created Example json encoded if no error.
// It uses ports.ExampleService Create.
func (hdl exampleHandler) post(w http.ResponseWriter, r *http.Request) {
	var dto dtos.Example
	if err := json.UnmarshalRead(r.Body, &dto); err != nil || dto.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, http.StatusText(http.StatusBadRequest))
		return
	}

	out, err := hdl.srv.Create(r.Context(), dto)
	if err != nil {
		hdl.loggr.LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	buf := bytes.Buffer{}
	if err = json.MarshalWrite(&buf, out); err != nil {
		hdl.loggr.LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(buf.Bytes())
}

// update is the handler to update an Example based on its id
// found in uri path. It returns updated Example json encoded
// if no error. It uses ports.ExampleService update.
func (hdl exampleHandler) update(w http.ResponseWriter, r *http.Request) {
	var dto dtos.Example
	if err := json.UnmarshalRead(r.Body, &dto); err != nil || dto.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, http.StatusText(http.StatusBadRequest))
		return
	}
	id, err := strconv.Atoi(r.PathValue(idPathValue))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, http.StatusText(http.StatusBadRequest))
		return
	}

	out, err := hdl.srv.Update(r.Context(), id, dto)
	if err != nil {
		hdl.loggr.LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	buf := bytes.Buffer{}
	if err = json.MarshalWrite(&buf, out); err != nil {
		hdl.loggr.LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(buf.Bytes())
}

// delete is the handler to update an Example based on its id
// found in uri path. It returns deleted Example json encoded
// if no error. It uses ports.ExampleService delete.
func (hdl exampleHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue(idPathValue))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, http.StatusText(http.StatusBadRequest))
		return
	}

	out, err := hdl.srv.Delete(r.Context(), id)
	if err != nil {
		hdl.loggr.LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	buf := bytes.Buffer{}
	if err = json.MarshalWrite(&buf, out); err != nil {
		hdl.loggr.LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
