package internalhttp

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/apex/log"
)

type MyHandler struct {
	db sql.DB
	l  *zap.Logger
}

type Server struct {
	*http.Server
}

type Application interface {
}

func NewServer(logger *zap.Logger, app Application) *Server {
	handler := &MyHandler{l: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("/set", requestLogger(handler, handler.SetEvent))
	mux.HandleFunc("/get", requestLogger(handler, handler.GetEvent))
	mux.HandleFunc("/getAll", requestLogger(handler, handler.GetEvents))
	mux.HandleFunc("/delete", requestLogger(handler, handler.DeleteEvent))
	mux.HandleFunc("/update", requestLogger(handler, handler.UpdateEvent))
	server := &http.Server{
		Addr:           ":8081",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe().Error())

	return &Server{}
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
		return err
	}
	return nil
}

func (m *MyHandler) SetEvent(resp http.ResponseWriter, req *http.Request) {
	time.Sleep(3 * time.Second)
}

func (m *MyHandler) GetEvent(resp http.ResponseWriter, req *http.Request) {
	time.Sleep(3 * time.Second)
}

func (m *MyHandler) GetEvents(resp http.ResponseWriter, req *http.Request) {
	time.Sleep(3 * time.Second)
}

func (m *MyHandler) DeleteEvent(resp http.ResponseWriter, req *http.Request) {
	time.Sleep(3 * time.Second)
}

func (m *MyHandler) UpdateEvent(writer http.ResponseWriter, request *http.Request) {
	time.Sleep(3 * time.Second)
}
