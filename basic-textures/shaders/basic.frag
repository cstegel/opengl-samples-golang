#version 410 core

in vec3 ourColor;
in vec2 TexCoord;

out vec4 color;

uniform sampler2D ourTexture0;
uniform sampler2D ourTexture1;

void main()
{
    // mix the two textures together (texture1 is colored with "ourColor")
    color = mix(texture(ourTexture0, TexCoord), texture(ourTexture1, TexCoord) * vec4(ourColor, 1.0f), 0.5);
}
