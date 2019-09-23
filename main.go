package main

import (
	"bufio"
	"fmt"
	"github.com/amitbasuri/linuxPathTraversal/system"
	"os"
	"strings"
)

const sessionClearInput = "session clear"

func main() {

	fmt.Println("Starting your application...")
	state := system.NewState()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		input = strings.TrimSpace(input)
		if input == sessionClearInput {
			state = system.NewState() // flush the state
			continue
		}
		msg := executeInput(input, state)
		fmt.Println(msg)
	}

	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}

}

// executeInput executes the Input to the system current state
func executeInput(input string, state *system.State) string {
	cmd, err := system.NewCmdFromStr(input)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return system.ExecuteCommandInCurrState(cmd, state)
}
