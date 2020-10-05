package hw04_lru_cache //nolint:golint,stylecheck
import (
	"fmt"
)

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
	// Place your code here
}

type list struct {
	Item      []ListItem
	FirstItem *ListItem
	LastItem  *ListItem
	Length    int
	// Place your code here
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *ListItem {
	return l.FirstItem
}

func (l *list) Back() *ListItem {
	return l.LastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := ListItem{}
	switch i := v.(type) {
	case int, cacheItem:
		newItem.Value = i
		newItem.Prev = nil
		newItem.Next = l.FirstItem
		if newItem.Next != nil {
			l.FirstItem.Prev = &newItem
		} else {
			l.LastItem = &newItem
		}
		l.FirstItem = &newItem
		l.Length++
	default:
		fmt.Printf("Unknown type %T!\n", v)
	}
	return &newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := ListItem{}
	switch i := v.(type) {
	case int, cacheItem:
		newItem.Value = i
		newItem.Next = nil
		newItem.Prev = l.LastItem
		if newItem.Prev != nil {
			l.LastItem.Next = &newItem
		} else {
			l.FirstItem = &newItem
		}
		l.LastItem = &newItem
		l.Length++
	default:
		fmt.Printf("\nUnknown type %T!\n", v)
	}
	return &newItem
}

func (l *list) Remove(i *ListItem) {
	if i != nil {
		switch {
		case i.Prev != nil && i.Next != nil:
			i.Prev.Next = i.Next
			i.Next.Prev = i.Prev
		case i.Prev == nil && i.Next != nil:
			l.FirstItem = i.Next.Prev
		case i.Prev != nil && i.Next == nil:
			l.LastItem = i.Prev
			l.LastItem.Next = nil
		case i.Prev == nil && i.Next == nil:
			l.FirstItem = nil
			l.LastItem = nil
		default:
			fmt.Printf("Unknown condition")
		}
		l.Length--
		return
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i != nil {
		if i == l.FirstItem {
			return
		}
		l.PushFront(i.Value)
		l.Remove(i)
	}
}

func NewList() List {
	return &list{}
}
