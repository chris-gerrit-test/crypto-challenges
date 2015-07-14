CC=gcc

INCLUDE=src/lib/

LIBCRYPTO=/usr/lib/libcrypto.dylib

# libcrypto is deprecated on Mac OS X
CFLAGS=-W -Wall -Wno-deprecated-declarations -std=c99

q1:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q1.c -o bin/q1
	bin/q1

q2:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q2.c -o bin/q2
	bin/q2

q3:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q3.c -o bin/q3
	bin/q3

q4:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q4.c -o bin/q4
	bin/q4 < data/set1/4.txt

q5:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q5.c -o bin/q5
	bin/q5

q6:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q6.c -o bin/q6
	bin/q6 < data/set1/6.txt

q7:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set1/q7.c -o bin/q7
	bin/q7 < data/set1/7.txt

set1: q1 q2 q3 q4 q5 q6 q7

all: set1