#ifndef _LIB_H
#define _LIB_H

#include <sys/time.h>

#define WIDTH 320
#define HEIGHT 240

/* used during calibration */
#define THRESH 1000
#define DIST 100

#define SIZE 5

void set_target(int r, int g, int b);
void reset();
void process_rgb(void *buf, int size);
void process_yuv(void *buf, int size);
void findmin(int *row, int *col);
void calibrate(void);

int elapsed(struct timeval *start, struct timeval *end);

#endif
