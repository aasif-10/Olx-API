package handlers

type PresignRequest struct {
	Files []PresignFile `json:"files"`
}

type PresignFile struct {
	ContentType string `json:"file"`
	SizeBytes   int64  `json:"size_bytes"`
}
