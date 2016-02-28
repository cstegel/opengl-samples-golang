package main

/*
Adapted from this tutorial: http://www.learnopengl.com/#!Getting-started/Camera

Shows how to create a basic controllable FPS camera
*/

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/opengl-samples-golang/basic-camera/gfx"
)

const windowWidth = 1280
const windowHeight = 720

// only using global variables because this is meant as a simple example
var cameraPos   = mgl32.Vec3{0.0, 0.0, 3.0}
var cameraFront = mgl32.Vec3{0.0, 0.0, -1.0}
var cameraUp    = mgl32.Vec3{0.0, 1.0, 0.0}

// vertices to draw 6 faces of a cube
var cubeVertices = []float32{
	// position        // texture position
	-0.5, -0.5, -0.5,  0.0, 0.0,
	 0.5, -0.5, -0.5,  1.0, 0.0,
	 0.5,  0.5, -0.5,  1.0, 1.0,
	 0.5,  0.5, -0.5,  1.0, 1.0,
	-0.5,  0.5, -0.5,  0.0, 1.0,
	-0.5, -0.5, -0.5,  0.0, 0.0,

	-0.5, -0.5,  0.5,  0.0, 0.0,
	 0.5, -0.5,  0.5,  1.0, 0.0,
	 0.5,  0.5,  0.5,  1.0, 1.0,
	 0.5,  0.5,  0.5,  1.0, 1.0,
	-0.5,  0.5,  0.5,  0.0, 1.0,
	-0.5, -0.5,  0.5,  0.0, 0.0,

	-0.5,  0.5,  0.5,  1.0, 0.0,
	-0.5,  0.5, -0.5,  1.0, 1.0,
	-0.5, -0.5, -0.5,  0.0, 1.0,
	-0.5, -0.5, -0.5,  0.0, 1.0,
	-0.5, -0.5,  0.5,  0.0, 0.0,
	-0.5,  0.5,  0.5,  1.0, 0.0,

	 0.5,  0.5,  0.5,  1.0, 0.0,
	 0.5,  0.5, -0.5,  1.0, 1.0,
	 0.5, -0.5, -0.5,  0.0, 1.0,
	 0.5, -0.5, -0.5,  0.0, 1.0,
	 0.5, -0.5,  0.5,  0.0, 0.0,
	 0.5,  0.5,  0.5,  1.0, 0.0,

	-0.5, -0.5, -0.5,  0.0, 1.0,
	 0.5, -0.5, -0.5,  1.0, 1.0,
	 0.5, -0.5,  0.5,  1.0, 0.0,
	 0.5, -0.5,  0.5,  1.0, 0.0,
	-0.5, -0.5,  0.5,  0.0, 0.0,
	-0.5, -0.5, -0.5,  0.0, 1.0,

	-0.5,  0.5, -0.5,  0.0, 1.0,
	 0.5,  0.5, -0.5,  1.0, 1.0,
	 0.5,  0.5,  0.5,  1.0, 0.0,
	 0.5,  0.5,  0.5,  1.0, 0.0,
	-0.5,  0.5,  0.5,  0.0, 0.0,
	-0.5,  0.5, -0.5,  0.0, 1.0,
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
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "basic camera", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow (go function bindings)
	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.SetKeyCallback(keyCallback)

	err = programLoop(window)
	if err != nil {
		log.Fatal(err)
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
	var stride int32 = 3*4 + 2*4
	var offset int = 0

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3*4

	// texture position
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)
	offset += 2*4

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO
}

func programLoop(window *glfw.Window) error {

	// the linked shader program determines how the data will be rendered
	vertShader, err := gfx.NewShaderFromFile("shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	fragShader, err := gfx.NewShaderFromFile("shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	program, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		return err
	}
	defer program.Delete()

	VAO := createVAO(cubeVertices, nil)
	texture0, err := gfx.NewTextureFromFile("../images/RTS_Crate.png",
	                                        gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
	if err != nil {
		panic(err.Error())
	}

	texture1, err := gfx.NewTextureFromFile("../images/trollface-transparent.png",
	                                        gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
	if err != nil {
		panic(err.Error())
	}

	// ensure that triangles that are "behind" others do not draw over top of them
	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		// poll events and call their registered callbacks
		glfw.PollEvents()

		// background color
		gl.ClearColor(0.2, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)  // depth buffer needed for DEPTH_TEST

		program.Use()

		// bind textures
		texture0.Bind(gl.TEXTURE0)
		texture0.SetUniform(program.GetUniformLocation("ourTexture0"))

		texture1.Bind(gl.TEXTURE1)
		texture1.SetUniform(program.GetUniformLocation("ourTexture1"))

		// cube rotation matrices
		rotateX   := (mgl32.Rotate3DX(mgl32.DegToRad(-60 * float32(glfw.GetTime()))))
		rotateY   := (mgl32.Rotate3DY(mgl32.DegToRad(-60 * float32(glfw.GetTime()))))
		rotateZ   := (mgl32.Rotate3DZ(mgl32.DegToRad(-60 * float32(glfw.GetTime()))))

		// creates perspective
		fov := float32(60.0)
		projectTransform := mgl32.Perspective(mgl32.DegToRad(fov), windowWidth/windowHeight, 0.1, 100.0)

		// Calculate camera transform
		// x/z are horizontal, y is vertical
		cameraTarget := cameraPos.Add(cameraFront)  // TODO-cs: why?

		cameraTransform := mgl32.LookAt(
			cameraPos.X(), cameraPos.Y(), cameraPos.Z(),
			cameraTarget.X(), cameraTarget.Y(), cameraTarget.Z(),
			cameraUp.X(), cameraUp.Y(), cameraUp.Z(),
		)

		gl.UniformMatrix4fv(program.GetUniformLocation("camera"), 1, false, &cameraTransform[0])
		gl.UniformMatrix4fv(program.GetUniformLocation("project"), 1, false,
		&projectTransform[0])

		gl.BindVertexArray(VAO)

		// draw each cube after all coordinate system transforms are bound
		for _, pos := range cubePositions {
			worldTranslate := mgl32.Translate3D(pos[0], pos[1], pos[2])
			worldTransform := (worldTranslate.Mul4(rotateX.Mul3(rotateY).Mul3(rotateZ).Mat4()))

			gl.UniformMatrix4fv(program.GetUniformLocation("world"), 1, false,
			                    &worldTransform[0])

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		gl.BindVertexArray(0)

		texture0.UnBind()
		texture1.UnBind()

		// end of draw loop

		// swap in the rendered buffer
		window.SwapBuffers()
	}

	return nil
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action,
	mods glfw.ModifierKey) {

	// When a user presses the escape key, we set the WindowShouldClose property to true,
	// which closes the application
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}

	var cameraSpeed float32 = 0.05
	if key == glfw.KeyW {
		cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
	}
	if key == glfw.KeyS {
		cameraPos = cameraPos.Sub(cameraFront.Mul(cameraSpeed))
	}
	if key == glfw.KeyA {
		cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
	}
	if key == glfw.KeyD {
		cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
	}
}
