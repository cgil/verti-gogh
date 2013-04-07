CFLAGS += -Wall -Wextra -Werror -Wno-unused-parameter -O2
PROGS = find calibrate calibrate_raw capture_raw_frames ping
LIBS = raw lib

all: $(PROGS)

build/%.o: %.c
	@mkdir -p $(@D)
	$(CC) $(CFLAGS) -c -MMD -MP -MF $(@:.o=.d) -MT $@ -o $@ $<

$(PROGS): %: build/%.o $(LIBS:%=build/%.o)
	$(CC) $(CFLAGS) -o $@ $^ -ljpeg -g

clean:
	rm -rf build $(PROGS)

-include $(PROGS:%=build/%.d) %(LIBS:%=build/%.d)
