package znet

import (
	"fmt"
	"math"
	"testing"
)

type NumMatrix struct {
	PreMatrix [][]int
}


func Constructor(matrix [][]int) NumMatrix {
	preMatrix := make([][]int, len(matrix)+1)
	for i:=0; i< len(preMatrix); i++ {
		row := make([]int, len(matrix[0])+1)
		preMatrix[i] = row
	}
	for i:=0; i< len(matrix); i++ {
		row := preMatrix[i+1]
		for j:=0; j<len(matrix[0]); j++ {
			row[j+1] = preMatrix[i][j+1] + preMatrix[i+1][j] - preMatrix[i][j] + matrix[i][j]
		}
		preMatrix[i+1] = row
	}
	return NumMatrix{
		PreMatrix: preMatrix,
	}
}


func (this *NumMatrix) SumRegion(row1 int, col1 int, row2 int, col2 int) int {
	row1++
	col1++
	row2++
	col2++
	return this.PreMatrix[row2+1][col2+1] + this.PreMatrix[row1][col1] - this.PreMatrix[row1][col2+1] - this.PreMatrix[row2+1][col1]
}

func TestMatrix(t *testing.T) {
	//m := make([][]int, 2)
	//row := make([]int, 2)
	//row[0] = 1
	//row[1] = 2
	//m[0] = row
	//row1 := make([]int, 2)
	//row1[0] = 3
	//row1[1] = 4
	//m[1] = row1
	//mr := Constructor(m)
	//fmt.Println(mr.SumRegion(0,0,1,1))
	//fmt.Println(findKthLargest([]int{3,2,3,1,2,4,5,5,6}, 4))
	fmt.Println(permute([]int{0,1,2,3,4,5,6,7}))
}

var result [][]int
var index int

func permute(nums []int) [][]int {
	result = [][]int{}
	used := make([]int, len(nums))
	track := make([]int, len(nums))
	backtrace(nums, track, used)
	return result
}

func backtrace(choices []int, route []int, used []int) {
	if index == len(choices) {
		tmp := make([]int, len(route))
		copy(tmp, route)
		result = append(result, tmp)
		return
	}

	for i, choice := range choices {
		if used[i] == 1 {
			continue
		}
		route[index] = choice
		used[i] = 1
		index++
		backtrace(choices, route, used)
		index--
		used[i] = 0
		route[index] = 0
	}
}

func coinChange(nums []int, target int) int {
	dp := make([]int, target+1)
	for i := 1; i <= target; i++ {
		dp[i] = math.MaxInt32
	}
	for _, num := range nums {
		for i := 0; i <= target-num; i++ {
			if dp[i] == math.MaxInt32 {
				continue
			}
			dp[i+num] = min(dp[i+num], dp[i]+1)
		}
	}
	if dp[target] == math.MaxInt32 {
		return -1
	}
	return dp[target]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func findKthLargest(nums []int, k int) int {
	quickSort(&nums, 0, len(nums)-1)

	return nums[len(nums)-k]
}

func quickSort(nums *[]int, l int, r int) {
	if l >= r {
		return
	}
	i, j, p := l, r, (*nums)[r]
	for i < j {
		for (*nums)[i] <= p && i < j {
			i++
		}
		for (*nums)[j] >= p  && i < j {
			j--
		}
		(*nums)[i],(*nums)[j] = (*nums)[j],(*nums)[i]
	}
	(*nums)[i],(*nums)[r] = (*nums)[r], (*nums)[i]
	quickSort(nums, l, i-1)
	quickSort(nums, i+1, r)
}