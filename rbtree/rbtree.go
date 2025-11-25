package rbtree

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// 아래 구현은 CLRS 교과서에 나오는 레드-블랙 트리(RBTree)를 그대로 옮긴 것이다.
// RBTree는 다음 규칙을 강제함으로써 삽입/조회/삭제를 모두 O(log n)에 수행한다.
//
//   1) 모든 노드는 빨강 또는 검정이다.
//   2) 루트는 항상 검정이다.
//   3) 어떤 노드에서 잎까지 내려가는 경로에도 빨강-빨강이 연속해서 나타나지 않는다.
//   4) 한 노드에서 내려가는 모든 경로는 같은 수의 검정 노드를 지난다(black height).
//
// 코드 곳곳에 과한 수준의 주석을 달아 알고리즘 흐름을 공부하기 쉽도록 했다.

// Color는 노드의 색 상태를 표현한다.
type Color bool

const (
	red   Color = true
	black Color = false
)

// Node는 트리의 한 정점을 표현한다. 실무 구현에서는 NIL 센티넬을 별도로 두지만,
// 여기서는 이해를 돕기 위해 nil 포인터를 잎으로 간주하고 보정 과정에서 검정으로 취급한다.
type Node struct {
	Key    string
	Value  interface{}
	Color  Color
	Parent *Node
	Left   *Node
	Right  *Node
}

// Tree 구조체는 루트 포인터와 원소 수를 추적하는 래퍼이다. 이 구조체에 연산 메서드를 붙여
// 회전/보정과 같은 내부 구현을 숨기고 API만 노출한다.
type Tree struct {
	root *Node
	size int
}

// New는 빈 RBTree를 만든다.
func New() *Tree {
	return &Tree{}
}

// Size는 현재 저장된 키 개수를 돌려준다.
func (t *Tree) Size() int {
	return t.size
}

// Root는 테스트나 예제에서 구조를 살펴볼 수 있도록 루트 포인터를 돌려준다.
func (t *Tree) Root() *Node {
	return t.root
}

// Search는 키를 가진 노드를 찾아 돌려준다. 일반적인 BST 탐색이므로 트리 구조를 바꾸지 않는다.
func (t *Tree) Search(key string) *Node {
	cur := t.root
	for cur != nil {
		switch {
		case key < cur.Key:
			cur = cur.Left
		case key > cur.Key:
			cur = cur.Right
		default:
			return cur
		}
	}
	return nil
}

// Insert는 키를 삽입한다. 단순화를 위해 중복 키는 무시하지만, 필요하다면 갯수 누적 등의 동작으로 확장할 수 있다.
func (t *Tree) Insert(key string, value interface{}) {
	var parent *Node
	cur := t.root

	// 먼저 일반 BST 삽입을 통해 부모 위치를 찾는다.
	for cur != nil {
		parent = cur
		switch {
		case key < cur.Key:
			cur = cur.Left
		case key > cur.Key:
			cur = cur.Right
		default:
			// 이미 존재하는 키면 값을 갱신하고 종료한다.
			cur.Value = value
			return
		}
	}

	// 삽입 노드는 항상 빨강으로 시작한다. 검정으로 넣으면 규칙 (4)가 깨질 수 있다.
	node := &Node{Key: key, Value: value, Color: red, Parent: parent}
	if parent == nil {
		t.root = node
	} else if node.Key < parent.Key {
		parent.Left = node
	} else {
		parent.Right = node
	}

	// 구조적 삽입 뒤 망가졌을 수 있는 규칙을 insertFixup으로 복원한다.
	t.insertFixup(node)
	t.size++
}

// Delete는 주어진 키를 삭제한다. 검정 노드를 제거하면 규칙 (2)(4)가 깨질 수 있으므로
// double black 개념을 사용해 위로 전파하면서 복구한다.
func (t *Tree) Delete(key string) bool {
	node := t.Search(key)
	if node == nil {
		return false
	}

	originalColor := node.Color
	var x, replacementParent *Node

	switch {
	case node.Left == nil:
		x = node.Right
		replacementParent = node.Parent
		t.transplant(node, node.Right)
	case node.Right == nil:
		x = node.Left
		replacementParent = node.Parent
		t.transplant(node, node.Left)
	default:
		// 후속 노드는 오른쪽 서브트리에서 가장 작은 값이다.
		successor := minimum(node.Right)
		originalColor = successor.Color
		x = successor.Right
		if successor.Parent == node {
			if x != nil {
				x.Parent = successor
			}
			replacementParent = successor
		} else {
			replacementParent = successor.Parent
			t.transplant(successor, successor.Right)
			successor.Right = node.Right
			successor.Right.Parent = successor
		}
		t.transplant(node, successor)
		successor.Left = node.Left
		successor.Left.Parent = successor
		successor.Color = node.Color
	}

	if originalColor == black {
		t.deleteFixup(x, replacementParent)
	}
	t.size--
	return true
}

// InOrder는 키를 정렬 순서대로 순회하며 fn을 호출한다. 테스트에서 구조를 확인할 때 유용하다.
func (t *Tree) InOrder(fn func(key string, value interface{})) {
	inOrder(t.root, fn)
}

// Print은 트리 구조를 들여쓰기 형태로 출력한다. w가 nil이면 stdout으로 대체한다.
func (t *Tree) Print(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	if t.root == nil {
		fmt.Fprintln(w, "(empty)")
		return
	}
	printNode(w, t.root, 0)
}

// PrintStdout은 편의를 위해 stdout으로 바로 출력한다.
func (t *Tree) PrintStdout() {
	t.Print(os.Stdout)
}

// insertFixup은 삽입으로 깨진 RB 규칙을 되돌린다. 빨강 부모-자식이 없어질 때까지 색을 바꾸거나 회전한다.
func (t *Tree) insertFixup(node *Node) {
	for node != t.root && colorOf(node.Parent) == red {
		if node.Parent == node.Parent.Parent.Left {
			uncle := node.Parent.Parent.Right
			switch colorOf(uncle) {
			case red:
				// Case 1: 부모와 삼촌이 모두 빨강이면 둘 다 검정으로 바꾸고 할아버지를 빨강으로 올린다.
				node.Parent.Color = black
				uncle.Color = black
				node.Parent.Parent.Color = red
				node = node.Parent.Parent
			default:
				if node == node.Parent.Parent.Right {
					// Case 2: 현재 노드가 오른쪽 자식이면 회전해서 Case 3으로 만들어 준다.
					node = node.Parent
					t.rotateLeft(node)
				}
				// Case 3: 현재 노드가 왼쪽 자식이므로 부모-할아버지 색을 뒤집고 오른쪽 회전한다.
				node.Parent.Color = black
				node.Parent.Parent.Color = red
				t.rotateRight(node.Parent.Parent)
			}
		} else {
			// 왼쪽/오른쪽만 뒤바꾼 대칭 케이스.
			uncle := node.Parent.Parent.Left
			switch colorOf(uncle) {
			case red:
				node.Parent.Color = black
				uncle.Color = black
				node.Parent.Parent.Color = red
				node = node.Parent.Parent
			default:
				if node == node.Parent.Left {
					node = node.Parent
					t.rotateRight(node)
				}
				node.Parent.Color = black
				node.Parent.Parent.Color = red
				t.rotateLeft(node.Parent.Parent)
			}
		}
	}
	t.root.Color = black
}

// deleteFixup은 검정 노드 삭제 후 생기는 double black을 제거한다.
// x가 nil일 수도 있으므로 parent를 함께 넘겨 nil 역참조를 피한다.
func (t *Tree) deleteFixup(x, parent *Node) {
	for (x != t.root) && colorOf(x) == black {
		if x == leftOf(parent) {
			sibling := rightOf(parent)
			if colorOf(sibling) == red {
				sibling.Color = black
				parent.Color = red
				t.rotateLeft(parent)
				sibling = rightOf(parent)
			}
			if colorOf(sibling.Left) == black && colorOf(sibling.Right) == black {
				sibling.Color = red
				x = parent
				parent = x.Parent
			} else {
				if colorOf(sibling.Right) == black {
					if sibling.Left != nil {
						sibling.Left.Color = black
					}
					sibling.Color = red
					t.rotateRight(sibling)
					sibling = rightOf(parent)
				}
				sibling.Color = colorOf(parent)
				parent.Color = black
				if sibling.Right != nil {
					sibling.Right.Color = black
				}
				t.rotateLeft(parent)
				x = t.root
				parent = nil
			}
		} else {
			sibling := leftOf(parent)
			if colorOf(sibling) == red {
				sibling.Color = black
				parent.Color = red
				t.rotateRight(parent)
				sibling = leftOf(parent)
			}
			if colorOf(sibling.Left) == black && colorOf(sibling.Right) == black {
				sibling.Color = red
				x = parent
				parent = x.Parent
			} else {
				if colorOf(sibling.Left) == black {
					if sibling.Right != nil {
						sibling.Right.Color = black
					}
					sibling.Color = red
					t.rotateLeft(sibling)
					sibling = leftOf(parent)
				}
				sibling.Color = colorOf(parent)
				parent.Color = black
				if sibling.Left != nil {
					sibling.Left.Color = black
				}
				t.rotateRight(parent)
				x = t.root
				parent = nil
			}
		}
	}
	if x != nil {
		x.Color = black
	}
}

// rotateLeft는 노드를 오른쪽 자식과 회전시킨다. 포인터만 바뀌므로 O(1)이다.
func (t *Tree) rotateLeft(node *Node) {
	right := node.Right
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Parent = node.Parent
	if node.Parent == nil {
		t.root = right
	} else if node == node.Parent.Left {
		node.Parent.Left = right
	} else {
		node.Parent.Right = right
	}
	right.Left = node
	node.Parent = right
}

// rotateRight는 rotateLeft의 좌우 대칭이다.
func (t *Tree) rotateRight(node *Node) {
	left := node.Left
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Parent = node.Parent
	if node.Parent == nil {
		t.root = left
	} else if node == node.Parent.Right {
		node.Parent.Right = left
	} else {
		node.Parent.Left = left
	}
	left.Right = node
	node.Parent = left
}

// transplant는 서브트리 u 자리에 v를 끼워 넣는다. 삭제 과정에서 부모 포인터를 깔끔하게 유지하기 위한 헬퍼다.
func (t *Tree) transplant(u, v *Node) {
	if u.Parent == nil {
		t.root = v
	} else if u == u.Parent.Left {
		u.Parent.Left = v
	} else {
		u.Parent.Right = v
	}
	if v != nil {
		v.Parent = u.Parent
	}
}

// 헬퍼 함수들 ---------------------------------------------------------------

func colorOf(node *Node) Color {
	if node == nil {
		return black
	}
	return node.Color
}

func leftOf(node *Node) *Node {
	if node == nil {
		return nil
	}
	return node.Left
}

func rightOf(node *Node) *Node {
	if node == nil {
		return nil
	}
	return node.Right
}

func minimum(node *Node) *Node {
	for node.Left != nil {
		node = node.Left
	}
	return node
}

func inOrder(node *Node, fn func(string, interface{})) {
	if node == nil {
		return
	}
	inOrder(node.Left, fn)
	fn(node.Key, node.Value)
	inOrder(node.Right, fn)
}

func printNode(w io.Writer, node *Node, depth int) {
	if node == nil {
		return
	}
	printNode(w, node.Right, depth+1)
	indent := strings.Repeat("  ", depth)
	fmt.Fprintf(w, "%s[%s] %s => %v\n", indent, colorString(node.Color), node.Key, node.Value)
	printNode(w, node.Left, depth+1)
}

func colorString(c Color) string {
	if c == red {
		return "R"
	}
	return "B"
}
