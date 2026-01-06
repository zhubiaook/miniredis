package client

import (
	"fmt"
	"strings"
)

const (
	TypeSimpleString = '+'
	TypeBulkString   = '$'
	TypeInteger      = ':'
	TypeArray        = '*'
	TypeError        = '-'
)

func Encode(args []string) []byte {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%c%d\r\n", TypeArray, len(args)))
	for _, arg := range args {
		b.WriteString(fmt.Sprintf("%c%d\r\n%s\r\n", TypeBulkString, len(arg), arg))
	}
	return []byte(b.String())
}
