#include <assert.h>
#include <math.h>
#include <limits.h>
#include <stdio.h>

#include "lib.h"
#include "colorspace.h"

static double H = 0;
static double S = 0;
static double L = 0;

static int nxt;
double buffer[HEIGHT][WIDTH];

void set_target(int r, int g, int b) {
  Rgb2Hsl(&H, &S, &L, r / 255.0, g / 255.0, b / 255.0);
}

void reset() {
  nxt = 0;
}

#define process(FETCH)                                  \
  double window[SIZE];                                  \
  int i;                                                \
  unsigned char *buf = _buf;                            \
  double sum = 0;                                          \
  for (i = 0; i < SIZE; i++)                            \
    window[i] = 0;                                      \
  for (i = 0; i < WIDTH; i++) {                         \
    double r, g, b;                                     \
    FETCH;                                              \
    double h, s, l;                                     \
    Rgb2Hsl(&h, &s, &l, r, g, b);                       \
    double x = h - H;                                   \
    double y = s - S;                                   \
    double z = l - L;                                   \
    double mag = x * x + y * y + z * z;                 \
    sum += mag - window[i % SIZE];                      \
    window[i % SIZE] = mag;                             \
    if (i >= SIZE) {                                    \
      buffer[nxt][i - SIZE] = sum;                      \
    }                                                   \
  }                                                     \
  nxt++;

#define YC(i, buf)  buf[i * 2]
#define CbC(i, buf) buf[(i & 0xfffffffe) * 2 + 1]
#define CrC(i, buf) buf[(i & 0xfffffffe) * 2 + 3]

void process_rgb(void *_buf, int amt) {
  assert(amt == 3 * WIDTH);
  process({
      r = buf[i * 3 + 0] / 255.0;
      g = buf[i * 3 + 1] / 255.0;
      b = buf[i * 3 + 2] / 255.0;
  });
}

void process_yuv(void *_buf, int amt) {
  assert(amt == 2 * WIDTH);
  process({
      double y  = buf[i * 2];
      double cb = buf[(i & 0xfffffffe) * 2 + 1];
      double cr = buf[(i & 0xfffffffe) * 2 + 3];
      Ycbcr2Rgb(&r, &g, &b, y, cb, cr);
  });
}

int elapsed(struct timeval *start, struct timeval *end) {
  return (((end->tv_sec - start->tv_sec) * 1000000) +
         (end->tv_usec - start->tv_usec))/1000;
}

void findmin(int *row, int *col) {
  assert(nxt == HEIGHT);

  int i, j;

  int mini = 0, minj = 0;
  double min = 10000000000000000;
  double windows[WIDTH - SIZE];
  for (i = 0; i < WIDTH - SIZE; i++)
    windows[i] = 0;
  for (i = 0; i < HEIGHT; i++) {
    for (j = 0; j < WIDTH - SIZE; j++) {
      assert(buffer[i][j] >= 0);
      if (i >= SIZE) {
        assert(windows[j] >= 0);
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

void calibrate() {
  int i, j, k;

  int windows[WIDTH - SIZE];
  for (i = 0; i < WIDTH - SIZE; i++)
    windows[i] = 0;

  struct {
    int x;
    int y;
    int score;
  } corners[4];
  for (i = 0; i < 4; i++)
    corners[i].score = INT_MAX;
  corners[0].x = 0;
  corners[0].y = 0;
  corners[1].x = WIDTH;
  corners[1].y = 0;
  corners[2].x = 0;
  corners[2].y = HEIGHT;
  corners[3].x = WIDTH;
  corners[3].y = HEIGHT;

  for (i = 0; i < HEIGHT; i++) {
    for (j = 0; j < WIDTH - SIZE; j++) {
      if (i >= SIZE) {
        int pix = windows[j];

        for (k = 0; k < 4; k++) {
          if (pix < corners[k].score) {
            int xdiff = corners[k].x - j;
            int ydiff = corners[k].y - i;
            if (xdiff * xdiff + ydiff * ydiff < DIST * DIST) {
              corners[k].x = j;
              corners[k].y = i;
              corners[k].score = pix;
            }
          }
        }
        windows[j] -= buffer[i - SIZE][j];
      }
      windows[j] += buffer[i][j];
    }
  }

  for (i = 0; i < 4; i++) {
    printf("%d %d\n", corners[i].x, corners[i].y);
  }
}
