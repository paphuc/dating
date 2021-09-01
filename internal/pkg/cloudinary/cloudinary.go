package cloudinary

import (
	"bytes"
	"context"

	cloudinary "github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/pkg/errors"
)

type Cloudinary struct {
	*cloudinary.Cloudinary
}

func New(CLOUDINARY_URL string) (*Cloudinary, error) {
	cld, err := cloudinary.NewFromURL(CLOUDINARY_URL)
	return &Cloudinary{
		cld,
	}, err
}

func (c *Cloudinary) UploadFile(ctx context.Context, fileBytes []byte, name string) (string, error) {
	uploadResult, err := c.Upload.Upload(
		context.Background(),
		bytes.NewReader(fileBytes),
		uploader.UploadParams{PublicID: name})
	return uploadResult.SecureURL, err
}

func (c *Cloudinary) DestroyFile(ctx context.Context, name string) (string, error) {
	//Find image
	destroyResult, err := c.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{PublicID: name, Invalidate: true})
	if destroyResult.Result == "not found" {
		//Find video
		destroyResult, err := c.Upload.Destroy(
			context.Background(),
			uploader.DestroyParams{PublicID: name, ResourceType: api.Video, Invalidate: true})

		if destroyResult.Result == "not found" {
			return destroyResult.Result, errors.New("not found")
		}
		return destroyResult.Result, err
	}
	return destroyResult.Result, err
}

func (c *Cloudinary) AssetFile(ctx context.Context, name string) (string, error) {
	//Find image
	assetResult, err := c.Admin.Asset(
		context.Background(),
		admin.AssetParams{PublicID: name, AssetType: api.Image})

	if assetResult.Error.Message != "" {
		//Find video
		assetResult, err = c.Admin.Asset(
			context.Background(),
			admin.AssetParams{PublicID: name, AssetType: api.Video})
		if assetResult.Error.Message != "" {
			return assetResult.SecureURL, errors.New(assetResult.Error.Message)
		}
		return assetResult.SecureURL, err
	}
	return assetResult.SecureURL, err
}
