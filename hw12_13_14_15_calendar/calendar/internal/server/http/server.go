package internalhttp

import (
	"context"
	"encoding/json"
	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Remneva/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func New() Server {
	return Server{}
}

type Server struct {
	app *app.App
	*http.Server
}

type Application interface {
}

type MyHandler struct {
	app *app.App
	ctx context.Context
}

func (m *MyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

func (s *Server) NewServer(mux *http.ServeMux, port string) (*Server, error) {
	server := &http.Server{
		Addr:           port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		s.app.Log.Error("ListenAndServe", zap.Error(err))
		return nil, errors.Wrap(err, "creating a new ServerTransport failed")
	}
	return &Server{}, nil
}
func NewHandler(ctx context.Context, app *app.App) (*MyHandler, *http.ServeMux) {
	handler := &MyHandler{ctx: ctx, app: app}
	mux := http.NewServeMux()
	mux.HandleFunc("/set", requestLoggerMiddleware(handler, handler.SetEvent))
	mux.HandleFunc("/get", requestLoggerMiddleware(handler, handler.GetEvent))
	mux.HandleFunc("/delete", requestLoggerMiddleware(handler, handler.DeleteEvent))
	mux.HandleFunc("/update", requestLoggerMiddleware(handler, handler.UpdateEvent))
	return handler, mux
}

func (s *Server) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.Shutdown(ctx)
	if err != nil {
		return errors.New("shutdown error")
	}
	return nil
}

func (m *MyHandler) SetEvent(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		m.app.Log.Info("BadRequest", zap.String("ReadAll", err.Error()))
	}
	rb := Event{}
	if err = json.Unmarshal(body, &rb); err != nil {
		m.app.Log.Error("Unmarshal error", zap.Error(err))
	}

	var eve storage.Event
	eve.Owner = rb.Owner
	eve.Title = rb.Title
	eve.Description = rb.Description
	eve.StartDate = rb.StartDate
	eve.EndDate = rb.EndDate
	eve.StartTime = rb.StartTime
	eve.EndTime = rb.EndTime

	id, err := m.app.Repo.AddEvent(m.ctx, eve)

	if err != nil {
		m.app.Log.Info("BadRequest", zap.String("error", err.Error()))
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(resp).Encode(id); err != nil {
		m.app.Log.Error("Encode error", zap.Error(err))
	}
}

func (m *MyHandler) GetEvent(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		m.app.Log.Info("BadRequest", zap.String("ReadAll", err.Error()))
	}

	r := JSONID{}
	id := r.ID
	if err = json.Unmarshal(body, &id); err != nil {
		m.app.Log.Error("Unmarshal error", zap.Error(err))
	}

	result, err := m.app.Repo.GetEvent(m.ctx, id)
	if err != nil {
		m.app.Log.Info("BadRequest", zap.String("error", err.Error()))
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
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(resp).Encode(&response); err != nil {
		m.app.Log.Error("Encode error", zap.Error(err))
	}
}

func (m *MyHandler) DeleteEvent(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		m.app.Log.Info("BadRequest", zap.String("ReadAll", err.Error()))
	}
	r := JSONID{}
	id := r.ID
	if err = json.Unmarshal(body, &id); err != nil {
		m.app.Log.Error("Unmarshal error", zap.Error(err))
	}
	err = m.app.Repo.DeleteEvent(m.ctx, id)
	if err != nil {
		m.app.Log.Info("BadRequest", zap.String("error", err.Error()))
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
}

func (m *MyHandler) UpdateEvent(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		m.app.Log.Info("BadRequest", zap.String("ReadAll", err.Error()))
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	rb := Event{}
	if err = json.Unmarshal(body, &rb); err != nil {
		m.app.Log.Error("Unmarshal error", zap.Error(err))
	}
	m.app.Log.Info("Update grpc method", zap.Int("req", int(rb.ID)))
	if rb.ID == 0 {
		m.app.Log.Info("BadRequest", zap.Int("ID can't be zero or nil value", int(rb.ID)))
		resp.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(resp).Encode("ID can't be zero or nil value"); err != nil {
			m.app.Log.Error("Encode error", zap.Error(err))
		}

		return
	}

	var eve storage.Event
	eve.ID = rb.ID
	eve.Owner = rb.Owner
	eve.Title = rb.Title
	eve.Description = rb.Description
	eve.StartDate = rb.StartDate
	eve.EndDate = rb.EndDate
	eve.StartTime = rb.StartTime
	eve.EndTime = rb.EndTime

	err = m.app.Repo.UpdateEvent(m.ctx, eve)
	if err != nil {
		m.app.Log.Info("BadRequest", zap.String("error", err.Error()))
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
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
