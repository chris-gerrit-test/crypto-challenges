CC=gcc

INCLUDE=src/lib/

LIBCRYPTO=/usr/lib/libcrypto.dylib

# libcrypto is deprecated on Mac OS X
CFLAGS=-W -Wall -Wno-deprecated-declarations -std=c99 -g

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

q8:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set1/q8.c -o bin/q8
	bin/q8 < data/set1/8.txt

q9:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q9.c -o bin/q9
	bin/q9

q10:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q10.c -o bin/q10
	bin/q10 < data/set2/10.txt

q11:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q11.c -o bin/q11
	bin/q11

q12:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q12.c -o bin/q12
	bin/q12

q13:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q13.c -o bin/q13
	bin/q13

q14:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q14.c -o bin/q14
	bin/q14

q15:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q15.c -o bin/q15
	bin/q15

q16:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q16.c -o bin/q16
	bin/q16

q17:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q17.c -o bin/q17
	bin/q17

q18:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q18.c -o bin/q18
	bin/q18

q19:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q19.c -o bin/q19
	bin/q19

q20:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q20.c -o bin/q20
	bin/q20

q21:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q21.c -o bin/q21
	bin/q21

q22:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q22.c -o bin/q22
	bin/q22

set1: q1 q2 q3 q4 q5 q6 q7 q8

set2: q9 q10 q11 q12 q13 q14 q15 q16

set3: q17 q18 q19 q20 q21 q22

all: set1 set2 set3