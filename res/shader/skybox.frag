#version 330

precision highp float;
precision highp sampler2D;

// Input vertex attributes (from vertex shader)
in vec3 fragPosition;
in vec4 fragColor;

// Input uniform values
uniform samplerCube texture0;
uniform vec4 colDiffuse;

// Output fragment color
out vec4 finalColor;

// Custom uniforms

void main() {
    // Texel color fetching from texture sampler
    vec4 texelColor = texture(texture0, fragPosition)*colDiffuse*fragColor;

    // Calculate final fragment color
    finalColor = texelColor;
}
