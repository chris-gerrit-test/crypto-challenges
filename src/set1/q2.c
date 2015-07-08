#include <stdio.h>

#include "encoding.c"
#include "math.c"

int main() {
	char in1s[] = "1c0111001f010100061a024b53535009181c";
	char *in2s = "686974207468652062756c6c277320657965";
	size_t max_bytes = num_bytes_from_hex(strlen(in1s));
	byte* in1 = calloc(max_bytes, sizeof(byte));
	byte* in2 = calloc(max_bytes, sizeof(byte));
	size_t num_bytes = hex_to_bytes(in1s, in1, max_bytes);
	hex_to_bytes(in2s, in2, max_bytes);

    printf(" %s\n^%s\n", in1s, in2s);

	repeated_xor(in1, num_bytes, in2, num_bytes, in1);
	bytes_to_hex(in1, num_bytes, in1s, strlen(in1s) + 1);

    printf("=%s\n", in1s);

    free(in1);
    free(in2);
}
