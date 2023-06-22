package fixtures

import (
	geo "github.com/gpabois/gostd/geojson"
	"github.com/gpabois/gostd/option"
)

func RandomPoint(props option.Option[geo.FeatureProperties]) geo.Feature {
	return geo.NewPoint(6.50, 4.2, props)
}
