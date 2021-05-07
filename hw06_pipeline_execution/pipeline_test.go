package hw06pipelineexecution

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}

func TestPipelineAnyData(t *testing.T) {
	stageFn := func(in In) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for v := range in {
				out <- v
			}
		}()
		return out
	}

	type User struct {
		Name string
	}

	in := make(Bi)
	data := []interface{}{1, true, "Hello", 3.14, User{"Vasya"}, int32(100)}
	go func() {
		for _, v := range data {
			in <- v
		}
		close(in)
	}()

	defineType := func(i interface{}) interface{} {
		switch i.(type) {
		case int, string, bool, float64, User:
			return i
		default:
			fmt.Println("unexpected type!")
			return nil
		}
	}

	result := make([]interface{}, 0, 10)
	for s := range ExecutePipeline(in, nil, stageFn, stageFn, stageFn) {
		result = append(result, defineType(s))
	}

	require.Equal(t, []interface{}{1, true, "Hello", 3.14, User{"Vasya"}, nil}, result)
}

func TestPipelineWithoutStages(t *testing.T) {
	stages := make([]Stage, 0)
	in := make(Bi)
	const testValue = "test"
	go func() {
		in <- testValue
		close(in)
	}()
	out := ExecutePipeline(in, nil, stages...)
	require.Equal(t, "test", <-out)
}

// Тест Алексея Бакина.
func TestPipelineDone(t *testing.T) {
	waitCh := make(chan struct{})
	defer close(waitCh)

	stageFn := func(in In) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for v := range in {
				<-waitCh
				out <- v
			}
		}()
		return out
	}

	in := make(Bi)
	const testValue = "test"
	go func() {
		in <- testValue
		close(in)
	}()

	doneCh := make(Bi)
	var resValue interface{}
	out := ExecutePipeline(in, doneCh, stageFn, stageFn, stageFn)
	close(doneCh)

	require.Eventually(t, func() bool {
		select {
		case resValue = <-out:
			return true
		default:
			return false
		}
	}, time.Second, time.Millisecond)

	require.Nil(t, resValue)
}
