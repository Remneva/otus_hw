package server

import "context"

type Stopper interface {
	Stop(context.Context) error
}
