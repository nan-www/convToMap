package ps

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"

	"github.com/nan-www/convToMap/ds"
)

const tag = "//go:generate convToMap"

type TemplateData struct {
	PackageName string
	Structs     []*Struct
}
type Struct struct {
	Name      string
	Fields    []Field
	ASTStruct *ast.StructType
}

type Field struct {
	Name     string
	TagName  string // 用于 map 的键名，通常是 JSON tag
	IsObj    bool
	IsPtrObj bool
	Type     string
}

func ParseMarkStruct(node *ast.File) *TemplateData {
	data := TemplateData{
		PackageName: node.Name.Name,
		Structs:     []*Struct{},
	}
	name2Node := make(map[string]*ds.Node[Struct])

	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE || genDecl.Doc == nil {
			return true
		}

		foundGenerate := false
		for _, comment := range genDecl.Doc.List {
			if strings.Contains(comment.Text, tag) {
				foundGenerate = true
				break
			}
		}
		if !foundGenerate {
			return true
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			currentStruct := Struct{
				Name:      typeSpec.Name.Name,
				Fields:    []Field{},
				ASTStruct: structType,
			}
			currentNode := &ds.Node[Struct]{
				Val:        &currentStruct,
				Children:   make([]*ds.Node[Struct], 0),
				ParentsNum: 0,
			}

			if cs, ok := name2Node[currentStruct.Name]; ok {
				cs.Val = &currentStruct
			} else {
				name2Node[currentStruct.Name] = currentNode
			}
			for _, field := range structType.Fields.List {
				// 如果有 inline 则构建依赖关系
				if len(field.Names) == 0 {
					//se, ok := field.Type.(*ast.SelectorExpr)
					// 先支持 inline 结构体在当前文件的情况
					if ident, ok := field.Type.(*ast.Ident); ok {
						parentStructName := ident.Name
						ps := name2Node[parentStructName]
						if ps == nil {
							name2Node[parentStructName] = &ds.Node[Struct]{
								Children:   make([]*ds.Node[Struct], 0),
								ParentsNum: 0,
							}
							ps = name2Node[parentStructName]
						}
						ps.Children = append(ps.Children, currentNode)
						currentNode.ParentsNum += 1
					}
				}
			}
			data.Structs = append(data.Structs, &currentStruct)
		}
		return false
	})
	for _, s := range data.Structs {
		node := name2Node[s.Name]
		if node.ParentsNum != 0 {
			continue
		}
		dfsProcessNode(node, name2Node)
	}
	return &data
}

func dfsProcessNode(node *ds.Node[Struct], name2Node map[string]*ds.Node[Struct]) {
	if node.ParentsNum < 0 {
		return
	}
	processStructField(node, name2Node)
	node.ParentsNum -= 1
	for _, child := range node.Children {
		child.ParentsNum -= 1
		if child.ParentsNum == 0 {
			// DFS多叉树
			dfsProcessNode(child, name2Node)
		}
	}
}

func processStructField(currentNode *ds.Node[Struct], name2Node map[string]*ds.Node[Struct]) {
	currentStruct := currentNode.Val
	for _, field := range currentStruct.ASTStruct.Fields.List {
		// 特殊处理 inline 字段
		if len(field.Names) == 0 {
			if ident, ok := field.Type.(*ast.Ident); ok {
				pn := name2Node[ident.Name]
				if pn == nil {
					fmt.Fprintf(os.Stderr, "Warning: Can't not find relevant struct node for inline struct: %s\n", ident.Name)
					panic("")
				}
				currentStruct.Fields = append(currentStruct.Fields, pn.Val.Fields...)
			}
			continue
		}
		var isObj bool
		var isPtrObj bool
		var typeStr string
		if ident, ok := field.Type.(*ast.Ident); ok {
			// 特殊处理结构体字段
			if ident.Obj != nil {
				isObj = true
				typeStr = ident.Obj.Name
			} else {
				typeStr = ident.Name
			}
		}

		if se, ok := field.Type.(*ast.StarExpr); ok {
			if ident, ok := se.X.(*ast.Ident); ok {
				if ident.Obj != nil {
					isPtrObj = true
					typeStr = ident.Obj.Name
				} else {
					typeStr = "*" + ident.Name
				}
			}
		}
		fieldName := field.Names[0].Name
		tagName := fieldName // 默认使用字段名

		// 尝试解析 tag
		if field.Tag != nil {
			// field.Tag.Value 是带有反引号的字符串，需要解析
			// 例如：`json:"id,omitempty"`
			tagString := strings.Trim(field.Tag.Value, "`")

			// 暂时只支持json tag
			if tag, found := reflectTag(tagString, "json"); found {
				// 忽略 ,omitempty 或其他选项
				tagName = strings.Split(tag, ",")[0]
				// 忽略 tag 中 "-" 的字段
				if tagName == "-" {
					continue
				}
			}
		}

		elems := Field{
			Name:    fieldName,
			TagName: tagName,
			Type:    typeStr,
		}
		if isPtrObj {
			elems.IsPtrObj = true
		} else if isObj {
			elems.IsObj = true
		}
		currentStruct.Fields = append(currentStruct.Fields, elems)
	}
}

func reflectTag(tagString, key string) (string, bool) {
	parts := strings.FieldsFunc(tagString, func(r rune) bool {
		return r == ' '
	})
	for _, part := range parts {
		if strings.HasPrefix(part, key+":") {
			// 找到 'json:"value"'，提取 value 部分
			value := strings.TrimPrefix(part, key+":")
			if len(value) > 1 && value[0] == '"' && value[len(value)-1] == '"' {
				return value[1 : len(value)-1], true
			}
		}
	}
	return "", false
}
