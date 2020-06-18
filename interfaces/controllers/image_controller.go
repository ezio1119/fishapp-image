package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/ezio1119/fishapp-image/conf"
	"github.com/ezio1119/fishapp-image/models"
	"github.com/ezio1119/fishapp-image/pb"
	"github.com/ezio1119/fishapp-image/usecase/interactor"
)

type imageController struct {
	imageInteractor interactor.ImageInteractor
}

func NewImageController(i interactor.ImageInteractor) *imageController {
	return &imageController{i}
}

func (c *imageController) BatchCreateImages(stream pb.ImageService_BatchCreateImagesServer) error {
	ctx, cancel := context.WithTimeout(stream.Context(), conf.C.Sv.TimeoutDuration)
	defer cancel()

	images := []*models.Image{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch x := req.Data.(type) {
		case *pb.BatchCreateImagesReq_Info:
			fmt.Println("BatchCreateImagesReq_Info")
			img := &models.Image{
				OwnerID:   x.Info.OwnerId,
				OwnerType: x.Info.OwnerType,
				Buf:       &bytes.Buffer{},
			}
			images = append(images, img)

		case *pb.BatchCreateImagesReq_Chunk:

			lastImg := images[len(images)-1]

			if _, err := lastImg.Buf.Write(x.Chunk); err != nil {
				return err
			}

		default:
			return fmt.Errorf("BatchCreateImages.Data has unexpected type %T", x)
		}
	}

	if err := c.imageInteractor.BatchCreateImages(ctx, images); err != nil {
		return err
	}

	imgsP, err := convListImageProto(images)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&pb.BatchCreateImagesRes{Images: imgsP})
}
