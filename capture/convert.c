#include <assert.h>
#include <stdio.h>
#include <unistd.h>
#include <jpeglib.h>

#include "lib.h"

int main() {
  struct jpeg_compress_struct cinfo;
  struct jpeg_error_mgr jerr;
  FILE * outfile;     /* target file */
  JSAMPROW row_pointer[1];    /* pointer to JSAMPLE row[s] */

  outfile = fopen("outfile.jpg", "wb");
  assert(outfile != NULL);
  jpeg_stdio_dest(&cinfo, outfile);

  cinfo.err = jpeg_std_error(&jerr);
  jpeg_create_compress(&cinfo);

  cinfo.client_data = (void*)&outfile;
  cinfo.image_width = WIDTH;
  cinfo.image_height = HEIGHT;
  cinfo.input_components = 3;     /* # of color components per pixel */
  cinfo.in_color_space = JCS_RGB; /* colorspace of input image */

  jpeg_set_defaults(&cinfo);
  jpeg_set_quality(&cinfo, 100, TRUE);
  jpeg_start_compress(&cinfo, TRUE);

  char input[WIDTH * 2];
  char lineBuffer[WIDTH * 3];
  int i, j;
  for (i = 0; i < HEIGHT; i++) {
    int b = read(STDIN_FILENO, input, sizeof(input));
    assert(b == sizeof(input));
    for (j = 0; j < WIDTH; j++) {
      int y = input[j * 2];
      int cB = input[(j & 0xfffffffe) * 2 + 1];
      int cR = input[(j & 0xfffffffe) * 2 + 3];

      int r   =   1.0 * y    + 0 * cB    + 1.402 * cR;
      int g   =   1.0 * y    - 0.344136 * cB - 0.714136 * cR;
      int b   =   1.0 * y    + 1.772 * cB    + 0 * cR;

      lineBuffer[j * 3 + 0] = r;
      lineBuffer[j * 3 + 1] = g;
      lineBuffer[j * 3 + 2] = b;
    }

    row_pointer[0] = (void*) lineBuffer;
    jpeg_write_scanlines(&cinfo, row_pointer, 1);
  }

  jpeg_finish_compress(&cinfo);
  fclose(outfile);
  jpeg_destroy_compress(&cinfo);
  return 0;
}
