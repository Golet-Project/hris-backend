package entities

import "hroost/shared/primitive"

type Company struct {
	Coordinate primitive.Coordinate `json:"coordinate,omitempty"`
	Address    primitive.String     `json:"address,omitempty"`
}
