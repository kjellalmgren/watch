// This code reproduces some of the concepts explained Liz Rice's talk
// "Building a container from scratch in Go", 2016.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

const usage = `usage:
    demo-container command [arguments]
commands:
    run: runs the first argument as an executable, using any
         subsequent arguments as the executable arguments.
`

const cmdRun = "run"

//
func main() {

	ctx := context.Background()

	// trap Ctrl+C and call cancel on the content
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		//
		defer func() {
			signal.Stop(c)
			cancel()
		}()
		go func() {
			select {
			case <-c:
				cancel()
			case <-ctx.Done():
			}
		}()
		//
		do, err := parseArgs()
		if err != nil {
			log.Println(err)
			fmt.Println(usage)
			os.Exit(1)
		}
		//
		if err := do(os.Args[2:]...); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("sleeping...")
		time.Sleep(5 * time.Second)
	}
}

//
func parseArgs() (action, error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("Watch needs a command")
	}
	switch os.Args[1] {
	case cmdRun:

		return run, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", os.Args[1])
	}
}

type action func(args ...string) error

func run(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("run command needs at least one argument")
	}
	cmd := exec.Command("run clear")
	cmd = exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//cmd.SysProcAttr = &syscall.SysProcAttr{
	// create new UTS namespace (since Linux 2.6.19)
	//	Cloneflags: syscall.CLONE_NEWUTS,
	//}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
