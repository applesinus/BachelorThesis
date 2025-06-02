package sat

import (
	tgs "BachelorThesis/engine/collision/resolving/TGS"
	"BachelorThesis/engine/constants"
	"BachelorThesis/engine/objects"
	"log"
	"math"
)

func SATNoParallel(aID, bID int, objectPool *[]objects.Object, resolveAlgorithm string) {
	switch (*objectPool)[aID].(type) {
	case *objects.Sphere:
		switch (*objectPool)[bID].(type) {
		case *objects.Sphere:
			satSphereSphere(aID, bID, objectPool, resolveAlgorithm)
		}

	default:
		log.Panicf("Unknown object type: %T", (*objectPool)[aID])
	}
}

func satSphereSphere(aID, bID int, objectPool *[]objects.Object, resolveAlgorithm string) {
	// Получаем объекты и проверяем, что это сферы
	objA := (*objectPool)[aID]
	sphereA, okA := objA.(*objects.Sphere) // Используем interface assertion с проверкой
	if !okA {
		log.Panicf("Object %d (ID: %s) is not a Sphere", aID, objA.GetId())
		return
	}

	objB := (*objectPool)[bID]
	sphereB, okB := objB.(*objects.Sphere)
	if !okB {
		log.Panicf("Object %d (ID: %s) is not a Sphere", bID, objB.GetId())
		return
	}

	// Получаем позиции
	posA, errA := sphereA.GetPosition()
	if errA != nil {
		log.Panicf("Error getting position of object %d: %v", aID, errA)
	}

	posB, errB := sphereB.GetPosition()
	if errB != nil {
		log.Panicf("Error getting position of object %d: %v", bID, errB)
	}

	radiusA := sphereA.GetRadius()
	radiusB := sphereB.GetRadius()
	sumRadii := radiusA + radiusB

	// Вектор от центра A к центру B - это и есть наша ось для проверки
	// (точнее, его направление после нормализации)
	axisX := posB.X - posA.X
	axisY := posB.Y - posA.Y
	axisZ := posB.Z - posA.Z

	// Квадрат расстояния между центрами
	distanceSq := axisX*axisX + axisY*axisY + axisZ*axisZ

	// Проверка столкновения (как и у тебя, это эффективно)
	// Если квадрат расстояния меньше или равен квадрату суммы радиусов, то есть столкновение
	if distanceSq <= sumRadii*sumRadii {
		// Столкновение обнаружено!

		// Теперь вычислим нормаль столкновения и глубину проникновения,
		// которые понадобятся для резолвера.

		distance := math.Sqrt(distanceSq)

		var penetrationDepth float64

		if distance == 0 {
			// Центры сфер совпадают. Это особый случай.
			// Глубина проникновения максимальна (сумма радиусов).
			// Нормаль можно выбрать произвольно, например, по оси X.
			// Такое обычно не должно происходить при корректном движении.
			penetrationDepth = sumRadii
		} else {
			// Нормаль столкновения - это нормализованный вектор от A к B (или от B к A)
			// Здесь нормаль будет указывать от A на B
			penetrationDepth = sumRadii - distance
		}

		// Убедимся, что глубина проникновения не отрицательная (из-за ошибок float)
		if penetrationDepth < 0 {
			penetrationDepth = 0
		}

		// "Проекции шаров на ось":
		// Для сферы A на ось (с началом в posA и направлением collisionNormal): интервал [-radiusA, radiusA]
		// Для сферы B на ось (с началом в posB и направлением -collisionNormal): интервал [-radiusB, radiusB]
		// Или, если привести к общему началу координат для оси:
		// Проекция центра A на ось: 0 (если ось из A в B)
		// Проекция центра B на ось: distance
		// Интервал A: [-radiusA, radiusA]
		// Интервал B: [distance - radiusB, distance + radiusB]
		// Они пересекаются, если radiusA >= distance - radiusB, что эквивалентно radiusA + radiusB >= distance.
		// Это условие мы уже проверили.

		// Вызываем резолвер
		switch resolveAlgorithm {
		case constants.PGS:
			switch constants.ResolveAlgoType {
			case constants.N:
				tgs.TGSNoParallel(aID, bID, objectPool)
			case constants.PNT:
				// TODO: pgs.PGSParallelWithNormalAndTangent(aID, bID, objectPool, collisionNormal, penetrationDepth, ...)
			default:
				log.Panicf("Unknown ResolveAlgoType for PGS: %s", constants.ResolveAlgoType)
			}

		case constants.NoAlgo:
			return // Просто выходим, если алгоритм не задан
		default:
			log.Panicf("Unknown resolve algorithm: %s", resolveAlgorithm)
		}
	}
}
