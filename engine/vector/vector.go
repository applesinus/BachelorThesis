package vector

import "math"

type Vector3D struct {
	X float64
	Y float64
	Z float64
}

func ZeroVector() *Vector3D {
	return &Vector3D{0, 0, 0}
}

func (v Vector3D) Add(v2 Vector3D) *Vector3D {
	return &Vector3D{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

func (v Vector3D) Sub(v2 Vector3D) *Vector3D {
	return &Vector3D{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
		Z: v.Z - v2.Z,
	}
}

func (v Vector3D) AddFloat(f float64) *Vector3D {
	return &Vector3D{
		X: v.X + f,
		Y: v.Y + f,
		Z: v.Z + f,
	}
}

func (v Vector3D) Length() float64 {
	return math.Sqrt(v.LengthSq())
}

func (v Vector3D) LengthSq() float64 {
	return (v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z)
}

func (v Vector3D) Normalize() *Vector3D {
	length := v.Length()
	if length == 0 {
		return &Vector3D{0, 0, 0}
	}
	return &Vector3D{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

func (v Vector3D) Mul(f float64) *Vector3D {
	return &Vector3D{
		X: v.X * f,
		Y: v.Y * f,
		Z: v.Z * f,
	}
}
