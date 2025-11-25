package rbtree

import (
	"bytes"
	"cmp"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestInsertAndSearch(t *testing.T) {
	tree := New[string, int]()
	values := []struct {
		key   string
		value int
	}{
		{"10", 10},
		{"3", 3},
		{"15", 15},
		{"7", 7},
		{"20", 20},
		{"1", 1},
		{"5", 5},
	}

	for _, v := range values {
		tree.Insert(v.key, v.value)
	}

	if tree.Size() != len(values) {
		t.Fatalf("expected size %d, got %d", len(values), tree.Size())
	}

	for _, v := range values {
		node := tree.Search(v.key)
		if node == nil || node.Key != v.key {
			t.Fatalf("missing key %q after insert", v.key)
		}
		if got := node.Value; got != v.value {
			t.Fatalf("key %q expected value %d got %d", v.key, v.value, got)
		}
	}

	tree.InOrder(func(key string, value int) {
		node := tree.Search(key)
		if node == nil {
			t.Fatalf("InOrder returned key %q but Search failed", key)
		}
		if node.Value != value {
			t.Fatalf("value mismatch for key %q", key)
		}
	})
}

func TestDelete(t *testing.T) {
	tree := New[string, string]()
	values := []string{"20", "15", "25", "10", "18", "8", "12", "16", "19"}

	for _, k := range values {
		tree.Insert(k, k)
	}

	toDelete := []string{"10", "15", "20"}
	for _, k := range toDelete {
		if !tree.Delete(k) {
			t.Fatalf("expected delete(%q) to succeed", k)
		}
		if tree.Search(k) != nil {
			t.Fatalf("key %q still searchable after delete", k)
		}
	}

	expectedSize := len(values) - len(toDelete)
	if tree.Size() != expectedSize {
		t.Fatalf("expected size %d, got %d", expectedSize, tree.Size())
	}

	assertRBProperties(t, tree)
}

func TestRBPropertiesRandom(t *testing.T) {
	tree := New[string, int]()
	const count = 1000
	var inserted []string
	seen := make(map[string]struct{})

	for i := 0; i < count; i++ {
		val := rand.Intn(10_000)
		key := strconv.Itoa(val)
		tree.Insert(key, val)
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			inserted = append(inserted, key)
		}
		assertRBProperties(t, tree)
	}

	rand.Shuffle(len(inserted), func(i, j int) {
		inserted[i], inserted[j] = inserted[j], inserted[i]
	})
	for _, k := range inserted[:len(inserted)/2] {
		if !tree.Delete(k) {
			t.Fatalf("delete(%q) unexpectedly failed", k)
		}
		assertRBProperties(t, tree)
	}

	var got []string
	tree.InOrder(func(key string, value int) {
		got = append(got, key)
	})
	if !sort.StringsAreSorted(got) {
		t.Fatalf("in-order walk not sorted: %v", got)
	}
}

func TestPrint(t *testing.T) {
	tree := New[string, int]()
	tree.Insert("b", 2)
	tree.Insert("a", 1)
	tree.Insert("c", 3)

	var buf bytes.Buffer
	tree.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "[B] b => 2") {
		t.Fatalf("expected root line, got %q", out)
	}
	if !strings.Contains(out, "[R] a => 1") {
		t.Fatalf("expected left child line, got %q", out)
	}

	var emptyBuf bytes.Buffer
	New[string, int]().Print(&emptyBuf)
	if strings.TrimSpace(emptyBuf.String()) != "(empty)" {
		t.Fatalf("empty tree should print (empty), got %q", emptyBuf.String())
	}
}

func assertRBProperties[K cmp.Ordered, V any](t *testing.T, tree *Tree[K, V]) {
	t.Helper()
	root := tree.Root()
	if root == nil {
		return
	}
	if root.Color != black {
		t.Fatalf("root must be black, got %v", root.Color)
	}
	checkNoRedRed(t, root)
	expectedBlackHeight := blackHeight(root)
	verifyBlackHeight(t, root, expectedBlackHeight, 0)
}

func checkNoRedRed[K cmp.Ordered, V any](t *testing.T, node *Node[K, V]) {
	if node == nil {
		return
	}
	if node.Color == red {
		if colorOf(node.Left) == red || colorOf(node.Right) == red {
			t.Fatalf("red node %v has red child", node.Key)
		}
	}
	checkNoRedRed(t, node.Left)
	checkNoRedRed(t, node.Right)
}

func blackHeight[K cmp.Ordered, V any](node *Node[K, V]) int {
	height := 0
	for node != nil {
		if node.Color == black {
			height++
		}
		node = node.Left
	}
	return height
}

func verifyBlackHeight[K cmp.Ordered, V any](t *testing.T, node *Node[K, V], expected, current int) {
	if node == nil {
		if current != expected {
			t.Fatalf("black height mismatch: expected %d got %d", expected, current)
		}
		return
	}
	if node.Color == black {
		current++
	}
	verifyBlackHeight(t, node.Left, expected, current)
	verifyBlackHeight(t, node.Right, expected, current)
}
