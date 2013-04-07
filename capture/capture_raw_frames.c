#include <assert.h>
#include <stdio.h>

#include "lib.h"
#include "raw.h"

static int process_image() {
  int row, col;
  findmin(&row, &col);
  printf("0 %d %d\n", row, col);
  fflush(stdout);
  return 1; /* keep going */
}

int main(int argc, char **argv) {
  process_raw(process_image);
  return 0;
}
