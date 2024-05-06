package snapshot

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/LaBatata101/goinsta/internal/gotextdiff"
)

type Snapshot struct {
	Loc     int
	path    string
	Name    string
	Source  string
	Content string
}

func (s Snapshot) Accept() error {
	if s.IsNew() {
		return os.Rename(s.path, strings.TrimSuffix(s.path, ".new"))
	}
	return nil
}

func (s Snapshot) Reject() {
	if s.IsNew() {
		os.Remove(s.path)
	}
}

// Compute the difference between the new snapshot (.snap.new) and the old snapshot (.snap).
// Return the diff string.
func (s Snapshot) Diff() string {
	oldSnapshotPath := strings.TrimSuffix(s.path, ".new")
	_, err := os.Stat(oldSnapshotPath)
	if s.IsNew() {
		if err == nil {
			// Don't need to handle error here, since, we already checkd that `oldSnapshotPath` exist.
			oldSnap, _ := Read(oldSnapshotPath)
			return gotextdiff.Unified(oldSnap.Content, s.Content)
		} else if errors.Is(err, fs.ErrNotExist) {
			return gotextdiff.Unified("", s.Content)
		}
	}
	return ""
}

func (s Snapshot) HasDifference() bool {
	oldSnapshotPath := strings.TrimSuffix(s.path, ".new")
	_, err := os.Stat(oldSnapshotPath)
	if s.IsNew() && err == nil {
		// Don't need to handle error here, since, we already checkd that `oldSnapshotPath` exist.
		oldSnap, _ := Read(oldSnapshotPath)
		edits := gotextdiff.Strings(oldSnap.Content, s.Content)
		return len(edits) > 0
	}
	return false
}

func (s Snapshot) IsNew() bool {
	return strings.HasSuffix(s.path, ".snap.new")
}

// Return the snapshot path relative to the `go.mod` directory.
func (s Snapshot) CleanPath() string {
	goModPath, found := findGoModPath(filepath.Dir(s.path))
	if found {
		return filepath.Join(filepath.Base(goModPath), strings.TrimPrefix(s.path, goModPath))
	}
	return s.path
}

func findGoModPath(dir string) (string, bool) {
	currentDir := dir
	for {
		modPath := filepath.Join(currentDir, "go.mod")
		_, err := os.Stat(modPath)
		if err == nil {
			return currentDir, true
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}
	return "", false
}

// Returns all the `.snap.new` snapshots in the current directory and its sub-directories.
func GetNewSnapshotPaths() ([]string, error) {
	var snapshots []string
	currentDir, err := os.Getwd()

	err = filepath.WalkDir(currentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q\n", path)
			return err
		}

		if strings.HasSuffix(path, ".snap.new") {
			snapshots = append(snapshots, path)
		}

		return nil
	})

	return snapshots, err
}

// Parses a snapshot file into the `Snapshot` struct.
// Returns an error if `snapshotPath` doesn't exist.
func Read(snapshotPath string) (Snapshot, error) {
	bytes, err := os.ReadFile(snapshotPath)
	content := strings.Split(string(bytes), "---")
	header := content[1]
	snapContent := strings.Trim(content[2], "\n")

	header = strings.Trim(header, "\n")
	headerLines := strings.Split(header, "\n")

	source := strings.TrimSpace(strings.Split(headerLines[0], ":")[1])
	assertionLine := strings.TrimSpace(strings.Split(headerLines[1], ":")[1])
	loc, err := strconv.ParseInt(assertionLine, 10, 0)
	if err != nil {
		panic("Failed to convert string to int")
	}

	var name string
	if strings.HasSuffix(snapshotPath, ".new") {
		name = strings.TrimSuffix(filepath.Base(snapshotPath), filepath.Ext(snapshotPath))
		name = strings.TrimSuffix(name, filepath.Ext(name))
		name = strings.ReplaceAll(name, "__", ".")
	} else {
		name = strings.ReplaceAll(strings.TrimSuffix(filepath.Base(snapshotPath), filepath.Ext(snapshotPath)), "__", ".")
	}

	return Snapshot{Source: source, Loc: int(loc), Content: snapContent, Name: name, path: snapshotPath}, err
}

// Writes a `.snap.new` snapshot to `path`.
func Write(path, snapshotName, content, source string, loc int) (Snapshot, error) {
	path = path + ".new"
	file, err := os.Create(path)
	defer file.Close()

	_, err = fmt.Fprintln(file, "---")
	_, err = fmt.Fprintf(file, "source: %s\n", source)
	_, err = fmt.Fprintf(file, "assertion_line: %d\n", loc)
	_, err = fmt.Fprintln(file, "---")
	_, err = fmt.Fprint(file, content)

	return Snapshot{Source: source, Loc: int(loc), Content: content, Name: snapshotName, path: path}, err
}

func RejectAll(paths []string) ([]Snapshot, error) {
	var rejectedSnaps []Snapshot
	for _, snapPath := range paths {
		snap, err := Read(snapPath)
		if err != nil {
			return rejectedSnaps, err
		}

		snap.Reject()
		rejectedSnaps = append(rejectedSnaps, snap)
	}

	return rejectedSnaps, nil
}

func AcceptAll(paths []string) ([]Snapshot, error) {
	var acceptSnaps []Snapshot
	for _, snapPath := range paths {
		snap, err := Read(snapPath)
		if err != nil {
			return acceptSnaps, err
		}

		snap.Accept()
		acceptSnaps = append(acceptSnaps, snap)
	}

	return acceptSnaps, nil
}
