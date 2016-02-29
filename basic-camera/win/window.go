package win

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

type Window struct {
	width int
	height int
	glfw *glfw.Window

	keysPressed [glfw.KeyLast]bool
	lastMouse mgl64.Vec2

}

// NewWindow assumes that glfw has already been initialized
func NewWindow(width, height int, title string) (*Window, error) {
	window := Window {
		width:width,
		height:height,
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfwWindow, err := glfw.CreateWindow(window.width, window.height, title, nil, nil)
	if err != nil {
		return nil, err
	}

	window.glfw = glfwWindow
	window.glfw.MakeContextCurrent()

	return &window, nil
}

func (window *Window) Width() int {
	return window.width
}

func (window *Window) Height() int {
	return window.height
}
