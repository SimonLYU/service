package MemoryModel

type Memory struct {
	Name     string
	DatabaseName string
	MemoList []string
}

type GetMemory struct {
	DatabaseName string
}