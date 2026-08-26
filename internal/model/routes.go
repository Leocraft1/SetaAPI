package model

type LineeRC struct {
	Linea        string  `db:"linea"`
	Rc           string  `db:"rc"`
	Disp_linea   *string `db:"disp_linea"`
	Disp_dest    *string `db:"disp_dest"`
	Desc         *string `db:"desc"`
	Still_exists bool    `db:"still_exists"`
}

type RLResponse struct {
	List []RouteListElement `json:"lines"`
}

type RouteListElement struct {
	Line        string    `json:"line"`
	Route_codes []RCEntry `json:"route_codes"`
}

type RCEntry struct {
	Rc        string  `json:"route_code"`
	Desc      *string `json:"description"`
	Disp_line *string `json:"display_line"`
	Disp_dest *string `json:"display_destination"`
	Exists    bool    `json:"exists"`
}

type StillExists struct {
	Rc           string `db:"rc"`
	Still_exists bool   `db:"still_exists"`
}
