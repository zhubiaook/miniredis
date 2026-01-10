package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/zhubiaook/miniredis/pkg/encoding"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:6379", "redis server address")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <command> [args...]\n", "mrctl")
		fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExamples:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  mrctl GET mykey\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  mrctl SET mykey value\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  mrctl -addr localhost:6379 PING\n")
		return
	}

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: failed to connect to %s: %v\n", *addr, err)
		return
	}
	defer conn.Close()

	if err := encoding.EncodeWrite(conn, args); err != nil {
		fmt.Printf("Error: failed to encode command: %v\n", err)
		return
	}

	var out []string
	if err := encoding.DecodeRead(conn, &out); err != nil {
		fmt.Printf("Error: failed to decode response: %v\n", err)
		return
	}

	fmt.Println(out)
}
