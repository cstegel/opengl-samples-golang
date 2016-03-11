#version 410 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 project;

uniform vec3 lightPos;  // only need one light for a basic example

out vec3 Normal;
out vec3 FragPos;
out vec3 LightPos;

void main()
{
    gl_Position = project * view * model * vec4(position, 1.0);

    // we transform positions and vectors to view space before performing lighting
    // calculations in the fragment shader so that we know that the viewer position is (0,0,0)
    // FragPos = vec3(model * vec4(position, 1.0));
    FragPos = vec3(view * model * vec4(position, 1.0));

    // LightPos = vec3(view * vec4(lightPos, 1.0));
    LightPos = vec3(view * vec4(lightPos, 1.0));

    // transform the normals to the view space
    // this is different from just multiplying by the model then view matrix since
    // normals can't translate and are changed by non-uniform scaling
    // instead we take the upper left 3x3 matrix of the transpose of the inverse of each transform
    // that we are transforming across
    // see here for more details: http://www.lighthouse3d.com/tutorials/glsl-tutorial/the-normal-matrix/
    mat3 normMatrix = mat3(transpose(inverse(view))) * mat3(transpose(inverse(model)));
    Normal = normMatrix * normal;
}
