package system

import (
	"fmt"
	"strings"
)

const (
	commandcd    = "cd"
	commandls    = "ls"
	commandmkdir = "mkdir"
	commandpwd   = "pwd"
	commandrm    = "rm"

	invalidInputErr = "ERR: CANNOT RECOGNIZE INPUT."
)

//CommandExecFuncMap maps the command to the execution function
var CommandExecFuncMap = map[string]cmdFunc{
	commandcd:    cd,
	commandls:    ls,
	commandmkdir: mkdir,
	commandpwd:   pwd,
	commandrm:    rm,
}

type Command struct { // cd home/docs
	Instruction string // cd
	Parameters  string // home/docs
}

func (cmd *Command) String() string {
	return cmd.Instruction + " " + cmd.Parameters
}

//validateInstruction validates an Instruction using CommandExecFuncMap
func validateInstruction(cmdInstruction string) error {

	if _, ok := CommandExecFuncMap[cmdInstruction]; !ok { // if its a valid command
		return fmt.Errorf(invalidInputErr)
	}

	return nil
}

// NewCmdFromStr will generate a Command from user input string
func NewCmdFromStr(input string) (*Command, error) {
	inputSplitArr := strings.Fields(input) // input = "mkdir a/b"

	if len(inputSplitArr) > 2 {
		return nil, fmt.Errorf(invalidInputErr)
	}
	cmdInstruction := inputSplitArr[0]
	if err := validateInstruction(cmdInstruction); err != nil { // validate instruction
		return nil, err
	}
	cmd := &Command{
		Instruction: cmdInstruction, // --> "mkdir"
	}

	if len(inputSplitArr) > 1 {
		cmd.Parameters = inputSplitArr[1] // --> "a/b"
	}

	return cmd, nil
}
