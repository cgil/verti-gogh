#include <assert.h>
#include <stdlib.h>
#include <stdio.h>

#include "lib.h"
#include "raw.h"

static int process_image() {
  int row, col;
  double ret = findmin(&row, &col);
  printf("%d %d %f\n", row, col, ret);
  fflush(stdout);

  static char buf[1024];
  if (fgets(buf, sizeof(buf), stdin) == NULL)
    return 0;
  return 1; /* keep going */
}

int main(int argc, char **argv) {
  assert(argc > 5);

  int c = strtol(argv[1], NULL, 16);
  set_target((c & 0x00ff0000) >> 16,
             (c & 0x0000ff00) >>  8,
             (c & 0x000000ff) >>  0);

  int topx = atoi(argv[2]);
  int topy = atoi(argv[3]);
  int botx = atoi(argv[4]);
  int boty = atoi(argv[5]);

  set_bounds(topx, topy, botx, boty);
  process_raw(process_image);
  return 0;
}
