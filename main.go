package main

import (
	"fmt"

	"github.com/EletricSaw/rbtree/rbtree"
)

func main() {
	tree := rbtree.New[string, string]()

	// 샘플 데이터 삽입
	samples := []struct {
		key   string
		value string
	}{
		{"k", "카카오"},
		{"g", "구글"},
		{"a", "애플"},
		{"n", "네이버"},
		{"b", "배달의민족"},
	}
	for _, s := range samples {
		tree.Insert(s.key, s.value)
	}

	// 하나를 삭제해본다.
	tree.Delete("g")

	// 검색 예시
	if node := tree.Search("n"); node != nil {
		fmt.Printf("key %q => %v\n", node.Key, node.Value)
	}

	fmt.Println("\n=== RBTree 구조 ===")
	tree.PrintStdout()
}
