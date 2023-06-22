package time_serie

import (
	"math"

	"github.com/gpabois/gostd/collection"
	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/ops"
	"github.com/gpabois/gostd/unit"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type TimeInterval = ops.Interval[int]

type Number interface {
	constraints.Integer | constraints.Float
}

type TimeIndex struct {
	Step     int           `serde:"t"`
	Sampling unit.Sampling `serde:"sampling"`
}

type TimeSeriePoint[Value any] struct {
	Step  int
	Value Value
}

func FromPair[Pair collection.IPair[int, Value], Value any](pair Pair) TimeSeriePoint[Value] {
	return TimeSeriePoint[Value]{
		Step:  pair.GetFirst(),
		Value: pair.GetSecond(),
	}
}

type TimeSerie[Value any] struct {
	Sampling unit.Sampling
	Begin    int
	End      int
	Points   []TimeSeriePoint[Value]
}

func (ts *TimeSerie[Value]) GetPeriod() int {
	return ts.Sampling.Period
}

func (ts *TimeSerie[Value]) GetUnit() string {
	return ts.Sampling.Unit
}

func (ts *TimeSerie[Value]) IterValues() iter.Iterator[Value] {
	return iter.Map(iter.IterSlice(&ts.Points), func(point TimeSeriePoint[Value]) Value {
		return point.Value
	})
}

func (ts *TimeSerie[Value]) ChangeUnit(newUnit string, aggregator func(iter.Chunk[Value]) Value) TimeSerie[Value] {
	newPeriod := int(1.0 / unit.GetUnitConversion(newUnit))
	nts := ts.ChangePeriod(newPeriod, aggregator)
	nts.Sampling.Unit = newUnit
	return nts
}

// Change the period of time serie
// If the new period is <= to the ts's period, returns a copy of it.
func (ts *TimeSerie[Value]) ChangePeriod(newPeriod int, aggregator func(iter.Chunk[Value]) Value) TimeSerie[Value] {
	chunkSize := int(math.Round(float64(newPeriod) / float64(ts.GetPeriod())))
	chunks := iter.ChunkEvery(
		ts.IterValues(),
		chunkSize,
	)

	// Get the offset
	offset := ts.Begin

	// Add the offset to the chunk enumeration
	points := iter.Map(iter.Enumerate(iter.Map(chunks, aggregator)), func(pair iter.Enumeration[Value]) iter.Enumeration[Value] {
		pair.First += offset
		return pair
	})

	return FromIter[Value](points, unit.Sampling{
		Unit:   ts.Sampling.Unit,
		Period: newPeriod,
	})
}

func FromIter[Value any, Pair collection.IPair[int, Value]](iterator iter.Iterator[Pair], sampling unit.Sampling) TimeSerie[Value] {
	// Transform each KV into Point
	points := iter.CollectToSlice[[]TimeSeriePoint[Value]](iter.Map(iterator, FromPair[Pair, Value]))

	// Order the points
	slices.SortFunc(points, func(a, b TimeSeriePoint[Value]) bool {
		return a.Step < b.Step
	})

	// Get the window of the time serie
	begin, end := points[0], points[len(points)-1]

	return TimeSerie[Value]{
		Sampling: sampling,
		Begin:    begin.Step,
		End:      end.Step,
		Points:   points,
	}
}
