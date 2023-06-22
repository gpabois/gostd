package geojson

import "github.com/gpabois/gostd/option"

type FeatureCollection struct {
	Type     string
	Features []Feature
}

type Feature struct {
	Type       string
	Geometry   Geometry
	Properties map[string]any
}

type Geometry struct {
	Type        string
	Coordinates []float64
}

type FeatureProperties map[string]any

func NewFeatureCollection(features ...Feature) FeatureCollection {
	return FeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}
}

func NewPoint(x float64, y float64, props option.Option[FeatureProperties]) Feature {
	return Feature{
		Type: "Feature",
		Geometry: Geometry{
			Type:        "Point",
			Coordinates: []float64{x, y},
		},
		Properties: props.UnwrapOr(func() FeatureProperties { return make(FeatureProperties) }),
	}
}
