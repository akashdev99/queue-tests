package queueimpl8

import (
    "fmt"
    "sync"
)

// Keeping below as var so it is possible to run the slice size bench tests with no coding changes.
var (
    // firstSliceSize holds the size of the first slice.
    firstSliceSize = 1

    // maxFirstSliceSize holds the maximum size of the first slice.
    maxFirstSliceSize = 16

    // maxInternalSliceSize holds the maximum size of each internal slice.
    maxInternalSliceSize = 128
)

// Queue represents an unbounded, dynamically growing FIFO queue.
// The zero value for queue is an empty queue ready to use.
type Queue[T any] struct {
    // Head points to the first node of the linked list.
    head *Node[T]

    // Tail points to the last node of the linked list.
    // In an empty queue, head and tail points to the same node.
    tail *Node[T]

    // Hp is the index pointing to the current first element in the queue
    // (i.e. first element added in the current queue values).
    hp int

    // Tp is the index pointing to the current last element in the queue
    // (i.e. last element added in the current queue values).
    tp int

    // Len holds the current queue values length.
    len int

    // lastSliceSize holds the size of the last created internal slice.
    lastSliceSize int

    //mutex locking to make queue operations thread safe
    mu sync.RWMutex
}

// Node represents a queue node.
// Each node holds a slice of user managed values.
type Node[T any] struct {
    // v holds the list of user added values in this node.
    v []T

    // n points to the next node in the linked list.
    n *Node[T]
}

// NewQueue returns an initialized queue.
func NewQueue[T any]() *Queue[T] {
    return new(Queue[T]).Init()
}

// Init initializes or clears queue q.
func (q *Queue[T]) Init() *Queue[T] {
    q.head = nil
    q.tail = nil

    q.hp = 0
    q.tp = 0

    q.len = 0
    return q
}

// Length returns the number of elements of queue q.
// The complexity is O(1).
func (q *Queue[T]) Length() int {
    q.mu.RLock()
    defer q.mu.RUnlock()
    length := q.len
    return length
}

// Front returns the first element of queue q or nil if the queue is empty.
// The second, bool result indicates whether a valid value was returned;
//
//  if the queue is empty, false will be returned.
//
// The complexity is O(1).
func (q *Queue[T]) Front() (interface{}, bool) {
    q.mu.RLock()
    headNode := q.head
    q.mu.RUnlock()

    if headNode == nil {
        return nil, false
    }

    return headNode.v[q.hp], true
}

// Push adds a value to the queue.
// The complexity is O(1).
func (q *Queue[T]) Push(v T) {
    q.mu.Lock()
    defer q.mu.Unlock()
    if q.head == nil {
        h := newNode[T](firstSliceSize)
        q.head = h
        q.tail = h
        q.lastSliceSize = maxFirstSliceSize
    } else if len(q.tail.v) >= q.lastSliceSize {
        n := newNode[T](maxInternalSliceSize)
        q.tail.n = n
        q.tail = n
        q.lastSliceSize = maxInternalSliceSize
    }

    q.tail.v = append(q.tail.v, v)
    q.len++
}

// Bulk push
func (q *Queue[T]) Enqueue(v []T) {
    var (
        vStartIndex     int
        vEndIndex       int
        vCount          = len(v)
        vRemainingItems = vCount
    )

    q.mu.Lock()
    defer q.mu.Unlock()
    //Push head node
    if q.head == nil {
        headNode := newNode[T](firstSliceSize)
        q.head = headNode
        q.tail = headNode
        q.lastSliceSize = maxFirstSliceSize

        if vCount > q.lastSliceSize {
            //setting max limit
            vEndIndex = q.lastSliceSize
        } else {
            vEndIndex = vRemainingItems
        }
    } else if q.tp < q.lastSliceSize {
        //If tail pointer is not at max slice size , then there is more space in node for adding
        //Fill node with available space before adding new Node
        availableSpace := q.lastSliceSize - q.tp

        if vCount > availableSpace {
            vEndIndex = availableSpace
        } else {
            vEndIndex = vRemainingItems
        }
    }

    if vStartIndex != vEndIndex {
        vRemainingItems = q.appendToTail(v, vStartIndex, vEndIndex, vCount)
    }

    //Add new node if more items remaining
    //Repeat until no remaining items
    for vRemainingItems > 0 {
        q.tp = 0
        node := newNode[T](maxInternalSliceSize)
        q.tail.n = node
        q.tail = node
        q.lastSliceSize = maxInternalSliceSize
        vStartIndex = vEndIndex

        if vRemainingItems > maxInternalSliceSize {
            vEndIndex += maxInternalSliceSize
        } else {
            vEndIndex += vRemainingItems
        }

        vRemainingItems = q.appendToTail(v, vStartIndex, vEndIndex, vCount)
    }
    q.len += vCount
}

func (q *Queue[T]) appendToTail(v []T, vStartIndex int, vEndIndex int, length int) (remaining int) {
    q.tail.v = append(q.tail.v, v[vStartIndex:vEndIndex]...)
    q.tp += vEndIndex - vStartIndex
    return length - vEndIndex
}

func (q *Queue[T]) Dequeue(count int) []T {
    q.mu.Lock()
    defer q.mu.Unlock()

    if q.head == nil {
        return nil
    }

    var (
        items              = make([]T, 0, count)
        vStartIndex        int
        vEndIndex          int
        maxSliceSize       int
        availableSliceSize int
    )

    for count > 0 && q.head != nil {
        maxSliceSize = len(q.head.v)
        availableSliceSize = maxSliceSize - q.hp
        vStartIndex = q.hp

        //If count of items to be popped is less than available values on node
        //Only pop up to count required
        if count < availableSliceSize {
            vEndIndex = q.hp + count
            q.hp = vEndIndex
            items = append(items, q.head.v[vStartIndex:vEndIndex]...)
            q.len -= count
            count = 0
            break
        }

        vEndIndex = q.hp + availableSliceSize
        items = append(items, q.head.v[vStartIndex:vEndIndex]...)
        count -= availableSliceSize

        //move to next node
        n := q.head.n
        q.head.n = nil // Avoid memory leaks
        q.head = n
        q.hp = 0
        q.len -= (vEndIndex - vStartIndex)
    }
    return items
}

// Pop retrieves and removes the current element from the queue.
// The second, bool result indicates whether a valid value was returned;
//
//  if the queue is empty, false will be returned.
//
// The complexity is O(1).

////////////////////////////////////////////////////////////////////////////////////////////////
//WARNING -  Single pop for generics only added for benchmarking , not for implementation
////////////////////////////////////////////////////////////////////////////////////////////////

func (q *Queue[T]) Pop() (T, bool) {
    q.mu.Lock()
    defer q.mu.Unlock()
    if q.head == nil {
        return *new(T), false
    }

    v := q.head.v[q.hp]
    // q.head.v[q.hp] = T{} // Avoid memory leaks
    q.len--
    q.hp++
    if q.hp >= len(q.head.v) {
        n := q.head.n
        q.head.n = nil // Avoid memory leaks
        q.head = n
        q.hp = 0
    }
    return v, true
}

// Print Queue Content
func (q *Queue[T]) Print() {
    q.mu.RLock()
    defer q.mu.RUnlock()
    node := q.head
    headNode := true
    if node == nil {
        fmt.Println("Queue is Empty")
        return
    }
    for node != nil {
        if headNode {
            //First node/Head node start from head pointer
            fmt.Printf("Slice %v , Length = %v \n ", node.v[q.hp:], len(node.v)-q.hp)
        } else {
            fmt.Printf("Slice %v , Length = %v \n ", node.v, len(node.v))
        }

        node = node.n
        headNode = false
    }
}

// Number of nodes in Linked List
func (q *Queue[T]) NoOfNodes() (count int) {
    q.mu.RLock()
    defer q.mu.RUnlock()
    node := q.head
    for node != nil {
        count++
        node = node.n
    }
    return count
}

// newNode returns an initialized node.
func newNode[T any](capacity int) *Node[T] {
    return &Node[T]{
        v: make([]T, 0, capacity),
    }
}
