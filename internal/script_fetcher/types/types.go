package types

import "io/fs"

type ScriptMetadata struct {
	BaseDir  string
	DirFiles []fs.DirEntry
}
