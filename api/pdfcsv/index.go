package pdfcsv

type UploadPdfBody struct {
	File UploadPfdFile `json:"file"`
}

type UploadPfdFile struct {
	LastModified int64  `json:"lastModified"`
	Name         string `json:"name"`
	Size         int16  `json:"size"`
	Type         string `json:"type"`
	File         string `json:"file"`
}
