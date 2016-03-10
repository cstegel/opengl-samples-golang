#version 410 core

in vec3 Normal;
in vec3 FragPos;
in vec3 LightPos;
out vec4 color;

uniform vec3 objectColor;
uniform vec3 lightColor;

void main()
{
	// affects diffuse and specular lighting
	float lightPower = 2.0f;

	// diffuse and specular intensity are affected by the amount of light they get based on how
	// far they are from a light source (inverse square of distance)
	float distToLight = length(LightPos - FragPos);
	float distIntensityDecay = 1.0f / pow(distToLight, 2);

	float ambientStrength = 0.05f;
	vec3 ambientLight = ambientStrength * lightColor;

	vec3 norm = normalize(Normal);
	vec3 dirToLight = normalize(LightPos - FragPos);
	float lightNormalDiff = max(dot(norm, dirToLight), 0.0);

	// diffuse light is greatest when surface is perpendicular to light (dot product)
	vec3 diffuse = lightNormalDiff * lightColor;
	vec3 diffuseLight = lightPower * diffuse * distIntensityDecay * lightColor;

	float specularStrength = 1.0f;
	int shininess = 64;
	vec3 viewPos = vec3(0.0f, 0.0f, 0.0f);
	vec3 dirToView = normalize(viewPos - FragPos);
	vec3 reflectDir = reflect(-dirToLight, norm);
	float spec = pow(max(dot(dirToView, reflectDir), 0.0), shininess);
	vec3 specularLight = lightPower * specularStrength * spec * distIntensityDecay * lightColor;

	vec3 result = (diffuseLight + specularLight + ambientLight) * objectColor;
	color = vec4(result, 1.0f);
}
