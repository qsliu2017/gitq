package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jmorganca/ollama/api"
)

const SYSTEM_PROMPT = `You are a command line tool to help developer transform natural language into git commands.
You should only output the git command, starting with '> ', without any additional information.
Example:
> git commit -m "Add a new feature"
`

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	go func() {
		<-signalChan
		cancel()
	}()

	if err := client.Generate(ctx,
		&api.GenerateRequest{
			Model:  "codellama",
			Prompt: strings.Join(os.Args[1:], " "),
			System: SYSTEM_PROMPT,
		},
		func(response api.GenerateResponse) error {
			if response.Done {
				fmt.Print("\n")
				response.Summary()
				return nil
			}

			fmt.Print(response.Response)
			return nil
		},
	); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		panic(err)
	}
}
