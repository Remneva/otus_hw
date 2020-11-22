package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func insertDone(in In, done Bi) Out {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case vv := <-in:
				if vv != nil {
					out <- vv
				} else {
					return
				}
			case <-done:
				return
			}
		}
	}()
	return out
}
func ExecutePipeline(in In, done Bi, stages ...Stage) Out {
	out := make(Out)
	out = stages[0](insertDone(in, done))
	for _, s := range stages[1:] {
		out = insertDone(s(insertDone(out, done)), done)
	}
	return out
}
