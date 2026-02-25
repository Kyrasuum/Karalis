#version 430

precision highp float;
precision highp sampler2D;

layout (triangles) in;
layout (triangle_strip) out;
layout (max_vertices = 256) out;

in vec4 vertPosition[];
in vec2 vertTexCoord[];
in vec4 vertColor[];
in vec3 vertNormal[];

#ifdef PORTAL_OBJ
in vec4 vertScrPos[];
#endif
#ifdef PORTAL_SCN
in float vertClip[];
#endif

// Output vertex attributes (to fragment shader)
out vec4 fragPosition;
out vec2 fragTexCoord;
out vec4 fragColor;
out vec3 fragNormal;

#ifdef PORTAL_OBJ
out vec4 fragScrPos;
#endif
#ifdef PORTAL_SCN
out float fragClip;
#endif

// Input uniform values
uniform mat4 mvp;
uniform mat4 matModel;
uniform mat4 matView;
uniform mat4 matProjection;

uniform sampler2D texture0;
uniform sampler2D texture1;
uniform sampler2D texture2;

// Heights are expected in normalized [0..1] to match your generation logic.
// Example: grass between ~0.41 and ~0.70 if seaLevel=0.38 and sandBand=0.03.
uniform float uGrassMinHeight;   // e.g. seaLevel + sandBand
uniform float uGrassMaxHeight;   // e.g. mountainStart

// Density controls spawn probability (0..1). You can also treat it as "blades per triangle" if you loop more.
uniform float uGrassDensity;     // e.g. 0.35
uniform float uMaxSlope;         // e.g. 0.65  (normal.y must be >= this to allow grass)

// Blade size controls
uniform float uBladeHeight;      // e.g. 0.18
uniform float uBladeHalfWidth;   // e.g. 0.02

// Deterministic random from world position + seed
uniform float uSeed;

// If your heightmap texture is higher-res and UVs are still 0..1, keep this at (1,1).
// If you atlas or pack tiles, you can scale UVs here.
uniform vec2 uHeightmapTexelScale; // usually vec2(1.0, 1.0)

float hash13(vec3 p) {
    return fract(sin(dot(p, vec3(12.9898, 78.233, 54.53))) * 43758.5453);
}
float clamp01(float x){ return clamp(x, 0.0, 1.0); }
float smoothBand(float a, float b, float x) {
    // smooth on/off with small feather
    float e = 0.02;
    return smoothstep(a, a+e, x) * (1.0 - smoothstep(b-e, b, x));
}

void emitBladeClip(vec4 clipPos, vec3 n, vec4 baseColor, float rnd) {
    float h = uBladeHeight * mix(0.7, 1.3, rnd);
    float w = uBladeHalfWidth;

    // simple triangle blade in clip-space (cheap)
    gl_Position = clipPos + vec4(+w, 0.0, 0.0, 0.0);
    fragTexCoord = vec2(0,0);
    fragNormal = n;
    fragColor = vec4(0.0, 0.30, 0.0, 1.0) * baseColor;
    EmitVertex();

    gl_Position = clipPos + vec4(0.0, h, 0.0, 0.0);
    fragTexCoord = vec2(1,0);
    fragNormal = n;
    fragColor = vec4(0.0, 0.55, 0.0, 1.0) * baseColor;
    EmitVertex();

    gl_Position = clipPos + vec4(-w, 0.0, 0.0, 0.0);
    fragTexCoord = vec2(1,1);
    fragNormal = n;
    fragColor = vec4(0.0, 0.30, 0.0, 1.0) * baseColor;
    EmitVertex();

    EndPrimitive();
}

void main() {
	const int bladesPerTri = 1;
	for (int b = 0; b < bladesPerTri; b++) {
		float r1 = hash13(vec3(float(gl_PrimitiveID), float(b), uSeed));
		float r2 = hash13(vec3(float(gl_PrimitiveID)+19.0, float(b)+7.0, uSeed));
		// Spawn position in CLIP space (robust with raylib’s pipeline)
		vec4 clipPos = mix(mix(gl_in[0].gl_Position, gl_in[1].gl_Position, r1), gl_in[2].gl_Position, r2);

		// Height sample
		float height01 = texture(texture0, vertTexCoord[0]).r;

		// spawn test
		float spawn = hash13(vec3(clipPos.xz, uSeed + float(b) * 13.0));
		if (spawn < uGrassDensity && height01 > uGrassMinHeight) {
			vec4 baseColor = mix(mix(vertColor[0], vertColor[1], r1), vertColor[2], r2);
			emitBladeClip(clipPos, vec3(0,1,0), baseColor, spawn);
		}
	}
}
