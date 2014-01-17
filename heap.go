//package heap
package main

import (
	"fmt"
	"strconv"
)

//~------------------------------------------------------------------------------------------------
//~INTERFACE DEFINITIONS---------------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
//Interface that combines the StringRep and Comparable interfaces
type Heapable interface {

	//Returns a string representation of the object
	String() string

	//Return 0 for equal
	//Return 1 if this is greater than value2
	//Return -1 if this is less than value2
	CompareTo(value2 interface{}) int
}

//Interface defining heuristic used to compare two types of the type passed into the constructor
type HeapHeuristic interface {

	//Compares the two objects passed
	//Returns 1 if val1 > val2
	//Returns -1 if val1 < val2
	//Returns 0 if val1 == val2
	Compare(val1, val2 Heapable) int
}

//Interface defining a Heap data structure
type Heap interface {

	//Constructor, takes variable and stores its type to ensure type safety
	Heap(heuristic HeapHeuristic)

	//Adds a variable to the heap if its type matches that of the heap's type
	Add(value Heapable)

	//Returns the element at the top of the heap
	Top() Heapable

	//Removes a specified element from the heap
	Remove(value Heapable) bool

	//Pops the element off the top of the heap
	Pop()

	//Does an in-order traversal of the heap and prints it
	PrintInOrder()
}


//~------------------------------------------------------------------------------------------------
//~STRUCT DECLARATIONS----------------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
//Node of a binary heap that holds a value and left and right children
type BinaryHeapNode struct {

	value Heapable
	left *BinaryHeapNode
	right *BinaryHeapNode
}


//BinaryHeap implementation
type BinaryHeap struct {

	//The root of the binary heap
	root *BinaryHeapNode

	//The heuristic that this heap is organized by
	heuristic HeapHeuristic
}

//MaxHeap heuristic which matches the HeapHeuristic interface
type MaxHeap struct {}

//~------------------------------------------------------------------------------------------------
//~MaxHeap FUNCTION DECLARATIONS-------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
//The Compare method for the MaxHeap implementation of the HeapHeuristic type
func (heuristic *MaxHeap) Compare(value1, value2 Heapable) int {

	if compareValue := value1.CompareTo(value2); compareValue > 0 {
		return 1
	} else if compareValue < 0 {
		return -1
	} else {
		return 0
	}
}

//~------------------------------------------------------------------------------------------------
//~BinaryHeapNode FUNCTION DECLARATIONS----------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
func (node *BinaryHeapNode) String() string {

	return node.value.String()
}

func (node *BinaryHeapNode) PrintFamily() string {

	returnString := "----"
	if node.right != nil {

		returnString += node.right.value.String()
	} else {

		returnString += "|"
	}

	returnString += "\n"

	returnString += node.value.String()
	returnString += "\n----"

	if node.left != nil {

		returnString += node.left.value.String()
	} else {

		returnString += "|"
	}
	returnString += "\n"

	return returnString
}

//~------------------------------------------------------------------------------------------------
//~BinaryHeap FUNCTION DECLARATIONS----------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
func (heap *BinaryHeap) Heap(heuristic HeapHeuristic) {

	heap.heuristic = heuristic
}

func (heap *BinaryHeap) Add(value Heapable) {

	newNode := new(BinaryHeapNode)
	newNode.value = value

	if heap.root != nil {

		heap.root = heap.AddHelper(heap.root, newNode)

	} else {

		heap.root = newNode
	}
}

//Recursive helper method to add nodes to the heap
func (heap *BinaryHeap) AddHelper(subRoot *BinaryHeapNode, newNode *BinaryHeapNode) *BinaryHeapNode {

	//if newNode belongs above subRoot
	if heap.heuristic.Compare(subRoot.value, newNode.value) < 0 {

		//TODO: Add balancing
		newNode = heap.AddHelper(newNode, subRoot)

		return newNode

	//newNode belongs below subRoot
	} else {

		//Make newNode subRoot's left child
		if subRoot.left == nil {

			subRoot.left = newNode

		//Make newNode subRoot's right child
		} else if subRoot.right == nil {

			subRoot.right = newNode

		//Recurse down
		} else {

			//TODO: Add node depth to BinaryHeapNode to keep heap balanced
			//go left
			subRoot.left = heap.AddHelper(subRoot.left, newNode)

			//go right
			//subRoot.right = heap.AddHelper(subRoot.right, newNode)
		}

		return subRoot
	}
}

func (heap *BinaryHeap) Top() Heapable {

	return heap.root.value
}

func (heap *BinaryHeap) Remove(value Heapable) bool {
 
	returnValue := false
	returnValue, heap.root = heap.RemoveHelper(heap.root, value)
	return returnValue
}

//Recursive helper method for Remove
func (heap *BinaryHeap) RemoveHelper(subRoot *BinaryHeapNode, value Heapable) (bool, *BinaryHeapNode) {

	//Descended to the bottom of a sub-tree
	if subRoot == nil {

		return false, subRoot

	//Descended too far down a sub-tree
	} else if heap.heuristic.Compare(subRoot.value, value) < 0 {

		return false, subRoot

	//Found element
	} else if heap.heuristic.Compare(subRoot.value, value) == 0 {

		return true, heap.PercolateUp(subRoot)		

		//return true, newSubRoot

	//The search continues
	} else {

		returnValue, subRoot := heap.RemoveHelper(subRoot.left, value)
		if !returnValue {

			returnValue, subRoot = heap.RemoveHelper(subRoot.right, value)
		}

		return returnValue, subRoot
	}
}

//Helper function to maintain heap structure in the case of Remove or Pop being called
//Takes the node being removed
func (heap *BinaryHeap) PercolateUp(subRoot *BinaryHeapNode) (*BinaryHeapNode) {

	if subRoot == nil {

		return nil
	}

	newSubRoot := subRoot

	//Left sub-tree exists
	if subRoot.left != nil && subRoot.right == nil {

		newSubRoot = subRoot.left

	//Right sub-tree exists
	} else if subRoot.left == nil && subRoot.right != nil {

		newSubRoot = subRoot.right

	//Both sub-trees exist
	} else if subRoot.left != nil && subRoot.right != nil {

		//Left child has a value greater than or equal to the right
		if heap.heuristic.Compare(subRoot.left.value, subRoot.right.value) >= 0 {

			newSubRoot = subRoot.left
			if newSubRoot.right != nil {

				oldRightChild := newSubRoot.right
				newSubRoot.right = subRoot.right
				heap.AddHelper(newSubRoot, oldRightChild)
			} else {

				newSubRoot.right = subRoot.right
			}

		//Right child has a value greater than the left
		} else if heap.heuristic.Compare(subRoot.left.value, subRoot.right.value) < 0 {

			newSubRoot = subRoot.right
			if newSubRoot.left != nil {

				oldLeftChild := newSubRoot.left
				newSubRoot.left = subRoot.left
				newSubRoot = heap.AddHelper(newSubRoot, oldLeftChild)
			} else {

				newSubRoot.left = subRoot.left
			}
		}

	//No sub-trees exist!
	} else {

		newSubRoot = nil
	}

	return newSubRoot
}

func (heap *BinaryHeap) Pop() Heapable {

	if heap.root == nil {

		return nil

	}

	returnValue := heap.root.value
	heap.root = heap.PercolateUp(heap.root)

	return returnValue
}

func (heap *BinaryHeap) PrintInOrder() {

	if heap.root == nil {
		fmt.Println("--|")
	} else {
		heap.PrintInOrderHelper(heap.root, 0)
	}
}


//Recursive helper method for the PrintInOrder method
func (heap *BinaryHeap) PrintInOrderHelper(subRoot *BinaryHeapNode, depth int) {

	if subRoot == nil {
		return
	}

	heap.PrintInOrderHelper(subRoot.right, depth + 1)
	for i := 0; i < depth; i++ {
		fmt.Printf("----")
	}
	if depth > 0 {
		fmt.Printf(" ")
	}
	//assert type here...somehow
	fmt.Println(subRoot.value.String())
	heap.PrintInOrderHelper(subRoot.left, depth + 1)
}

//~------------------------------------------------------------------------------------------------
//~Test Type TYPEDEF----------------------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
type integer int

func (num integer) String() string {
	return strconv.Itoa(int(num))
}

func (num integer) Int() int {
	return int(num)
}

func (num1 integer) CompareTo(num2 interface{}) int {
	switch num2.(type) {
		case integer:
			if int(num1) > int(num2.(integer)) {
				return 1
			} else if int(num1) < int(num2.(integer)) {
				return -1
			} else {
				return 0
			}
		default:
			fmt.Println("VALUE PASSED TO integer CompareTo IS INVALID!")
			return -99
	}
}

//~------------------------------------------------------------------------------------------------
//~MAIN--------------------------------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
//~------------------------------------------------------------------------------------------------
func main () {

	//Does this call the constructor? Don't think so...but maybe....
	myHeap := new(BinaryHeap)
	myHeuristic := new(MaxHeap)
	
	//Set heap to hold integers (but does not add the number 5...there must be a cleaner way to do this
	myHeap.Heap(myHeuristic)
	myHeap.PrintInOrder()
	fmt.Println("CREATED HEAP\n\n")
	
	myHeap.Add(integer(5))
	myHeap.PrintInOrder()
	fmt.Println("ADDED 5\n\n")
	
	myHeap.Add(integer(55))
	myHeap.PrintInOrder()
	fmt.Println("ADDED 55\n\n")
	
	myHeap.Add(integer(13))
	myHeap.PrintInOrder()
	fmt.Println("ADDED 13\n\n")
	
	myHeap.Add(integer(-9))
	myHeap.PrintInOrder()
	fmt.Println("ADDED -9\n\n")
	
	myHeap.Add(integer(0))
	myHeap.PrintInOrder()
	fmt.Println("ADDED 0\n\n")

	myHeap.Remove(integer(55))
	myHeap.PrintInOrder()
	fmt.Println("REMOVED 55\n\n")

	poppedVal := myHeap.Pop()
	myHeap.PrintInOrder()
	if poppedVal != nil {
		fmt.Println("POPPED ", poppedVal.String(), "\n\n")
	} else {
		fmt.Println("POPPED nil\n\n")
	}

	poppedVal = myHeap.Pop()
	myHeap.PrintInOrder()
	if poppedVal != nil {
		fmt.Println("POPPED ", poppedVal.String(), "\n\n")
	} else {
		fmt.Println("POPPED nil\n\n")
	}

	poppedVal = myHeap.Pop()
	myHeap.PrintInOrder()
	if poppedVal != nil {
		fmt.Println("POPPED ", poppedVal.String(), "\n\n")
	} else {
		fmt.Println("POPPED nil\n\n")
	}

	poppedVal = myHeap.Pop()
	myHeap.PrintInOrder()
	if poppedVal != nil {
		fmt.Println("POPPED ", poppedVal.String(), "\n\n")
	} else {
		fmt.Println("POPPED nil\n\n")
	}

	poppedVal = myHeap.Pop()
	myHeap.PrintInOrder()
	if poppedVal != nil {
		fmt.Println("POPPED ", poppedVal.String(), "\n\n")
	} else {
		fmt.Println("POPPED nil\n\n")
	}

	poppedVal = myHeap.Pop()
	myHeap.PrintInOrder()
	if poppedVal != nil {
		fmt.Println("POPPED ", poppedVal.String(), "\n\n")
	} else {
		fmt.Println("POPPED nil\n\n")
	}

	poppedVal = myHeap.Pop()
	myHeap.PrintInOrder()
	if poppedVal != nil {
		fmt.Println("POPPED ", poppedVal.String(), "\n\n")
	} else {
		fmt.Println("POPPED nil\n\n")
	}
	
	
	myHeap.Add(integer(82))
	myHeap.PrintInOrder()
	fmt.Println("ADDED 82\n\n")
	
	myHeap.Add(integer(99))
	myHeap.PrintInOrder()
	fmt.Println("ADDED 99\n\n")
	
	myHeap.Add(integer(-33))
	myHeap.PrintInOrder()
	fmt.Println("ADDED -33\n\n")
	
	myHeap.Add(integer(-55))
	myHeap.PrintInOrder()
	fmt.Println("ADDED -55\n\n")
	
	myHeap.Add(integer(-2))
	myHeap.PrintInOrder()
	fmt.Println("ADDED -2\n\n")
	
	myHeap.Add(integer(0))
	myHeap.PrintInOrder()
	fmt.Println("ADDED 0\n\n")
	//*/
	fmt.Println("hey look I ran alright, LEAVE ME ALONE!!!!!")
}