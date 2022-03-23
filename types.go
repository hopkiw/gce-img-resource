package gceimgresource

//

type Source struct {
	Project string `json:"project"`
	// Family is optional, it filters
	Family string `json:"family"`
	// TODO: we should support regexp being optional - if it's not provided, we just return images sorted by creation date.
	Regexp string `json:"regexp"`
}

type Version struct {
	Name string `json:"name"`
	// version is either the semver in the regexp, or it is the creation date (in unix epoch?)
	Version string `json:"version,omitempty"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
