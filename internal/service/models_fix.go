package service

import "setaapi/internal/repository"

func GetModels() []string {
	return repository.GetModelsDistinct();
}