package system

import "fmt"

// Directory a single Directory
type Directory struct {
	name             string
	ParentDirectory  *Directory
	childDirectories []*Directory
}

//Validate Validates a Directory
func (dir *Directory) Validate() error {
	if dir.name == ".." || dir.name == "" {
		return fmt.Errorf("invalid Directory name")
	}

	return nil
}

//NewDirectory returns New Directory
func NewDirectory(parentDir *Directory, name string) *Directory {

	childDirectories := make([]*Directory, 0)
	dir := &Directory{
		name:             name,
		childDirectories: childDirectories,
		ParentDirectory:  parentDir,
	}
	return dir
}

func (dir *Directory) String() string {
	return fmt.Sprintf("%v", dir.name)
}

//getChild returns the child of a parent with the index of Directory.childDirectories
func getChild(parentDir *Directory, childDirName string) (*Directory, int) {
	for index, dir := range parentDir.childDirectories {
		if dir.name == childDirName {
			return dir, index
		}
	}
	return nil, 0
}

// addChildDirectory  adds a Child Directory
func (dir *Directory) addChildDirectory(childDirName string) {

	childDir := NewDirectory(dir, childDirName)
	dir.childDirectories = append(dir.childDirectories, childDir)

}

// removeChildDirectory
func (dir *Directory) removeChildDirectory(childDirName string) {

	childDir, index := getChild(dir, childDirName)
	if childDir != nil {
		dir.childDirectories = append(dir.childDirectories[:index], dir.childDirectories[index+1:]...)
	}

}
