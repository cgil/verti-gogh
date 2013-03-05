CFLAGS += -Wall -Wextra -Werror -Wno-unused-parameter
all: find capture

find: find.c lib.c
	$(CC) $(CFLAGS) -O2 -o $@ $^ -ljpeg -g -DWIDTH=640 -DHEIGHT=360

capture: capture_raw_frames.c lib.c
	$(CC) $(CFLAGS) -O2 -o $@ $^ -g -DWIDTH=320 -DHEIGHT=240

clean:
	rm -f find
