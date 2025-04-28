package vector

type Angle3D struct {
	X float64
	Y float64
	Z float64
}

func ZeroAngle() *Angle3D {
	return &Angle3D{0, 0, 0}
}

func (a Angle3D) Add(a2 Angle3D) *Angle3D {
	return &Angle3D{
		X: a.X + a2.X,
		Y: a.Y + a2.Y,
		Z: a.Z + a2.Z,
	}
}

func (a Angle3D) AddFloat(f float64) *Angle3D {
	return &Angle3D{
		X: a.X + f,
		Y: a.Y + f,
		Z: a.Z + f,
	}
}

func (a *Angle3D) Normalize() {
	xOffset := int(a.X / 360)
	if xOffset != 0 {
		a.X = a.X - 360*float64(xOffset)
	}

	yOffset := int(a.Y / 360)
	if yOffset != 0 {
		a.Y = a.Y - 360*float64(yOffset)
	}

	zOffset := int(a.Z / 360)
	if zOffset != 0 {
		a.Z = a.Z - 360*float64(zOffset)
	}
}
