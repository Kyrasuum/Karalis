#version 430
in vec3 vN;
in float vFade;

out vec4 fragColor;

void main() {
    // Simple lambert-ish
    vec3 N = normalize(vN);
    vec3 L = normalize(vec3(0.4, 1.0, 0.2));
    float ndl = clamp(dot(N, L)*0.5 + 0.5, 0.0, 1.0);

    vec3 grass = vec3(0.18, 0.45, 0.20);
    vec3 col = grass * (0.55 + 0.65*ndl) * vFade;

    fragColor = vec4(col, 1.0);
}
