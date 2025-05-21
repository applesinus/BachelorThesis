package objects

import (
	"BachelorThesis/engine/vector"
)

type BoundingBox struct {
	Min *vector.Vector3D
	Max *vector.Vector3D
}

type Object interface {
	Update()
	GetBoundingBox() (*BoundingBox, error)

	SetPosition(vector.Vector3D)
	GetPosition() (*vector.Vector3D, error)

	ApplyVelocity(vector.Vector3D) error
	GetVelocity() (*vector.Vector3D, error)

	SetAngle(vector.Angle3D)
	GetAngle() (*vector.Angle3D, error)

	ApplyRotation(vector.Angle3D) error
	GetRotation() (*vector.Angle3D, error)

	GetId() string
}
