package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/zigbee-s/go-todo-cli/api"
)

func main() {
	add := flag.Bool("add", false, "add a new todo")
	complete := flag.Int("complete", 0, "mark a todo as completed")
	del := flag.Int("del", 0, "delete a todo")
	list := flag.Bool("list", false, "list all todos")

	flag.Parse()

	switch {
	case *add:

		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = api.Add(task)
		if err != nil {
			panic(err)
		}

	case *complete > 0:
		err := api.Complete(*complete)
		if err != nil {
			panic(err)
		}

	case *del > 0:
		err := api.Delete(*del)
		if err != nil {
			panic(err)
		}

	case *list:
		err := api.List()
		if err != nil {
			panic(err)
		}

	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(0)

	}
}

func getInput(r io.Reader, args ...string) (string, error) {

	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed")
	}

	return text, nil

}
