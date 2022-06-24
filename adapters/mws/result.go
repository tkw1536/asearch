package mws

type Result struct {
	Total    int64 `json:"total"`
	TookInMS int   `json:"time"`

	Variables []QueryVariable `json:"qvars"`

	MathWebSearchIDs []int64 `json:"ids,omitempty"`
	//  Hits             []*Hit  `json:"hits,omitempty"`
}

type QueryVariable struct {
	Name  string `json:"name"`  // name of the variable
	XPath string `json:"xpath"` // xpath of the variable relative to the root
}
