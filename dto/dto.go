package dto

// ResumeUploadResponse defines the response structure for /resume-upload.
type ResumeUploadResponse struct {
    ChunkIndex  int    `json:"chunk_index"`
    TotalChunks int    `json:"total_chunks"`
    FileName    string `json:"file_name"`
}

// UploadChunkResponse defines the response structure for /upload-chunk.
type UploadChunkResponse struct {
    Message   string `json:"message"`
    RequestID string `json:"request_id"`
}
