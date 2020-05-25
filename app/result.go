package app

import (
	"github.com/gin-gonic/gin"
	"sort"
)

type searchResultSet map[int]float32 // caseId -> weight

var (
	emptyResponse = pairList{}.toResponse()
)

// A data structure to hold a caseId/value pair.
type pair struct {
	caseId int
	weight float32
}

// A slice of Pairs that implements sort.Interface to sort by weight.
type pairList []pair

func (p pairList) Len() int {
	return len(p)
}

func (p pairList) Less(i, j int) bool {
	return p[i].weight < p[j].weight
}

func (p pairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Intersect search results, nil set stands for **full set**
func (s searchResultSet) merge(t searchResultSet) searchResultSet {
	if s == nil {
		return t
	}
	if t == nil {
		return s
	}
	if len(t) < len(s) {
		return t.merge(s)
	}
	// Assume len(s) <= len(t)
	res := searchResultSet{}
	for id, w1 := range s {
		if w2, contains := t[id]; contains {
			res[id] = w1 + w2 // TODO maybe other operations
		}
	}
	return res
}

// A function to turn a map into a pairList, then sort and return it.
func (s searchResultSet) sortMapByValue() pairList {
	p := make(pairList, len(s))
	i := 0
	for k, v := range s {
		p[i] = pair{k, v}
		i += 1
	}
	sort.Sort(sort.Reverse(p))
	return p
}

// To ByteArray
func (p pairList) toByteArray() []byte {
	ids := make([]byte, len(p) * 3)
	for i, pr := range p {
		ids[i * 3 + 0] = byte((pr.caseId >>  0) & 0xff)
		ids[i * 3 + 1] = byte((pr.caseId >>  8) & 0xff)
		ids[i * 3 + 2] = byte((pr.caseId >> 16) & 0xff)
	}
	return ids
}

// To http response body
func (p pairList) toResponse() gin.H {
	ids := make([]int, len(p))
	for i, pr := range p {
		ids[i] = pr.caseId
	}
	return gin.H{
		"count":  len(p),
		"result": ids,
	}
}
