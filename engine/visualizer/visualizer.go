package visualizer

import (
	"BachelorThesis/engine/objects"
	"log"
	"math"
	"runtime"

	hg "github.com/harfang3d/harfang-go"
)

const (
	title = "Bachelor Thesis Visualization"
)

func Start(algorithm string, endChan chan bool, pool []objects.Object) {
	hg.InputInit()
	hg.WindowSystemInit()

	//подготовка окна для рисования
	var resX int32 = 1280
	var resY int32 = 720
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang Sample", resX, resY, hg.RFMSAA4X)
	pipeline := hg.CreateForwardPipelineWithShadowMapResolution(4096)

	//создание конвейера отрисовки
	res := hg.NewPipelineResources()
	vtxLayout := hg.VertexLayoutPosFloatNormUInt8()
	sphereMdl := hg.CreateSphereModel(vtxLayout, 0.1, 8, 16)
	sphereRef := res.AddModel("sphere", sphereMdl)

	//будем использовать предкомпилированные шейдеры
	shader := hg.LoadPipelineProgramRefFromFile("resources_compiled/core/shader/default.hps", res, hg.GetForwardPipelineInfo())
	//определение материала (для цвета задается RGB от 0 до 1)
	sphereMat := hg.CreateMaterialWithValueName0Value0ValueName1Value1(shader, "uDiffuseColor", hg.NewVec4WithXYZ(1, 0, 0), "uSpecularColor", hg.NewVec4WithXYZ(1, 0.8, 0))

	//создания.и конфигурация фонового цвета
	scene := hg.NewScene()
	scene.GetCanvas().SetColor(hg.NewColorWithRGB(0.1, 0.1, 0.1))
	scene.GetEnvironment().SetAmbient(hg.NewColorWithRGB(0.1, 0.1, 0.1))

	//установка камеры
	cam := hg.CreateCamera(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(15.5, 5, -6), hg.NewVec3WithXYZ(0.4, -1.2, 0)), 0.01, 100)
	scene.SetCurrentCamera(cam)

	//добавление света
	hg.CreateSpotLightWithDiffuseDiffuseIntensitySpecularSpecularIntensityPriorityShadowTypeShadowBias(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(-8.8, 21.7, -8.8), hg.Deg3(60, 45, 0)), 0, hg.Deg(5), hg.Deg(30), hg.ColorGetWhite(), 1, hg.ColorGetWhite(), 1, 0, hg.LSTMap, 0.000005)

	//создаем объекты сцены (сферы)
	rows := [][]*hg.Transform{}
	for z := float32(-100.0); z < 100.0; z += 2.0 {
		row := []*hg.Transform{}
		for x := float32(-100.0); x < 100.0; x += 2.0 {
			node := hg.CreateObjectWithSliceOfMaterials(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(x*0.1, 0.1, z*0.1)), sphereRef, hg.GoSliceOfMaterial{sphereMat})
			row = append(row, node.GetTransform())
		}
		rows = append(rows, row)
	}
	angle := 0.0
	rect := hg.NewIntRectWithSxSyExEy(0, 0, resX, resY)

	log.Printf("Visualization objects created")

	//основной цикл
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		dt := hg.TickClock()
		angle += float64(hg.TimeToSecF(dt))

		//поворот объектов сцены (узлов) на угол
		for j, row := range rows {
			rowY := math.Cos(angle + float64(j)*0.1)
			for i, trs := range row {
				/*pos := trs.GetPos()
				pos.SetY(float32(0.1 * (rowY*math.Sin(angle+float64(i)*0.1)*6 + 6.5)))
				trs.SetPos(pos)
				*/
				p := hg.NewVec3()
				p, _ = trs.GetPosRot()
				p.SetY(float32(0.1 * (rowY*math.Sin(angle+float64(i)*0.1)*6 + 6.5)))
				trs.SetPos(p)

			}
		}

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

	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
