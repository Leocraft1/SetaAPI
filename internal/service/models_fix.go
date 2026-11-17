package service

import (
	"setaapi/internal/repository/mezzi"
)

func GetModels() []string {
	return mezzi.GetModelsDistinct()
}
