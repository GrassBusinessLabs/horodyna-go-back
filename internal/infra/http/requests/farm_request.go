package requests

import (
	"boilerplate/internal/domain"
)

type PointRequest struct {
	Lat float64 `json:"lat" validate:"required"`
	Lng float64 `json:"lng" validate:"required"`
}
type PointsRequest struct {
	UpperLeftPoint   PointRequest `json:"upper_left_point" validate:"required"`
	BottomRightPoint PointRequest `json:"bottom_right_point" validate:"required"`
	Category         string       `json:"category"`
}

type FarmRequest struct {
	Name      *string `json:"name"`
	City      string  `json:"city" validate:"required"`
	Address   string  `json:"address" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

func (p PointsRequest) ToDomainModel() (interface{}, error) {
	return domain.Points{
		UpperLeftPoint:   p.UpperLeftPoint.ToDomainModel(),
		BottomRightPoint: p.BottomRightPoint.ToDomainModel(),
		Category:         p.Category,
	}, nil
}

func (m PointRequest) ToDomainModel() domain.Point {
	return domain.Point{
		Lat: m.Lat,
		Lng: m.Lng,
	}
}

func (m FarmRequest) ToDomainModel() (interface{}, error) {
	return domain.Farm{
		Name:      m.Name,
		City:      m.City,
		Address:   m.Address,
		Longitude: m.Longitude,
		Latitude:  m.Latitude,
	}, nil
}
