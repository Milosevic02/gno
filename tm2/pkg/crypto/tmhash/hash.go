package tmhash

import (
	"crypto/sha256"
	"hash"
	"sync"
)

const (
	Size      = sha256.Size
	BlockSize = sha256.BlockSize
)

// HashPool provides a pool of hash.Hash instances to reduce memory allocations
var hashPool = sync.Pool{
	New: func() interface{} {
		return sha256.New()
	},
}

// New returns a new hash.Hash from the pool.
func New() hash.Hash {
	return hashPool.Get().(hash.Hash)
}

// ReturnHash returns a hash.Hash to the pool.
func ReturnHash(h hash.Hash) {
	h.Reset()
	hashPool.Put(h)
}

// Sum returns the SHA256 of the bz.
func Sum(bz []byte) []byte {
	h := hashPool.Get().(hash.Hash)
	defer func() {
		h.Reset()
		hashPool.Put(h)
	}()

	h.Write(bz)
	result := make([]byte, Size)
	h.Sum(result[:0])
	return result
}

//-------------------------------------------------------------

const (
	TruncatedSize = 20
)

type sha256trunc struct {
	sha256 hash.Hash
}

func (h sha256trunc) Write(p []byte) (n int, err error) {
	return h.sha256.Write(p)
}

func (h sha256trunc) Sum(b []byte) []byte {
	shasum := h.sha256.Sum(b)
	return shasum[:TruncatedSize]
}

func (h sha256trunc) Reset() {
	h.sha256.Reset()
}

func (h sha256trunc) Size() int {
	return TruncatedSize
}

func (h sha256trunc) BlockSize() int {
	return h.sha256.BlockSize()
}

// truncatedHashPool provides a pool of truncated hash instances
var truncatedHashPool = sync.Pool{
	New: func() interface{} {
		return sha256trunc{
			sha256: sha256.New(),
		}
	},
}

// NewTruncated returns a new hash.Hash from the truncated pool.
func NewTruncated() hash.Hash {
	return truncatedHashPool.Get().(hash.Hash)
}

// ReturnTruncatedHash returns a truncated hash to the pool.
func ReturnTruncatedHash(h hash.Hash) {
	h.Reset()
	truncatedHashPool.Put(h)
}

// SumTruncated returns the first 20 bytes of SHA256 of the bz.
func SumTruncated(bz []byte) []byte {
	h := hashPool.Get().(hash.Hash)
	defer func() {
		h.Reset()
		hashPool.Put(h)
	}()

	h.Write(bz)
	result := make([]byte, sha256.Size)
	h.Sum(result[:0])
	return result[:TruncatedSize]
}
