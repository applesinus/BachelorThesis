package visualizer

import (
	"BachelorThesis/engine/objects"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

const (
	title = "Bachelor Thesis Visualization"
)

func Start(algorithm string, endChan chan bool, pool []objects.Object) {
	a := app.App()
	scene := core.NewNode()

	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	light := light.NewDirectional(&math32.Color{1, 1, 1}, 1)
	scene.Add(light)

	geom := geometry.NewBox(1, 1, 1)
	mat := material.NewStandard(&math32.Color{1, 0, 0})
	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)

	a.Run(scene)

	/*for {
		select {
		case <-endChan:
			return

		default:
			offset++
			rl.BeginDrawing()
			rl.ClearBackground(rl.RayWhite)

			header = fmt.Sprintf("FPS: %v\ni=%v", rl.GetFPS(), offset)
			rl.DrawText(header, int32(offset), 0, 20, rl.DarkGray)

			rl.EndDrawing()
		}
	}*/
}
