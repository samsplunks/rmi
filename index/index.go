package index

import (
	"fmt"
	"math"
	"rmi/linear"
	"sort"
)

type kv struct {
	key    float64
	offset int
}

/*
LearnedIndex is an index structure that use inference to locate keys
*/
type LearnedIndex struct {
	M            *linear.RegressionModel
	rawdataTable []float64
	sortedTable  []*kv
	Len          int
	maxError     int
}

/*
New return an LearnedIndex fitted over the dataset with a linear regression algorythm
*/
func New(dataset []float64) *LearnedIndex {
	var sortedTable []*kv
	for i, row := range dataset {
		sortedTable = append(sortedTable, &kv{key: row, offset: i})
	}
	sort.SliceStable(sortedTable, func(i, j int) bool { return sortedTable[i].key < sortedTable[j].key })

	x, y := linear.Cdf(dataset)
	len := len(dataset)
	m := linear.Fit(x, y)
	maxErr := 0
	for i, k := range x {
		guessPosition := math.Round(scale(m.Predict(k), len))
		truePosition := math.Round(scale(y[i], len))
		residual := math.Sqrt(math.Pow(truePosition-guessPosition, 2))
		if float64(maxErr) < residual {
			maxErr = int(residual)
		}

	}
	return &LearnedIndex{M: m, Len: len, maxError: maxErr, sortedTable: sortedTable}
}

/*
scale return the CDF value x datasetLen -1 to get back the position in a sortedTable
*/
func scale(cdfVal float64, datasetLen int) float64 {
	return cdfVal*float64(datasetLen) - 1
}

/*
GuessIndex return the predicted position of the key in the index
and upper / lower positions' search interval
*/
func (idx *LearnedIndex) GuessIndex(key float64) (guess, lower, upper int) {
	guess = int(math.Round(scale(idx.M.Predict(key), idx.Len)))
	if guess < 0 {
		guess = 0
	} else if guess > idx.Len-1 {
		guess = idx.Len - 1
	}
	lower = guess - idx.maxError
	if lower < 0 {
		lower = 0
	}
	upper = guess + idx.maxError
	if upper > idx.Len-1 {
		upper = idx.Len - 1
	}
	return guess, lower, upper
}

/*
Lookup return the first offset of the key or err if the key is not found in the index
*/
func (idx *LearnedIndex) Lookup(key float64) (offset int, err error) {
	guess, lower, upper := idx.GuessIndex(key)

	if 0 <= guess && guess < idx.Len {
		if idx.sortedTable[guess].key == key {
			return idx.sortedTable[guess].offset, nil
		} else if idx.sortedTable[guess].key < key {
			return binarySearch(key, idx.sortedTable[guess+1:upper+1])
		} else {
			return binarySearch(key, idx.sortedTable[lower:guess])
		}
	}

	return -1, fmt.Errorf("The following key <%f> is not found in the index", key)
}

/*
binarySearch implementation is for finding the leftmost element
*/
func binarySearch(key float64, searchSpace []*kv) (offeset int, err error) {
	L := 0
	R := len(searchSpace) - 1
	nIter := 0

	for L <= R {
		m := int(math.Floor(float64((L + R) / 2)))
		if searchSpace[m].key < key {
			L = m + 1
		} else if searchSpace[m].key > key {
			R = m - 1
		} else {
			return searchSpace[m].offset, nil
		}
		nIter++
	}
	return -1, fmt.Errorf("The following key <%f> is not found in the index", key)
}
