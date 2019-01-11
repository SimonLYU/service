package MemoryModel

type Memory struct {
	Name     string
	DatabaseName string
	MemoList []string
}

type SingleMemory struct {
	Name     string
	DatabaseName string
	Memo string
}

type GetMemory struct {
	DatabaseName string
}