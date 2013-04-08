#include <assert.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>

#include "lib.h"
#include "raw.h"

static int process_image() {
  int row, col;
  findmin(&row, &col);
  printf("%d %d\n", col, row);
  return 0; /* stop */
}

int main(int argc, char **argv) {
  assert(argc > 1);
  int c = strtol(argv[1], NULL, 16);
  set_target((c & 0x00ff0000) >> 16,
             (c & 0x0000ff00) >>  8,
             (c & 0x000000ff) >>  0);
  process_raw(process_image);
  return 0;
}
