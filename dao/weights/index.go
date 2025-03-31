package weights

type Weights struct {
	ID      int     `json:"id"`
	Dcr     float64 `json:"dcr"`
	DnrDpmo float64 `json:"dnr_dpmo"`
	Ce      float64 `json:"ce"`
	Pod     float64 `json:"pod"`
	Cc      float64 `json:"cc"`
	Dex     float64 `json:"dex"`
}
