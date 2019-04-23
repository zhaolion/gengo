//go:generate marshal-gen -i github.com/zhaolion/gengo/example/marshal-gen/model

package model

type T1 struct {
	Byte byte
	//Int8    int8 //TODO: int8 becomes byte in SnippetWriter
	Int16   int16
	Int32   int32
	Int64   int64
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64
	Str     string
}

type Inner interface {
	Function() float64
	DeepCopyInner() Inner
}

type T2 struct {
	I []Inner
}

type T3 struct {
	Byte map[string]byte
	//Int8    map[string]int8 //TODO: int8 becomes byte in SnippetWriter
	Int16        map[string]int16
	Int32        map[string]int32
	Int64        map[string]int64
	Uint8        map[string]uint8
	Uint16       map[string]uint16
	Uint32       map[string]uint32
	Uint64       map[string]uint64
	Float32      map[string]float32
	Float64      map[string]float64
	StringPtr    map[string]*string
	StringPtrPtr map[string]**string
	Map          map[string]map[string]string
	MapPtr       map[string]*map[string]string
	Slice        map[string][]string
	SlicePtr     map[string]*[]string
	Struct       map[string]T1
	StructPtr    map[string]*T2
}
