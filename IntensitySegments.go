package IntentsityToolkit

import (
	"github.com/emirpasic/gods/v2/trees/redblacktree"
	"math"
)

// IntensitySegments 强度区间段
// notice:
//
//	这个区间段struct，用了一个最小值做标记位，所以，可以支持的区间段值为 [math.MinInt64+1, math.MaxInt64)
type IntensitySegments struct {
	// 内部实现是区间红黑数，key是区间段的左端点
	// 利用红黑树特性，可以再O(logN)时间内进行查找和更新
	underlying *redblacktree.Tree[int64, *Segment]
}

func NewIntensitySegments() *IntensitySegments {
	s := &IntensitySegments{
		underlying: redblacktree.New[int64, *Segment](),
	}
	s.underlying.Put(math.MinInt64, &Segment{math.MinInt64, math.MinInt64, 0})
	return s
}

func (s *IntensitySegments) Add(from, to, amount int64) {
	panic("implement me")
}

func (s *IntensitySegments) Set(from, to, amount int64) {
	panic("implement me")
}

func (s *IntensitySegments) ToString() {
	panic("implement me")
}

// baseOperate 处理区间段更新，
// amountUpdater 用于处理 重复区间段的更新逻辑，因为set和add的逻辑不同，固抽出这个function交由各自实现
func (s *IntensitySegments) baseOperate(from, to, amount int64, amountUpdater func(amount int64) int64) {
	panic("implement me")
}
