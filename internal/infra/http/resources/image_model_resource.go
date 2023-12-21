package resources

import "boilerplate/internal/domain"

type ImageMDto struct {
	Name     string `json:"name"`
	Data     string `json:"data"`
	Entity   string `json:"entity"`
	EntityId uint64 `json:"entity_id"`
}

type ImagesMDto struct {
	Items []ImageMDto `json:"items"`
	Total uint64      `json:"total"`
	Pages uint        `json:"pages"`
}

func (d ImageMDto) DomainToDto(imageM domain.Image) ImageMDto {
	return ImageMDto{
		Name:     imageM.Name,
		Data:     imageM.Data,
		Entity:   imageM.Entity,
		EntityId: imageM.EntityId,
	}
}
