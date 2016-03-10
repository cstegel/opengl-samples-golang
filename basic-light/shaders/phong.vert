#version 410 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 project;

out vec3 Normal;
out vec3 FragPos;

void main()
{
    gl_Position = project * view * model * vec4(position, 1.0);
    FragPos = vec3(model * vec4(position, 1.0));

    // transform the normals to the model
    // this is different from just multiplying by the model matrix since
    // normals can't translate and are changed by non-uniform scaling
    // instead we take the upper left 3x3 matrix of the transpose of the inverse model matrix
    // see here for more details: http://www.lighthouse3d.com/tutorials/glsl-tutorial/the-normal-matrix/
    mat3 normMatrix = mat3(transpose(inverse(model)));
    Normal = normMatrix * normal;
}
