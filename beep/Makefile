CC ?= cc
# miniaudio.h is a huge single-header library: -O0 keeps builds fast on
# Raspberry Pi / Termux while audio playback quality is unaffected.
# Override with `make CFLAGS=-O2` if you want it optimized.
CFLAGS ?= -O0

UNAME_M := $(shell uname -m)

# 64-bit atomics are built in on aarch64; on 32-bit ARM (e.g. old
# Termux installs) they need libatomic.
ifeq ($(UNAME_M),armv7l)
    LIBS += -latomic
endif

all: beep

beep: beep.c miniaudio.h
	$(CC) $(CFLAGS) beep.c -o beep -lm $(LIBS)

clean:
	rm -f beep

.PHONY: all clean
