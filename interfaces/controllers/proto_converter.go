package controllers

import (
	"github.com/ezio1119/fishapp-image/models"
	"github.com/ezio1119/fishapp-image/pb"
	"github.com/golang/protobuf/ptypes"
)

func convImageProto(i *models.Image) (*pb.Image, error) {
	cAt, err := ptypes.TimestampProto(i.CreatedAt)
	if err != nil {
		return nil, err
	}
	uAt, err := ptypes.TimestampProto(i.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &pb.Image{
		Id:        i.ID,
		Name:      i.Name,
		OwnerId:   i.OwnerID,
		OwnerType: i.OwnerType,
		CreatedAt: cAt,
		UpdatedAt: uAt,
	}, nil
}

func convListImageProto(list []*models.Image) ([]*pb.Image, error) {
	listI := make([]*pb.Image, len(list))
	for i, image := range list {
		iProto, err := convImageProto(image)
		if err != nil {
			return nil, err
		}
		listI[i] = iProto
	}
	return listI, nil
}
