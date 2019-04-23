# gengo

一些 Golang 代码自动生成工具

## marshal-gen

自动生成 `struct` 的`json marshaler/unmarshaler` & `stringer` 方法

```
MarshalJSON() ([]byte, error)

UnmarshalJSON(data []byte) error

String() string
``` 

**安装 marshal-gen**

```
go install github.com/zhaolion/gengo/marshal-gen
```

**在你的模型文件中添加 go generate 注释**

- 参考 [example](example/marshal-gen/model/model.go)
- `//go:generate marshal-gen -i github.com/zhaolion/gengo/example/marshal-gen/model`

**go generate it**

```
go generate ./...
```
