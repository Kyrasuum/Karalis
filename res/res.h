#include "raylib.h"
#include "tinyobj_loader.h"

extern char* GetData(char* file, char* dir);
Model* LoadOBJ(char *fileName, char *fileText);
void ProcessMaterialsOBJ(Material *rayMaterials, tinyobj_material_t *materials, int materialCount);
extern unsigned int rlGetTextureIdDefault(void);
