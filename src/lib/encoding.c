#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

typedef uint8_t byte;

char hex_digits[16] = {
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'
};


int bytes_to_hex(byte *bytes, size_t num_bytes, char *buf, size_t buf_size) {
	size_t num_chars = num_bytes * 2;
	if (num_chars + 1 > buf_size) {
		return -1;
	}
	buf[num_chars] = '\0';
	for (int i = 0; i < num_bytes; i++) {
		byte b = bytes[i];
		buf[2 * i] = hex_digits[b >> 4];
		buf[2 * i + 1] = hex_digits[b & 0xf];
	}
	return num_chars;
}

int hex_to_bytes(char* hex, byte *buf, size_t buf_size) {
	size_t num_chars = strlen(hex);
	size_t num_bytes = (strlen(hex) + 1) / 2;
	if (num_bytes > buf_size) {
		return -1;
	}
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
			buf[num_bytes - i / 2 - 1] = digit;
		} else {
			buf[num_bytes - i / 2 - 1] |= digit << 4;
		}
	}
	return num_bytes;
}

void print_bytes(byte* bytes, size_t num_bytes) {
	printf("Printing %zd bytes:", num_bytes);
	for (int i = 0; i < num_bytes; i++) {
		printf(" %u", bytes[i]);
	}
	printf("\n");
}