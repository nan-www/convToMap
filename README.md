# convToMap

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

一个强大的 Go 代码生成工具，用于自动生成 struct 与 `map[string]any` 之间的双向转换方法。

## ✨ 特性

- 🚀 **双向转换**：自动生成 `ToStringMap()` 和 `Map2Struct()` 方法
- 🎯 **智能处理**：支持嵌套结构体、指针字段、内联结构体
- 🏷️ **JSON Tag 支持**：自动识别 `json` tag，支持 `omitempty` 和 `-` 忽略标记
- 🔄 **递归转换**：自动处理嵌套的自定义结构体
- 📦 **零依赖运行时**：生成的代码无需额外依赖

## 📦 安装

```bash
go install github.com/nan-www/convToMap@latest
```

## ⚠️ 限制和注意事项

1. **跨包结构体**：目前不支持内联来自不同包的结构体（普通嵌套支持）
2. **复杂类型**：不支持 slice、array、map 等复杂类型的自动转换
3. **类型断言**：`Map2Struct` 使用类型断言，运行时类型不匹配会 panic
4. **同文件要求**：内联的结构体必须在同一个文件中定义

## 🚀 快速开始

### 1. 在你的结构体上添加注释

在需要生成转换方法的结构体前添加 `//go:generate convToMap` 注释：

```go
package example

//go:generate convToMap example.go
type User struct {
    ID       int     `json:"id"`
    Name     string  `json:"name,omitempty"`
    Email    string  `json:"email"`
    Age      *int    `json:"age,omitempty"`
    Profile  Profile `json:"profile"`
}

//go:generate convToMap example.go
type Profile struct {
    Bio    string `json:"bio"`
    Avatar string `json:"avatar,omitempty"`
}
```

### 2. 运行代码生成

```bash
convToMap example.go
```

或者使用 `go generate`：

```bash
go generate ./...
```

### 3. 使用生成的方法

生成的代码会创建两个文件：
- `example_generated_0.go` - 包含 `ToStringMap()` 方法
- `example_generated_1.go` - 包含 `Map2Struct()` 方法

#### Struct 转 Map

```go
user := &User{
    ID:    1,
    Name:  "Alice",
    Email: "alice@example.com",
    Profile: Profile{
        Bio:    "Software Engineer",
        Avatar: "avatar.jpg",
    },
}

// 转换为 map
m := user.ToStringMap()
// m = map[string]any{
//     "id": 1,
//     "name": "Alice",
//     "email": "alice@example.com",
//     "profile": map[string]any{
//         "bio": "Software Engineer",
//         "avatar": "avatar.jpg",
//     },
// }
```

#### Map 转 Struct

```go
m := map[string]any{
    "id":    1,
    "name":  "Alice",
    "email": "alice@example.com",
    "profile": map[string]any{
        "bio":    "Software Engineer",
        "avatar": "avatar.jpg",
    },
}

user := &User{}
user.Map2Struct(m)
// user 现在包含了 map 中的所有数据
```

## 📖 功能详解

### 支持的字段类型

| 类型 | 说明 | 示例 |
|------|------|------|
| 基本类型 | int, int32, int64, float32, float64, string, bool | `ID int` |
| 指针类型 | 基本类型的指针 | `Age *int` |
| 结构体 | 自定义结构体 | `Profile Profile` |
| 结构体指针 | 自定义结构体指针 | `Profile *Profile` |
| 内联结构体 | 使用 `json:",inline"` 标记 | `BaseModel json:",inline"` |

### JSON Tag 支持

- **字段重命名**：`json:"custom_name"` - 在 map 中使用自定义键名
- **omitempty**：`json:"name,omitempty"` - 零值时不添加到 map（仅 ToStringMap）
- **忽略字段**：`json:"-"` - 完全忽略该字段
- **内联**：`json:",inline"` - 将嵌套结构体的字段展平到父级

### 零值处理

`ToStringMap()` 方法会智能处理零值：

- **string**: 空字符串不会添加到 map（如果有 omitempty）
- **int/int32/int64/float32/float64**: 0 值不会添加到 map（如果有 omitempty）
- **指针**: nil 指针不会添加到 map
- **结构体**: 始终调用其 `ToStringMap()` 方法

### 嵌套结构体

工具会自动处理嵌套的自定义结构体：

```go
//go:generate convToMap example.go
type Company struct {
    Name    string `json:"name"`
    Address Address `json:"address"`
}

//go:generate convToMap example.go
type Address struct {
    Street string `json:"street"`
    City   string `json:"city"`
}
```

生成的代码会递归调用嵌套结构体的转换方法。

### 内联结构体

支持使用 `json:",inline"` 标记的内联结构体：

```go
//go:generate convToMap example.go
type BaseModel struct {
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

//go:generate convToMap example.go
type User struct {
    BaseModel `json:",inline"`
    ID        int    `json:"id"`
    Name      string `json:"name"`
}
```

内联字段会被展平到父结构体的 map 中。

## 📝 完整示例

查看 [unit_test/example.go](./unit_test/example.go) 和 [unit_test/example_test.go](./unit_test/example_test.go) 获取完整的使用示例。

```go
package unit_test

//go:generate convToMap example.go
type Example struct {
    FooPtr       *Foo              `json:"fooPtr"`
    Foo          Foo               `json:"foo,omitempty"`
    ID           int               `json:"id,omitempty"`
    Name         string            `json:"name,omitempty"`
    Float        float64           `json:"float64,omitempty"`
    Ignore       map[string]string `json:"-"`
    PtrInt       *int64            `json:"ptrInt,omitempty"`
    InlineStruct `json:",inline"`
}

//go:generate convToMap example.go
type InlineStruct struct {
    A string
    B int
}

//go:generate convToMap example.go
type Foo struct {
    Bar string `json:"bar"`
}
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🔗 相关链接

- [Go AST 文档](https://pkg.go.dev/go/ast)
- [Go Generate 文档](https://go.dev/blog/generate)

## 💡 使用场景

- **API 开发**：在 HTTP handler 中快速转换请求/响应
- **数据库操作**：与 NoSQL 数据库（如 MongoDB）交互
- **配置管理**：动态配置的序列化/反序列化
- **测试**：快速构造测试数据
- **数据迁移**：在不同数据格式之间转换
