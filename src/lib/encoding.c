#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

typedef uint8_t byte;

char hex_digits[16] = {
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'
};

char* hex_from_bytes(byte *bytes, size_t num_bytes) {
	size_t num_chars = num_bytes * 2;
	char* string = calloc(num_chars + 1, sizeof(char));
	char c;
	for (int i = 0; i < num_bytes; i++) {
		byte b = bytes[i];
		string[2 * i] = hex_digits[b >> 4];
		string[2 * i + 1] = hex_digits[b & 0xf];
	}
	return string;
}

void bytes_from_hex(char* hex, byte **bytes, size_t *num_bytes) {
	size_t num_chars = strlen(hex);
	*num_bytes = (strlen(hex) + 1) / 2;
	*bytes = calloc(*num_bytes, sizeof(byte));
	for (int i = 0; i < num_chars; ++i) {
		char c = hex[num_chars - i - 1];
		char digit = 0;
		if (c >= '0' && c <= '9') {
			digit = c - '0';
		} else if (c >= 'A' && c <= 'F') {
			digit = 10 + (c - 'A');
		} else if (c >= 'a' && c <= 'f') {
			digit = 10 + (c - 'a');
		}
		if (i % 2 == 0) {
			(*bytes)[*num_bytes - i / 2 - 1] = digit;
		} else {
			(*bytes)[*num_bytes - i / 2 - 1] |= digit << 4;
		}
	}
}

void print_bytes(byte* bytes, size_t num_bytes) {
	printf("Printing %zd bytes:", num_bytes);
	for (int i = 0; i < num_bytes; i++) {
		printf(" %u", bytes[i]);
	}
	printf("\n");
}