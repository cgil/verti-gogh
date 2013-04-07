/*
 *  V4L2 video capture example
 *
 *  This program can be used and distributed without restrictions.
 *
 *      This program is provided with the V4L2 API
 * see http://linuxtv.org/docs.php for more information
 */


#include "lib.h"

static void process_image(void *p, int size) {
  assert(size == WIDTH * HEIGHT * 2);
  reset();

  char *buf = p;
  int i;
  for (i = 0; i < HEIGHT; i++) {
    process_yuv(buf + i * WIDTH * 2, WIDTH * 2);
  }

  int row, col;
  findmin(&row, &col);
  printf("0 %d %d\n", row, col);
  fflush(stdout);
}

int main(int argc, char **argv) {
  process_raw(process_image);
  return 0;
}
