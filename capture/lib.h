#ifndef _LIB_H
#define _LIB_H

#include <sys/time.h>

#define R 0x00
#define G 0x39
#define B 0x89
#define SIZE 20

/* http://en.wikipedia.org/wiki/YCbCr */
#define Y  ((int) (16 + 65.481 * R + 128.553 * G + 24.966 * B))
#define Cb ((int) (128 + -37.797 * R - 74.203 * G + 112.0 * B))
#define Cr ((int) (128 + 112.0 * R - 93.786 * G + 18.214 * B))

void reset();
void process_rgb(void *buf, int size);
void process_yuv(void *buf, int size);
void findmin(int *row, int *col);

int elapsed(struct timeval *start, struct timeval *end);

#endif
