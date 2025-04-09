package tiers

type Tiers struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	FanPlus float64 `json:"fan_plus"`
	Fan     float64 `json:"fan"`
	Great   float64 `json:"great"`
	Fair    float64 `json:"fair"`
}
