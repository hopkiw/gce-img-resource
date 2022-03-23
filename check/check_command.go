package check

import (
	"context"
	"fmt"

	gceimgresource "github.com/hopkiw/gce-img-resource"
	"google.golang.org/api/compute/v1"
)

type Command struct{}

func NewCommand() *Command {
	return &Command{}
}

/*
A check request. The 'project' field is required.

{
  "source": {
		"project": "some-project",
		"family": "some-family",
		"regexp": "rhel-8-v([0-9]+).*",
  },
  "version": { "name": "rhel-8-v20220322" }
}
*/

type Request struct {
	Source  gceimgresource.Source  `json:"source"`
	Version gceimgresource.Version `json:"version"`
}

type Response []gceimgresource.Version

func (command *Command) Run(request Request) (Response, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return Response{}, err
	}

	call := computeService.Images.List(request.Source.Project) //.OrderBy("creationTimestamp")
	if request.Source.Family != "" {
		call = call.Filter(fmt.Sprintf("family = %s", request.Source.Family))
	}
	if request.Version.Name == "" {
		call = call.OrderBy("creationTimestamp desc")
	}

	var is []*compute.Image
	var pt string
	for il, err := call.PageToken(pt).Do(); ; il, err = call.PageToken(pt).Do() {
		if err != nil {
			return Response{}, err
		}
		is = append(is, il.Items...)
		if il.NextPageToken == "" || request.Version.Name == "" {
			break
		}
		pt = il.NextPageToken
	}

	var response Response
	var start bool
	for _, image := range is {
		if request.Version.Name == "" {
			// No version specified, return only the latest image.
			response = append(response, gceimgresource.Version{Name: image.Name})
			break
		}
		if image.Name == request.Version.Name {
			start = true
		}
		if start {
			response = append(response, gceimgresource.Version{Name: image.Name})
		}
	}

	return response, nil
}
