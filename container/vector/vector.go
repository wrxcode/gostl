package vector

import (
	"errors"
	"fmt"
	. "github.com/liyue201/gostl/container"
)

var ErrOutOffRange = errors.New("out off range")
var ErrEmpty = errors.New("vector is empty")
var ErrInvalidIterator = errors.New("invalid iterator")

type Vector struct {
	data []interface{}
}

func New(capacity int) *Vector {
	return &Vector{data: make([]interface{}, 0, capacity)}
}

func NewFromVector(other *Vector) *Vector {
	this := &Vector{data: make([]interface{}, other.Size(), other.Capacity())}
	for i := range other.data {
		this.data[i] = other.data[i]
	}
	return this
}

func (this *Vector) Size() int {
	return len(this.data)
}

func (this *Vector) Capacity() int {
	return cap(this.data)
}

func (this *Vector) Empty() bool {
	if len(this.data) == 0 {
		return true
	}
	return false
}

func (this *Vector) PushBack(val interface{}) {
	this.data = append(this.data, val)
}

func (this *Vector) SetAt(index int, val interface{}) error {
	if index < 0 || index >= this.Size() {
		return ErrOutOffRange
	}
	this.data[index] = val
	return nil
}

func (this *Vector) InsertAt(index int, val interface{}) error {
	if index < 0 || index > this.Size() {
		return ErrOutOffRange
	}
	this.data = append(this.data, val)
	for i := len(this.data) - 1; i > index; i-- {
		this.data[i] = this.data[i-1]
	}
	this.data[index] = val
	return nil
}

func (this *Vector) EraseAt(index int) error {
	return this.EraseIndexRange(index, index+1)
}

func (this *Vector) EraseIndexRange(first, last int) error {
	if first > last {
		return nil
	}
	if first < 0 || last > this.Size() {
		return ErrOutOffRange
	}

	left := this.data[:first]
	right := this.data[last:]
	this.data = append(left, right...)
	return nil
}

//At returns the value at index, returns nil if index out off range .
func (this *Vector) At(index int) interface{} {
	if index < 0 || index > this.Size() {
		return nil
	}
	return this.data[index]
}

//At returns the first value of the vector, returns nil if the vector is empty.
func (this *Vector) Front() interface{} {
	return this.At(0)
}

//At returns the last value of the vector, returns nil if the vector is empty.
func (this *Vector) Back() interface{} {
	return this.At(this.Size() - 1)
}

//At returns the last value of the vector and erase it, returns nil if the vector is empty.
func (this *Vector) PopBack() interface{} {
	if this.Empty() {
		return nil
	}
	val := this.Back()
	this.data = this.data[:len(this.data)-1]
	return val
}

func (this *Vector) Reserve(capacity int) {
	if cap(this.data) >= capacity {
		return
	}
	data := make([]interface{}, this.Size(), capacity)
	for i := 0; i < len(this.data); i++ {
		data[i] = this.data[i]
	}
	this.data = data
}

func (this *Vector) ShrinkToFit() {
	if len(this.data) == cap(this.data) {
		return
	}
	len := this.Size()
	data := make([]interface{}, len, len)
	for i := 0; i < len; i++ {
		data[i] = this.data[i]
	}
	this.data = data
}

func (this *Vector) Clear() {
	this.data = this.data[:0]
}

func (this *Vector) Data() [] interface{} {
	return this.data
}

func (this *Vector) Begin() Iterator {
	return &VectorIterator{vec: this, curIndex: 0}
}

func (this *Vector) End() Iterator {
	return &VectorIterator{vec: this, curIndex: this.Size()}
}

func (this *Vector) RBegin() ReverseIterator {
	return &VectorReverseIterator{vec: this, curIndex: this.Size() - 1}
}

func (this *Vector) REnd() ReverseIterator {
	return &VectorReverseIterator{vec: this, curIndex: -1}
}

func (this *Vector) Insert(iter Iterator, val interface{}) Iterator {
	this.InsertAt(iter.(*VectorIterator).curIndex, val)
	return iter
}

func (this *Vector) Erase(iter Iterator) Iterator {
	this.EraseAt(iter.(*VectorIterator).curIndex)
	return iter
}

func (this *Vector) EraseRange(first, last Iterator) Iterator {
	from := first.(*VectorIterator).curIndex
	to := last.(*VectorIterator).curIndex
	this.EraseIndexRange(from, to)
	return &VectorIterator{vec: this, curIndex: from}
}

func (this *Vector) Resize(size int) {
	if size >= this.Size() {
		return
	}
	this.data = this.data[:size]
}

func (this *Vector) Swap(other *Vector) {
	this.data, other.data = other.data, this.data
}

func (this *Vector) String() string {
	return fmt.Sprintf("%v", this.data)
}

type VectorIterator struct {
	vec      *Vector
	curIndex int
}

func (this *VectorIterator) Next() Iterator {
	index := this.curIndex + 1
	if index > this.vec.Size() {
		index = this.vec.Size()
	}
	return &VectorIterator{vec: this.vec, curIndex: index}
}

func (this *VectorIterator) Value() interface{} {
	val := this.vec.At(this.curIndex)
	return val
}

func (this *VectorIterator) Set(val interface{}) error {
	return this.vec.SetAt(this.curIndex, val)
}

func (this *VectorIterator) Equal(other Iterator) bool {
	otherItr, ok := other.(*VectorIterator)
	if !ok {
		return false
	}
	if this.vec == otherItr.vec && otherItr.curIndex == this.curIndex {
		return true
	}
	return false
}

type VectorReverseIterator struct {
	vec      *Vector
	curIndex int
}

func (this *VectorReverseIterator) Next() ReverseIterator {
	index := this.curIndex - 1
	if index < -1 {
		index = -1
	}
	return &VectorReverseIterator{vec: this.vec, curIndex: index}
}

func (this *VectorReverseIterator) Set(val interface{}) error {
	return this.vec.SetAt(this.curIndex, val)
}

func (this *VectorReverseIterator) Value() interface{} {
	return this.vec.At(this.curIndex)
}

func (this *VectorReverseIterator) Equal(other ReverseIterator) bool {
	otherItr, ok := other.(*VectorReverseIterator)
	if !ok {
		return false
	}
	if this.vec == otherItr.vec && otherItr.curIndex == this.curIndex {
		return true
	}
	return false
}