package app

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/shatylos/ffmpeg-screenshots/tools/apperrors"
	"github.com/shatylos/ffmpeg-screenshots/tools/logger"
)

// Storage persists screenshots by their output name and reads them back.
// Implementations must be safe for concurrent use.
type Storage interface {
	Save(output string, image []byte) error
	Get(output string) ([]byte, bool)
}

// NewStorage returns the Storage implementation selected by config.Storage.
func NewStorage(config Config) (Storage, error) {
	switch config.Storage {
	case StorageDisk:
		return NewDiskStorage(config.Outputdir), nil
	case StorageMemory:
		return NewMemoryStorage(), nil
	default:
		return nil, apperrors.New("unknown storage type: %q", config.Storage)
	}
}

// DiskStorage writes screenshots to files under a base directory and reads
// them back from disk on demand.
type DiskStorage struct {
	dir string
}

func NewDiskStorage(dir string) *DiskStorage {
	return &DiskStorage{dir: dir}
}

func (s *DiskStorage) Save(output string, image []byte) error {
	path := filepath.Join(s.dir, output)
	if err := os.WriteFile(path, image, 0o644); err != nil {
		return apperrors.Wrap(err, "failed to write screenshot %s", path)
	}
	return nil
}

func (s *DiskStorage) Get(output string) ([]byte, bool) {
	data, err := os.ReadFile(filepath.Join(s.dir, output))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			logger.PrintError(apperrors.Wrap(err, "failed to read screenshot %s", output))
		}
		return nil, false
	}
	return data, true
}

// MemoryStorage keeps the latest screenshot for each output in memory.
type MemoryStorage struct {
	mu     sync.RWMutex
	images map[string][]byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{images: map[string][]byte{}}
}

func (s *MemoryStorage) Save(output string, image []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.images[output] = image
	return nil
}

func (s *MemoryStorage) Get(output string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.images[output]
	return data, ok
}
