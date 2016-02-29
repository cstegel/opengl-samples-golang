package cam

import (
	"log"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type CursorListener interface {
	// dx, dy
	UpdateCursor(float64, float64)
}

type CursorManager interface {
	RegisterCursorListener(*CursorListener)
}

type MovementListener interface {
	UpdateMovement(float64, ...Direction)
}

type MovementManager interface {
	RegisterMovementListener(*MovementListener)
}

type Direction int
const (
	FORWARD Direction = iota
	BACKWARD Direction = iota
	LEFT Direction = iota
	RIGHT Direction = iota
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
}

func NewFpsCamera(position, worldUp mgl32.Vec3, yaw, pitch float64) (*FpsCamera) {
	cam := FpsCamera {
		moveSpeed: 5.00,
		cursorSensitivity: 0.05,
		pitch: pitch,
		yaw: yaw,
		pos: position,
		up: mgl32.Vec3{0, 1, 0},
		worldUp: worldUp,
	}

	return &cam
}

func (camera *FpsCamera) updateVectors() {
	// x, y, z
	camera.front[0] = float32(math.Cos(mgl64.DegToRad(camera.pitch)) * math.Cos(mgl64.DegToRad(camera.yaw)))
	camera.front[1] = float32(math.Sin(mgl64.DegToRad(camera.pitch)))
	camera.front[2] = float32(math.Cos(mgl64.DegToRad(camera.pitch)) * math.Sin(mgl64.DegToRad(camera.yaw)))
	camera.front = camera.front.Normalize()

	// Gram-Schmidt process to figure out right and up vectors
	camera.right = camera.worldUp.Cross(camera.front).Normalize()
	camera.up = camera.right.Cross(camera.front).Normalize()
}

// UpdatePosition updates this camera's position by giving directions that
// the camera is to travel in and for how long
func (camera *FpsCamera) UpdatePosition(dTime float64, directions ...Direction) {
	adjustedSpeed := float32(dTime * camera.moveSpeed)

	for _, dir := range directions {
		switch dir {
		case FORWARD:
			camera.pos = camera.pos.Add(camera.front.Mul(adjustedSpeed))
		case BACKWARD:
			camera.pos = camera.pos.Sub(camera.front.Mul(adjustedSpeed))
		case LEFT:
			camera.pos = camera.pos.Sub(camera.front.Cross(camera.up).Normalize().Mul(adjustedSpeed))
		case RIGHT:
			camera.pos = camera.pos.Add(camera.front.Cross(camera.up).Normalize().Mul(adjustedSpeed))
		default:
			log.Output(1, string(dir) + " is an invalid 'Direction'")
		}
	}
}

// UpdateCursor updates the direction of the camera by giving it delta x/y values
// that came from a cursor input device
func (camera *FpsCamera) UpdateCursor(dx, dy float64) {
	dx *= camera.cursorSensitivity
	dy *= -camera.cursorSensitivity // reversed since y goes from bottom to top

	camera.pitch += dy
	if camera.pitch > 89.0 {
		camera.pitch = 89.0
	} else if camera.pitch < -89.0 {
		camera.pitch = -89.0
	}

	camera.yaw = math.Mod(camera.yaw + dx, 360)
	camera.updateVectors()
}

// GetCameraTransform gets the matrix to transform from world coordinates to
// this camera's coordinates
func (camera *FpsCamera) GetCameraTransform() mgl32.Mat4 {
	cameraTarget := camera.pos.Add(camera.front)

	return mgl32.LookAt(
		camera.pos.X(), camera.pos.Y(), camera.pos.Z(),
		cameraTarget.X(), cameraTarget.Y(), cameraTarget.Z(),
		camera.up.X(), camera.up.Y(), camera.up.Z(),
	)
}
