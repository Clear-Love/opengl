#version 410
layout (location = 0) in vec3 vert;
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
out vec4 vertColor;
void main() {
    vertColor = vec4(0.46, 0.51, 0.64, 1.0);
    gl_Position = projection * camera * model * vec4(vert, 10.0);
}