#version 410 core

// special fragment shader that is not affected by lighting
// useful for debugging like showing locations of lights

out vec4 color;

void main()
{
	color = vec4(1.0f); // color white
}
