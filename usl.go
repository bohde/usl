package usl

import (
	"math"

	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/gonum/stat"
)

type Model struct {
	Sigma  float64
	Kappa  float64
	Lambda float64
	R2     float64
}

func (m Model) Throughput(n float64) float64 {
	return (m.Lambda * n) / (1 + (m.Sigma * (n - 1)) + (m.Kappa * n * (n - 1)))
}

func (m Model) MaxConcurrency() float64 {
	return math.Sqrt((1 - m.Sigma) / m.Kappa)
}

func (m Model) MaxThroughput() float64 {
	return m.Throughput(m.MaxConcurrency())
}

func Fit(points [][2]float64) (Model, error) {
	ys := make([]float64, len(points))
	xs := make([]float64, len(points))

	for i := range points {
		xs[i] = points[i][0]
		ys[i] = points[i][1]
	}

	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			ssRes := 0.0
			m := Model{
				Sigma:  x[0],
				Kappa:  x[1],
				Lambda: x[2],
			}

			for i := range ys {
				actualY := ys[i]
				testY := m.Throughput(xs[i])
				err := (testY - actualY)
				ssRes += (err * err)
			}

			return ssRes
		},
	}

	settings := optimize.Settings{
		Converger: &optimize.FunctionConverge{
			Iterations: 5000,
		},
	}

	result, err := optimize.Minimize(problem, []float64{0.01, 0.1, 1000}, &settings, &optimize.NelderMead{})
	if err != nil && err != optimize.ErrLinesearcherFailure {
		return Model{}, err
	}

	mean := stat.Mean(ys, nil)

	ss := float64(0)
	for i := range ys {
		v := ys[i] - mean
		ss += (v * v)
	}

	m := Model{
		Sigma:  result.X[0],
		Kappa:  result.X[1],
		Lambda: result.X[2],
		R2:     1 - (result.F / ss),
	}

	return m, nil
}
