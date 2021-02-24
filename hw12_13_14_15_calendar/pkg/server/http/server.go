package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/app"
	srv "github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/server"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/pkg/storage"
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
		m.respond(err, resp)
		return
	}
	rb := Event{}
	if err = json.Unmarshal(body, &rb); err != nil {
		m.respond(err, resp)
		return
	}
	eve := parseToEventStorageStruct(rb)
	id, err := m.app.Create(m.ctx, eve)
	if err != nil {
		m.respond(err, resp)
		return
	}
	r := JSONID{}
	r.ID = id
	m.respond(r, resp)
}

func (m *MyHandler) getEvent(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		m.respond(err, resp)
		return
	}

	r := JSONID{}
	id := r.ID
	if err = json.Unmarshal(body, &id); err != nil {
		m.respond(err, resp)
		return
	}
	var result storage.Event
	result, err = m.app.Get(m.ctx, id)
	if err != nil {
		m.respond(err, resp)
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
	m.respond(response, resp)
}

func (m *MyHandler) deleteEvent(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		m.respond(err, resp)
		return
	}
	r := JSONID{}
	id := r.ID
	if err = json.Unmarshal(body, &id); err != nil {
		m.respond(err, resp)
		return
	}
	err = m.app.Delete(m.ctx, id)
	if err != nil {
		m.respond(err, resp)
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
		m.respond(err, resp)
		return
	}
	rb := Event{}
	if err = json.Unmarshal(body, &rb); err != nil {
		m.respond(err, resp)
	}
	m.log.Info("Update http method", zap.Int("req", int(rb.ID)))
	if rb.ID == 0 {
		m.respond("ID can't be zero or nil value", resp)
		return
	}
	eve := parseToEventStorageStruct(rb)
	err = m.app.Update(m.ctx, eve)
	if err != nil {
		m.respond(err, resp)
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

func (m *MyHandler) respond(data interface{}, resp http.ResponseWriter) {
	switch i := data.(type) {
	case Event, JSONID:
		resp.WriteHeader(200)
		if err := json.NewEncoder(resp).Encode(&i); err != nil {
			m.log.Error("Encode error", zap.Error(err))
		}
	case error:
		m.writeResponse(i.Error(), 400, resp)
	case string:
		m.writeResponse(i, 400, resp)
	}
}

func (m *MyHandler) writeResponse(message string, code int, resp http.ResponseWriter) {
	m.log.Error("BadRequest", zap.String("error", message))
	resp.WriteHeader(code)
	if err := json.NewEncoder(resp).Encode(message); err != nil {
		m.log.Error("Encode error", zap.Error(err))
	}
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
