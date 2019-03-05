package main

type backendHeap []*httpBackend

func (b backendHeap) Len() int           { return len(b) }
func (b backendHeap) Less(i, j int) bool { return b[i].weight > b[j].weight }
func (b backendHeap) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func (b *backendHeap) Push(x *httpBackend) {
	*b = append(*b, x)
}

// Add timeout if no backends in heap given

func (b *backendHeap) Pop() *httpBackend {
	old := *b
	n := len(old)
	x := old[n-1]
	*b = old[0 : n-1]
	return x
}