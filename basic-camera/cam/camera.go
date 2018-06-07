package cam

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"

	"github.com/cstegel/opengl-samples-golang/basic-camera/win"
)

type FpsCamera struct {
	// Camera options
	moveSpeed float64
	cursorSensitivity float64

	// Eular Angles
	pitch float64
	yaw float64

	// Camera attributes
	pos mgl32.Vec3
	front mgl32.Vec3
	up mgl32.Vec3
	right mgl32.Vec3
	worldUp mgl32.Vec3

	inputManager *win.InputManager
}

func NewFpsCamera(position, worldUp mgl32.Vec3, yaw, pitch float64, im *win.InputManager) (*FpsCamera) {
	cam := FpsCamera {
		moveSpeed: 5.00,
		cursorSensitivity: 0.05,
		pitch: pitch,
		yaw: yaw,
		pos: position,
		up: mgl32.Vec3{0, 1, 0},
		worldUp: worldUp,
		inputManager: im,
	}

	return &cam
}

func (c *FpsCamera) Update(dTime float64) {
	c.updatePosition(dTime)
	c.updateDirection()
}

// UpdatePosition updates this camera's position by giving directions that
// the camera is to travel in and for how long
func (c *FpsCamera) updatePosition(dTime float64) {
	adjustedSpeed := float32(dTime * c.moveSpeed)

	if c.inputManager.IsActive(win.PLAYER_FORWARD) {
		c.pos = c.pos.Add(c.front.Mul(adjustedSpeed))
	}
	if c.inputManager.IsActive(win.PLAYER_BACKWARD) {
		c.pos = c.pos.Sub(c.front.Mul(adjustedSpeed))
	}
	if c.inputManager.IsActive(win.PLAYER_LEFT) {
		c.pos = c.pos.Sub(c.front.Cross(c.up).Normalize().Mul(adjustedSpeed))
	}
	if c.inputManager.IsActive(win.PLAYER_RIGHT) {
		c.pos = c.pos.Add(c.front.Cross(c.up).Normalize().Mul(adjustedSpeed))
	}
}

// UpdateCursor updates the direction of the camera by giving it delta x/y values
// that came from a cursor input device
func (c *FpsCamera) updateDirection() {
	dCursor := c.inputManager.CursorChange()

	dx := -c.cursorSensitivity * dCursor[0]
	dy := c.cursorSensitivity * dCursor[1]

	c.pitch += dy
	if c.pitch > 89.0 {
		c.pitch = 89.0
	} else if c.pitch < -89.0 {
		c.pitch = -89.0
	}

	c.yaw = math.Mod(c.yaw + dx, 360)
	c.updateVectors()
}

func (c *FpsCamera) updateVectors() {
	// x, y, z
	c.front[0] = float32(math.Cos(mgl64.DegToRad(c.pitch)) * math.Cos(mgl64.DegToRad(c.yaw)))
	c.front[1] = float32(math.Sin(mgl64.DegToRad(c.pitch)))
	c.front[2] = float32(math.Cos(mgl64.DegToRad(c.pitch)) * math.Sin(mgl64.DegToRad(c.yaw)))
	c.front = c.front.Normalize()

	// Gram-Schmidt process to figure out right and up vectors
	c.right = c.worldUp.Cross(c.front).Normalize()
	c.up = c.right.Cross(c.front).Normalize()
}

// GetCameraTransform gets the matrix to transform from world coordinates to
// this camera's coordinates.
func (camera *FpsCamera) GetTransform() mgl32.Mat4 {
	cameraTarget := camera.pos.Add(camera.front)

	return mgl32.LookAt(
		camera.pos.X(), camera.pos.Y(), camera.pos.Z(),
		cameraTarget.X(), cameraTarget.Y(), cameraTarget.Z(),
		camera.up.X(), camera.up.Y(), camera.up.Z(),
	)
}
