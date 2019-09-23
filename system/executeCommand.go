package system

import (
	"fmt"
	"strings"
)

//CommandExecSuccesMsg is a map of msg on Success
var CommandExecSuccesMsg = map[string]string{
	commandcd:    "SUCC: REACHED",
	commandls:    "SUCC: ls",
	commandmkdir: "SUCC: CREATED",
	commandpwd:   "SUCC: pwd",
	commandrm:    "SUCC: REMOVED",
}

//CommandExecErrorMsg is a map of msg on Error
var CommandExecErrorMsg = map[string]string{
	commandcd:    "ERR: INVALID PATH",
	commandls:    "ERR: ls",
	commandmkdir: "ERR: DIRECTORY ALREADY EXISTS",
	commandpwd:   "ERR: pwd",
	commandrm:    "ERR: DIRECTORY / PATH",
}

type cmdFunc func(cmd *Command, state *State) string

//ExecuteCommandInCurrState will the call the required func for a given command using CommandExecFuncMap
func ExecuteCommandInCurrState(cmd *Command, state *State) string {
	return CommandExecFuncMap[cmd.Instruction](cmd, state)
}

// cd changes directory
func cd(cmd *Command, state *State) string {
	changedCurrDir := false
	initialStateCurrDir := state.CurrentDirectory

	if cmd.Parameters == "/" || cmd.Parameters == "" {
		state.CurrentDirectory = state.RootDirectory // switch to root dir
		changedCurrDir = true
	}

	// trim left by "/"
	if strings.HasPrefix(cmd.Parameters, "/") {
		state.CurrentDirectory = state.RootDirectory // switch to root dir
		cmd.Parameters = strings.TrimLeft(cmd.Parameters, "/")
	}

	goTodirNames := strings.Split(cmd.Parameters, "/") // "dir1/dir2" --> ["dir1", "dir2"]
	for _, goTodirName := range goTodirNames {

		// goto parent
		if goTodirName == ".." {
			if state.CurrentDirectory.ParentDirectory != nil {
				state.CurrentDirectory = state.CurrentDirectory.ParentDirectory
				changedCurrDir = true
				continue
			}
			return fmt.Sprintf(CommandExecErrorMsg[cmd.Instruction]) // its root directory
		}

		// goto children
		for _, childDir := range state.CurrentDirectory.childDirectories {
			if childDir.name == goTodirName {
				state.CurrentDirectory = childDir
				changedCurrDir = true
				break
			}
		}
		// if does not match any child dir
		if !changedCurrDir {
			break
		}
	}

	if !changedCurrDir { // was unable to change CurrDir --> error
		state.CurrentDirectory = initialStateCurrDir
		return fmt.Sprintf(CommandExecErrorMsg[cmd.Instruction])
	}

	return fmt.Sprintf(CommandExecSuccesMsg[cmd.Instruction])
}

// ls is used to list all child directories
func ls(cmd *Command, state *State) string {
	dirs := "DIRS:"
	for _, dir := range state.CurrentDirectory.childDirectories {
		dirs = dirs + " " + dir.name
	}
	return dirs
}

//mkdir is used to create new directory
func mkdir(cmd *Command, state *State) string {
	initialStateCurrDir := state.CurrentDirectory

	// trim left by "/"
	if strings.HasPrefix(cmd.Parameters, "/") {
		state.CurrentDirectory = state.RootDirectory // switch to root dir
		cmd.Parameters = strings.TrimLeft(cmd.Parameters, "/")
	}

	goToDirNames := strings.Split(cmd.Parameters, "/")
	for level, goToDirName := range goToDirNames {
		// goto parent
		if goToDirName == ".." {
			if state.CurrentDirectory.ParentDirectory != nil {
				state.CurrentDirectory = state.CurrentDirectory.ParentDirectory
				continue
			}
			return fmt.Sprintf(CommandExecErrorMsg[cmd.Instruction]) // its root directory
		}

		childDir, _ := getChild(state.CurrentDirectory, goToDirName)

		// if childDir with same name is there and it is the LAST goToDirName
		if childDir != nil && level == len(goToDirNames)-1 {
			return fmt.Sprintf(CommandExecErrorMsg[cmd.Instruction])
		} else if childDir != nil { //  if childDir with same name  and it is NOT the LAST goToDirName
			state.CurrentDirectory = childDir //goto children
		} else {
			state.CurrentDirectory.addChildDirectory(goToDirName)
		}

	}
	state.CurrentDirectory = initialStateCurrDir
	return fmt.Sprintf(CommandExecSuccesMsg[cmd.Instruction])
}

//pwd gives the present directory
func pwd(cmd *Command, state *State) string {
	initialStateCurrDir := state.CurrentDirectory

	pwd := state.CurrentDirectory.name
	for state.CurrentDirectory != state.RootDirectory && state.RootDirectory != nil {
		pwd = state.CurrentDirectory.ParentDirectory.name + "/" + pwd
		state.CurrentDirectory = state.CurrentDirectory.ParentDirectory
	}
	state.CurrentDirectory = initialStateCurrDir

	if pwd == "" {
		return "PATH: /"
	}
	return "PATH: " + pwd
}

func rm(cmd *Command, state *State) string {
	initialStateCurrDir := state.CurrentDirectory

	// trim left by "/"
	if strings.HasPrefix(cmd.Parameters, "/") {
		state.CurrentDirectory = state.RootDirectory // switch to root dir
		cmd.Parameters = strings.TrimLeft(cmd.Parameters, "/")
	}

	goToDirNames := strings.Split(cmd.Parameters, "/")
	for level, goToDirName := range goToDirNames {

		// goto parent
		if goToDirName == ".." {
			return fmt.Sprintf(CommandExecErrorMsg[cmd.Instruction])
		}

		childDir, _ := getChild(state.CurrentDirectory, goToDirName)

		// if childDir with same name is there and it is the LAST goToDirName
		if childDir != nil && level == len(goToDirNames)-1 {
			state.CurrentDirectory.removeChildDirectory(goToDirName) // remove it
		} else if childDir == nil { //no such childDir
			return fmt.Sprintf(CommandExecErrorMsg[cmd.Instruction])
		} else {
			state.CurrentDirectory = childDir //goto children as  it is NOT the LAST goToDirName
		}

	}
	state.CurrentDirectory = initialStateCurrDir
	return fmt.Sprintf(CommandExecSuccesMsg[cmd.Instruction])
}
