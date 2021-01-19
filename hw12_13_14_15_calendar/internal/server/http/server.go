package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/app"
	srv "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/server"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	log    *zap.Logger
}

var _ srv.Stopper = (*Server)(nil)

type MyHandler struct {
	app *app.App
	ctx context.Context
	log *zap.Logger
}

func NewHTTP(ctx context.Context, app *app.App, l *zap.Logger, port string) (*Server, error) {
	_, mux := newHandler(ctx, app, l)
	srv, err := newServer(mux, port, l)
	if err != nil {
		return nil, fmt.Errorf("server creation error: %w", err)
	}
	return srv, nil
}
func newServer(mux *http.ServeMux, port string, log *zap.Logger) (*Server, error) { //nolint
	log.Info("http is running...")
	server := &http.Server{
		Addr:           port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	srv := &Server{
		server: server,
		log:    log,
	}
	return srv, nil
}
func newHandler(ctx context.Context, app *app.App, l *zap.Logger) (*MyHandler, *http.ServeMux) {
	handler := &MyHandler{ctx: ctx, app: app, log: l}
	mux := http.NewServeMux()

	mux.HandleFunc("/set", handler.requestLoggerMiddleware(headerSetter(handler.setEvent)))
	mux.HandleFunc("/get", handler.requestLoggerMiddleware(headerSetter(handler.getEvent)))
	mux.HandleFunc("/delete", handler.requestLoggerMiddleware(headerSetter(handler.deleteEvent)))
	mux.HandleFunc("/update", handler.requestLoggerMiddleware(headerSetter(handler.updateEvent)))
	return handler, mux
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("http server shutdown")
	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}
	return nil
}

func (m *MyHandler) setEvent(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		m.log.Error("BadRequest", zap.Error(err))
		resp.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(resp).Encode(err.Error()); err != nil {
			m.log.Error("Encode error", zap.Error(err))
		}
		return
	}
	rb := Event{}
	if err = json.Unmarshal(body, &rb); err != nil {
		m.log.Error("Unmarshal error", zap.Error(err))
		resp.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(resp).Encode(err.Error()); err != nil {
			m.log.Error("Encode error", zap.Error(err))
		}
		return
	}
	eve := parseToEventStorageStruct(rb)
	var id int64
	id, err = m.app.Create(m.ctx, eve)
	if err != nil {
		m.log.Info("BadRequest", zap.String("error", err.Error()))
		resp.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(resp).Encode(err.Error()); err != nil {
			m.log.Error("Encode error", zap.Error(err))
		}
		return
	}
	resp.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(resp).Encode(id); err != nil {
		m.log.Error("Encode error", zap.Error(err))
	}
}

func (m *MyHandler) getEvent(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		m.log.Error("BadRequest", zap.Error(err))
		resp.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(resp).Encode(err.Error()); err != nil {
			m.log.Error("Encode error", zap.Error(err))
		}
		return
	}

	r := JSONID{}
	id := r.ID
	if err = json.Unmarshal(body, &id); err != nil {
		m.log.Error("Unmarshal error", zap.Error(err))
		resp.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(resp).Encode(err.Error()); err != nil {
			m.log.Error("Encode error", zap.Error(err))
		}
		return
	}
	var result storage.Event
	result, err = m.app.Get(m.ctx, id)
	if err != nil {
		m.log.Info("BadRequest", zap.String("error", err.Error()))
		if err = json.NewEncoder(resp).Encode(err.Error()); err != nil {
			m.log.Error("Encode error", zap.Error(err))
		}
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	response := Event{
		ID:          result.ID,
		Owner:       result.Owner,
		Title:       result.Title,
		Description: result.Description,
		StartDate:   result.StartDate,
		StartTime:   result.StartTime,
		EndDate:     result.EndDate,
		EndTime:     result.EndTime,
	}
	resp.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(resp).Encode(&response); err != nil {
		m.log.Error("Encode error", zap.Error(err))
	}
}

func (m *MyHandler) deleteEvent(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		m.log.Error("BadRequest", zap.Error(err))
		resp.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(resp).Encode(err.Error()); err != nil {
			m.log.Error("Encode error", zap.Error(err))
		}
		return
	}
	r := JSONID{}
	id := r.ID
	if err = json.Unmarshal(body, &id); err != nil {
		m.log.Error("Unmarshal error", zap.Error(err))
		resp.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(resp).Encode(err.Error()); err != nil {
			m.log.Error("Encode error", zap.Error(err))
		}
		return
	}
	err = m.app.Delete(m.ctx, id)
	if err != nil {
		m.log.Info("BadRequest", zap.String("error", err.Error()))
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(resp).Encode("OK"); err != nil {
		m.log.Error("Encode error", zap.Error(err))
	}
}

func (m *MyHandler) updateEvent(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		m.log.Info("BadRequest", zap.String("ReadAll", err.Error()))
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	rb := Event{}
	if err = json.Unmarshal(body, &rb); err != nil {
		m.log.Error("Unmarshal error", zap.Error(err))
	}
	m.log.Info("Update http method", zap.Int("req", int(rb.ID)))
	if rb.ID == 0 {
		m.log.Info("BadRequest", zap.Int("ID can't be zero or nil value", int(rb.ID)))
		resp.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(resp).Encode("ID can't be zero or nil value"); err != nil {
			m.log.Error("Encode error", zap.Error(err))
		}
		return
	}
	eve := parseToEventStorageStruct(rb)
	err = m.app.Update(m.ctx, eve)
	if err != nil {
		m.log.Info("BadRequest", zap.String("error", err.Error()))
		resp.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(resp).Encode(err.Error()); err != nil {
			m.log.Error("Encode error", zap.Error(err))
		}
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func parseToEventStorageStruct(req Event) storage.Event {
	var eve storage.Event
	eve.ID = req.ID
	eve.Owner = req.Owner
	eve.Title = req.Title
	eve.Description = req.Description
	eve.StartDate = req.StartDate
	eve.EndDate = req.EndDate
	eve.StartTime = req.StartTime
	eve.EndTime = req.EndTime
	return eve
}

type Event struct {
	ID          int64     `json:"ID"`
	Owner       int64     `json:"Owner"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	StartDate   string    `json:"StartDate"`
	StartTime   time.Time `json:"StartTime"`
	EndDate     string    `json:"EndDate"`
	EndTime     time.Time `json:"EndTime"`
}

type JSONID struct {
	ID int64 `json:"ID"`
}

func (m *MyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}
