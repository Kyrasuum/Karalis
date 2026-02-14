#version 430

// Input vertex attributes
in vec3 vertexPosition;
in vec4 vertexColor;

// Input uniform values
uniform mat4 mvp;
uniform mat4 matView;
uniform mat4 matProjection;

// Output vertex attributes (to fragment shader)
out vec3 fragPosition;
out vec4 fragColor;

// Custom Uniforms

void main() {
    mat4 view = mat4(mat3(matView));
    // Calculate final vertex position
    gl_Position = matProjection * view * vec4(vertexPosition, 1.0);

    // Send vertex attributes to fragment shader
    fragPosition = vertexPosition;
    fragColor = vertexColor;
}
