package vector

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

func (v Vector3D) AddFloat(f float64) *Vector3D {
	return &Vector3D{
		X: v.X + f,
		Y: v.Y + f,
		Z: v.Z + f,
	}
}
