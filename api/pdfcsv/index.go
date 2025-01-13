package pdfcsv

type UploadPdfBody struct {
	File UploadPfdFile `json:"file"`
}

type UploadPfdFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}
