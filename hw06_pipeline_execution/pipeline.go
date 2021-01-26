package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	//var input Bi
	input := make(Bi)
	go func() {
		//input = make(Bi)
		defer close(input)
		for item := range in{
			select {
			case <-done:
				return
			case input <- item:
			}
		}
	}()
	var out Out
	for _, stage := range stages {
		if out == nil {
			out = stage(input)
		} else {
			out = stage(out)
		}
	}
	return out
}
