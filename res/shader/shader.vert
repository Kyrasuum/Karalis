#version 430

// Input vertex attributes
in vec3 inPosition;
in vec2 inTexCoord;
in vec3 inNormal;
in vec4 inColor;
in vec4 inTangent;
in vec4 inTexCoord2;

#ifdef GRASS_OBJ
layout(std430, binding = 1) buffer InstanceBuf {
    struct GrassInstance {
        vec4 pos_height;
        vec4 n_rot;
    } instances[];
};
#endif

// Input uniform values
uniform mat4 mvp;
uniform mat4 matModel;
uniform mat4 matView;
uniform mat4 matProjection;
uniform mat4 matNormal;

// Output vertex attributes (to next shader)
out vec4 vertPosition;
out vec2 vertTexCoord;
out vec4 vertColor;
out vec3 vertNormal;

#ifdef PORTAL_OBJ
out vec4 vertScrPos;
#endif
#ifdef PORTAL_SCN
out float vertClip;
#endif

// Custom Uniforms
uniform vec3 portalPos = vec3(0,0,0);
uniform vec3 portalNorm = vec3(0,0,0);
uniform mat4 portalMat = mat4(1.0);

mat2 rot2(float a){
    float s = sin(a), c = cos(a);
    return mat2(c,-s,s,c);
}

void main() {
    // grass based positioning
    #ifdef GRASS_OBJ
    uint id = uint(gl_InstanceID);
    vec3 basePos = instances[id].pos_height.xyz;
    float h      = instances[id].pos_height.w;
    float rot    = instances[id].n_rot.w;

    // rotate around Y in local grass space
    local.xz = rot2(rot) * inPosition.xz;
    // scale height
    local.y *= h;
    inPosition = basePos + inPosition
    #endif

    // Calculate final vertex position
    gl_Position = mvp*portalMat*portalMat*vec4(inPosition, 1.0);

    // Send vertex attributes to next shader
    vertPosition = matModel*vec4(inPosition, 1.0);
    vertTexCoord = inTexCoord;
    vertColor = inColor;
    vertNormal = normalize(vec3(matNormal*vec4(inNormal, 1.0)));

    #ifdef PORTAL_OBJ
    vertScrPos = gl_Position;
    #endif
    #ifdef PORTAL_SCN
    vertClip = dot((vec3(vertPosition) - portalPos), portalNorm);
    #endif
}
