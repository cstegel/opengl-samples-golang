#version 410 core

out vec4 color;

uniform vec3 objectColor;
uniform vec3 lightColor;

void main()
{
    // the color of the light "reflects" off the object
    color = vec4(objectColor * lightColor, 1.0f);
}
