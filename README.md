#intensityToolKit
> a toolkit for intensity


#### 思路
1. 使用红黑树维护区间段，key为区间段的左端点，value为区间段具体信息
2. 区间段更新的时候，找到最左侧区间，然后和当前区间进行比较。然后再和右侧区间进行比较。
   1. 迭代寻找区间的时候，将所有有重叠的区间都抽象理解为了三部分，增加可读性，然后再在处理完所有空间后进行合并

#### usage
```shell
go get -v github.com/knoxgao67/IntensityToolkit
```

```go
// create a new intensity segments
 s := NewIntensitySegments()

 // add segment to s
 s.Add(10, 30, 1)

 // set segment
 s.Set(15,35,2)

 // query intensity with idx
 s.Query(20)
```