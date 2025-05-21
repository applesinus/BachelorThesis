package visualizer

import (
	"BachelorThesis/engine/constants"
	"BachelorThesis/engine/objects"
	"BachelorThesis/engine/vector"
	"log"
	"math/rand"
	"runtime"
	"time"

	hg "github.com/harfang3d/harfang-go"
)

const (
	title = "Bachelor Thesis Visualization"
)

type obj struct {
	transform *hg.Transform
	object    objects.Object
}

func Start(algorithm string, endChan chan bool, pool *[]objects.Object) {

	win, pipeline, scene := prepareScene()

	defer hg.RenderShutdown()
	defer hg.DestroyWindow(win)

	//создание конвейера отрисовки

	res := hg.NewPipelineResources()

	sphereRef, sphereMat := createSphereRefAndMat(res)

	//добавление света
	hg.CreateSpotLightWithDiffuseDiffuseIntensitySpecularSpecularIntensityPriorityShadowTypeShadowBias(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(-8.8, 21.7, -8.8), hg.Deg3(60, 45, 0)), 0, hg.Deg(5), hg.Deg(30), hg.ColorGetWhite(), 1, hg.ColorGetWhite(), 1, 0, hg.LSTMap, 0.000005)

	rect := hg.NewIntRectWithSxSyExEy(0, 0, constants.WindowWidth, constants.WindowHeight)

	rendererPool := make(map[string]*obj, 0)
	for _, object := range *pool {
		rendererPool[object.GetId()] = &obj{
			transform: newSphere(scene, sphereRef, sphereMat),
			object:    object,
		}
	}

	//основной цикл
	frame := 0

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		frame++

		newSphereInPool(pool, &rendererPool, scene, sphereRef, sphereMat)

		if frame%100 == 0 {
			/*
				log.Printf("frame %d", frame)
				for _, object := range rendererPool {
					log.Printf("object %s", object.object.GetId())
				}
				log.Println()
			*/
		}

		for _, object := range rendererPool {
			updateObjectOnRenderer(object)
		}

		for _, a := range rendererPool {
			aTime, err := time.Parse(time.RFC3339, a.object.GetId())
			if err != nil {
				continue
			}

			for _, b := range rendererPool {
				bTime, err := time.Parse(time.RFC3339, b.object.GetId())
				if err != nil {
					continue
				}

				if a != b && aTime.Before(bTime) {
					if _, ok := a.object.(*objects.Sphere); !ok {
						continue
					}
					if _, ok := b.object.(*objects.Sphere); !ok {
						continue
					}
					a.object.(*objects.Sphere).Collide(b.object.(*objects.Sphere))
				}
			}
		}

		dt := hg.TickClock()

		//обновление сцены
		scene.Update(dt)

		viewID := uint16(0)
		//отправка сцены в pipeline
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontal(&viewID, scene, rect, true, pipeline, res)

		//обновление экрана из буфера
		hg.Frame()
		hg.UpdateWindow(win)
		runtime.GC()

		select {
		case <-endChan:
			return
		default:
			continue
		}
	}
}

func createSphereRefAndMat(res *hg.PipelineResources) (*hg.ModelRef, *hg.Material) {
	vtxLayout := hg.VertexLayoutPosFloatNormUInt8()
	sphereMdl := hg.CreateSphereModel(vtxLayout, 1, 8, 16)

	shader := hg.LoadPipelineProgramRefFromFile("resources_compiled/core/shader/default.hps", res, hg.GetForwardPipelineInfo())

	return res.AddModel("sphere", sphereMdl),
		hg.CreateMaterialWithValueName0Value0ValueName1Value1(
			shader,
			"uDiffuseColor",
			hg.NewVec4WithXYZ(1, 0, 0),
			"uSpecularColor",
			hg.NewVec4WithXYZ(1, 0.8, 0),
		)
}

func newSphereInPool(pool *[]objects.Object, rendererPool *map[string]*obj, scene *hg.Scene, sphereRef *hg.ModelRef, sphereMat *hg.Material) {
	id := time.Now().Format(time.RFC3339)
	if _, ok := (*rendererPool)[id]; ok {
		return
	}

	posX := float64(rand.Intn(10)-5) / 10
	posY := float64(rand.Intn(10)-5) / 10
	posZ := float64(rand.Intn(10)-5) / 10

	forceX := float64(rand.Intn(10)-5) / 1000
	forceY := float64(rand.Intn(10)-5) / 1000
	forceZ := float64(rand.Intn(10)-5) / 1000

	sphere := objects.NewSphere(1, id)
	sphere.SetPosition(vector.Vector3D{X: posX, Y: posY, Z: posZ})
	sphere.ApplyVelocity(vector.Vector3D{X: forceX, Y: forceY, Z: forceZ})

	*pool = append(*pool, &sphere)

	sphereRenderer := &obj{
		transform: newSphere(scene, sphereRef, sphereMat),
		object:    &sphere,
	}

	pos, err := sphere.GetPosition()
	if err != nil {
		log.Printf("error: %v", err)
		pos = vector.ZeroVector()
	}
	sphereRenderer.transform.SetPos(hg.NewVec3WithXYZ(float32(pos.X), float32(pos.Y), float32(pos.Z)))

	(*rendererPool)[sphere.GetId()] = sphereRenderer
}

func newSphere(scene *hg.Scene, sphereRef *hg.ModelRef, sphereMat *hg.Material) *hg.Transform {
	node := hg.CreateObjectWithSliceOfMaterials(
		scene,
		hg.TranslationMat4(hg.NewVec3WithXYZ(0.1, 0.1, 0.1)),
		sphereRef,
		hg.GoSliceOfMaterial{sphereMat},
	)

	return node.GetTransform()
}

func updateObjectOnRenderer(object *obj) {
	object.object.Update()

	pos, err := object.object.GetPosition()
	if err != nil {
		log.Printf("error: %v", err)
		pos = vector.ZeroVector()
	}
	object.transform.SetPos(hg.NewVec3WithXYZ(float32(pos.X), float32(pos.Y), float32(pos.Z)))

	angle, err := object.object.GetAngle()
	if err != nil {
		log.Printf("error: %v", err)
		angle = vector.ZeroAngle()
	}
	object.transform.SetRot(hg.NewVec3WithXYZ(float32(angle.X), float32(angle.Y), float32(angle.Z)))
}

func prepareScene() (*hg.Window, *hg.ForwardPipeline, *hg.Scene) {
	hg.InputInit()
	hg.WindowSystemInit()

	//подготовка окна для рисования
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags(title, constants.WindowWidth, constants.WindowHeight, hg.RFMSAA4X)
	pipeline := hg.CreateForwardPipelineWithShadowMapResolution(4096)

	//создания.и конфигурация фонового цвета
	scene := hg.NewScene()
	scene.GetCanvas().SetColor(hg.NewColorWithRGB(0.1, 0.1, 0.1))
	scene.GetEnvironment().SetAmbient(hg.NewColorWithRGB(0.1, 0.1, 0.1))

	//установка камеры
	cam := hg.CreateCamera(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(15.5, 5, -6), hg.NewVec3WithXYZ(0.4, -1.2, 0)), 0.01, 100)
	scene.SetCurrentCamera(cam)

	return win, pipeline, scene
}
