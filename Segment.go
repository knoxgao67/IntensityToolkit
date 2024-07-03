package IntentsityToolkit

// Segment 表示区间段，左闭右开,[left,right)
// amount表示区间段的值.
type Segment struct {
	From   int64
	To     int64
	Amount int64
}

func NewSegment(from, to, amount int64) *Segment {
	if from >= to {
		return nil
	}
	return &Segment{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

type SegmentOperatorCache struct {
	createSegs []*Segment         // 标记后续需要create的区间
	deleteKeys map[int64]struct{} // 标记后续需要delete的区间key
}

func NewSegmentOperatorCache() *SegmentOperatorCache {
	return &SegmentOperatorCache{
		createSegs: make([]*Segment, 0),
		deleteKeys: make(map[int64]struct{}),
	}
}

func (so *SegmentOperatorCache) Create(list ...*Segment) {
	so.createSegs = append(so.createSegs, list...)
}

func (so *SegmentOperatorCache) Delete(list ...int64) {
	for _, item := range list {
		so.deleteKeys[item] = struct{}{}
	}
}

// Merge 合并需要创建的区间段.
func (so *SegmentOperatorCache) Merge() {
	// 去掉nil，在Add和Set的时候我们可能会加入nil
	segs := make([]*Segment, 0, len(so.createSegs))
	for _, item := range so.createSegs {
		if item == nil {
			continue
		}
		segs = append(segs, item)
	}

	// merge
	so.createSegs = segs
	segs = make([]*Segment, 0, len(so.createSegs))
	cur := so.createSegs[0]

	for i := 1; i < len(so.createSegs); i++ {
		seg := so.createSegs[i]
		if cur.To >= seg.From && cur.Amount == seg.Amount {
			cur.To = seg.To
			// seg合并进cur，seg可以删除了
			so.Delete(seg.From)
		} else {
			segs = append(segs, cur)
			cur = seg
		}
	}
	segs = append(segs, cur)

	so.createSegs = segs

	for _, item := range so.createSegs {
		delete(so.deleteKeys, item.From)
	}
}
