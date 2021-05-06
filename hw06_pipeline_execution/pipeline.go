package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	result := in
	for _, stage := range stages {
		result = doneStage(done, stage(result))
	}

	return result
}

func doneStage(done In, in In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for {
				_, ok := <-in
				if !ok {
					break
				}
			}
		}()

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-done:
					return
				}
			}
		}
	}()
	return out
}
