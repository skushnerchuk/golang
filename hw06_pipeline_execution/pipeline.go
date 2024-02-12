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
	out := make(chan any)
	go func() {
		defer close(out)

		sc := stageWrapper(done)
		next := in
		for _, s := range stages {
			if s != nil {
				next = s(sc(next))
			}
		}
		last := sc(next)
		for v := range last {
			out <- v
		}
	}()
	return out
}

func stageWrapper(done In) Stage {
	return func(in In) Out {
		out := make(Bi)
		go func() {
			defer close(out)

			for {
				select {
				case <-done:
					go drain(in)
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					out <- v
				}
			}
		}()
		return out
	}
}

func drain(ch Out) {
	for range ch {
	}
}
