# gengo

一些 Golang 代码自动生成工具

## marshal-gen

自动生成 `struct` 一些方法

```
MarshalJSONBinary() ([]byte, error)

UnmarshalJSONBinary(data []byte) error

String() string
``` 

> 选择不实现 以下方法是因为:    
> 生成 `MarshalJSON` 方法存在嵌套结构体出现 `goroutine stack exceeds` 问题  
> 生成 `UnmarshalJSON` 方法存在嵌套结构体出现 `goroutine stack exceeds` 问题   

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
