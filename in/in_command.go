package in

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	gceimgresource "github.com/hopkiw/gce-img-resource"
	"google.golang.org/api/compute/v1"
)

type Command struct{}

func NewCommand() *Command {
	return &Command{}
}

type Request struct {
	Source  gceimgresource.Source  `json:"source"`
	Version gceimgresource.Version `json:"version"`
	Params  Params                 `json:"params"`
}

type Params struct {
	SkipDownload string `json:"skip_download"`
}

type Response struct {
	Version  gceimgresource.Version    `json:"version"`
	Metadata []gceimgresource.Metadata `json:"metadata"`
}

func (command *Command) Run(destinationDir string, request Request) (Response, error) {
	err := os.MkdirAll(destinationDir, 0755)
	if err != nil {
		return Response{}, err
	}

	url := fmt.Sprintf("projects/%s/global/images/%s", request.Source.Project, request.Version.Name)

	ctx := context.Background()
	computeService, err := compute.NewService(ctx)

	image, err := computeService.Images.Get(request.Source.Project, request.Version.Name).Do()
	if err != nil {
		return Response{}, err
	}

	if err := writeOutput(destinationDir, "name", request.Version.Name); err != nil {
		return Response{}, err
	}
	if err := writeOutput(destinationDir, "version", request.Version.Version); err != nil {
		return Response{}, err
	}
	if err := writeOutput(destinationDir, "url", url); err != nil {
		return Response{}, err
	}

	return Response{
		Version: gceimgresource.Version{
			Name:    request.Version.Name,
			Version: request.Version.Version,
		},
		Metadata: []gceimgresource.Metadata{
			{
				Name:  "image_id",
				Value: fmt.Sprintf("%d", image.Id),
			},
			{
				Name:  "description",
				Value: image.Description,
			},
			{
				Name:  "creation_timestamp",
				Value: image.CreationTimestamp,
			},
		},
	}, nil
}

func writeOutput(destinationDir, filename string, content string) error {
	return ioutil.WriteFile(filepath.Join(destinationDir, filename), []byte(content), 0644)
}
