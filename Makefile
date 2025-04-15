BIN = exec


# LIBS += -lm

SRC += main.go events.go queue.go

all: build

build:
	go build -o ${BIN} ${SRC}

run:
	./${BIN}


# PHONY all run
