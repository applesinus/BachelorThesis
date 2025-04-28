package objects

import (
	"BachelorThesis/engine/vector"
	"fmt"
)

type Sphere struct {
	radius float64
	id     string

	position *vector.Vector3D
	velocity *vector.Vector3D

	angle    *vector.Angle3D
	rotation *vector.Angle3D

	boundingBox *BoundingBox
}

func (s *Sphere) Update() {
	s.SetPosition(*s.position.Add(*s.velocity))
	s.SetAngle(*s.angle.Add(*s.rotation))

	// sphere bounding box
	s.boundingBox = &BoundingBox{
		Min: s.position.AddFloat(-s.radius),
		Max: s.position.AddFloat(s.radius),
	}
}

func (s *Sphere) GetBoundingBox() (*BoundingBox, error) {
	if s.boundingBox == nil {
		return nil, fmt.Errorf("bounding box of %s is not set", s.id)
	}

	return s.boundingBox, nil
}

// Standart object behavior

func NewSphere(radius float64, id string) Sphere {
	return Sphere{
		id:     id,
		radius: radius,

		position: vector.ZeroVector(),
		velocity: vector.ZeroVector(),

		angle:    vector.ZeroAngle(),
		rotation: vector.ZeroAngle(),

		boundingBox: &BoundingBox{
			Min: vector.ZeroVector().AddFloat(-radius),
			Max: vector.ZeroVector().AddFloat(radius),
		},
	}
}

func (s *Sphere) GetId() string {
	return s.id
}

func (s *Sphere) SetPosition(position vector.Vector3D) {
	s.position = &position
}

func (s *Sphere) GetPosition() (*vector.Vector3D, error) {
	if s.position == nil {
		return nil, fmt.Errorf("position of %s is not set", s.id)
	}

	return s.position, nil
}

func (s *Sphere) AppllyVelocity(velocity vector.Vector3D) error {
	if s.velocity == nil {
		return fmt.Errorf("velocity of %s is not set", s.id)
	}

	s.velocity = s.velocity.Add(velocity)
	return nil
}

func (s *Sphere) GetVelocity() (*vector.Vector3D, error) {
	if s.velocity == nil {
		return nil, fmt.Errorf("velocity of %s is not set", s.id)
	}

	return s.velocity, nil
}

func (s *Sphere) SetAngle(angle vector.Angle3D) {
	s.angle = &angle
	s.angle.Normalize()
}

func (s *Sphere) GetAngle() (*vector.Angle3D, error) {
	if s.angle == nil {
		return nil, fmt.Errorf("angle of %s is not set", s.id)
	}

	return s.angle, nil
}

func (s *Sphere) ApplyRotation(rotation vector.Angle3D) error {
	if s.rotation == nil {
		return fmt.Errorf("rotation of %s is not set", s.id)
	}

	s.rotation = s.rotation.Add(rotation)
	return nil
}

func (s *Sphere) GetRotation() (*vector.Angle3D, error) {
	if s.rotation == nil {
		return nil, fmt.Errorf("rotation of %s is not set", s.id)
	}

	return s.rotation, nil
}
