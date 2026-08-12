package service

import (
	"setaapi/internal/model"
	"setaapi/internal/repository"
)

func GetStops() []model.Stop {
	return repository.GetStops()
}
