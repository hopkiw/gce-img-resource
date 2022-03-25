package check

import (
	"context"
	"fmt"

	gceimgresource "github.com/hopkiw/gce-img-resource"
	"google.golang.org/api/compute/v1"
)

/*
{
  "source": {
		"project": "some-project",
		"family": "some-family",
		"regexp": "rhel-8-v([0-9]+).*",
  },
  "version": { "name": "rhel-8-v20220322" }
}
*/

// Request is the input of a resource check.
type Request struct {
	Source  gceimgresource.Source  `json:"source"`
	Version gceimgresource.Version `json:"version"`
}

// Response is the output of a resource check.
type Response []gceimgresource.Version

// Run performs a check for image versions.
func Run(request Request) (Response, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return Response{}, err
	}

	call := computeService.Images.List(request.Source.Project)
	if request.Source.Family != "" {
		call = call.Filter(fmt.Sprintf("family = %s", request.Source.Family))
	}

	var is []*compute.Image
	var pt string
	for il, err := call.PageToken(pt).Do(); ; il, err = call.PageToken(pt).Do() {
		if err != nil {
			return Response{}, err
		}
		is = append(is, il.Items...)
		if il.NextPageToken == "" {
			break
		}
		pt = il.NextPageToken
	}

	// Use this for correct encoding of empty list.
	response := Response{}

	if request.Version.Name == "" && len(is) > 0 {
		// No version specified, return only the latest image.
		image := is[len(is)-1]
		response = append(response, gceimgresource.Version{Name: image.Name})
		return response, nil
	}

	var start bool
	for _, image := range is {
		if image.Name == request.Version.Name {
			// Start appending from the matching version.
			start = true
		}
		if start {
			response = append(response, gceimgresource.Version{Name: image.Name})
		}
	}

	return response, nil
}
