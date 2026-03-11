#version 430

in vec2 fragTexCoord;
in vec4 fragColor;

out vec4 finalColor;

uniform sampler2D uUnderwaterMask;
uniform vec4 uWaterColor;   // rgba, e.g. (0.0, 0.32, 0.42, 0.35)
uniform float uStrength;    // overall multiplier 0..1

void main()
{
    float mask = texture(uUnderwaterMask, fragTexCoord).r;

    // optional softening
    mask = smoothstep(0.05, 0.95, mask);

    float alpha = mask * uWaterColor.a * uStrength;

    finalColor = vec4(uWaterColor.rgb, alpha) * fragColor;
}
