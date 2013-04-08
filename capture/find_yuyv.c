#include <assert.h>
#include <stdio.h>

#include "lib.h"

int main(int argc, char *argv[]) {
  assert(argc > 1);
  FILE *f = fopen(argv[1], "r");
  assert(f != NULL);

  set_target(0x98, 0x64, 0x7a);

  int i;
  char buf[WIDTH * 2];
  for (i = 0; i < HEIGHT; i++) {
    fread(buf, sizeof(buf), 1, f);
    assert(!feof(f));
    process_yuv(buf, sizeof(buf));
  }
  fclose(f);

  int row, col;
  findmin(&row, &col);
  printf("0 %d %d\n", col, row);
  return 0;
}
