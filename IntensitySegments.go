package IntentsityToolkit

import (
	"fmt"
	"math"

	"github.com/emirpasic/gods/v2/trees/redblacktree"
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
	s.baseOperate(from, to, amount, func(curAmount int64) int64 {
		if (amount > 0 && curAmount > 0 && curAmount+amount < 0) || (amount < 0 && curAmount < 0 && amount+curAmount > 0) {
			panic("overflow as the amount is too large")
		}
		return curAmount + amount
	})
}

func (s *IntensitySegments) Set(from, to, amount int64) {
	s.baseOperate(from, to, amount, func(curAmount int64) int64 {
		return amount
	})
}

func (s *IntensitySegments) ToString() {
	fmt.Println(s.GetAllValues())
}

func (s *IntensitySegments) GetAllValues() [][]int64 {
	list := make([][]int64, 0, s.underlying.Size()) // 去除哨兵位置
	if s.underlying.Size() <= 1 {
		return [][]int64{}
	}
	iter := s.underlying.Iterator()
	iter.Next() // 这里的目的是跳过哨兵
	prev := iter.Value()
	for iter.Next() {
		if prev.From != math.MinInt64 && prev.To < iter.Value().From {
			list = append(list, []int64{prev.To, 0})
		}
		list = append(list, []int64{iter.Value().From, iter.Value().Amount})
		prev = iter.Value()
	}
	iter.Prev()
	list = append(list, []int64{iter.Value().To, 0})
	return list
}

func (s *IntensitySegments) Query(idx int64) int64 {
	curNode, ok := s.underlying.Floor(idx)
	if !ok {
		panic("error, can't find the segment")
	}
	if curNode.Value.To < idx {
		return 0
	}
	return curNode.Value.Amount
}

// baseOperate 处理区间段更新，
// amountUpdater 用于处理 重复区间段的更新逻辑，因为set和add的逻辑不同，固抽出这个function交由各自实现.
func (s *IntensitySegments) baseOperate(from, to, amount int64, amountUpdater func(amount int64) int64) {
	if from >= to {
		panic("param error, from must small than to")
	}

	// 在tree里面所有小于等于from的keys中找到最大的那一个
	curNode, ok := s.underlying.Floor(from)
	if !ok { // 如果没有找到
		// 因为插入了最小的节点，这里一定不能为false，如果出现false是代码的bug
		panic("something error, there must have a small node")
	}

	iter := s.underlying.IteratorAt(curNode) // 迭代器

	segOper := newSegmentOperatorCache()

	// curNode是前置区间段
	if curNode.Value.To > from {
		seg := curNode.Value
		// 因为区间合并判断太过复杂，代码可读性很差，这里逻辑做了一定抽象处理
		// 将oldSeg 和 newSeg进行合并，会得到三部分区间段，即重叠，重叠前，重叠后。对应下面添加的三个区间段
		// 如果这三部分区间段部分不存在，NewSegment会返回nil，或许会进行过滤
		segOper.Create(
			NewSegment(seg.From, from, seg.Amount),                       // 重叠 区间前半部分
			NewSegment(from, min(to, seg.To), amountUpdater(seg.Amount)), // 重叠 区间
			NewSegment(to, seg.To, seg.Amount),                           // 重叠 区间后半部分
		)
		segOper.Delete(seg.From)
	}

	tmpFrom := max(from, curNode.Value.To)
	// 遍历后面的区间段，直到后面的区间段的 from 大于 新区间段的 to
	for iter.Next() {
		seg := iter.Value()
		if seg.From >= to {
			break
		}

		segOper.Create(
			NewSegment(tmpFrom, seg.From, amount),
			NewSegment(seg.From, min(to, seg.To), amountUpdater(seg.Amount)),
			NewSegment(to, seg.To, seg.Amount),
		)
		segOper.Delete(seg.From)

		// 这里更新的意思是 tmpFrom以前的已经都处理了
		tmpFrom = seg.To
	}

	if tmpFrom < to {
		segOper.Create(NewSegment(tmpFrom, to, amountUpdater(0)))
	}

	segOper.Merge()
	for key := range segOper.deleteKeys {
		s.underlying.Remove(key)
	}
	for _, item := range segOper.createSegs {
		s.underlying.Put(item.From, item)
	}

	s.merge(from, to)
}

// merge 合并[prev] [from,to) [next]区间段
// 这个和segmentOperate不同在于 它是在树内操作，主要是用于合并已经写入的区间.
func (s *IntensitySegments) merge(from, to int64) {
	startNode, ok := s.underlying.Ceiling(from)
	if !ok {
		// 如果找不到 则返回
		return
	}
	iter := s.underlying.IteratorAt(startNode)
	iter.Prev()
	prev := iter.Value()
	var needDeleteKeys []int64
	for iter.Next() {
		seg := iter.Value()
		// should merge
		if seg.From == prev.To && seg.Amount == prev.Amount {
			needDeleteKeys = append(needDeleteKeys, seg.From)
			prev.To = seg.To
			continue
		}
		if seg.From > to {
			break
		}
		prev = seg
	}

	// 将重复的segment删除
	for _, key := range needDeleteKeys {
		s.underlying.Remove(key)
	}

	// 最前和最后的0需要移除
	for s.underlying.Size() > 1 {
		iter := s.underlying.Iterator()
		iter.First()
		iter.Next()
		if iter == nil || iter.Value().Amount != 0 {
			break
		}
		s.underlying.Remove(iter.Value().From)
	}

	for s.underlying.Size() > 1 {
		iter := s.underlying.Iterator()
		iter.Last()
		if iter.Value() == nil || iter.Value().Amount != 0 {
			break
		}
		s.underlying.Remove(iter.Value().From)
	}
}
