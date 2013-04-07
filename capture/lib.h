#ifndef _LIB_H
#define _LIB_H

#include <sys/time.h>

#define WIDTH 320
#define HEIGHT 240

/* used during calibration */
#define THRESH 1000
#define DIST 100

#define R 0x00
#define G 0x00
#define B 0x00
#define SIZE 5

/* http://en.wikipedia.org/wiki/YCbCr */
#define Y  ((int) (0.299 * R + 0.587 * G + 0.114 * B))
#define Cb ((int) (128 - 0.168736 * R - 0.331264 * G + 0.5 * B))
#define Cr ((int) (128 + 0.5 * R - 0.418688 * G - 0.081312 * B))

void reset();
void process_rgb(void *buf, int size);
void process_yuv(void *buf, int size);
void findmin(int *row, int *col);
void calibrate(void);

int elapsed(struct timeval *start, struct timeval *end);

#endif
