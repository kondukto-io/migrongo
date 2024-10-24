package types

import "io/fs"

type ScriptMetadata struct {
	DirFiles []fs.DirEntry
}
