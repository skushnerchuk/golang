package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, s := range stages {
		in = s(manager(in, done))
	}
	return in
}

func manager(in In, done In) Out {
	output := make(Bi)

	go func() {
		defer close(output)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				output <- v
			case <-done:
				return
			}
		}
	}()

	return output
}
