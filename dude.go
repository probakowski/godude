package godude

import (
	"fmt"
	"github.com/dterei/gotsc"
	"github.com/montanaflynn/stats"
	"math"
	"math/rand"
	"runtime"
	"time"
)

func Measure[T, R any](m func(t T) R, genRand, genFix func() T, iters int) float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ticks := make([]float64, iters)
	tests := make([]T, iters)
	classes := make([]int, iters)
	for i := 0; i < iters; i++ {
		classes[i] = r.Intn(2)
		if classes[i] == 0 {
			tests[i] = genRand()
		} else {
			tests[i] = genFix()
		}
	}
	s := Stat{}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	tsc := gotsc.TSCOverhead()

	for i := 0; i < iters; i++ {
		t := tests[i]
		st := gotsc.BenchStart()
		m(t)
		end := gotsc.BenchEnd()
		ticks[i] = float64(end - st - tsc)
	}

	p90, _ := stats.Percentile(ticks, 90)
	fmt.Println(p90)
	for i := 10; i < len(ticks); i++ {
		if ticks[i] <= p90 {
			s.push(ticks[i], classes[i])
		}
	}

	return s.compute()
}

type Stat struct {
	mean  [2]float64
	m2    [2]float64
	count [2]float64
}

func (s *Stat) push(ticks float64, class int) {
	s.count[class]++
	delta := ticks - s.mean[class]
	s.mean[class] += delta / s.count[class]
	s.m2[class] += delta * (ticks - s.mean[class])
}

func (s *Stat) compute() float64 {
	v0 := s.m2[0] / (s.count[0] - 1)
	v1 := s.m2[1] / (s.count[1] - 1)
	num := s.mean[0] - s.mean[1]
	vv0 := v0 / s.count[0]
	vv1 := v1 / s.count[1]
	ss := vv0 + vv1
	den := math.Sqrt(ss)
	t := num / den
	return math.Abs(t)
}
