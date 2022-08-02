#version 330 core

layout (location = 0) in vec3 vert;
layout (location = 1) in vec3 color;

out vec3 aColor;
void main() 
{
    gl_Position = vec4(vert, 1.0);
    aColor = color;
}