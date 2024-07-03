# intensityToolKit
> a toolkit for intensity


#### 思路
1. 使用红黑树维护区间段，key为区间段的左端点，value为区间段具体信息
2. `Add`和`Set`操作抽象`baseOperate`用于完成核心逻辑处理。`baseOperate`负责完成区间段查找，拆分，合并，删除等逻辑。
3. `Add`和`Set`操作只负责角色amount如何更新
4. `baseOperate`核心逻辑如下：
    - 查找区间段，根据`form`找到最左的区间段
    - 拆分区间段，从找的区间段开始向右遍历，找到所有与当前区间段有交集的区间段，将其拆分为三部分:重叠左，重叠部分，重叠右.(TODO 补充更详细说明)
    - 合并区间段，将所有需要新增的区间段合并
    - 删除区间段，删除不再需要的区间段

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

