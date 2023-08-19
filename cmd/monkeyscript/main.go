package main

import (
	"fmt"
	"os"
	"os/user"

	"git.tigh.dev/tigh-latte/monkeyscript/repl"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("unable to identify: %w", err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Println("Feel free to type commands")

	repl.Start(os.Stdin, os.Stdout)
	return nil
}
