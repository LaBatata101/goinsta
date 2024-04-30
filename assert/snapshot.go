package assert

import (
	"errors"
	"gotextdiff/myers"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/LaBatata101/goinsta/internal/snapshot"
	"github.com/LaBatata101/goinsta/internal/ui"
	"github.com/sanity-io/litter"
)

const snapshotDirPath = "testdata/snapshots"

func getParentCallerFuncName() (string, string, int) {
	pc, sourceFile, loc, ok := runtime.Caller(2)
	if !ok {
		log.Fatal("Failed to get parent caller function name")
	}

	fn := runtime.FuncForPC(pc)
	return fn.Name(), sourceFile, loc
}

func Snapshot(t *testing.T, value any) {
	if _, err := os.Stat(snapshotDirPath); errors.Is(err, fs.ErrNotExist) {
		err := os.MkdirAll(snapshotDirPath, 0755)
		if err != nil {
			t.Fatal("An error ocurred while creating the snapshot directory: ", err)
		}
	}

	callerFuncName, sourceFile, loc := getParentCallerFuncName()
	callerFuncName = filepath.Base(callerFuncName)
	snapshotName := strings.ReplaceAll(callerFuncName, ".", "__") + ".snap"
	snapshotPath := filepath.Join(snapshotDirPath, snapshotName)
	snapshotFullPath, err := filepath.Abs(snapshotPath)
	if err != nil {
		t.Fatal("An error ocurred while creating absulute path for snapshot file: ", err)
	}

	sq := litter.Options{
		HidePrivateFields: false,
	}
	newContent := sq.Sdump(value) + "\n"

	_, err = os.Stat(snapshotPath)
	if err == nil {
		// Don't need to handle error here, since, we already checkd that `snapshotPath` exist.
		snap, _ := snapshot.Read(snapshotPath)
		edits := myers.ComputeEdits(newContent, snap.Content)

		// Only show summary if the snapshot content was changed
		if len(edits) > 0 {
			snap, err := snapshot.Write(snapshotFullPath, callerFuncName, newContent, sourceFile, loc)
			if err != nil {
				t.Fatal("An error ocurred while creating new snapshot file: ", err)
			}

			ui.RenderSnapshotSummary(&snap)
			t.Fail()
		}
	} else if errors.Is(err, fs.ErrNotExist) {
		snap, err := snapshot.Write(snapshotFullPath, callerFuncName, newContent, sourceFile, loc)
		if err != nil {
			t.Fatal("An error ocurred while creating new snapshot file: ", err)
		}

		// TODO: keep this as a log?
		infoLog := log.New(os.Stdout, ui.BoldText.Render("INFO: "), 0)
		infoLog.Printf("%s %s", ui.GreenText.Render("stored new snapshot"),
			ui.GreenText2Underlined.Render(snapshotFullPath+".new"))

		ui.RenderSnapshotSummary(&snap)
		t.Fail()
	} else {
		t.Fatal("An error ocurred while checking snapshot file: ", err)
	}
}
