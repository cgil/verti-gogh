#include <assert.h>
#include <stdio.h>

#include "lib.h"
#include "raw.h"

static int process_image() {
  calibrate();
  return 0; /* stop */
}

int main(int argc, char **argv) {
  process_raw(process_image);
  return 0;
}
