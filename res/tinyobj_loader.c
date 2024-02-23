#ifndef TINYOBJ_LOADER_C_IMPLEMENTATION
#define TINYOBJ_LOADER_C_IMPLEMENTATION
#include "tinyobj_loader.h"
#include "stddef.h"
#include "stdio.h"
#include "string.h"
#include "errno.h"
#include "assert.h"

int tinyobj_parse_and_index_mtl_file2(tinyobj_material_t **materials_out,
                                            unsigned int *num_materials_out,
                                            const char *filename,
                                            hash_table_t* material_table,
                                            const char *workingDir) {
  tinyobj_material_t material;
  unsigned int buffer_size = 128;
  char *linebuf;
  FILE *fp;
  unsigned int num_materials = 0;
  tinyobj_material_t *materials = NULL;
  int has_previous_material = 0;
  const char *line_end = NULL;

  if (materials_out == NULL) {
    return TINYOBJ_ERROR_INVALID_PARAMETER;
  }

  if (num_materials_out == NULL) {
    return TINYOBJ_ERROR_INVALID_PARAMETER;
  }

  (*materials_out) = NULL;
  (*num_materials_out) = 0;

  fp = ReadData(filename, workingDir);
  if (!fp) {
    fprintf(stderr, "TINYOBJ: Error reading file '%s': %s (%d)\n", filename, strerror(errno), errno);
    return TINYOBJ_ERROR_FILE_OPERATION;
  }

  /* Create a default material */
  initMaterial(&material);

  linebuf = (char*)TINYOBJ_MALLOC(buffer_size);
  while (NULL != dynamic_fgets(&linebuf, &buffer_size, fp)) {
    const char *token = linebuf;

    line_end = token + strlen(token);

    /* Skip leading space. */
    token += strspn(token, " \t");

    assert(token);
    if (token[0] == '\0') continue; /* empty line */

    if (token[0] == '#') continue; /* comment line */

    /* new mtl */
    if ((0 == strncmp(token, "newmtl", 6)) && IS_SPACE((token[6]))) {
      char namebuf[4096];

      /* flush previous material. */
      if (has_previous_material) {
        materials = tinyobj_material_add(materials, num_materials, &material);
        num_materials++;
      } else {
        has_previous_material = 1;
      }

      /* initial temporary material */
      initMaterial(&material);

      /* set new mtl name */
      token += 7;
#ifdef _MSC_VER
      sscanf_s(token, "%s", namebuf, (unsigned)_countof(namebuf));
#else
      sscanf(token, "%s", namebuf);
#endif
      material.name = my_strdup(namebuf, (unsigned int) (line_end - token));

      /* Add material to material table */
      if (material_table)
        hash_table_set(material.name, num_materials, material_table);

      continue;
    }

    /* ambient */
    if (token[0] == 'K' && token[1] == 'a' && IS_SPACE((token[2]))) {
      float r, g, b;
      token += 2;
      parseFloat3(&r, &g, &b, &token);
      material.ambient[0] = r;
      material.ambient[1] = g;
      material.ambient[2] = b;
      continue;
    }

    /* diffuse */
    if (token[0] == 'K' && token[1] == 'd' && IS_SPACE((token[2]))) {
      float r, g, b;
      token += 2;
      parseFloat3(&r, &g, &b, &token);
      material.diffuse[0] = r;
      material.diffuse[1] = g;
      material.diffuse[2] = b;
      continue;
    }

    /* specular */
    if (token[0] == 'K' && token[1] == 's' && IS_SPACE((token[2]))) {
      float r, g, b;
      token += 2;
      parseFloat3(&r, &g, &b, &token);
      material.specular[0] = r;
      material.specular[1] = g;
      material.specular[2] = b;
      continue;
    }

    /* transmittance */
    if (token[0] == 'K' && token[1] == 't' && IS_SPACE((token[2]))) {
      float r, g, b;
      token += 2;
      parseFloat3(&r, &g, &b, &token);
      material.transmittance[0] = r;
      material.transmittance[1] = g;
      material.transmittance[2] = b;
      continue;
    }

    /* ior(index of refraction) */
    if (token[0] == 'N' && token[1] == 'i' && IS_SPACE((token[2]))) {
      token += 2;
      material.ior = parseFloat(&token);
      continue;
    }

    /* emission */
    if (token[0] == 'K' && token[1] == 'e' && IS_SPACE(token[2])) {
      float r, g, b;
      token += 2;
      parseFloat3(&r, &g, &b, &token);
      material.emission[0] = r;
      material.emission[1] = g;
      material.emission[2] = b;
      continue;
    }

    /* shininess */
    if (token[0] == 'N' && token[1] == 's' && IS_SPACE(token[2])) {
      token += 2;
      material.shininess = parseFloat(&token);
      continue;
    }

    /* illum model */
    if (0 == strncmp(token, "illum", 5) && IS_SPACE(token[5])) {
      token += 6;
      material.illum = parseInt(&token);
      continue;
    }

    /* dissolve */
    if ((token[0] == 'd' && IS_SPACE(token[1]))) {
      token += 1;
      material.dissolve = parseFloat(&token);
      continue;
    }
    if (token[0] == 'T' && token[1] == 'r' && IS_SPACE(token[2])) {
      token += 2;
      /* Invert value of Tr(assume Tr is in range [0, 1]) */
      material.dissolve = 1.0f - parseFloat(&token);
      continue;
    }

    /* ambient texture */
    if ((0 == strncmp(token, "map_Ka", 6)) && IS_SPACE(token[6])) {
      token += 7;
      material.ambient_texname = my_strdup(token, (unsigned int) (line_end - token));
      continue;
    }

    /* diffuse texture */
    if ((0 == strncmp(token, "map_Kd", 6)) && IS_SPACE(token[6])) {
      token += 7;
      material.diffuse_texname = my_strdup(token, (unsigned int) (line_end - token));
      continue;
    }

    /* specular texture */
    if ((0 == strncmp(token, "map_Ks", 6)) && IS_SPACE(token[6])) {
      token += 7;
      material.specular_texname = my_strdup(token, (unsigned int) (line_end - token));
      continue;
    }

    /* specular highlight texture */
    if ((0 == strncmp(token, "map_Ns", 6)) && IS_SPACE(token[6])) {
      token += 7;
      material.specular_highlight_texname = my_strdup(token, (unsigned int) (line_end - token));
      continue;
    }

    /* bump texture */
    if ((0 == strncmp(token, "map_bump", 8)) && IS_SPACE(token[8])) {
      token += 9;
      material.bump_texname = my_strdup(token, (unsigned int) (line_end - token));
      continue;
    }

    /* alpha texture */
    if ((0 == strncmp(token, "map_d", 5)) && IS_SPACE(token[5])) {
      token += 6;
      material.alpha_texname = my_strdup(token, (unsigned int) (line_end - token));
      continue;
    }

    /* bump texture */
    if ((0 == strncmp(token, "bump", 4)) && IS_SPACE(token[4])) {
      token += 5;
      material.bump_texname = my_strdup(token, (unsigned int) (line_end - token));
      continue;
    }

    /* displacement texture */
    if ((0 == strncmp(token, "disp", 4)) && IS_SPACE(token[4])) {
      token += 5;
      material.displacement_texname = my_strdup(token, (unsigned int) (line_end - token));
      continue;
    }

    /* @todo { unknown parameter } */
  }

  fclose(fp);

  if (material.name) {
    /* Flush last material element */
    materials = tinyobj_material_add(materials, num_materials, &material);
    num_materials++;
  }

  (*num_materials_out) = num_materials;
  (*materials_out) = materials;

  if (linebuf) {
    TINYOBJ_FREE(linebuf);
  }

  return TINYOBJ_SUCCESS;
}

int tinyobj_parse_mtl_file2(tinyobj_material_t **materials_out,
                           unsigned int *num_materials_out,
                           const char *filename,
                           const char *workingDir) {
  return tinyobj_parse_and_index_mtl_file2(materials_out, num_materials_out, filename, NULL, workingDir);
}


int tinyobj_parse_obj2(tinyobj_attrib_t *attrib, tinyobj_shape_t **shapes,
                      unsigned int *num_shapes, tinyobj_material_t **materials_out,
                      unsigned int *num_materials_out, const char *buf, unsigned int len,
                      unsigned int flags, const char* workingDir) {
  LineInfo *line_infos = NULL;
  Command *commands = NULL;
  unsigned int num_lines = 0;

  unsigned int num_v = 0;
  unsigned int num_vn = 0;
  unsigned int num_vt = 0;
  unsigned int num_f = 0;
  unsigned int num_faces = 0;

  int mtllib_line_index = -1;

  tinyobj_material_t *materials = NULL;
  unsigned int num_materials = 0;

  hash_table_t material_table;

  if (len < 1) return TINYOBJ_ERROR_INVALID_PARAMETER;
  if (attrib == NULL) return TINYOBJ_ERROR_INVALID_PARAMETER;
  if (shapes == NULL) return TINYOBJ_ERROR_INVALID_PARAMETER;
  if (num_shapes == NULL) return TINYOBJ_ERROR_INVALID_PARAMETER;
  if (buf == NULL) return TINYOBJ_ERROR_INVALID_PARAMETER;
  if (materials_out == NULL) return TINYOBJ_ERROR_INVALID_PARAMETER;
  if (num_materials_out == NULL) return TINYOBJ_ERROR_INVALID_PARAMETER;

  tinyobj_attrib_init(attrib);
   /* 1. Find '\n' and create line data. */
  {
    unsigned int i;
    unsigned int end_idx = len;
    unsigned int prev_pos = 0;
    unsigned int line_no = 0;
    unsigned int last_line_ending = 0;

    /* Count # of lines. */
    for (i = 0; i < end_idx; i++) {
      if (is_line_ending(buf, i, end_idx)) {
        num_lines++;
        last_line_ending = i;
      }
    }
    /* The last char from the input may not be a line
     * ending character so add an extra line if there
     * are more characters after the last line ending
     * that was found. */
    if (end_idx - last_line_ending > 0) {
        num_lines++;
    }

    if (num_lines == 0) return TINYOBJ_ERROR_EMPTY;

    line_infos = (LineInfo *)TINYOBJ_MALLOC(sizeof(LineInfo) * num_lines);

    /* Fill line infos. */
    for (i = 0; i < end_idx; i++) {
      if (is_line_ending(buf, i, end_idx)) {
        line_infos[line_no].pos = prev_pos;
        line_infos[line_no].len = i - prev_pos;

// ---- QUICK BUG FIX : https://github.com/raysan5/raylib/issues/3473
        if ( i > 0 && buf[i-1] == '\r' ) line_infos[line_no].len--;
// --------

        prev_pos = i + 1;
        line_no++;
      }
    }
    if (end_idx - last_line_ending > 0) {
      line_infos[line_no].pos = prev_pos;
      line_infos[line_no].len = end_idx - 1 - last_line_ending;
    }
  }

  commands = (Command *)TINYOBJ_MALLOC(sizeof(Command) * num_lines);

  create_hash_table(HASH_TABLE_DEFAULT_SIZE, &material_table);

  /* 2. parse each line */
  {
    unsigned int i = 0;
    for (i = 0; i < num_lines; i++) {
      int ret = parseLine(&commands[i], &buf[line_infos[i].pos],
                          line_infos[i].len, flags & TINYOBJ_FLAG_TRIANGULATE);
      if (ret) {
        if (commands[i].type == COMMAND_V) {
          num_v++;
        } else if (commands[i].type == COMMAND_VN) {
          num_vn++;
        } else if (commands[i].type == COMMAND_VT) {
          num_vt++;
        } else if (commands[i].type == COMMAND_F) {
          num_f += commands[i].num_f;
          num_faces += commands[i].num_f_num_verts;
        }

        if (commands[i].type == COMMAND_MTLLIB) {
          mtllib_line_index = (int)i;
        }
      }
    }
  }

  /* line_infos are not used anymore. Release memory. */
  if (line_infos) {
    TINYOBJ_FREE(line_infos);
  }

  /* Load material(if exits) */
  if (mtllib_line_index >= 0 && commands[mtllib_line_index].mtllib_name &&
      commands[mtllib_line_index].mtllib_name_len > 0) {
    char *filename = my_strndup(commands[mtllib_line_index].mtllib_name,
                                commands[mtllib_line_index].mtllib_name_len);

    int ret = tinyobj_parse_and_index_mtl_file2(&materials, &num_materials, filename, &material_table, workingDir);

    if (ret != TINYOBJ_SUCCESS) {
      /* warning. */
      fprintf(stderr, "TINYOBJ: Failed to parse material file '%s': %d\n", filename, ret);
    }

    TINYOBJ_FREE(filename);

  }

  /* Construct attributes */

  {
    unsigned int v_count = 0;
    unsigned int n_count = 0;
    unsigned int t_count = 0;
    unsigned int f_count = 0;
    unsigned int face_count = 0;
    int material_id = -1; /* -1 = default unknown material. */
    unsigned int i = 0;

    attrib->vertices = (float *)TINYOBJ_MALLOC(sizeof(float) * num_v * 3);
    attrib->num_vertices = (unsigned int)num_v;
    attrib->normals = (float *)TINYOBJ_MALLOC(sizeof(float) * num_vn * 3);
    attrib->num_normals = (unsigned int)num_vn;
    attrib->texcoords = (float *)TINYOBJ_MALLOC(sizeof(float) * num_vt * 2);
    attrib->num_texcoords = (unsigned int)num_vt;
    attrib->faces = (tinyobj_vertex_index_t *)TINYOBJ_MALLOC(sizeof(tinyobj_vertex_index_t) * num_f);
    attrib->face_num_verts = (int *)TINYOBJ_MALLOC(sizeof(int) * num_faces);

    attrib->num_faces = (unsigned int)num_faces;
    attrib->num_face_num_verts = (unsigned int)num_f;

    attrib->material_ids = (int *)TINYOBJ_MALLOC(sizeof(int) * num_faces);

    for (i = 0; i < num_lines; i++) {
      if (commands[i].type == COMMAND_EMPTY) {
        continue;
      } else if (commands[i].type == COMMAND_USEMTL) {
        /* @todo
           if (commands[t][i].material_name &&
           commands[t][i].material_name_len > 0) {
           std::string material_name(commands[t][i].material_name,
           commands[t][i].material_name_len);

           if (material_map.find(material_name) != material_map.end()) {
           material_id = material_map[material_name];
           } else {
        // Assign invalid material ID
        material_id = -1;
        }
        }
        */
        if (commands[i].material_name &&
           commands[i].material_name_len >0)
        {
          /* Create a null terminated string */
          char* material_name_null_term = (char*) TINYOBJ_MALLOC(commands[i].material_name_len + 1);
          memcpy((void*) material_name_null_term, (const void*) commands[i].material_name, commands[i].material_name_len);
          material_name_null_term[commands[i].material_name_len] = 0;

          if (hash_table_exists(material_name_null_term, &material_table))
            material_id = (int)hash_table_get(material_name_null_term, &material_table);
          else
            material_id = -1;

          TINYOBJ_FREE(material_name_null_term);
        }
      } else if (commands[i].type == COMMAND_V) {
        attrib->vertices[3 * v_count + 0] = commands[i].vx;
        attrib->vertices[3 * v_count + 1] = commands[i].vy;
        attrib->vertices[3 * v_count + 2] = commands[i].vz;
        v_count++;
      } else if (commands[i].type == COMMAND_VN) {
        attrib->normals[3 * n_count + 0] = commands[i].nx;
        attrib->normals[3 * n_count + 1] = commands[i].ny;
        attrib->normals[3 * n_count + 2] = commands[i].nz;
        n_count++;
      } else if (commands[i].type == COMMAND_VT) {
        attrib->texcoords[2 * t_count + 0] = commands[i].tx;
        attrib->texcoords[2 * t_count + 1] = commands[i].ty;
        t_count++;
      } else if (commands[i].type == COMMAND_F) {
        unsigned int k = 0;
        for (k = 0; k < commands[i].num_f; k++) {
          tinyobj_vertex_index_t vi = commands[i].f[k];
          int v_idx = fixIndex(vi.v_idx, v_count);
          int vn_idx = fixIndex(vi.vn_idx, n_count);
          int vt_idx = fixIndex(vi.vt_idx, t_count);
          attrib->faces[f_count + k].v_idx = v_idx;
          attrib->faces[f_count + k].vn_idx = vn_idx;
          attrib->faces[f_count + k].vt_idx = vt_idx;
        }

        for (k = 0; k < commands[i].num_f_num_verts; k++) {
          attrib->material_ids[face_count + k] = material_id;
          attrib->face_num_verts[face_count + k] = commands[i].f_num_verts[k];
        }

        f_count += commands[i].num_f;
        face_count += commands[i].num_f_num_verts;
      }
    }
  }

  /* 5. Construct shape information. */
  {
    unsigned int face_count = 0;
    unsigned int i = 0;
    unsigned int n = 0;
    unsigned int shape_idx = 0;

    const char *shape_name = NULL;
    unsigned int shape_name_len = 0;
    const char *prev_shape_name = NULL;
    unsigned int prev_shape_name_len = 0;
    unsigned int prev_shape_face_offset = 0;
    unsigned int prev_face_offset = 0;
    tinyobj_shape_t prev_shape = {NULL, 0, 0};

    /* Find the number of shapes in .obj */
    for (i = 0; i < num_lines; i++) {
      if (commands[i].type == COMMAND_O || commands[i].type == COMMAND_G) {
        n++;
      }
    }

    /* Allocate array of shapes with maximum possible size(+1 for unnamed
     * group/object).
     * Actual # of shapes found in .obj is determined in the later */
    (*shapes) = (tinyobj_shape_t*)TINYOBJ_MALLOC(sizeof(tinyobj_shape_t) * (n + 1));

    for (i = 0; i < num_lines; i++) {
      if (commands[i].type == COMMAND_O || commands[i].type == COMMAND_G) {
        if (commands[i].type == COMMAND_O) {
          shape_name = commands[i].object_name;
          shape_name_len = commands[i].object_name_len;
        } else {
          shape_name = commands[i].group_name;
          shape_name_len = commands[i].group_name_len;
        }

        if (face_count == 0) {
          /* 'o' or 'g' appears before any 'f' */
          prev_shape_name = shape_name;
          prev_shape_name_len = shape_name_len;
          prev_shape_face_offset = face_count;
          prev_face_offset = face_count;
        } else {
          if (shape_idx == 0) {
            /* 'o' or 'g' after some 'v' lines. */
            (*shapes)[shape_idx].name = my_strndup(
                                                   prev_shape_name, prev_shape_name_len); /* may be NULL */
            (*shapes)[shape_idx].face_offset = prev_shape.face_offset;
            (*shapes)[shape_idx].length = face_count - prev_face_offset;
            shape_idx++;

            prev_face_offset = face_count;

          } else {
            if ((face_count - prev_face_offset) > 0) {
              (*shapes)[shape_idx].name =
                my_strndup(prev_shape_name, prev_shape_name_len);
              (*shapes)[shape_idx].face_offset = prev_face_offset;
              (*shapes)[shape_idx].length = face_count - prev_face_offset;
              shape_idx++;
              prev_face_offset = face_count;
            }
          }

          /* Record shape info for succeeding 'o' or 'g' command. */
          prev_shape_name = shape_name;
          prev_shape_name_len = shape_name_len;
          prev_shape_face_offset = face_count;
        }
      }
      if (commands[i].type == COMMAND_F) {
        face_count++;
      }
    }

    if ((face_count - prev_face_offset) > 0) {
      unsigned int length = face_count - prev_shape_face_offset;
      if (length > 0) {
        (*shapes)[shape_idx].name =
          my_strndup(prev_shape_name, prev_shape_name_len);
        (*shapes)[shape_idx].face_offset = prev_face_offset;
        (*shapes)[shape_idx].length = face_count - prev_face_offset;
        shape_idx++;
      }
    } else {
      /* Guess no 'v' line occurrence after 'o' or 'g', so discards current
       * shape information. */
    }

    (*num_shapes) = shape_idx;
  }

  if (commands) {
    TINYOBJ_FREE(commands);
  }

  destroy_hash_table(&material_table);

  (*materials_out) = materials;
  (*num_materials_out) = num_materials;

  return TINYOBJ_SUCCESS;
}

#endif /* TINYOBJ_LOADER_C_IMPLEMENTATION */
