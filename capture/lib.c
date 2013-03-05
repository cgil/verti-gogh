#include <assert.h>
#include <limits.h>

#include "lib.h"

static int nxt;
static int buffer[HEIGHT][WIDTH];

void reset() {
  nxt = 0;
}

#define process(XC, YC, ZC, X, Y, Z)                    \
  int window[SIZE];                                     \
  int i;                                                \
  unsigned char *buf = _buf;                            \
  int sum = 0;                                          \
  for (i = 0; i < SIZE; i++)                            \
    window[i] = 0;                                      \
  for (i = 0; i < WIDTH; i++) {                         \
    int x = XC(i, buf) - X;                             \
    int y = YC(i, buf) - Y;                             \
    int z = ZC(i, buf) - Z;                             \
    int mag = x * x + y * y + z * z;                    \
    sum += mag - window[i % SIZE];                      \
    window[i % SIZE] = mag;                             \
    if (i >= SIZE) {                                    \
      buffer[nxt][i - SIZE] = sum;                      \
    }                                                   \
  }                                                     \
  nxt++;

#define RC(i, buf) buf[i * 3 + 0]
#define GC(i, buf) buf[i * 3 + 1]
#define BC(i, buf) buf[i * 3 + 2]

#define YC(i, buf)  buf[i * 2]
#define CbC(i, buf) buf[(i & 0xfffffffe) * 2 + 1]
#define CrC(i, buf) buf[(i & 0xfffffffe) * 2 + 3]

void process_rgb(void *_buf, int amt) {
  assert(amt == 3 * WIDTH);
  process(RC, GC, BC, R, G, B);
}

void process_yuv(void *_buf, int amt) {
  assert(amt == 2 * WIDTH);
  process(YC, CrC, CbC, Y, Cr, Cb);
}

int elapsed(struct timeval *start, struct timeval *end) {
  return (((end->tv_sec - start->tv_sec) * 1000000) + 
                  (end->tv_usec - start->tv_usec))/1000;
}

void findmin(int *row, int *col) {
  assert(nxt == HEIGHT);

  int i, j;
  
  int mini = 0, minj = 0;
  int min = INT_MAX;
  int windows[WIDTH - SIZE];
  for (i = 0; i < WIDTH - SIZE; i++)
    windows[i] = 0;
  for (i = 0; i < HEIGHT; i++) {
    for (j = 0; j < WIDTH - SIZE; j++) {
      if (i >= SIZE) {
        if (windows[j] < min) {
          min = windows[j];
          mini = i;
          minj = j;
        }
        windows[j] -= buffer[i - SIZE][j];
      }
      windows[j] += buffer[i][j];
    }
  }
  *row = mini;
  *col = minj;
}
