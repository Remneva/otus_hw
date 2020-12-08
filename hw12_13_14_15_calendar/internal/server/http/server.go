package internalhttp //nolint:golint,stylecheck

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/apex/log"
	"go.uber.org/zap"
)

type MyHandler struct {
	l *zap.Logger
}

type Server struct {
	*http.Server
}

type Application interface {
}

func NewServer(logger *zap.Logger, app Application) *Server {
	handler := &MyHandler{l: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("/set", requestLoggerMiddleware(handler, handler.SetEvent))
	mux.HandleFunc("/get", requestLoggerMiddleware(handler, handler.GetEvent))
	mux.HandleFunc("/getAll", requestLoggerMiddleware(handler, handler.GetEvents))
	mux.HandleFunc("/delete", requestLoggerMiddleware(handler, handler.DeleteEvent))
	mux.HandleFunc("/update", requestLoggerMiddleware(handler, handler.UpdateEvent))
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
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.Shutdown(ctx)
	if err != nil {
		return errors.New("shutdown error")
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
