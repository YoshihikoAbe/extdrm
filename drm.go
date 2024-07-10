package extdrm

import (
	"bytes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"math/big"
	"os"
	"path"

	"github.com/dgryski/go-camellia"
)

type DrmConfig struct {
	KeySalt       []byte        `json:"key_salt"`
	PathSalt      []byte        `json:"path_salt"`
	Magic         *big.Int      `json:"magic"`
	DisableCRC    bool          `json:"disable_crc"`
	SaltGenerator SaltGenerator `json:"salt_generator"`
}

type DrmFile struct {
	io.Reader
	io.Closer
	Path string
}

type MetadataEntry struct {
	DPath string
	SPath string
	SSha1 string
}

func (entry *MetadataEntry) UnmarshalJSON(b []byte) error {
	m := map[string]any{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			continue
		}
		switch k {
		case "pathh":
			fallthrough
		case "hashed_path":
			fallthrough
		case "dpath":
			entry.DPath = s

		case "path":
			fallthrough
		case "spath":
			entry.SPath = s

		case "sha1":
			fallthrough
		case "ssha1":
			entry.SSha1 = s
		}
	}
	return nil
}

type Metadata struct {
	Files []MetadataEntry `json:"files"`
}

func ReadFS(config DrmConfig, root string) (chan DrmFile, error) {
	ch := make(chan DrmFile, 1)
	state := &readState{
		DrmConfig: config,
		root:      root,
		ch:        ch,

		sha384: sha512.New384(),
	}
	if err := state.run(); err != nil {
		return nil, err
	}
	return ch, nil
}

type readState struct {
	DrmConfig
	root string
	ch   chan DrmFile

	sha384 hash.Hash
}

func (state *readState) run() error {
	metadata, err := state.readMetadata()
	if err != nil {
		return err
	}

	go func() {
		for _, entry := range metadata.Files {
			f, err := state.openFile(entry.DPath)
			if err != nil {
				fmt.Println(entry.SPath+":", "skipping file:", err)
				continue
			}

			state.ch <- DrmFile{
				Reader: state.makeReader(f, entry.SPath),
				Path:   entry.SPath,
				Closer: f,
			}
		}
		close(state.ch)
	}()
	return nil
}

func (state *readState) readMetadata() (*Metadata, error) {
	const filename = "/__metadata.metatxt"

	fp, err := state.openFile(state.hashPath(filename))
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	b, err := io.ReadAll(state.makeReader(fp, filename))
	if err != nil {
		return nil, err
	}
	mdata := &Metadata{}
	if err := json.Unmarshal(b, mdata); err != nil {
		return nil, err
	}

	// the first file written to the channel is always
	// the metadata file
	state.ch <- DrmFile{
		Reader: bytes.NewReader(b),
		Closer: io.NopCloser(nil),
		Path:   filename,
	}

	return mdata, nil
}

func (state *readState) openFile(filename string) (*os.File, error) {
	return os.Open(path.Join(state.root, filename))
}

func (state *readState) makeReader(rd io.Reader, spath string) io.Reader {
	key, iv := state.deriveKey(spath)
	block, _ := camellia.New(key)
	rd = &cipher.StreamReader{
		S: newCTR(block, state.Magic, iv),
		R: rd,
	}

	if state.DisableCRC {
		return rd
	}
	return &crcSkipReader{
		rd: rd,
	}
}

func (state *readState) deriveKey(spath string) (key []byte, iv *big.Int) {
	salt := state.SaltGenerator.KeySalt(state.KeySalt, spath)
	state.sha384.Write(salt)
	state.sha384.Write([]byte(spath))
	sum := state.sha384.Sum(nil)

	key = sum[:32]
	reverseByteSlice(key)
	iv = big.NewInt(0).SetBytes(sum[32:])

	state.sha384.Reset()
	return
}

func (state *readState) hashPath(spath string) string {
	salt := state.SaltGenerator.PathSalt(state.PathSalt, spath)
	hash := sha1.New()
	hash.Write(salt)
	hash.Write([]byte(spath))
	sum := hex.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("/%c/%c/%c/%s", sum[0], sum[2], sum[4], sum)
}

func reverseByteSlice(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}
