#version 410 core

layout (location = 0) in vec3 position;

out vec2 TexCoord;

uniform mat4 world;
uniform mat4 camera;
uniform mat4 project;

void main()
{
    gl_Position = project * camera * world * vec4(position, 1.0);
}
