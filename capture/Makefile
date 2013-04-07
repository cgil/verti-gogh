CFLAGS += -Wall -Wextra -Werror -Wno-unused-parameter -O2
all: find capture calibrate calibrate_raw

build/%.o: %.c
	@mkdir -p $(@D)
	$(CC) $(CFLAGS) -c -MMD -MP -MF $(@:.o=.d) -MT $@ -o $@ $<

find: build/find.o build/lib.o
	$(CC) -o $@ $^ -ljpeg -g

calibrate: build/calibrate.o build/lib.o
	$(CC) -o $@ $^ -ljpeg -g

calibrate_raw: build/calibrate_raw.o build/lib.o build/raw.o
	$(CC) -o $@ $^

capture: build/capture_raw_frames.o build/lib.o build/raw.o
	$(CC) -o $@ $^

clean:
	rm -rf build find capture

-include build/find.d biuld/lib.d build/capture_raw_frames.d build/calibrate.d \
	 build/raw.d
