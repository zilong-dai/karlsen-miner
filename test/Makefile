CC=gcc
CFLAGS=-I../target/release/
LDFLAGS=-L../target/release/ -lkls
OBJS=test.o
TARGET=test

all: $(OBJS)
	$(CC) -o $(TARGET) $(OBJS) $(LDFLAGS)

$(OBJS): test.c
	$(CC) -c test.c $(CFLAGS)

run: all
	./test

clean:
	rm -f $(OBJS) $(TARGET)
