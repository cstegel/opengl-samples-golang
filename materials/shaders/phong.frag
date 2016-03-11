#version 410 core

struct Material {
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
	float shininess;
};

struct Light {
	vec3 position;

	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
};

in vec3 Normal;
in vec3 FragPos;
in vec3 LightPos;
out vec4 color;

uniform Material material;
uniform Light light;

void main()
{
	// ambient
	vec3 ambient = light.ambient * material.ambient;

	// diffuse
	vec3 norm = normalize(Normal);
	vec3 dirToLight = normalize(LightPos - FragPos);
	float lightNormalDiff = max(dot(norm, dirToLight), 0.0);
	vec3 diffuse = light.diffuse * (material.diffuse * lightNormalDiff);

	// specular
	vec3 viewPos = vec3(0.0f, 0.0f, 0.0f);
	vec3 dirToView = normalize(viewPos - FragPos);
	vec3 reflectDir = reflect(-dirToLight, norm);
	float spec = pow(max(dot(dirToView, reflectDir), 0.0), material.shininess);
	vec3 specular = light.specular * (spec * material.specular);

	vec3 result = diffuse + specular + ambient;
	color = vec4(result, 1.0f);
}
