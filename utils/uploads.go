package utils

import (
    "sync"
)

// UploadTracker keeps track of ongoing uploads for file names
type UploadTracker struct {
    mu       sync.Mutex
    uploads  map[string]int
}

// NewUploadTracker initializes a new UploadTracker
func NewUploadTracker() *UploadTracker {
    return &UploadTracker{
        uploads: make(map[string]int),
    }
}

// IncrementChunkCount increments the chunk count for a given file name
func (ut *UploadTracker) IncrementChunkCount(fileName string, totalChunks int) {
    ut.mu.Lock()
    defer ut.mu.Unlock()
    if count, exists := ut.uploads[fileName]; exists {
        ut.uploads[fileName] = count + 1
    } else {
        ut.uploads[fileName] = 1
    }
}

// CheckCompletion checks if all chunks are uploaded
func (ut *UploadTracker) CheckCompletion(fileName string, totalChunks int) bool {
    ut.mu.Lock()
    defer ut.mu.Unlock()
    return ut.uploads[fileName] >= totalChunks
}

// ResetUpload resets the upload count for a file name
func (ut *UploadTracker) ResetUpload(fileName string) {
    ut.mu.Lock()
    defer ut.mu.Unlock()
    delete(ut.uploads, fileName)
}
