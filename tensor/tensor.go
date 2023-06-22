package tensor

import (
	"github.com/gpabois/gostd/iter"
	"golang.org/x/exp/slices"
)

// Yale-format encoded matrix
type SparseMatrix[Data any] struct {
	V           []Data // Length == NNZ
	RowIndex    []int
	ColumnIndex []int // Length == NNZ

	RowLength    int
	ColumnLength int
}

func NewSparseMatrix[Data iter.Iterable[Element], Element any, Value any](rowLength int, columnLength int, data Data, indexer func(el Element) (row int, col int), valuer func(el Element) Value) SparseMatrix[Value] {
	elements := iter.CollectToSlice[[]Element](data.Iter())
	NNZ := len(elements)
	var mat SparseMatrix[Value]

	// Initialise parameters
	mat.RowLength = rowLength
	mat.ColumnLength = columnLength
	mat.V = make([]Value, NNZ)
	mat.ColumnIndex = make([]int, NNZ)
	mat.RowIndex = make([]int, rowLength+1)

	dict := make(map[int]*[]Element)

	// Group by rows
	for _, el := range elements {
		row, _ := indexer(el)
		if column, ok := dict[row]; ok {
			*column = append(*column, el)
		} else {
			dict[row] = &[]Element{el}
		}
	}

	rowOffset := 0
	for row, group := range dict {
		slices.SortFunc(*group, func(i, j Element) bool {
			_, coli := indexer(i)
			_, colj := indexer(j)

			return coli < colj
		})

		rowStart := rowOffset
		rowEnd := rowStart + len(*group)

		columnIndexSlice := make([]int, len(*group))
		values := make([]Value, len(*group))

		for i, el := range *group {
			_, col := indexer(el)
			columnIndexSlice[i] = col
			values[i] = valuer(el)
		}

		// Set the RowIndex, V and ColumnIndex values
		mat.RowIndex[row] = rowStart
		slices.Replace(mat.V, rowStart, rowEnd-1, values...)
		slices.Replace(mat.ColumnIndex, rowStart, rowEnd-1, columnIndexSlice...)

		rowOffset += len(*group)
	}

	return mat
}

func (mat *SparseMatrix[Data]) GetOrDefault(row int, col int) Data {
	var data Data

	if row >= len(mat.RowIndex) {
		return data
	}

	rowStart := mat.RowIndex[row]
	rowEnd := mat.RowIndex[row+1]

	nnzColumnIndexes := mat.ColumnIndex[rowStart:rowEnd]
	nnzColumnValues := mat.V[rowStart:rowEnd]

	if colIndex := slices.Index(nnzColumnIndexes, col); colIndex != -1 {
		return nnzColumnValues[colIndex]
	}

	return data
}

// Windowed Sparse Matrix
type WSMBounds struct {
	RowOffset    int
	RowLength    int
	ColumnOffset int
	ColumnLength int
}
type WSM[Data any] struct {
	Window SparseMatrix[Data]
	Bounds WSMBounds
}

func NewWSM[Data iter.Iterable[Element], Element any, Value any](rowOffset, rowLength, columnOffset, columnLength int, data Data, indexer func(el Element) (row int, col int), valuer func(el Element) Value) WSM[Value] {

	window := NewSparseMatrix(rowLength, columnLength, data, func(el Element) (int, int) {
		row, col := indexer(el)
		return row - rowOffset, col - columnOffset
	}, valuer)

	return WSM[Value]{
		Window: window,
		Bounds: WSMBounds{
			RowOffset:    rowOffset,
			RowLength:    rowLength,
			ColumnOffset: rowOffset,
			ColumnLength: columnLength,
		},
	}
}
