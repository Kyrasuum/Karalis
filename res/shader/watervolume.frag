#version 330

in vec3 vWorldPos;
in vec3 vVertPos;
in vec2 vUv;
in vec3 vWaveNormal;
in float vWaveHeight;

out vec4 finalColor;

void main() {
    vec3 col = vec3(1.0);
    if (vVertPos.y >= 0.0) {
    	col = vec3(0.0);
    }

    finalColor = vec4(col, 1.0);
}
