package main

/*
Adapted from this tutorial: http://www.learnopengl.com/#!Getting-started/Hello-Triangle
*/

import (
	"log"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const windowWidth = 800
const windowHeight = 600

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
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Hello advanced!", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow (go function bindings)
	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.SetKeyCallback(keyCallback)

	programLoop(window)
}

var vertexShaderSource = `
#version 410 core

layout (location = 0) in vec3 position;

void main()
{
    gl_Position = vec4(position.x, position.y, position.z, 1.0);
}
`

var fragmentShaderBlueSource = `
#version 410 core

out vec4 color;

void main()
{
    color = vec4(0.0f, 0.0f, 1.0f, 1.0f);
}
`

var fragmentShaderRedSource = `
#version 410 core

out vec4 color;

void main()
{
    color = vec4(1.0f, 0.0f, 0.0f, 1.0f);
}
`

type getGlParam func(uint32, uint32, *int32)
type getInfoLog func(uint32, int32, *int32, *uint8)

func checkGlError(glObject uint32, errorParam uint32, getParamFn getGlParam,
	getInfoLogFn getInfoLog, failMsg string) {

	var success int32
	getParamFn(glObject, errorParam, &success)
	if success != 1 {
		var infoLog [512]byte
		getInfoLogFn(glObject, 512, nil, (*uint8)(unsafe.Pointer(&infoLog)))
		log.Fatalln(failMsg, "\n", string(infoLog[:512]))
	}
}

func checkShaderCompileErrors(shader uint32) {
	checkGlError(shader, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog,
		"ERROR::SHADER::COMPILE_FAILURE")
}

func checkProgramLinkErrors(program uint32) {
	checkGlError(program, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog,
		"ERROR::PROGRAM::LINKING_FAILURE")
}

func compileShader(shaderCode string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, (**uint8)(unsafe.Pointer(&shaderCode)), nil)
	gl.CompileShader(shader)
	checkShaderCompileErrors(shader)
	return shader
}

/*
 * Link the provided shaders in the order they were given and return the linked program.
 */
func linkShaders(shaders []uint32) uint32 {
	program := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}
	gl.LinkProgram(program)
	checkProgramLinkErrors(program)

	return program
}

/*
 * Creates the Vertex Array Object for a triangle.
 */
func createTriangleVAO(vertices []float32, indices []uint32) uint32 {
	if len(vertices) > 256 || len(indices) > 256 {
		panic("slices are too big!")
	}
	verticesArray := [256]float32{}
	copy(verticesArray[:], vertices)

	indicesArray := [256]uint32{}
	copy(indicesArray[:], indices)

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
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(verticesArray)), unsafe.Pointer(&verticesArray), gl.STATIC_DRAW)

	// copy indices into element buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(indicesArray)), unsafe.Pointer(&indicesArray), gl.STATIC_DRAW)

	// specify the format of our vertex input
	// (shader) input 0
	// vertex has size 3
	// vertex items are of type FLOAT
	// do not normalize (already done)
	// stride of 3 * sizeof(float) (separation of vertices)
	// offset of where the position data starts (0 for the beginning)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO
}

func programLoop(window *glfw.Window) {

	// the linked shader program determines how the data will be rendered
	vertexShader := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShaderBlue := compileShader(fragmentShaderBlueSource, gl.FRAGMENT_SHADER)
	shaderProgramBlue := linkShaders([]uint32{vertexShader, fragmentShaderBlue})
	fragmentShaderRed := compileShader(fragmentShaderRedSource, gl.FRAGMENT_SHADER)
	shaderProgramRed := linkShaders([]uint32{vertexShader, fragmentShaderRed})

	// shader objects are not needed after they are linked into a program object
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShaderBlue)
	gl.DeleteShader(fragmentShaderRed)

	vertices1 := []float32{
		0.2, 0.2, 0.0,    // top right
		0.2, -0.8, 0.0,   // bottom right
		-0.8, -0.8, 0.0,  // bottom left
		-0.8, 0.2, 0.0,   // top left
	}

	indices1 := []uint32{
		0, 1, 3,  // first triangle
		1, 2, 3,  // second triangle
	}

	VAO1 := createTriangleVAO(vertices1, indices1)

	vertices2 := []float32{
		0.2, 0.6, 0.0,    // top
		0.6, -0.2, 0.0,   // bottom right
		-0.2, -0.2, 0.0,  // bottom left
	}

	indices2 := []uint32{
		0, 1, 2,  // only triangle
	}

	VAO2 := createTriangleVAO(vertices2, indices2)

	for !window.ShouldClose() {
		// poll events and call their registered callbacks
		glfw.PollEvents()

		// perform rendering
		gl.ClearColor(0.2, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// draw loop

		// draw rectangle
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		gl.UseProgram(shaderProgramRed)
		gl.BindVertexArray(VAO1)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, unsafe.Pointer(nil))
		gl.BindVertexArray(0)

		// draw triangle
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		gl.UseProgram(shaderProgramBlue)
		gl.BindVertexArray(VAO2)
		gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, unsafe.Pointer(nil))
		gl.BindVertexArray(0)

		// end of draw loop

		// swap in the rendered buffer
		window.SwapBuffers()
	}
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action,
	mods glfw.ModifierKey) {

	// When a user presses the escape key, we set the WindowShouldClose property to true,
	// which closes the application
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}
