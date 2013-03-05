CFLAGS += -Wall -Wextra -Werror -Wno-unused-parameter -O2
all: find capture

build/%.o: %.c
	@mkdir -p $(@D)
	$(CC) $(CFLAGS) -c -MMD -MP -MF $(@:.o=.d) -MT $@ -o $@ $<

find: build/find.o build/lib.o
	$(CC) -o $@ $^ -ljpeg -g

capture: build/capture_raw_frames.o build/lib.o
	$(CC) -o $@ $^

clean:
	rm -rf build find capture

-include build/find.d biuld/lib.d build/capture_raw_frames.d
