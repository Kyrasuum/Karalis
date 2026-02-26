#version 430
// raylib provides these for custom shaders
in vec3 vertexPosition;
in vec3 vertexNormal;
in vec2 vertexTexCoord;

in mat4 instanceTransform; 

out vec3 vN;
out float vFade;

uniform float uTime;
uniform mat4 mvp;
uniform mat4 uRot;

// Visible blades (compacted) at binding=1
struct Blade { vec4 posH; };
layout(std430, binding = 1) readonly buffer VisibleBlades { Blade blades[]; };

void main() {
    uint iid = uint(gl_InstanceID);
    vec4 posH = blades[iid].posH;

    vec3 base = posH.xyz;
    float H = posH.w;

    // Simple wind (cheap)
    float w = sin(uTime*1.7 + base.x*0.35 + base.z*0.35) * 0.12;
    vec3 bend = vec3(w, 0.0, w*0.6);

    vec3 local = (uRot * vec4(vertexPosition, 1.0)).xyz;

    // Bend more near the tip
    float tip = clamp(local.y / max(H, 0.001), 0.0, 1.0);
    local += bend * (tip*tip);

    vec3 worldPos = base + local;

    gl_Position = mvp * instanceTransform * vec4(worldPos, 1.0);

    vN = vertexNormal;
    vFade = 1.0 - tip*0.15; // tiny darkening at tip
}
