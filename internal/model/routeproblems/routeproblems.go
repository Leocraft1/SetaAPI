package routeproblems

type Problem struct {
	Num         string `json:"num"`
	HasProblems bool   `json:"has_problems"`
	SiteCode    int    `json:"site_code"`
}

type ProblemCodesResponse struct {
	Problem []Problem `json:"codes"`
}
