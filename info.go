package main

// Struct holding file/directory info
type Info struct {
    directory   bool
    owner       string
    group       string
    size        uint64
    name        string
    children    []*Info
    parent      *Info
}

// Implement sort interface [Len(), Less(i, j), Swap(i, j)] for list of Info's
// sort by size
type bySize []*Info

func (i bySize) Len() int {
    return len(i)
}

func (i bySize) Less(a, b int) bool {
    // largest is first
    return i[a].size > i[b].size
}

func (i bySize) Swap(a, b int) {
    i[a], i[b] = i[b], i[a]
}

// sort by size
type byName []*Info

func (i byName) Len() int {
    return len(i)
}

func (i byName) Less(a, b int) bool {
    // largest is first
    return i[a].name < i[b].name
}

func (i byName) Swap(a, b int) {
    i[a], i[b] = i[b], i[a]
}

// sort by owner
type byOwner []*Info

func (i byOwner) Len() int {
    return len(i)
}

func (i byOwner) Less(a, b int) bool {
    // largest is first
    return i[a].owner < i[b].owner
}

func (i byOwner) Swap(a, b int) {
    i[a], i[b] = i[b], i[a]
}

// sort by size
type byGroup []*Info

func (i byGroup) Len() int {
    return len(i)
}

func (i byGroup) Less(a, b int) bool {
    // largest is first
    return i[a].group < i[b].group
}

func (i byGroup) Swap(a, b int) {
    i[a], i[b] = i[b], i[a]
}



