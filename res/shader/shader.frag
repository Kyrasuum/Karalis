#version 330

precision highp float;
precision highp sampler2D;

// Input vertex attributes (from vertex shader)
in vec4 fragPosition;
in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 fragNormal;
in vec4 fragScrPos;
in float fragClip;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;

// Output fragment color
out vec4 finalColor;

// Custom uniforms
uniform float portalObj = 0.0;
uniform float portalScn = 0.0;

vec2 CorrectUV(vec4 pos) {
    vec2 ndx = pos.xy / pos.w;
    vec2 uv = 0.5 * ndx + 0.5;

    return uv;
}

void main() {
    if (portalScn < 1 && fragClip < 0.0) {
        discard;
    }

    // Texel color fetching from texture sampler
    vec4 texelColor;
    if (portalObj < 1) {
        texelColor = texture(texture0, fragTexCoord)*colDiffuse*fragColor;
    } else {
        vec2 uv = CorrectUV(fragScrPos);
        texelColor = texture(texture0, uv)*colDiffuse*fragColor;
    }

    // Calculate final fragment color
    finalColor = texelColor;
}
