package resources

import "boilerplate/internal/domain"

type ImageMDto struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Entity   string `json:"entity"`
	EntityId uint64 `json:"entity_id"`
}

type ImagesMDto struct {
	Items []ImageMDto `json:"data"`
}

func (d ImageMDto) DomainToDtoMass(images []domain.Image) ImagesMDto {
	imgsDto := make([]ImageMDto, len(images))
	for i, item := range images {
		imgsDto[i] = ImageMDto{}.DomainToDto(item)
	}

	return ImagesMDto{Items: imgsDto}
}

func (d ImageMDto) DomainToDto(imageM domain.Image) ImageMDto {
	return ImageMDto{
		Id:       imageM.Id,
		Name:     imageM.Name,
		Entity:   imageM.Entity,
		EntityId: imageM.EntityId,
	}
}
