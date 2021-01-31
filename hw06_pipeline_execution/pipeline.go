package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	transmitter := func(in Bi, done In, out Out) {
		for {
			select {
			case <-done:
				close(in)
				return
			default:
				select {
				case <-done:
					close(in)
					return
				case tmp, ok := <-out:
					if ok {
						in <- tmp
					} else {
						close(in)
						return
					}
				}
			}
		}
	}
	//var wg = sync.WaitGroup{}
	for _, stage := range stages {
		//wg.Add(1)
		//stage := stage
		//go func() {
		//	defer wg.Done()
		//	input := make(Bi)
		//	wg.Add(1)
		//	go func() {
		//		defer wg.Done()
		//		transmitter(input, done, in)
		//	}()
		//	in = stage(input)
		//}()
		input := make(Bi)
		go transmitter(input, done, in)
		in = stage(input)
	}
	//wg.Wait()
	return in
}
