#include <assert.h>
#include <stdio.h>
#include <string.h>

#include "jpeglib.h"
#include "lib.h"

void decode(char *filename) {
  struct jpeg_decompress_struct cinfo;
  struct jpeg_error_mgr err;
  cinfo.err = jpeg_std_error(&err);

  FILE *infile = fopen(filename, "rb");
  assert(infile != NULL);

  jpeg_create_decompress(&cinfo);
  jpeg_stdio_src(&cinfo, infile);

  jpeg_read_header(&cinfo, TRUE);
  jpeg_start_decompress(&cinfo);

  int row_stride = cinfo.output_width * cinfo.output_components;
  JSAMPARRAY buffer = (*cinfo.mem->alloc_sarray)
    ((j_common_ptr) &cinfo, JPOOL_IMAGE, row_stride, 1);
  assert(buffer != NULL);

  while (cinfo.output_scanline < cinfo.output_height) {
    jpeg_read_scanlines(&cinfo, buffer, 1);
    process_rgb(buffer[0], row_stride);
  }

  jpeg_finish_decompress(&cinfo);
  jpeg_destroy_decompress(&cinfo);
  fclose(infile);
}

int main(int argc, char *argv[]) {
  assert(argc > 1);
  reset();
  decode(argv[1]);

  int row, col;
  findmin(&row, &col);
  printf("0 %d %d\n", row, col);
  return 0;
}
