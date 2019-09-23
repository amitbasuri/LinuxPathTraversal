package system

import "fmt"

type State struct {
	RootDirectory    *Directory
	CurrentDirectory *Directory
}

func (s *State) String() string {
	return fmt.Sprintf("Current Dir: %s", s.CurrentDirectory)
}

// NewState returns a new State
func NewState() *State {
	childDirectories := make([]*Directory, 0)
	rootDir := &Directory{
		name:             "",
		childDirectories: childDirectories,
		ParentDirectory:  nil, //rootDir dir parent will always be nil
	}
	return &State{RootDirectory: rootDir, CurrentDirectory: rootDir}
}
