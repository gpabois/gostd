package slippy_map

import (
	"math"

	"github.com/gpabois/gostd/geojson"
	"github.com/gpabois/gostd/ops"
	"github.com/gpabois/gostd/result"
)

type TileIndex struct {
	X int
	Y int
	Z int
}

func (tile TileIndex) Right() TileIndex {
	return TileIndex{
		X: tile.X + 1,
		Y: tile.Y,
		Z: tile.Z,
	}
}

func (tile TileIndex) Left() TileIndex {
	return TileIndex{
		X: tile.X - 1,
		Y: tile.Y,
		Z: tile.Z,
	}
}

func (tile TileIndex) Up() TileIndex {
	return TileIndex{
		X: tile.X,
		Y: tile.Y - 1,
		Z: tile.Z,
	}
}

func (tile TileIndex) Down() TileIndex {
	return TileIndex{
		X: tile.X,
		Y: tile.Y + 1,
		Z: tile.Z,
	}
}

func (tile TileIndex) UpLeft() TileIndex {
	return tile.Up().Left()
}

func (tile TileIndex) DownRight() TileIndex {
	return tile.Down().Right()
}

func (tile TileIndex) ZoomOut(out int) TileIndex {
	return TileIndex{
		X: int(float64(tile.X) * math.Pow(-2.0, float64(out))),
		Y: int(float64(tile.Y) * math.Pow(-2.0, float64(out))),
		Z: tile.Z - out,
	}
}

func (tile TileIndex) ZoomIn(in int) TileIndex {
	return TileIndex{
		X: int(float64(tile.X) * math.Pow(2.0, float64(in))),
		Y: int(float64(tile.Y) * math.Pow(2.0, float64(in))),
		Z: tile.Z + in,
	}
}

func (tile TileIndex) Upscale(in int) TileBounds {
	upLeft := tile.ZoomIn(4)
	downRight := tile.DownRight().ZoomIn(4).UpLeft()

	return Bounds(upLeft, downRight)
}

func TileIndexFromLatLng(lat float64, lng float64, zoom int) TileIndex {
	n := math.Exp2(float64(zoom))
	x := int(math.Floor((lng + 180.0) / 360.0 * n))

	if float64(x) >= n {
		x = int(n - 1)
	}

	y := int(math.Floor((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * n))

	return TileIndex{
		X: x,
		Y: y,
		Z: zoom,
	}
}

func (t TileIndex) IntoLatLng() (float64, float64) {
	n := math.Pi - 2.0*math.Pi*float64(t.Y)/math.Exp2(float64(t.Z))
	lat := 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	lng := float64(t.X)/math.Exp2(float64(t.Z))*360.0 - 180.0
	return lat, lng
}

// Returns Tile Index from Point[lng, lat]
func TileIndexFromGeometry(geometry geojson.Geometry, zoom int) result.Result[TileIndex] {
	lng := geometry.Coordinates[0]
	lat := geometry.Coordinates[1]

	return result.Success(TileIndexFromLatLng(lat, lng, zoom))
}

type TileBounds struct {
	UpLeft    TileIndex
	DownRight TileIndex
}

func Bounds(upLeft TileIndex, downRight TileIndex) TileBounds {
	return TileBounds{
		UpLeft:    upLeft,
		DownRight: downRight,
	}
}

func (bounds TileBounds) DX() int {
	return bounds.MaxX() - bounds.MinX()
}

func (bounds TileBounds) DY() int {
	return bounds.MaxY() - bounds.MinY()
}

func (bounds TileBounds) MinX() int {
	return ops.Min(bounds.UpLeft.X, bounds.DownRight.X)
}

func (bounds TileBounds) MinY() int {
	return ops.Min(bounds.UpLeft.Y, bounds.DownRight.Y)
}

func (bounds TileBounds) MaxX() int {
	return ops.Max(bounds.UpLeft.X, bounds.DownRight.X)
}

func (bounds TileBounds) MaxY() int {
	return ops.Max(bounds.UpLeft.Y, bounds.DownRight.Y)
}

func (bounds TileBounds) MaxZ() int {
	return ops.Max(bounds.UpLeft.Z, bounds.DownRight.Z)
}

func (bounds TileBounds) MinZ() int {
	return ops.Min(bounds.UpLeft.Z, bounds.DownRight.Z)
}
