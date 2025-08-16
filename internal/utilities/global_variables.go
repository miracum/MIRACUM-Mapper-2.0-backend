package utilities

import (
	"sync"
)

type ImportStatus struct {
	Progress int // 0-100
	Running  bool
	Error    error
}

var (
	isImporting  bool
	importLock   sync.Mutex
	importStatus ImportStatus = ImportStatus{Progress: 0, Running: false, Error: nil}
)

func TryImporting() bool {
	importLock.Lock()
	defer importLock.Unlock()

	if isImporting {
		return false
	}

	isImporting = true
	return true
}

func DoneImporting() {
	importLock.Lock()
	defer importLock.Unlock()

	isImporting = false
}

func SetImportBegin() {
	importStatus.Progress = 0
	importStatus.Running = true
	importStatus.Error = nil
}

func SetImportProgress(progress int) {
	importStatus.Progress = progress
}

func SetImportDone(err error) {
	if err == nil {
		importStatus.Progress = 100
	}
	importStatus.Running = false
	importStatus.Error = err
}

func GetImportStatus() ImportStatus {
	return importStatus
}
