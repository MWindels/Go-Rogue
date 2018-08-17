package adt

import (
	"sort"
	"fmt"
	"math"
	"github.com/mwindels/go-rogue/geom"
)

//Could go elsewhere...
const (
	MinPartitionWidth float64 = 3.0
	MinPartitionHeight float64 = 3.0
)

const (
	Closed uint = iota
	SemiOpen
	Open
)

const (
	noAction uint = iota
	openNode
	closeNode
)

type BSPNode struct {
	area geom.Rectangle	//from the upper left point (inclusive) to the lower right point (exclusive)
	traversability uint
	left *BSPNode	//less than the partition
	right *BSPNode	//greater than or equal to the partition
}

func (node BSPNode) Area() geom.Rectangle {
	return node.area
}

func (node BSPNode) Traversability() uint {
	return node.traversability
}

func (node BSPNode) Left() *BSPNode {
	return node.left
}

func (node BSPNode) Right() *BSPNode {
	return node.right
}

type BSPTree struct {
	root *BSPNode
	depth uint		//effectively the lowest level of the tree (starts at zero (root) sice a nodeless bsp tree doesn't make sense, as it wouldn't be partitioning anything)
}

func (tree BSPTree) Root() *BSPNode {
	return tree.root
}

func (tree BSPTree) Depth() uint {
	return tree.depth
}

//Determines whether or not an element from a given set of partitions can create a spatial partition which obeys the dimensional minimums (but only in the partition dimension, because it is assumed that the space has already been verified to be at least the minimum dimensions).
//Starts at the median partition, and starts searching around it.  Returns the index and the associated partitioned spaces if a suitable partition is found, or -1 and arealess rectanlges otherwise.
func hasSufficientPartition(space geom.Rectangle, sortedPartitions []geom.Point, dimension uint) (int, geom.Rectangle, geom.Rectangle) {
	partitionMin, partitionMax := 0, len(sortedPartitions) - 1
	partitionMedian := func() int {return (partitionMin + partitionMax) / 2}
	
	if len(sortedPartitions) > 0 {
		var minPartitionSize float64
		var getLeftSize, getRightSize (func() float64)
		var getLeftSpace, getRightSpace (func() geom.Rectangle)
		if dimension % 2 == 0 {
			minPartitionSize = MinPartitionWidth
			getLeftSize = func() float64 {return sortedPartitions[partitionMedian()].X - space.UpperLeft().X}
			getRightSize = func() float64 {return space.UpperRight().X - sortedPartitions[partitionMedian()].X}
			getLeftSpace = func() geom.Rectangle {return geom.InitRectangle(space.UpperLeft().X, space.UpperLeft().Y, getLeftSize(), space.LowerLeft().Y - space.UpperLeft().Y)}
			getRightSpace = func() geom.Rectangle {return geom.InitRectangle(sortedPartitions[partitionMedian()].X, space.UpperLeft().Y, getRightSize(), space.LowerLeft().Y - space.UpperLeft().Y)}
		}else{
			minPartitionSize = MinPartitionHeight
			getLeftSize = func() float64 {return sortedPartitions[partitionMedian()].Y - space.UpperLeft().Y}
			getRightSize = func() float64 {return space.LowerLeft().Y - sortedPartitions[partitionMedian()].Y}
			getLeftSpace = func() geom.Rectangle {return geom.InitRectangle(space.UpperLeft().X, space.UpperLeft().Y, space.UpperRight().X - space.UpperLeft().X, getLeftSize())}
			getRightSpace = func() geom.Rectangle {return geom.InitRectangle(space.UpperLeft().X, sortedPartitions[partitionMedian()].Y, space.UpperRight().X - space.UpperLeft().X, getRightSize())}
		}
		
		for partitionMin <= partitionMax {
			if getLeftSize() < minPartitionSize || getLeftSize() == 0.0 {	//deals with things that lie right on the partition (e.g. are equal to 0), will ignore them for now (they might be useful in future partitions in the opposite dimension)
				partitionMin = partitionMedian() + 1
			}else if getRightSize() < minPartitionSize || getRightSize() == 0.0 {	//likewise
				partitionMax = partitionMedian() - 1
			}else{
				return partitionMedian(), getLeftSpace(), getRightSpace()
			}
		}
	}
	return -1, geom.InitRectangle(0.0, 0.0, 0.0, 0.0), geom.InitRectangle(0.0, 0.0, 0.0, 0.0)
}

//Recursively builds a BSP Tree (assumes the inital space is larger than the minimum dimensions).
func constructBSPTree(space geom.Rectangle, partitions []geom.Point, depth uint) (*BSPNode, uint) {
	leftDepth, rightDepth := depth, depth
	node := BSPNode{area: space, traversability: Closed, left: nil, right: nil}
	if len(partitions) > 0 {
		var sortedPartitions []geom.Point
		var sortFunction (func(int, int) bool)
		if depth % 2 == 0 {
			sortFunction = func(i, j int) bool {return sortedPartitions[i].X < sortedPartitions[j].X}
		}else{
			sortFunction = func(i, j int) bool {return sortedPartitions[i].Y < sortedPartitions[j].Y}
		}
		sortedPartitions = make([]geom.Point, len(partitions), cap(partitions))
		copy(sortedPartitions, partitions)
		sort.Slice(sortedPartitions, sortFunction)
		
		partitionIndex, leftSpace, rightSpace := hasSufficientPartition(space, sortedPartitions, depth % 2)
		if partitionIndex >= 0 {
			node.left, leftDepth = constructBSPTree(leftSpace, sortedPartitions[:partitionIndex], depth + 1)
			node.right, rightDepth = constructBSPTree(rightSpace, sortedPartitions[partitionIndex:], depth + 1)
		}
	}
	return &node, uint(math.Max(float64(leftDepth), float64(rightDepth)))
}

func InitBSPTree(initialArea geom.Rectangle, partitions []geom.Point) BSPTree {
	if initialArea.UpperRight().X - initialArea.UpperLeft().X < MinPartitionWidth {
		panic(fmt.Sprintf("The width of the given space is less than the minimum of %d.", MinPartitionWidth))
	}
	if initialArea.LowerLeft().Y - initialArea.UpperLeft().Y < MinPartitionHeight {
		panic(fmt.Sprintf("The height of the given space is less than the minimum of %d.", MinPartitionHeight))
	}
	rootNode, treeDepth := constructBSPTree(initialArea, partitions, 0)
	return BSPTree{root: rootNode, depth: treeDepth}
}

//Based on a left and right node, deduces the traversability of a node constituted of those two nodes.
func deduceTraversability(left, right BSPNode) uint {
	if left.traversability == Open && right.traversability == Open {
		return Open
	}else if left.traversability == Closed && right.traversability == Closed {
		return Closed
	}else{
		return SemiOpen
	}
}

//The FSM which opens and closes the nodes of a BSPTree randomly.
func (node *BSPNode) recursivelyRandomizeTraversability(openFunc, closeFunc (func (float64) bool), state, depth, maxDepth uint) {
	relativeDepth := float64(depth) / float64(maxDepth)
	if state == noAction {
		if openFunc(relativeDepth) {
			state = openNode
		}else if closeFunc(relativeDepth) {
			state = closeNode
		}
	}else if state == openNode {
		if closeFunc(relativeDepth) {
			state = closeNode
		}
	}
	
	if node.left != nil && node.right != nil {
		node.left.recursivelyRandomizeTraversability(openFunc, closeFunc, state, depth + 1, maxDepth)
		node.right.recursivelyRandomizeTraversability(openFunc, closeFunc, state, depth + 1, maxDepth)
		
		node.traversability = deduceTraversability(*(node.left), *(node.right))
	}else{
		if state == openNode {
			node.traversability = Open
		}else if state == closeNode {
			node.traversability = Closed
		}
	}
}

//Randomizes the traversability of nodes in a BSPTree.
func (tree BSPTree) RandomizeTraversability(openFunc, closeFunc (func (float64) bool)) {
	tree.root.recursivelyRandomizeTraversability(openFunc, closeFunc, noAction, 0, tree.depth)
}