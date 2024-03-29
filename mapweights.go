package goutils

import (
	"log/slog"
	"math/rand"
)

type MapWeights struct {
	MapWeights  map[string]int
	TotalWeight int
	DefaultName string
}

func NewMapWeights() *MapWeights {
	return &MapWeights{
		MapWeights: make(map[string]int),
	}
}

func (mapWeights *MapWeights) SetDefaultIsMaxWeight() {
	maxw := 0
	for k, v := range mapWeights.MapWeights {
		if v > maxw {
			maxw = v

			mapWeights.DefaultName = k
		}
	}
}

func (mapWeights *MapWeights) AddWeight(name string, weight int, isDefault bool) error {
	_, isok := mapWeights.MapWeights[name]
	if isok {
		Error("MapWeights.AddWeight",
			slog.String("name", name),
			Err(ErrInvalidNameInMapWeights))

		return ErrInvalidNameInMapWeights
	}

	mapWeights.MapWeights[name] = weight

	mapWeights.TotalWeight += weight

	if isDefault {
		mapWeights.DefaultName = name
	}

	return nil
}

func (mapWeights *MapWeights) Rand() string {
	cr := rand.Int() % mapWeights.TotalWeight
	for k, v := range mapWeights.MapWeights {
		if cr < v {
			return k
		}

		cr -= v
	}

	Error("MapWeights.Rand",
		slog.Int("cr", cr))

	return mapWeights.DefaultName
}
