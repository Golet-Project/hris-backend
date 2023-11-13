package entities

import "hris/module/shared/primitive"

type Company struct {
	Coordinate primitive.Coordinate `json:"coordinate,omitempty"`
	Address    primitive.String     `json:"address,omitempty"`
}
