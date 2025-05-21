package objects

import (
	"BachelorThesis/engine/vector"
	"fmt"
	"math"
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

	// temp
	if s.position.X > 10 || s.position.X < -10 {
		s.velocity.X = -s.velocity.X
	}
	if s.position.Y > 10 || s.position.Y < -10 {
		s.velocity.Y = -s.velocity.Y
	}
	if s.position.Z > 10 || s.position.Z < -10 {
		s.velocity.Z = -s.velocity.Z
	}
}

// temp
func (s *Sphere) Collide(other *Sphere) {
	// 0. Оптимизация через boundingBox (если есть и используется)
	// if s.boundingBox.Min.X > other.boundingBox.Max.X || s.boundingBox.Max.X < other.boundingBox.Min.X {
	// 	return
	// }
	// // ... и так далее для Y и Z

	// 1. Вектор от центра s к центру other (ось столкновения)
	collisionAxis := other.position.Sub(*s.position)
	distSq := collisionAxis.LengthSq() // Квадрат расстояния

	// 2. Проверка на реальное столкновение по радиусам
	sumRadii := s.radius + other.radius
	if distSq == 0 { // Центры совпадают - особая ситуация, избегаем деления на ноль
		// Можно немного растолкнуть их в случайном направлении или вернуть
		// Например, сдвинуть other немного в произвольном направлении
		// other.position = other.position.Add(Vector3D{X: 0.01, Y: 0, Z: 0})
		return
	}
	if distSq > sumRadii*sumRadii { // Нет столкновения
		return
	}

	// 3. Нормаль столкновения (единичный вектор вдоль оси столкновения)
	// Нормализуем collisionAxis. Если collisionAxis.LengthSq() == 0, Normalize вернет нулевой вектор
	// что приведет к нулевому impulseScalar, так что это безопасно.
	normal := collisionAxis.Normalize()
	if normal.LengthSq() < 1e-9 { // Дополнительная проверка на случай, если Normalize вернул почти нулевой вектор
		return // Если нормаль нулевая (например, центры очень близки), дальнейшие расчеты бессмысленны
	}

	// 4. Относительная скорость
	relativeVelocity := other.velocity.Sub(*s.velocity)

	// 5. Скорость сближения вдоль нормали
	velAlongNormal := Dot(*relativeVelocity, *normal)

	// 6. Если скорости вдоль нормали положительны, значит шары уже удаляются друг от друга
	if velAlongNormal > 0 {
		return
	}

	// 7. Коэффициент восстановления (0 для абсолютно неупругого, 1 для абсолютно упругого)
	cor := 0.9 // У вас был 0.9, это почти упругий удар

	// 8. Расчет скаляра импульса.
	// Массы равны 1, поэтому формула упрощается: j = -(1+e) * v_rel_normal / (1/m1 + 1/m2)
	// Для m1=m2=1, (1/m1 + 1/m2) = 2.
	// j = -(1+e) * velAlongNormal / 2
	// velAlongNormal отрицателен для сближения, поэтому j будет положительным.
	impulseScalar := -(1 + cor) * velAlongNormal / 2.0 // Делим на 2, т.к. две массы участвуют

	// 9. Применяем импульс к скоростям (массы равны 1)
	// v1_new = v1_old - (impulseScalar/m1) * normal
	// v2_new = v2_old + (impulseScalar/m2) * normal
	// Так как m1=m2=1:
	// Вектор импульса для s: normal.Mul(-impulseScalar)
	// Вектор импульса для other: normal.Mul(impulseScalar)

	// s получает импульс ПРОТИВ нормали (normal указывает от s к other)
	s.velocity = s.velocity.Sub(*normal.Mul(impulseScalar))
	// other получает импульс ПО нормали
	other.velocity = other.velocity.Add(*normal.Mul(impulseScalar))

	// 10. (Опционально, но очень рекомендуется) Разрешение проникновения (статическое разрешение)
	// Если шары проникли друг в друга, их нужно немного растолкнуть,
	// чтобы они касались, а не пересекались.
	// Это важно для предотвращения "залипания" или некорректных повторных коллизий.
	distance := math.Sqrt(distSq) // Фактическое расстояние
	penetrationDepth := sumRadii - distance
	if penetrationDepth > 0 {
		// Коэффициент для смягчения "выталкивания", чтобы не было дрожания
		// Можно использовать 0.2-0.8. Или просто 1.0 для полного выталкивания за один шаг.
		// Для простоты здесь 0.5 для каждого шара (в сумме 1.0)
		// Либо можно двигать только один, или пропорционально массам (но массы у нас 1)
		correctionAmount := penetrationDepth / 2.0 // Каждый шар отодвигается на половину глубины проникновения
		correctionVector := normal.Mul(correctionAmount)

		s.position = s.position.Sub(*correctionVector)
		other.position = other.position.Add(*correctionVector)
	}
}

func Dot(v1 vector.Vector3D, v2 vector.Vector3D) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
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

func (s *Sphere) ApplyVelocity(velocity vector.Vector3D) error {
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
