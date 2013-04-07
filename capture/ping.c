#include <assert.h>
#include <stdio.h>
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
  process_raw(process_image);
  return 0;
}
