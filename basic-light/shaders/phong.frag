#version 410 core

in vec3 Normal;
in vec3 FragPos;
out vec4 color;

uniform vec3 objectColor;
uniform vec3 lightColor;
uniform vec3 lightPos;  // only need one light for a basic example

void main()
{
	float ambientStrength = 0.1f;
	vec3 ambientLight = ambientStrength * lightColor;

	float lightPower = 1.0f;
	float distToLight = length(lightPos - FragPos);

	vec3 norm = normalize(Normal);
	vec3 dirToLight = normalize(lightPos - FragPos);
	float lightNormalDiff = max(dot(norm, dirToLight), 0.0);

	// amount of light is greatest when surface is perpendicular to to light (dot product)
	// intensity decays with inverse square of distance
	float distIntensityDecay = 1.0f / pow(distToLight, 2);
	float lightIntensity = lightPower * distIntensityDecay;
	vec3 diffuseLight = (lightNormalDiff * lightColor) * lightIntensity;

	vec3 result = (ambientLight + diffuseLight) * objectColor;
	color = vec4(result, 1.0f);
}
