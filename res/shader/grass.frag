#version 430
in vec4 fragPosition;
in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 fragNormal;
in float vFade;


out vec4 finalColor;

uniform sampler2D texture0;

void main() {
    // Simple lambert-ish
    vec3 N = normalize(fragNormal);
    vec3 L = normalize(vec3(0.4, 1.0, 0.2));
    float ndl = clamp(dot(N, L)*0.5 + 0.5, 0.0, 1.0);

    vec3 col = vec3(0.18, 0.45, 0.20) * (0.55 + 0.65*ndl) * vFade;

    finalColor = vec4(col, 1.0);
}
