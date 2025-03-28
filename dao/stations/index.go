package stations

type Station struct {
	ID      int    `json:"id"`
	Station string `json:"station"`
	Fan     int    `json:"fan"`
	Great   int    `json:"great"`
	Fair    int    `json:"fair"`
}
