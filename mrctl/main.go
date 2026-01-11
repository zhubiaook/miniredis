package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:6379", "redis server address")
	flag.Parse()

	opt := Options{
		Addr: *addr,
	}
	client, err := NewClient(opt)
	if err != nil {
		fmt.Println(err)
		return
	}

	// non-interactive mode
	args := flag.Args()
	if len(args) > 0 {
		out, err := client.Do(args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println(out)
		return
	}

	// interactive mode
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.ToUpper(line) == "EXIT" {
			fmt.Println("Goodbye!")
			break
		}

		args := strings.Fields(line)
		list, err := client.Do(args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		if len(list) == 1 {
			fmt.Println(list[0])
			continue
		}
		for i, v := range list {
			fmt.Printf("%d) %q\n", i+1, v)
		}
	}

	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}
}
