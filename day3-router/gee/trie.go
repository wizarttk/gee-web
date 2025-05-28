package gee

import (
	"fmt"
	"strings"
)

// node 结构体表示路由树的节点
type node struct {
	pattern  string  // 待匹配路由，完整的路由模式，例如 /p/:lang/doc
	part     string  // 当前节点表示的路由部分，例如 :lang
	children []*node // 子节点切片，表示该节点的下级路由，例如 [doc, tutorial, intro]
	isWild   bool    // 标识该节点是否是通配符节点（动态节点），part 含有 : 或 * 时为true(通配节点)
}

// String方法 用于返回节点的字符串表示，便于调试和日志记录。
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// TODO:额外的，教程里没有讲的
//
// travel 用于遍历Trie树的所有节点，并将满足条件的节点（即那些有非空pattern的节点）添加到一个列表中
// n 是当前的节点     list是一个指向节点切片的指针，允许在函数内部直接修改传入的切片
// 使用指针的原因是为了能够在函数内部修改原始切片，而不是仅仅在函数内对其进行操作的副本。
// 1. 切片是引用类型，意味着切片本身包含指向底层数组的指针、长度和容量。当你将切片作为参数传递给函数时，实际上是传递了一个切片的副本。
// 2. 这个副本仍然指向相同的底层数组，但它的长度和容量是独立的。
// 3. 在函数中内部append扩展这个切片时，如果扩展导致分配新的底层数组，会改变副本指向的底层数组，但不会改变原切片指向的底层数组
// 4. 使用切片的指针可以避免这中情况，通过传递指向切片的指针，在函数内部可以使用 *list 直接修改原始切片，而不是操作其副本。
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)

	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert方法用于向前缀树插入一个新的路由（将一个路由规则拆分成各个部分parts后插入前缀树）
// - pattern：完整的路由路径，例如 "/p/:lang/doc"
// - parts：将路由路径按"/"分割后的各部分,例如 ["", "p", "lang", "doc"]
// - height：当前处理的深度or当前处理路径部分的索引（从0开始）
func (n *node) insert(pattern string, parts []string, height int) {
	// 递归终止条件：当 height 等于 parts的长度时，表示已处理完所有路由部分，设置节点的 pattern。
	if len(parts) == height {
		n.pattern = pattern // 将当前节点的pattern设置为pattern,表示该节点对应一个完整的路由模式
		return              // 结束递归
	}

	part := parts[height] // 取出当前需要处理的路径部分
	// 查找当前节点的子节点中是否有匹配当前路由部分的节点。
	child := n.matchChild(part)
	if child == nil {
		// 如果没有匹配的子节点，创建一个新的子节点。
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'} // 判断是否为动态节点。
		n.children = append(n.children, child)                              // 新创建的节点添加到当前节点的子节点列表中。
	}
	// 递归插入子节点。
	child.insert(pattern, parts, height+1)
}

// search 方法用于在前缀树中查找与给定路径匹配的节点。
// - parts：将请求路径按 "/" 分割后的各部分。
// - height：当前处理的深度（从0开始）。
// 返回值：匹配的节点（如果存在）和路由参数的映射。
func (n *node) search(parts []string, height int) *node {
	// 递归终止条件：如果已处理完所有路由部分，或当前节点为通配符节点，且存在完整的路由匹配。
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// 获取当前路由部分。
	part := parts[height]
	// 查找当前节点的所有子节点中与当前路由部分匹配的节点。
	children := n.matchChildren(part)

	for _, child := range children {
		// 递归搜索子节点。
		result := child.search(parts, height+1)
		if result != nil { // 如果找到匹配节点，立即返回
			return result
		}
	}

	return nil // 遍历完所有可能的节点未找到匹配，返回nil
}
