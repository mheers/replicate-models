package replicate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/replicate/replicate-go"
	"github.com/sirupsen/logrus"
)

type Replicate struct {
	client  *replicate.Client
	version string
	token   string
}

func NewReplicate(token, version string, opts ...replicate.ClientOption) (*Replicate, error) {
	opts = append(opts, replicate.WithToken(token))
	client, err := replicate.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return &Replicate{client: client, version: version, token: token}, nil
}

// waitUntilFinished waits until the prediction is finished and returns the updated prediction // TODO: add max retries and timeout
func (r *Replicate) waitUntilFinished(prediction *replicate.Prediction) (*replicate.Prediction, error) {
	succeeded := prediction.Status == "succeeded"
	for !succeeded {
		logrus.Infof("waiting for prediction %s to finish", prediction.ID)
		logrus.Debugf("prediction: %v", prediction)
		var err error
		prediction, err = r.client.GetPrediction(context.Background(), prediction.ID)
		if err != nil {
			return nil, err
		}
		succeeded = prediction.Status == "succeeded"
		if !succeeded {
			time.Sleep(1 * time.Second)
		}
	}

	return prediction, nil
}

func (r *Replicate) download(p *replicate.Prediction, dstPath string, fileEntry int) error {
	p, err := r.waitUntilFinished(p)
	if err != nil {
		return err
	}

	url := ""
	switch p.Output.(type) {
	case string:
		url = p.Output.(string)
	case []interface{}:
		urls := p.Output.([]interface{})
		url = urls[fileEntry].(string)
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	folder := path.Dir(dstPath)
	err = os.MkdirAll(folder, 0777)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (r *Replicate) GetOutputFromID(id string) (any, error) {
	p, err := r.client.GetPrediction(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("GetOutputFromID: error for id %s: %v", id, err)
	}

	p, err = r.waitUntilFinished(p)
	if err != nil {
		return nil, err
	}

	return p.Output, nil
}

func (r *Replicate) DownloadFromID(id, dstPath string, fileEntry int) error {
	p, err := r.client.GetPrediction(context.Background(), id)
	if err != nil {
		return fmt.Errorf("DownloadFromID: error for id %s to %s: %v", id, dstPath, err)
	}
	return r.download(p, dstPath, fileEntry)
}

func (r *Replicate) CreateFileFromBytes(ctx context.Context, data []byte, options *replicate.CreateFileOptions) (*replicate.File, error) {
	f, err := r.client.CreateFileFromBytes(ctx, data, options)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *Replicate) DeleteFile(ctx context.Context, fileID string) error {
	err := r.client.DeleteFile(ctx, fileID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Replicate) CreatePrediction(ctx context.Context, input replicate.PredictionInput, webhook *replicate.Webhook, stream bool) (*replicate.Prediction, error) {
	prediction, err := r.client.CreatePrediction(ctx, r.version, input, webhook, stream)
	if err != nil {
		return nil, err
	}

	return prediction, nil
}

func (r *Replicate) CreatePredictionWithModel(ctx context.Context, modelOwner string, modelName string, input replicate.PredictionInput, webhook *replicate.Webhook, stream bool) (*replicate.Prediction, error) {
	prediction, err := r.client.CreatePredictionWithModel(ctx, modelOwner, modelName, input, webhook, stream)
	if err != nil {
		return nil, err
	}

	return prediction, nil
}

func (r *Replicate) Run(ctx context.Context, input replicate.PredictionInput, webhook *replicate.Webhook, stream bool) (replicate.PredictionOutput, error) {
	prediction, err := r.client.Run(ctx, r.version, input, webhook)
	if err != nil {
		return nil, err
	}

	return prediction, nil
}
