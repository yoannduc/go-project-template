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

type exampleHandler struct {
	srv ports.ExampleService
}

func New(srv ports.ExampleService) exampleHandler {
	return exampleHandler{
		srv: srv,
	}
}

func (hdl exampleHandler) LoadRoutes(mux *http.ServeMux) {
	mux.Handle(http.MethodGet+" "+routeGetAll, http.HandlerFunc(hdl.GetAllWithSearch))
	mux.Handle(http.MethodGet+" "+routeGetByID, http.HandlerFunc(hdl.GetByID))
	mux.Handle(http.MethodPost+" "+routeCreate, http.HandlerFunc(hdl.Create))
	mux.Handle(http.MethodPut+" "+routeUpdate, http.HandlerFunc(hdl.Update))
	mux.Handle(http.MethodDelete+" "+routeDelete, http.HandlerFunc(hdl.Delete))
}

func (hdl exampleHandler) GetAllWithSearch(w http.ResponseWriter, r *http.Request) {
	var out []dtos.Example
	var err error
	if search := r.URL.Query().Get("search"); search != "" {
		out, err = hdl.srv.GetByLabelContaining(r.Context(), search)
	} else {
		out, err = hdl.srv.GetAll(r.Context())
	}
	if err != nil {
		logger.Get().LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	buf := bytes.Buffer{}
	if err = json.MarshalWrite(&buf, out); err != nil {
		logger.Get().LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (hdl exampleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue(idPathValue))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, http.StatusText(http.StatusBadRequest))
		return
	}

	out, err := hdl.srv.GetByID(r.Context(), id)
	if err != nil {
		logger.Get().LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	buf := bytes.Buffer{}
	if err = json.MarshalWrite(&buf, out); err != nil {
		logger.Get().LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (hdl exampleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto dtos.Example
	if err := json.UnmarshalRead(r.Body, &dto); err != nil || dto.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, http.StatusText(http.StatusBadRequest))
		return
	}

	out, err := hdl.srv.Post(r.Context(), dto)
	if err != nil {
		logger.Get().LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	buf := bytes.Buffer{}
	if err = json.MarshalWrite(&buf, out); err != nil {
		logger.Get().LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

func (hdl exampleHandler) Update(w http.ResponseWriter, r *http.Request) {
	var dto dtos.Example
	if err := json.UnmarshalRead(r.Body, &dto); err != nil || dto.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, http.StatusText(http.StatusBadRequest))
		return
	}
	id, err := strconv.Atoi(r.PathValue(idPathValue))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, http.StatusText(http.StatusBadRequest))
		return
	}

	out, err := hdl.srv.Patch(r.Context(), id, dto)
	if err != nil {
		logger.Get().LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	buf := bytes.Buffer{}
	if err = json.MarshalWrite(&buf, out); err != nil {
		logger.Get().LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

func (hdl exampleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue(idPathValue))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, http.StatusText(http.StatusBadRequest))
		return
	}

	out, err := hdl.srv.Delete(r.Context(), id)
	if err != nil {
		logger.Get().LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}

	buf := bytes.Buffer{}
	if err = json.MarshalWrite(&buf, out); err != nil {
		logger.Get().LogAttrs(r.Context(), logger.LevelError, "", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
