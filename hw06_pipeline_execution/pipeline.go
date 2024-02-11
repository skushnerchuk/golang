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

	out := make(Bi)
	prev := in
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		bufChannel := make(Bi)
		go channelWrapper(bufChannel, prev, done)
		next := make(Bi)
		go stageRunner(stage, bufChannel, next, done)
		prev = next
	}

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-prev:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()

	return out
}

func channelWrapper(bufChannel Bi, in In, done In) {
	defer close(bufChannel)
	for v := range in {
		select {
		case bufChannel <- v:
		case <-done:
			return
		}
	}
}

func stageRunner(stage Stage, in In, next Bi, done In) {
	defer close(next)
	for v := range stage(in) {
		select {
		case next <- v:
		case <-done:
			return
		}
	}
}
