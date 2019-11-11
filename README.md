# gengo

一些 Golang 代码自动生成工具

## install

```
go get -u github.com/zhaolion/gengo/...
``` 

## gomock generate doc

根据项目里面的接口自动生成带 `//go:generate` [golang/mock](https://github.com/golang/mock) 文件，方便维护你的 gomock 代码

install:
```
go install github.com/zhaolion/gengo/cmd/gomock
```

usage:

```
gomock -packagesBase={your package go path base} -packagesDir="." -targetDoc="mock/doc.go" -targetPrefix="" -targetPackage="mock"
go generate ./mock/...
```


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
go install github.com/zhaolion/gengo/cmd/autogen/marshal-gen
```

**在你的模型文件中添加 go generate 注释**

- 参考 [example](example/marshal-gen/model/model.go)
- `//go:generate marshal-gen -i github.com/zhaolion/gengo/example/marshal-gen/model`

**go generate it**

```
go generate ./...
```

## deepcoy-gen

自动生成 `struct` 一些方法

```
func (t T) DeepCopy() T 

or

func (t *T) DeepCopy() *T
``` 

```
func (t T) DeepCopyInto(t *T)

or

func (t *T) DeepCopyInto(t *T)
``` 

**安装 marshal-gen**

```
go install github.com/zhaolion/gengo/cmd/autogen/deepcopy-gen
```

**在你的包文件中添加 go generate 注释**

- 参考 [example](example/deepcopy-gen/model/doc.go)

```
// model test model structs
//
// gengo:deepcopy=package
//
package model

//go:generate deepcopy-gen -i github.com/zhaolion/gengo/example/marshal-gen/model
```
