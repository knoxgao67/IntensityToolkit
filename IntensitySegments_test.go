package IntentsityToolkit

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntensitySegments_Add_Normal1(t *testing.T) {
	s := NewIntensitySegments()
	checkSegmentsValue(t, s, [][]int64{})

	s.Add(10, 30, 1)
	checkSegmentsValue(t, s, [][]int64{{10, 1}, {30, 0}})

	s.Add(20, 40, 1)
	checkSegmentsValue(t, s, [][]int64{{10, 1}, {20, 2}, {30, 1}, {40, 0}})

	s.Add(10, 40, -2)
	checkSegmentsValue(t, s, [][]int64{{10, -1}, {20, 0}, {30, -1}, {40, 0}})
}

func TestIntensitySegments_Add_Normal2(t *testing.T) {
	s := NewIntensitySegments()

	s.Add(10, 30, 1)
	checkSegmentsValue(t, s, [][]int64{{10, 1}, {30, 0}})

	s.Add(20, 40, 1)
	checkSegmentsValue(t, s, [][]int64{{10, 1}, {20, 2}, {30, 1}, {40, 0}})

	s.Add(10, 40, -1)
	checkSegmentsValue(t, s, [][]int64{{20, 1}, {30, 0}})

	s.Add(10, 40, -1)
	checkSegmentsValue(t, s, [][]int64{{10, -1}, {20, 0}, {30, -1}, {40, 0}})

	s.Add(10, 20, 1)
	checkSegmentsValue(t, s, [][]int64{{30, -1}, {40, 0}})

	s.Add(30, 40, 1)
	checkSegmentsValue(t, s, [][]int64{})
}

func TestIntensitySegments_Add_Normal3(t *testing.T) {
	s := NewIntensitySegments()

	s.Add(10, 30, 1)
	checkSegmentsValue(t, s, [][]int64{{10, 1}, {30, 0}})

	s.Add(40, 50, 1)
	checkSegmentsValue(t, s, [][]int64{{10, 1}, {30, 0}, {40, 1}, {50, 0}})

	s.Add(30, 40, 1)
	checkSegmentsValue(t, s, [][]int64{{10, 1}, {50, 0}})

	s.Add(10, 50, -1)
	checkSegmentsValue(t, s, [][]int64{})
}

func TestIntensitySegments_SetAndQuery_Normal1(t *testing.T) {
	s := NewIntensitySegments()
	checkSegmentsValue(t, s, [][]int64{})

	s.Add(10, 30, 1)
	checkSegmentsValue(t, s, [][]int64{{10, 1}, {30, 0}})

	s.Add(20, 40, 1)
	checkSegmentsValue(t, s, [][]int64{{10, 1}, {20, 2}, {30, 1}, {40, 0}})

	s.Set(10, 40, 0)
	checkSegmentsValue(t, s, [][]int64{})

	s.Add(10, 30, 1)
	s.Add(20, 40, 1)

	s.Set(50, 100, 1)
	checkSegmentsValue(t, s, [][]int64{{10, 1}, {20, 2}, {30, 1}, {40, 0}, {50, 1}, {100, 0}})

	require.EqualValues(t, 1, s.Query(15))
	require.EqualValues(t, 0, s.Query(-10))
	require.EqualValues(t, 0, s.Query(1000))
	require.EqualValues(t, 1, s.Query(35))
	require.EqualValues(t, 0, s.Query(45))

}

func checkSegmentsValue(t *testing.T, s *IntensitySegments, expected [][]int64) {
	require.EqualValues(t, expected, s.GetAllValues())
}

func BenchmarkNewSegment(b *testing.B) {
	s := NewIntensitySegments()
	for i := 0; i < b.N; i++ {
		s.Add(1, 7, 10)
	}
}
