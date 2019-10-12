package usl

import (
	"math"

	"gonum.org/v1/gonum/optimize"
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
	problem := optimize.Problem{
		Func: func(x []float64) float64 {
			ssRes := 0.0
			m := Model{
				Sigma:  x[0],
				Kappa:  x[1],
				Lambda: x[2],
			}

			for i := range points {
				actualY := points[i][1]
				testY := m.Throughput(points[i][0])
				err := testY - actualY
				ssRes += err * err
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

	sum := float64(0)
	for i := range points {
		sum += points[i][1]
	}

	mean := sum / float64(len(points))

	ss := float64(0)
	for i := range points {
		v := points[i][1] - mean
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
