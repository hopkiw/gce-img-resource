package gceimgresource

// Source is the configuration specifying which images to return.
type Source struct {
	Project string `json:"project"`
	// Family limits matching images to those with the specified family. Optional.
	Family string `json:"family"`
	// Regexp defines a regular expression to find semver embedded in image names. Optional.
	//Regexp string `json:"regexp"` // TODO: Not yet implemented.
}

// Version represents a single image version.
type Version struct {
	Name string `json:"name"`
	// Version is used for ordering returned images. Defaults to the image creation timestamp. If regexp is
	// specified, Version will instead contain the parsed semver.
	Version string `json:"version,omitempty"`
}
