package main

import (
	"flag"
	"fmt"

	"github.com/zhubiaook/miniredis/mrctl/client"
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

	opts := client.Options{
		Addr: *addr,
	}

	cli, err := client.NewClient(opts)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: failed to connect to %s: %v\n", *addr, err)
		return
	}
	defer cli.Close()

	resp, err := cli.Do(args...)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: command failed: %v\n", err)
		return
	}
	fmt.Println(string(resp))
}
