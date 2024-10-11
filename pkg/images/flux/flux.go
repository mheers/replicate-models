package flux

import (
	"context"
	"fmt"

	iReplicate "github.com/mheers/replicate-models/pkg/replicate"
	"github.com/replicate/replicate-go"
	"github.com/sirupsen/logrus"
)

type Flux struct {
	r *iReplicate.Replicate
}

type CreateImage struct {
	AspectRatio          string `json:"aspectRatio"`
	DisableSafetyChecker bool   `json:"disableSafetyChecker"`
	OutputFormat         string `json:"outputFormat"`
	OutputQuality        int    `json:"outputQuality"`
}

func NewFlux(token string) (*Flux, error) {
	r, err := iReplicate.NewReplicate(token, "")
	if err != nil {
		return nil, fmt.Errorf("image create: error creating new replicate: %v", err)
	}

	return &Flux{r: r}, nil
}

func (f *Flux) Create(prompt, dstFile string, options *CreateImage) error {
	if options == nil {
		options = &CreateImage{}
	}
	if options.AspectRatio == "" {
		options.AspectRatio = "1:1"
	}
	if options.OutputFormat == "" {
		options.OutputFormat = "webp"
	}
	if options.OutputQuality == 0 {
		options.OutputQuality = 80
	}

	logrus.Infof("image create: prompt: %s", prompt)
	logrus.Infof("image create: dstFile: %s", dstFile)
	logrus.Infof("image create: options: %v", options)

	p, err := f.r.CreatePredictionWithModel(context.Background(), "black-forest-labs", "flux-schnell", replicate.PredictionInput{
		"prompt":         prompt,
		"num_outputs":    1,
		"aspect_ratio":   options.AspectRatio,
		"output_format":  options.OutputFormat,
		"output_quality": options.OutputQuality,
	}, nil, false)
	if err != nil {
		return fmt.Errorf("image create: error creating prediction with model: %v", err)
	}

	// TODO: use webhooks to get notified when the prediction is finished
	return f.r.DownloadFromID(p.ID, dstFile, 0) // TODO: add possibility to download all outputs
}
