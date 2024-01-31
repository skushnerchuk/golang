package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil || len(stages) == 0 {
		return nil
	}

	for _, s := range stages {
		if s == nil {
			continue
		}
		in = s(stageRunner(in, done))
	}
	return in
}

func stageRunner(in In, done In) Out {
	output := make(Bi)

	go func() {
		defer close(output)
		for {
			select {
			case v := <-in:
				if v == nil {
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
