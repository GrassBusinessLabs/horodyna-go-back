package resources

import "boilerplate/internal/domain"

type ImageMDto struct {
	Id       uint64 `json:"id"`
	Title    string `json:"title"`
	Entity   string `json:"entity"`
	EntityId uint64 `json:"entity_id"`
}

type ImagesMDto struct {
	Items []ImageMDto `json:"items"`
	Total uint64      `json:"total"`
	Pages uint        `json:"pages"`
}

func (d ImageMDto) DomainToDto(imageM domain.ImageModel) ImageMDto {
	return ImageMDto{
		Id:       imageM.Id,
		Title:    imageM.Title,
		Entity:   imageM.Entity,
		EntityId: imageM.EntityId,
	}
}
