package shader

import (
	"unsafe"
	"log"
	"io/ioutil"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Program struct {
	Program uint32
}

func (program *Program) Close() {
	gl.DeleteProgram(program.Program)
}

func (program *Program) Use() {
	gl.UseProgram(program.Program)
}

func NewProgram(vertShaderFile string, fragShaderFile string) *Program {
	shaderProgram := &Program{}
	vertShader := compileShader(vertShaderFile, gl.VERTEX_SHADER)
	fragShader := compileShader(fragShaderFile, gl.FRAGMENT_SHADER)
	shaderProgram.Program = linkShaders([]uint32{vertShader, fragShader})

	gl.DeleteShader(vertShader)
	gl.DeleteShader(fragShader)
	return shaderProgram
}

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

func checkShaderCompileErrors(shader uint32, shaderFile string) {
	checkGlError(shader, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog,
		"ERROR::SHADER::COMPILE_FAILURE::" + shaderFile)
}

func checkProgramLinkErrors(program uint32) {
	checkGlError(program, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog,
		"ERROR::PROGRAM::LINKING_FAILURE")
}

func compileShader(shaderFile string, shaderType uint32) uint32 {
	shaderSource, err := ioutil.ReadFile(shaderFile)
	if err != nil {
		log.Fatal(err)
	}

	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, gl.Str(shaderSource), nil)
	gl.CompileShader(shader)
	checkShaderCompileErrors(shader, shaderFile)
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
