// Package common содержит общие типы данных, константы и структуры
package common

const (
	KB = 1 << (10 * (iota + 1))
	MB
)

type Uint24 uint32
