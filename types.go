package gceimgresource

type Source struct {
	Project string `json:"project"`
	Family  string `json:"family"`
	Regexp  string `json:"regexp"`
}

type Version struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
