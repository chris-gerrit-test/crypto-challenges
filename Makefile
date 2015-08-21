CC=gcc

INCLUDE=src/lib/

LIBCRYPTO=-lcrypto

LIBGMP=-lgmp

LIB_SRC=src/lib/sha1.c src/lib/md4.c

# libcrypto is deprecated on Mac OS X
CFLAGS=-W -Wall -Wno-deprecated-declarations -std=c99 -g

q1:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q1.c $(LIB_SRC) -o bin/q1
	bin/q1

q2:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q2.c $(LIB_SRC) -o bin/q2
	bin/q2

q3:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q3.c $(LIB_SRC) -o bin/q3
	bin/q3

q4:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q4.c $(LIB_SRC) -o bin/q4
	bin/q4 < data/set1/4.txt

q5:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q5.c $(LIB_SRC) -o bin/q5
	bin/q5

q6:
	$(CC) $(CFLAGS) -I $(INCLUDE) src/set1/q6.c $(LIB_SRC) -o bin/q6
	bin/q6 < data/set1/6.txt

q7:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set1/q7.c $(LIB_SRC) -o bin/q7
	bin/q7 < data/set1/7.txt

q8:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set1/q8.c $(LIB_SRC) -o bin/q8
	bin/q8 < data/set1/8.txt

q9:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q9.c $(LIB_SRC) -o bin/q9
	bin/q9

q10:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q10.c $(LIB_SRC) -o bin/q10
	bin/q10 < data/set2/10.txt

q11:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q11.c $(LIB_SRC) -o bin/q11
	bin/q11

q12:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q12.c $(LIB_SRC) -o bin/q12
	bin/q12

q13:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q13.c $(LIB_SRC) -o bin/q13
	bin/q13

q14:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q14.c $(LIB_SRC) -o bin/q14
	bin/q14

q15:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q15.c $(LIB_SRC) -o bin/q15
	bin/q15

q16:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set2/q16.c $(LIB_SRC) -o bin/q16
	bin/q16

q17:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q17.c $(LIB_SRC) -o bin/q17
	bin/q17

q18:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q18.c $(LIB_SRC) -o bin/q18
	bin/q18

q19:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q19.c $(LIB_SRC) -o bin/q19
	bin/q19

q20:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q20.c $(LIB_SRC) -o bin/q20
	bin/q20

q21:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q21.c $(LIB_SRC) -o bin/q21
	bin/q21

q22:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q22.c $(LIB_SRC) -o bin/q22
	bin/q22

q23:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q23.c $(LIB_SRC) -o bin/q23
	bin/q23

q24:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set3/q24.c $(LIB_SRC) -o bin/q24
	bin/q24

q25:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set4/q25.c $(LIB_SRC) -o bin/q25
	bin/q25 < data/set1/7.txt

q26:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set4/q26.c $(LIB_SRC) -o bin/q26
	bin/q26

q27:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set4/q27.c $(LIB_SRC) -o bin/q27
	bin/q27

q28:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set4/q28.c $(LIB_SRC) -o bin/q28
	bin/q28

q29:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set4/q29.c $(LIB_SRC) -o bin/q29
	bin/q29

q30:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set4/q30.c $(LIB_SRC) -o bin/q30
	bin/q30

q31:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set4/q31s.c $(LIB_SRC) -o bin/q31s
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBCRYPTO) src/set4/q31c.c $(LIB_SRC) -o bin/q31c
	@"bin/q31s" & PID=$$!; sleep 0.25; bin/q31c; kill "$$PID";

q32: q31

q33:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBGMP) src/set5/q33.c $(LIB_SRC) -o bin/q33
	bin/q33

q34:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBGMP) $(LIBCRYPTO) src/set5/q34.c $(LIB_SRC) -o bin/q34
	bin/q34

q35:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBGMP) $(LIBCRYPTO) src/set5/q35.c $(LIB_SRC) -o bin/q35
	bin/q35

q36:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBGMP) $(LIBCRYPTO) src/set5/q36.c $(LIB_SRC) -o bin/q36
	bin/q36

q37:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBGMP) $(LIBCRYPTO) src/set5/q37.c $(LIB_SRC) -o bin/q37
	bin/q37

q38:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBGMP) $(LIBCRYPTO) src/set5/q38.c $(LIB_SRC) -o bin/q38
	bin/q38

q39:
	$(CC) $(CFLAGS) -I $(INCLUDE) $(LIBGMP) $(LIBCRYPTO) src/set5/q39.c $(LIB_SRC) -o bin/q39
	bin/q39

set1: q1 q2 q3 q4 q5 q6 q7 q8

set2: q9 q10 q11 q12 q13 q14 q15 q16

set3: q17 q18 q19 q20 q21 q22 q23 q24

set4: q25 q26 q27 q28 q29 q30 q31 q32

set5: q33 q34 q35 q36 q37 q38 q39

all: set1 set2 set3 set4 set5