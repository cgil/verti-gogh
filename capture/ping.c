#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lib.h"
#include "raw.h"

static int first = 1;

static int process_image() {
  if (!first) {
    int row, col;
    findmin(&row, &col);
    printf("%d %d\n", row, col);
    fflush(stdout);
  }
  first = 0;

  /* wait on stdin for something to happen */
  char buf[1024];
  size_t ret = fread(buf, 1, sizeof(buf), stdin);

  if (strncmp(buf, "stop\n", ret) == 0)
    return 0;
  return 1;
}

int main(int argc, char **argv) {
  if (argc > 1) {
    int c = atoi(argv[1]);
    set_target((c & 0x00ff0000) >> 16,
               (c & 0x0000ff00) >>  8,
               (c & 0x000000ff) >>  0);
  }
  process_raw(process_image);
  return 0;
}