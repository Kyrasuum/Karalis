#version 330

in vec3 vWorldPos;
in vec2 vUv;
in vec3 vWaveNormal;
in float vWaveHeight;

out vec4 finalColor;

uniform vec3 uCameraPos;
uniform float uTime;
uniform sampler2D texture0;

// Appearance controls
uniform vec4 uWaterColor;      // base tint
uniform float uFresnelPower;   // 2..6 typical
uniform float uSpecPower;      // 32..256 typical
uniform float uDetailStrength; // 0..1
// base look


float clamp01(float x){ return clamp(x, 0.0, 1.0); }

void main() {
    if (texture(texture0, vUv).r > 0.33+vWaveHeight) {
        discard;
    }

    vec3 N = normalize(vWaveNormal);
    vec3 V = normalize(uCameraPos - vWorldPos);

    // Fresnel
    float ndv = clamp01(dot(N, V));
    float fres = pow(1.0 - ndv, uFresnelPower);

    // Sun spec
    vec3 L = normalize(vec3(0.35, 1.0, 0.15));
    vec3 H = normalize(L + V);
    float spec = pow(max(dot(N, H), 0.0), uSpecPower);

    vec3 base = uWaterColor.rgb;

    // Make waves read: modulate color by facing + a little animated detail
    float face = 0.25 + 0.75 * ndv; // brighter when looking down a bit
    float detail = 0.0;
    if (uDetailStrength > 0.001) {
        detail = sin((vWorldPos.x + vWorldPos.z) * 0.35 + uTime * 1.2) * 0.5 + 0.5; // 0..1
        detail = (detail - 0.5) * 0.20 * uDetailStrength; // small +/- shift
    }

    vec3 reflectTint = vec3(0.55, 0.70, 0.90);
    vec3 col = mix(base * face, reflectTint, fres);

    col += spec * 0.9;
    col += detail;

    finalColor = vec4(col, uWaterColor.a);
}
