package interactor

import (
	"context"

	"github.com/ezio1119/fishapp-image/models"
	"github.com/ezio1119/fishapp-image/usecase/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type imageInteractor struct {
	db                *gorm.DB
	imageUploaderRepo repo.ImageUploaderRepo
}

func NewImageInteractor(db *gorm.DB, ur repo.ImageUploaderRepo) *imageInteractor {
	return &imageInteractor{db, ur}
}

type ImageInteractor interface {
	BatchCreateImages(ctx context.Context, images []*models.Image) error
	BatchDeleteImages(ctx context.Context, id []int64) error
	DeleteImagesByOwnerID(ownerType int64, ownerID int64) error
}

func (i *imageInteractor) BatchCreateImages(ctx context.Context, images []*models.Image) error {
	for _, img := range images {
		img.Name = uuid.New().String()

		if err := resizeImage(img); err != nil {
			return err
		}

		if err := i.db.Create(img).Error; err != nil {
			return err
		}

		if err := i.imageUploaderRepo.UploadImage(ctx, img.Buf, img.Name); err != nil {
			return err
		}
	}

	return nil

}

func (i *imageInteractor) BatchDeleteImages(ctx context.Context, id []int64) error {
	return i.db.Where("id IN (?)", id).Delete(&models.Image{}).Error
}

func (i *imageInteractor) DeleteImagesByOwnerID(ownerType int64, ownerID int64) error {
	return i.db.Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).Delete(&models.Image{}).Error
}
