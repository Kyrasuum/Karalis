#version 330

in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec3 vertexNormal;

out vec3 vWorldPos;
out vec2 vUv;
out vec3 vWaveNormal;
out float vWaveHeight;

uniform mat4 mvp;
uniform mat4 matModel;
uniform sampler2D texture0;
uniform float uTime;

// Wave params (base)
uniform float uWaveAmp;      // deep-water amplitude (e.g. 0.10)
uniform float uWaveFreq;     // base frequency (e.g. 0.18)
uniform float uWaveSpeed;    // (e.g. 0.6)


// Small helper: compute wave height + deriv given an amplitude scale
vec3 waveHeightAndDeriv(vec2 p, float ampScale) {
    float t = uTime * uWaveSpeed;

    vec2 d1 = normalize(vec2(1.0, 0.2));
    vec2 d2 = normalize(vec2(-0.4, 1.0));
    vec2 d3 = normalize(vec2(0.7, -0.8));

    float k1 = 6.2831853 * uWaveFreq;
    float k2 = 6.2831853 * (uWaveFreq * 1.37);
    float k3 = 6.2831853 * (uWaveFreq * 0.73);

    float a1 = uWaveAmp * 0.55 * ampScale;
    float a2 = uWaveAmp * 0.30 * ampScale;
    float a3 = uWaveAmp * 0.20 * ampScale;

    float ph1 = k1 * dot(d1, p) + t * 1.00;
    float ph2 = k2 * dot(d2, p) + t * 1.31;
    float ph3 = k3 * dot(d3, p) + t * 0.84;

    float h = a1*sin(ph1) + a2*sin(ph2) + a3*sin(ph3);

    vec2 dh =
        a1*cos(ph1) * k1 * d1 +
        a2*cos(ph2) * k2 * d2 +
        a3*cos(ph3) * k3 * d3;

    return vec3(h, dh.x, dh.y);
}

float clamp01(float x){ return clamp(x, 0.0, 1.0); }

void main() {
    vec4 world = matModel * vec4(vertexPosition, 1.0);
    vUv = vertexTexCoord;

    // Sample normalized depth in [0..1]

    // Turn depth into an amplitude multiplier:
    float depth = texture(texture0, vUv).r;
    float d = max(0.33 - depth, 0.0);
    float ampScale = pow(d, 0.9);

    // Optional: reduce wave frequency in shallow water too (looks nicer)
    // Uncomment if desired:
    // float localFreqScale = mix(0.65, 1.0, ampScale);
    // ...and multiply uWaveFreq by localFreqScale above (requires refactor).
    // We'll keep it simple for now.

    vec2 p = world.xz;
    vec3 wh = waveHeightAndDeriv(p, ampScale);

    // Displace
    vWaveHeight = wh.x;
    world.y += wh.x;

    // Normal from derivatives
    vWaveNormal = normalize(vec3(-wh.y, 1.0, -wh.z));
    vWorldPos = world.xyz;

    // Displace in model-y for mvp convenience (flat plane assumption)
    gl_Position = mvp * vec4(vertexPosition + vec3(0.0, wh.x, 0.0), 1.0);
}
