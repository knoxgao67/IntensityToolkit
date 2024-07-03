package IntentsityToolkit

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSegmentOperatorCache_Merge(t *testing.T) {
	segOper := newSegmentOperatorCache()
	segOper.Create(
		NewSegment(10, 20, 1),
		NewSegment(20, 30, 1),
		NewSegment(30, 40, 1),
		NewSegment(40, 50, 1))
	segOper.Merge()
	require.EqualValues(t, []*Segment{{10, 50, 1}}, segOper.createSegs)
	require.NotEmpty(t, 3, len(segOper.deleteKeys))

	segOper = newSegmentOperatorCache()
	segOper.Create(
		NewSegment(10, 20, 1),
		NewSegment(20, 30, 2),
		NewSegment(30, 40, 1),
		NewSegment(40, 50, 1))
	segOper.Merge()
	require.EqualValues(t, []*Segment{
		{10, 20, 1},
		{20, 30, 2},
		{30, 50, 1},
	}, segOper.createSegs)
	require.NotEmpty(t, 1, len(segOper.deleteKeys))
}
