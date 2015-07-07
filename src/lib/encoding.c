#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

typedef uint8_t byte;

char hex_digits[16] = {
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'
};

/*
	The byte's two high bits will be ignored.
*/
char byte_to_base64_digit(byte b) {
	b = b & 0x3f;
	if (b < 26) return 'A' + b;
	if (b < 52) return 'a' + (b - 26);
	if (b < 62) return '0' + (b - 52);
	return b == 62 ? '=' : '/';
}

int num_base64_digits(size_t num_bytes) {
	return ((num_bytes + 2) / 3) * 4;
}

int num_hex_digits(size_t num_bytes) {
	return num_bytes * 2;
}

int num_bytes_from_hex(size_t num_hex_digits) {
	return (num_hex_digits + 1) / 2;
}

int bytes_to_base64(byte *bytes, size_t num_bytes, char *buf, size_t buf_size) {
	size_t num_chars = num_base64_digits(num_bytes);  // 4 base64 chars fit in 3 bytes
	if (num_chars + 1 > buf_size) {
		return -1;
	}
	buf[num_chars] = '\0';
	for (int i = 0; i < num_bytes; i+= 3) {
		uint32_t group;  // the 24-bit group of bits we are handling now
		// Read in from input
		group = bytes[i] << 16;
		if (i + 1 < num_bytes) {
			group |= (bytes[i + 1] << 8);
		}
		if (i + 2 < num_bytes) {
			group |= bytes[i + 2];
		}
		// Write out to output
		if (i + 2 < num_bytes) {
			buf[i / 3 * 4 + 3] = byte_to_base64_digit(group & 0x3f);
		} else { // missing one or more bytes
			buf[i / 3 * 4 + 3] = '=';
		}
		group = group >> 6;
		if (i + 1 < num_bytes) {
			buf[i / 3 * 4 + 2] = byte_to_base64_digit(group & 0x3f);
		} else {  // missing two bytes
			buf[i / 3 * 4 + 2] = '=';
		}
		group = group >> 6;
		buf[i / 3 * 4 + 1] = byte_to_base64_digit(group & 0x3f);
		group = group >> 6;
		buf[i / 3 * 4] = byte_to_base64_digit(group & 0x3f);
	}
	return num_chars;
}

int bytes_to_hex(byte *bytes, size_t num_bytes, char *buf, size_t buf_size) {
	size_t num_chars = num_hex_digits(num_bytes);
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
	size_t num_bytes = num_bytes_from_hex(num_chars);
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