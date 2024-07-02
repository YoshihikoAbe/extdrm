package extdrm

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path"
	"time"
)

type VerifyResult struct {
	Time time.Time

	Broken       []string `json:"broken"`
	Missing      []string `json:"missing"`
	TotalBroken  int      `json:"total_broken"`
	TotalMissing int      `json:"total_missing"`
	TotalFiles   int      `json:"total_files"`
}

func VerifyFS(root string, metadata *Metadata) (*VerifyResult, error) {
	result := &VerifyResult{
		Time:    time.Now(),
		Broken:  []string{},
		Missing: []string{},
	}
	hash := sha1.New()
	result.TotalFiles = len(metadata.Files)
	for _, entry := range metadata.Files {
		rd, err := os.Open(path.Join(root, entry.SPath))
		if err != nil {
			result.Missing = append(result.Missing, entry.SPath)
			result.TotalMissing++
			continue
		}

		if _, err := io.Copy(hash, rd); err != nil {
			rd.Close()
			return nil, err
		}
		rd.Close()

		if entry.SSha1 != hex.EncodeToString(hash.Sum(nil)) {
			result.Broken = append(result.Broken, entry.SPath)
			result.TotalBroken++
		}
		hash.Reset()
	}
	return result, nil
}
