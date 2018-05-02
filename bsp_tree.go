package main

import (
	"sort"
	"fmt"
)

const (
	minPartitionWidth float64 = 3.0
	minPartitionHeight float64 = 3.0
)

/*const (
	closed uint = iota
	semiOpen
	open
)*/

type bspNode struct {
	area rectangle	//from the upper left point (inclusive) to the lower right point (exclusive)
	//traversability uint
	left *bspNode	//less than the partition
	right *bspNode	//greater than or equal to the partition
}

//Determines whether or not an element from a given set of partitions can create a spatial partition which obeys the dimensional minimums (but only in the partition dimension, because it is assumed that the space has already been verified to be at least the minimum dimensions).
//Starts at the median partition, and starts searching around it.  Returns the index and the associated partitioned spaces if a suitable partition is found, or -1 and arealess rectanlges otherwise.
func hasSufficientPartition(space rectangle, sortedPartitions []point, dimension uint) (int, rectangle, rectangle) {
	partitionMin, partitionMax := 0, len(sortedPartitions) - 1
	partitionMedian := func() int {return (partitionMin + partitionMax) / 2}
	
	if len(sortedPartitions) > 0 {
		var minPartitionSize float64
		var getLeftSize, getRightSize (func() float64)
		var getLeftSpace, getRightSpace (func() rectangle)
		if dimension % 2 == 0 {
			minPartitionSize = minPartitionWidth
			getLeftSize = func() float64 {return sortedPartitions[partitionMedian()].x - space.corners[upperLeft].x}
			getRightSize = func() float64 {return space.corners[upperRight].x - sortedPartitions[partitionMedian()].x}
			getLeftSpace = func() rectangle {return initRectangle(space.corners[upperLeft].x, space.corners[upperLeft].y, getLeftSize(), space.corners[lowerLeft].y - space.corners[upperLeft].y)}
			getRightSpace = func() rectangle {return initRectangle(sortedPartitions[partitionMedian()].x, space.corners[upperLeft].y, getRightSize(), space.corners[lowerLeft].y - space.corners[upperLeft].y)}
		}else{
			minPartitionSize = minPartitionHeight
			getLeftSize = func() float64 {return sortedPartitions[partitionMedian()].y - space.corners[upperLeft].y}
			getRightSize = func() float64 {return space.corners[lowerLeft].y - sortedPartitions[partitionMedian()].y}
			getLeftSpace = func() rectangle {return initRectangle(space.corners[upperLeft].x, space.corners[upperLeft].y, space.corners[upperRight].x - space.corners[upperLeft].x, getLeftSize())}
			getRightSpace = func() rectangle {return initRectangle(space.corners[upperLeft].x, sortedPartitions[partitionMedian()].y, space.corners[upperRight].x - space.corners[upperLeft].x, getRightSize())}
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
	return -1, initRectangle(0.0, 0.0, 0.0, 0.0), initRectangle(0.0, 0.0, 0.0, 0.0)
}

//Recursively builds a BSP Tree (assumes the inital space is larger than the minimum dimensions).
func constructBSPTree(space rectangle, partitions []point, depth uint) *bspNode {
	node := bspNode{area: space, left: nil, right: nil}
	if len(partitions) > 0 {
		var sortedPartitions []point
		var sortFunction (func(int, int) bool)
		if depth % 2 == 0 {
			sortFunction = func(i, j int) bool {return sortedPartitions[i].x < sortedPartitions[j].x}
		}else{
			sortFunction = func(i, j int) bool {return sortedPartitions[i].y < sortedPartitions[j].y}
		}
		sortedPartitions = make([]point, len(partitions), cap(partitions))
		copy(sortedPartitions, partitions)
		sort.Slice(sortedPartitions, sortFunction)
		
		partitionIndex, leftSpace, rightSpace := hasSufficientPartition(space, sortedPartitions, depth % 2)
		if partitionIndex >= 0 {
			node.left = constructBSPTree(leftSpace, sortedPartitions[:partitionIndex], depth + 1)
			node.right = constructBSPTree(rightSpace, sortedPartitions[partitionIndex:], depth + 1)
		}
	}
	return &node
}

//Really need to refactor some time in the furture so this can become it's own package, with this function being the only one visible...
func initBSPTree(initial_area rectangle, partitions []point) bspNode {
	if initial_area.corners[upperRight].x - initial_area.corners[upperLeft].x < minPartitionWidth {
		panic(fmt.Sprintf("The width of the given space is less than the minimum of %d.", minPartitionWidth))
	}
	if initial_area.corners[lowerLeft].y - initial_area.corners[upperLeft].y < minPartitionHeight {
		panic(fmt.Sprintf("The height of the given space is less than the minimum of %d.", minPartitionHeight))
	}
	return *constructBSPTree(initial_area, partitions, 0)
}