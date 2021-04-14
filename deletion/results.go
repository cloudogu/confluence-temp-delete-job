package deletion

import "os"

// Results keeps statistics about the deletion process.
type Results struct {
	deleted       int
	deletedSizeKB int64
	failed        int
	skipped       int
}

// PrintStats prints deletion statistics as one-liner.
func (r *Results) PrintStats() {
	log.Infof("objects deleted: %d, skipped: %d, failed: %d", r.deleted, r.skipped, r.failed)
}

func (r *Results) fail(path string, err error) {
	log.Debugf("failed: %s with error '%v'", path, err)
	r.failed++
}

func (r *Results) pass(path string, info os.FileInfo) {
	sizeKB := info.Size() / 1024
	log.Debugf("deleted: %s (%d KB)", path, sizeKB)

	r.deleted++
	r.deletedSizeKB += sizeKB
}

func (r *Results) skip(path string) {
	log.Debugf("skipped: %s", path)
	r.skipped++
}
