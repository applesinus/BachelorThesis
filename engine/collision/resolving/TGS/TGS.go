package tgs

import (
	"BachelorThesis/engine/objects"
	"BachelorThesis/engine/vector"
	"log"
	"math"
)

const (
	TGS_ITERATIONS = 10
	RESTITUTION    = 0.5
	SLOP           = 0.001
	BAUMGARTE_BIAS = 0.0
)

func TGSNoParallel(aID, bID int, objectPool *[]objects.Object) {
	if aID == bID {
		log.Panicf("Object %d (ID: %s) is the same as object %d (ID: %s)", aID, (*objectPool)[aID].GetId(), bID, (*objectPool)[bID].GetId())
	}

	objA := (*objectPool)[aID]
	objB := (*objectPool)[bID]

	sphereA, okA := objA.(*objects.Sphere)
	if !okA {
		log.Printf("TGS: Объект %d (ID: %s) не является сферой, пропуск разрешения коллизии.", aID, objA.GetId())
		return
	}
	sphereB, okB := objB.(*objects.Sphere)
	if !okB {
		log.Printf("TGS: Объект %d (ID: %s) не является сферой, пропуск разрешения коллизии.", bID, objB.GetId())
		return
	}

	// Получаем свойства объектов
	posA, errA := sphereA.GetPosition()
	velA, errVelA := sphereA.GetVelocity()
	radiusA := sphereA.GetRadius()
	// Предполагается, что у Sphere есть метод GetMass() float64.
	// Если нет, или масса = 0, то по умолчанию масса 1.0 (для движущихся объектов)
	// или 0.0 (для статических/бесконечно массивных объектов).
	massA := 1.0 // Значение по умолчанию, заменить на sphereA.GetMass(), если доступно

	if errA != nil || errVelA != nil {
		log.Printf("TGS: Ошибка получения свойств для объекта А (%s): %v, %v", sphereA.GetId(), errA, errVelA)
		return
	}

	posB, errB := sphereB.GetPosition()
	velB, errVelB := sphereB.GetVelocity()
	radiusB := sphereB.GetRadius()
	massB := 1.0 // Значение по умолчанию, заменить на sphereB.GetMass(), если доступно

	if errB != nil || errVelB != nil {
		log.Printf("TGS: Ошибка получения свойств для объекта B (%s): %v, %v", sphereB.GetId(), errB, errVelB)
		return
	}

	// Вычисляем вектор от B к A (направление нормали, отталкивающей B от A)
	deltaPos := posA.Sub(*posB)

	distanceSq := deltaPos.LengthSq()
	sumRadii := radiusA + radiusB
	sumRadiiSq := sumRadii * sumRadii

	// Если объекты не проникают или только касаются, разрешение не требуется для этого шага.
	if distanceSq > sumRadiiSq {
		return
	}

	distance := math.Sqrt(distanceSq)
	if distance < 1e-9 { // Если центры совпадают или очень близки, выбираем произвольную нормаль
		deltaPos = &vector.Vector3D{X: 1, Y: 0, Z: 0}
		distance = 1.0 // Чтобы избежать деления на ноль
	}

	normal := deltaPos.Normalize() // Нормализованный вектор от B к A

	// Вычисляем обратные массы для определения эффективной массы
	// Если масса равна 0, это означает бесконечную массу (статический объект), поэтому обратная масса равна 0.
	invMassA := 0.0
	if massA > 1e-9 { // Проверка на ненулевую массу
		invMassA = 1.0 / massA
	}
	invMassB := 0.0
	if massB > 1e-9 { // Проверка на ненулевую массу
		invMassB = 1.0 / massB
	}

	// Если оба объекта статичны, импульс применить нельзя
	if invMassA == 0 && invMassB == 0 {
		return
	}

	// Сумма обратных масс, используемая в знаменателе формулы импульса
	effectiveMassInverse := invMassA + invMassB

	// Глубина проникновения
	penetration := sumRadii - distance

	// Итерации TGS
	for i := 0; i < TGS_ITERATIONS; i++ {
		// Вычисляем текущую относительную скорость вдоль нормали
		relativeVelocity := velA.Sub(*velB).Dot(*normal)

		// Вычисляем смещение для позиционной коррекции (стабилизация Баумгарте)
		// Это помогает предотвратить "проваливание" объектов, если они глубоко проникли
		bias := 0.0
		if penetration > SLOP {
			bias = (BAUMGARTE_BIAS / float64(TGS_ITERATIONS)) * (penetration - SLOP)
		}

		// Желаемая относительная скорость после столкновения (учитывая восстановление и смещение)
		// Если объекты уже расходятся, не применяем дальнейший импульс для проникновения
		if relativeVelocity >= 0 && bias == 0 {
			// Объекты уже разделяются или находятся в покое и нет проникновения. Импульс не нужен.
			continue
		}

		// Вычисляем величину импульса (lambda_change), необходимого для разрешения
		// J = -( (1 + e) * v_rel_normal + bias ) / (1/m_A + 1/m_B)
		impulseMagnitude := -((1+RESTITUTION)*relativeVelocity + bias) / effectiveMassInverse

		// Применяем импульсы к текущим скоростям
		// Это "последовательная" часть алгоритма: обновленные скорости используются немедленно.

		// Импульс, применяемый к объекту A (вдоль нормали)
		impulseVecA := normal.Mul(impulseMagnitude * invMassA)
		newVelA := velA.Add(*impulseVecA)

		// Импульс, применяемый к объекту B (в противоположном направлении от нормали)
		impulseVecB := normal.Mul(impulseMagnitude * invMassB)
		newVelB := velB.Sub(*impulseVecB)

		// Обновляем скорости объектов в пуле
		err := sphereA.ApplyVelocity(*newVelA)
		if err != nil {
			log.Printf("TGS: Ошибка применения скорости к объекту А (%s): %v", sphereA.GetId(), err)
		}
		err = sphereB.ApplyVelocity(*newVelB)
		if err != nil {
			log.Printf("TGS: Ошибка применения скорости к объекту B (%s): %v", sphereB.GetId(), err)
		}

		// Обновляем локальные копии скоростей для следующей итерации внутри цикла
		velA = newVelA
		velB = newVelB
	}
}
