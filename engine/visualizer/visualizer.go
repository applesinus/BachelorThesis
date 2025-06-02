package visualizer

import (
	"BachelorThesis/engine/constants"
	"BachelorThesis/engine/objects"
	st "BachelorThesis/engine/singletone"
	"BachelorThesis/engine/vector"
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	hg "github.com/harfang3d/harfang-go"
)

const (
	title        = "Bachelor Thesis Visualization"
	sphereRadius = 1
)

type obj struct {
	transform *hg.Transform
	object    objects.Object
}

func Start(engineSingletone *st.Engine, cancel context.CancelFunc) {
	// scene setup
	win, pipeline, scene, cam := prepareScene()

	// pipeline setup
	res := hg.NewPipelineResources()

	sphereRef, shader := createSphereRefAndRes(res)

	// light setup
	hg.CreateSpotLightWithDiffuseDiffuseIntensitySpecularSpecularIntensityPriorityShadowTypeShadowBias(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(-8.8, 21.7, -8.8), hg.Deg3(60, 45, 0)), 0, hg.Deg(5), hg.Deg(30), hg.ColorGetWhite(), 1, hg.ColorGetWhite(), 1, 0, hg.LSTMap, 0.000005)

	rect := hg.NewIntRectWithSxSyExEy(0, 0, constants.WindowWidth, constants.WindowHeight)

	// initial objects adding to the renderer
	engineSingletone.Mute()

	rendererPool := make(map[string]*obj, 0)
	for _, object := range *engineSingletone.ObjectPool {
		sphereMat := hg.CreateMaterialWithValueName0Value0ValueName1Value1(
			shader,
			"uDiffuseColor",
			hg.NewVec4WithXYZ(float32(rand.Intn(100))/100, float32(rand.Intn(100))/100, float32(rand.Intn(100))/100),
			"uSpecularColor",
			hg.NewVec4WithXYZ(1, 0.8, 0),
		)

		rendererPool[object.GetId()] = &obj{
			transform: newSphere(scene, sphereRef, sphereMat),
			object:    object,
		}
	}

	runtime.GC()
	engineSingletone.Unmute()

	defer func() {
		if cam != nil && cam.IsValid() {
			if cam.HasCamera() {
				cameraComponent := cam.GetCamera()
				if cameraComponent != nil && cameraComponent.IsValid() {
					cam.RemoveCamera()
					scene.DestroyCamera(cameraComponent)
				}
			}
			scene.DestroyNode(cam)
			cam = nil
		}

		if pipeline != nil {
			hg.DestroyForwardPipeline(pipeline)
			pipeline = nil
		}

		if scene != nil {
			scene.Clear()
			scene.GarbageCollect()
			scene = nil
		}

		hg.RenderShutdown()

		if win != nil {
			hg.DestroyWindow(win)
			win = nil
		}

		hg.WindowSystemShutdown()
		hg.InputShutdown()

		runtime.GC()

		log.Print("visualizer shutdown")
	}()

	// main loop
	frame := 0

	for i := 0; i < 1024; i++ {
		newSphereInPool_TEMP(engineSingletone, &rendererPool, scene, sphereRef, shader)
	}
	log.Printf("Objects in pool on start: %d", len(*engineSingletone.ObjectPool))
	log.Println()
	timer := time.Now()

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		frame++

		if time.Since(timer) > time.Second*60 {
			log.Printf("Objects in pool: %d", len(*engineSingletone.ObjectPool))
			log.Printf("FPM: %d", frame)

			if len(*engineSingletone.ObjectPool) >= 32768 {
				log.Printf("Simulation is too long, stopping...")
				cancel()
				return
			}

			addingCount := len(*engineSingletone.ObjectPool)
			for i := 0; i < addingCount; i++ {
				newSphereInPool_TEMP(engineSingletone, &rendererPool, scene, sphereRef, shader)
			}

			log.Println()
			timer = time.Now()
			frame = 0
		}

		for _, object := range rendererPool {
			updateObjectOnRenderer(object)
		}

		engineSingletone.ProcessCollisions()

		dt := hg.TickClock()
		scene.Update(dt)

		viewID := uint16(0)
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontal(&viewID, scene, rect, true, pipeline, res)

		hg.Frame()
		hg.UpdateWindow(win)

		select {
		case <-engineSingletone.Context.Done():
			return
		default:
			continue
		}
	}

	select {
	case <-engineSingletone.Context.Done():
		return
	default:
		cancel()
	}
}

func createSphereRefAndRes(res *hg.PipelineResources) (*hg.ModelRef, *hg.PipelineProgramRef) {
	vtxLayout := hg.VertexLayoutPosFloatNormUInt8()
	sphereMdl := hg.CreateSphereModel(vtxLayout, sphereRadius, 8, 16)

	return res.AddModel("sphere", sphereMdl),
		hg.LoadPipelineProgramRefFromFile("resources_compiled/core/shader/default.hps", res, hg.GetForwardPipelineInfo())
}

func newSphereInPool_TEMP(engineSingletone *st.Engine, rendererPool *map[string]*obj, scene *hg.Scene, sphereRef *hg.ModelRef, shader *hg.PipelineProgramRef) {
	id := fmt.Sprintf("%s_%d", time.Now().Format(time.RFC3339), len(*engineSingletone.ObjectPool))

	posX := float64(rand.Intn(50) - 25)
	posY := float64(rand.Intn(50) - 25)
	posZ := float64(rand.Intn(50) - 25)

	maxInitialSpeed := 0.1                             // Например, максимальная начальная скорость 0.1
	forceX := (rand.Float64()*2 - 1) * maxInitialSpeed // Диапазон [-maxInitialSpeed, maxInitialSpeed)
	forceY := (rand.Float64()*2 - 1) * maxInitialSpeed
	forceZ := (rand.Float64()*2 - 1) * maxInitialSpeed

	sphere := objects.NewSphere(sphereRadius, id)
	sphere.SetPosition(vector.Vector3D{X: posX, Y: posY, Z: posZ})
	sphere.ApplyVelocity(vector.Vector3D{X: forceX, Y: forceY, Z: forceZ})

	engineSingletone.AddObject(&sphere)

	sphereMat := hg.CreateMaterialWithValueName0Value0ValueName1Value1(
		shader,
		"uDiffuseColor",
		hg.NewVec4WithXYZ(float32(rand.Intn(100))/100, float32(rand.Intn(100))/100, float32(rand.Intn(100))/100),
		"uSpecularColor",
		hg.NewVec4WithXYZ(1, 0.8, 0),
	)

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

func prepareScene() (*hg.Window, *hg.ForwardPipeline, *hg.Scene, *hg.Node) {
	hg.InputInit()
	hg.WindowSystemInit()

	//подготовка окна для рисования
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags(title, constants.WindowWidth, constants.WindowHeight, hg.RFMSAA4X|hg.RFVSync)
	pipeline := hg.CreateForwardPipelineWithShadowMapResolution(1024)

	//создания.и конфигурация фонового цвета
	scene := hg.NewScene()
	scene.GetCanvas().SetColor(hg.NewColorWithRGB(0.1, 0.1, 0.1))
	scene.GetEnvironment().SetAmbient(hg.NewColorWithRGB(0.1, 0.1, 0.1))

	//установка камеры
	//cam := hg.CreateCamera(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(155, 50, -60), hg.NewVec3WithXYZ(0.4, -1.2, 0)), 0.01, 10000)
	cam := hg.CreateCamera(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(15, 5, -6), hg.NewVec3WithXYZ(0.4, -1.2, 0)), 0.01, 10000)
	scene.SetCurrentCamera(cam)

	return win, pipeline, scene, cam
}
