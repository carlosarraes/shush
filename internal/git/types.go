package git

type GitStatus struct {
	IsRepo  bool
	RootDir string
}

type FileStatus int

const (
	StatusStaged FileStatus = iota
	StatusUnstaged
	StatusUntracked
)

type FileChange struct {
	Path       string
	Status     FileStatus
	LineRanges []LineRange
}

type LineRange struct {
	Start int
	End   int
}
