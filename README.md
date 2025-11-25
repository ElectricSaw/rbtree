# rbtree

교육용 레드-블랙 트리(RBTree) 구현입니다. CLRS 알고리즘을 Go로 옮겨와 삽입/삭제/검색/중위순회와
함께, 트리 구조를 자체적으로 출력하는 기능까지 제공합니다. 모든 핵심 단계마다 한글 주석을 붙여
자료구조 동작을 공부하기 쉽도록 했습니다.

## 특징

- **제네릭 타입 지원**: Go 1.18+ 제네릭을 사용하여 모든 정렬 가능한 타입(int, string, float64 등)을 키로 사용 가능
- **타입 안전성**: 컴파일 타임에 타입 체크를 통해 런타임 오류 방지
- **어디서든 사용 가능**: 다른 프로젝트에서 쉽게 임포트하여 사용할 수 있는 독립적인 패키지
- **간단한 API**: 직관적이고 사용하기 쉬운 인터페이스

## 포함 내용
- `rbtree` 패키지: 제네릭 타입을 지원하는 RBTree 구현 (키: `cmp.Ordered`, 값: `any`)
- 상세 주석과 헬퍼 함수(`insertFixup`, `deleteFixup`, 회전 등)
- `Print`, `PrintStdout`으로 트리 구조 시각화
- 풍부한 테스트(`go test ./...`)로 불변식 검증
- `main.go`: 간단한 샘플 데이터 삽입/삭제/검색 후 트리 출력

## 사용 예시

```go
package main

import (
    "fmt"
    "github.com/EletricSaw/rbtree/rbtree"
)

func main() {
    // 문자열 키, 문자열 값
    tree1 := rbtree.New[string, string]()
    tree1.Insert("apple", "사과")
    tree1.Insert("banana", "바나나")
    
    // 정수 키, 정수 값
    tree2 := rbtree.New[int, int]()
    tree2.Insert(10, 100)
    tree2.Insert(20, 200)
    
    // 검색
    if node := tree1.Search("apple"); node != nil {
        fmt.Printf("Found: %s => %s\n", node.Key, node.Value)
    }
    
    // 삭제
    tree1.Delete("banana")
    
    // 순회
    tree1.InOrder(func(key string, value string) {
        fmt.Printf("%s: %s\n", key, value)
    })
}
```

## 실행 방법
```bash
go test ./...
go run .
```

## 라이선스
[MIT](./LICENSE)