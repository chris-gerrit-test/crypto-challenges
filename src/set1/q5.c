#include <stdio.h>

#include "encoding.c"
#include "math.c"

int main() {
	char msg[] = "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal";
	char *key = "ICE";
	repeated_xor((byte*)msg, strlen(msg), (byte*)key, strlen(key), (byte*)msg);
	size_t hex_digits = num_hex_digits(strlen(msg));
	char* buf = calloc(hex_digits + 1, sizeof(char));
	bytes_to_hex((byte*)msg, strlen(msg), buf, hex_digits + 1);
    printf("%s\n", buf);
    free(buf);
}
