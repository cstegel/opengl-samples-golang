package main

/*
http://www.learnopengl.com/#!Lighting/Materials

Shows basic materials with phong lighting
*/

import (
	"log"
	"runtime"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/opengl-samples-golang/materials/gfx"
	"github.com/opengl-samples-golang/materials/win"
	"github.com/opengl-samples-golang/materials/cam"
)

// vertices to draw 6 faces of a cube
var cubeVertices = []float32{
	// position        // normal vector
	-0.5, -0.5, -0.5,  0.0,  0.0, -1.0,
	 0.5, -0.5, -0.5,  0.0,  0.0, -1.0,
	 0.5,  0.5, -0.5,  0.0,  0.0, -1.0,
	 0.5,  0.5, -0.5,  0.0,  0.0, -1.0,
	-0.5,  0.5, -0.5,  0.0,  0.0, -1.0,
	-0.5, -0.5, -0.5,  0.0,  0.0, -1.0,

	-0.5, -0.5,  0.5,  0.0,  0.0, 1.0,
	 0.5, -0.5,  0.5,  0.0,  0.0, 1.0,
	 0.5,  0.5,  0.5,  0.0,  0.0, 1.0,
	 0.5,  0.5,  0.5,  0.0,  0.0, 1.0,
	-0.5,  0.5,  0.5,  0.0,  0.0, 1.0,
	-0.5, -0.5,  0.5,  0.0,  0.0, 1.0,

	-0.5,  0.5,  0.5, -1.0,  0.0,  0.0,
	-0.5,  0.5, -0.5, -1.0,  0.0,  0.0,
	-0.5, -0.5, -0.5, -1.0,  0.0,  0.0,
	-0.5, -0.5, -0.5, -1.0,  0.0,  0.0,
	-0.5, -0.5,  0.5, -1.0,  0.0,  0.0,
	-0.5,  0.5,  0.5, -1.0,  0.0,  0.0,

	 0.5,  0.5,  0.5,  1.0,  0.0,  0.0,
	 0.5,  0.5, -0.5,  1.0,  0.0,  0.0,
	 0.5, -0.5, -0.5,  1.0,  0.0,  0.0,
	 0.5, -0.5, -0.5,  1.0,  0.0,  0.0,
	 0.5, -0.5,  0.5,  1.0,  0.0,  0.0,
	 0.5,  0.5,  0.5,  1.0,  0.0,  0.0,

	-0.5, -0.5, -0.5,  0.0, -1.0,  0.0,
	 0.5, -0.5, -0.5,  0.0, -1.0,  0.0,
	 0.5, -0.5,  0.5,  0.0, -1.0,  0.0,
	 0.5, -0.5,  0.5,  0.0, -1.0,  0.0,
	-0.5, -0.5,  0.5,  0.0, -1.0,  0.0,
	-0.5, -0.5, -0.5,  0.0, -1.0,  0.0,

	-0.5,  0.5, -0.5,  0.0,  1.0,  0.0,
	 0.5,  0.5, -0.5,  0.0,  1.0,  0.0,
	 0.5,  0.5,  0.5,  0.0,  1.0,  0.0,
	 0.5,  0.5,  0.5,  0.0,  1.0,  0.0,
	-0.5,  0.5,  0.5,  0.0,  1.0,  0.0,
	-0.5,  0.5, -0.5,  0.0,  1.0,  0.0,
}

var cubePositions = [][]float32 {
	{ 0.0,  0.0,  -3.0},
	{ 2.0,  5.0, -15.0},
	{-1.5, -2.2, -2.5 },
	{-3.8, -2.0, -12.3},
	{ 2.4, -0.4, -3.5 },
	{-1.7,  3.0, -7.5 },
	{ 1.3, -2.0, -2.5 },
	{ 1.5,  2.0, -2.5 },
	{ 1.5,  0.2, -1.5 },
	{-1.3,  1.0, -1.5 },
}

func init() {
	// GLFW event handling must be run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to inifitialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Initialize Glow (go function bindings)
	if err := gl.Init(); err != nil {
		panic(err)
	}

	window := win.NewWindow(1280, 720, "Materials")

	err := programLoop(window)
	if err != nil {
		log.Fatalln(err)
	}
}

/*
 * Creates the Vertex Array Object for a triangle.
 * indices is leftover from earlier samples and not used here.
 */
func createVAO(vertices []float32, indices []uint32) uint32 {

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32;
	gl.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// size of one whole vertex (sum of attrib sizes)
	var stride int32 = 3*4 + 3*4
	var offset int = 0

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3*4

	// normal
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)
	offset += 3*4

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO
}

func programLoop(window *win.Window) error {

	// the linked shader program determines how the data will be rendered
	vertShader, err := gfx.NewShaderFromFile("shaders/phong.vert", gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	fragShader, err := gfx.NewShaderFromFile("shaders/phong.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	program, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		return err
	}
	defer program.Delete()

	lightFragShader, err := gfx.NewShaderFromFile("shaders/light.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	// special shader program so that lights themselves are not affected by lighting
	lightProgram, err := gfx.NewProgram(vertShader, lightFragShader)
	if err != nil {
		return err
	}

	VAO := createVAO(cubeVertices, nil)
	lightVAO := createVAO(cubeVertices, nil)

	// ensure that triangles that are "behind" others do not draw over top of them
	gl.Enable(gl.DEPTH_TEST)

	camera := cam.NewFpsCamera(mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 1, 0}, -90, 0, window.InputManager())

	for !window.ShouldClose() {

		// swaps in last buffer, polls for window events, and generally sets up for a new render frame
		window.StartFrame()

		// update camera position and direction from input evevnts
		camera.Update(window.SinceLastFrame())

		// background color
		gl.ClearColor(0, 0, 0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)  // depth buffer needed for DEPTH_TEST


		// cube rotation matrices
		rotateX   := (mgl32.Rotate3DX(mgl32.DegToRad(-45 * float32(glfw.GetTime()))))
		rotateY   := (mgl32.Rotate3DY(mgl32.DegToRad(-45 * float32(glfw.GetTime()))))
		rotateZ   := (mgl32.Rotate3DZ(mgl32.DegToRad(-45 * float32(glfw.GetTime()))))

		// creates perspective
		fov := float32(60.0)
		projectTransform := mgl32.Perspective(mgl32.DegToRad(fov),
		                                      float32(window.Width())/float32(window.Height()),
		                                      0.1,
		                                      100.0)

		camTransform := camera.GetTransform()
		lightPos := mgl32.Vec3{0.6, 1, 0.1}
		lightTransform := mgl32.Translate3D(lightPos.X(), lightPos.Y(), lightPos.Z()).Mul4(
		                                    mgl32.Scale3D(0.2, 0.2, 0.2))

		program.Use()
		gl.UniformMatrix4fv(program.GetUniformLocation("view"), 1, false, &camTransform[0])
		gl.UniformMatrix4fv(program.GetUniformLocation("project"), 1, false,
		                    &projectTransform[0])

		gl.BindVertexArray(VAO)

		// draw each cube after all coordinate system transforms are bound

		// obj is colored, light is white
		gl.Uniform3f(program.GetUniformLocation("material.ambient"), 1.0, 0.5, 0.31)
		gl.Uniform3f(program.GetUniformLocation("material.diffuse"), 1.0, 0.5, 0.31)
		gl.Uniform3f(program.GetUniformLocation("material.specular"), 0.5, 0.5, 0.5)
		gl.Uniform1f(program.GetUniformLocation("material.shininess"), 32.0)

		lightColor := mgl32.Vec3{
			float32(math.Sin(glfw.GetTime() * 1)),
			float32(math.Sin(glfw.GetTime() * 0.35)),
			float32(math.Sin(glfw.GetTime() * 0.65)),
		}

		diffuseColor := mgl32.Vec3{
			0.5 * lightColor[0],
			0.5 * lightColor[1],
			0.5 * lightColor[2],
		}
		ambientColor := mgl32.Vec3{
			0.2 * lightColor[0],
			0.2 * lightColor[1],
			0.2 * lightColor[2],
		}

		gl.Uniform3f(program.GetUniformLocation("light.ambient"),
		             ambientColor[0], ambientColor[1], ambientColor[2])
		gl.Uniform3f(program.GetUniformLocation("light.diffuse"),
		             diffuseColor[0], diffuseColor[1], diffuseColor[2])
		gl.Uniform3f(program.GetUniformLocation("light.specular"), 1.0, 1.0, 1.0)
		gl.Uniform3f(program.GetUniformLocation("light.position"), lightPos.X(), lightPos.Y(), lightPos.Z())

		for _, pos := range cubePositions {

			// turn the cubes into rectangular prisms for more fun
			worldTranslate := mgl32.Translate3D(pos[0], pos[1], pos[2])
			worldTransform := worldTranslate.Mul4(
				rotateX.Mul3(rotateY).Mul3(rotateZ).Mat4(),
			)

			gl.UniformMatrix4fv(program.GetUniformLocation("model"), 1, false,
			                    &worldTransform[0])

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}
		gl.BindVertexArray(0)

		// Draw the light obj after the other boxes using its separate shader program
		// this means that we must re-bind any uniforms
		lightProgram.Use()
		gl.BindVertexArray(lightVAO)
		gl.UniformMatrix4fv(lightProgram.GetUniformLocation("model"), 1, false, &lightTransform[0])
		gl.UniformMatrix4fv(lightProgram.GetUniformLocation("view"), 1, false, &camTransform[0])
		gl.UniformMatrix4fv(lightProgram.GetUniformLocation("project"), 1, false, &projectTransform[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		gl.BindVertexArray(0)

		// end of draw loop
	}

	return nil
}
