CFLAGS += -Wall -Wextra -Werror -Wno-unused-parameter -O2 -g
PROGS = find calibrate calibrate_raw capture_raw_frames ping yuyv find_yuyv \
	convert find_raw
LIBS = raw lib colorspace

all: $(PROGS)

build/%.o: %.c
	@mkdir -p $(@D)
	$(CC) $(CFLAGS) -c -MMD -MP -MF $(@:.o=.d) -MT $@ -o $@ $<

$(PROGS): %: build/%.o $(LIBS:%=build/%.o)
	$(CC) $(CFLAGS) -o $@ $^ -ljpeg -g -lm

clean:
	rm -rf build $(PROGS)

-include $(PROGS:%=build/%.d) %(LIBS:%=build/%.d)
