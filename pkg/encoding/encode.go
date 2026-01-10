package encoding

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	TypeSimpleString = '+'
	TypeBulkString   = '$'
	TypeInteger      = ':'
	TypeArray        = '*'
	TypeError        = '-'
)

func EncodeWrite(out io.Writer, in []string) error {
	inLen := len(in)
	if inLen == 0 {
		return fmt.Errorf("invalid length")
	}

	w := bufio.NewWriter(out)
	if _, err := w.WriteString(fmt.Sprintf("%c%d\r\n", TypeArray, inLen)); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	for _, v := range in {
		if _, err := w.WriteString(fmt.Sprintf("%c%d\r\n%s\r\n", TypeBulkString, len(v), v)); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}
	return w.Flush()
}

func DecodeRead(in io.Reader, out *[]string) error {
	rst := make([]string, 0)

	r := bufio.NewReader(in)
	b, err := r.ReadByte()
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	switch b {
	case TypeSimpleString:
		token, err := r.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}
		token = token[:len(token)-2]
		rst = append(rst, token)
	case TypeArray:
		token, err := r.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}
		token = token[:len(token)-2]
		arrLen, err := strconv.Atoi(token)
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}
		for range arrLen {
			token, err := r.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read error: %w", err)
			}
			token = token[:len(token)-2]
			if token[0] != TypeBulkString {
				return fmt.Errorf("read error: %w", err)
			}
			bukLen, err := strconv.Atoi(token[1:])
			if err != nil {
				return fmt.Errorf("read error: %w", err)
			}
			buk := make([]byte, bukLen+2)
			if _, err := r.Read(buk); err != nil {
				return fmt.Errorf("read error: %w", err)
			}
			buk = buk[:bukLen]
			rst = append(rst, string(buk))
		}
	default:
		return fmt.Errorf("invalid type")
	}

	*out = rst
	return nil
}
