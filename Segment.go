package IntentsityToolkit

// Segment 表示区间段，左闭右开,[left,right)
// amount表示区间段的值
type Segment struct {
	From   int64
	To     int64
	Amount int64
}
