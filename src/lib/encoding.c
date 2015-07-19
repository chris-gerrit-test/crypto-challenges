#ifndef ENCODING_C
#define ENCODING_C

#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

char hex_digits[16] = {
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'
};

/*
	The byte's two high bits will be ignored.
*/
char byte_to_base64_digit(char b) {
	b = b & 0x3f;
	if (b < 26) return 'A' + b;
	if (b < 52) return 'a' + (b - 26);
	if (b < 62) return '0' + (b - 52);
	return b == 62 ? '=' : '/';
}

char base64_digit_to_byte(char digit) {
	if (digit >= 'A' && digit <= 'Z') {
		return digit - 'A';
	}
	if (digit >= 'a' && digit <= 'z') {
		return digit - 'a' + 26;
	}
	if (digit >= '0' && digit <= '9') {
		return digit - '0' + 52;
	}
	if (digit == '+') return 62;
	return 63;
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

int num_bytes_from_base64(size_t num_base64_digits) {
	return ((num_base64_digits + 3) / 4) * 3;
}

int bytes_to_base64(char *bytes, size_t num_bytes, char *buf, size_t buf_size) {
	size_t num_chars = num_base64_digits(num_bytes);  // 4 base64 chars fit in 3 bytes
	if (num_chars + 1 > buf_size) {
		return -1;
	}
	buf[num_chars] = '\0';
	for (size_t i = 0; i < num_bytes; i+= 3) {
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

int bytes_to_hex(char *bytes, size_t num_bytes, char *buf, size_t buf_size) {
	size_t num_chars = num_hex_digits(num_bytes);
	if (num_chars + 1 > buf_size) {
		return -1;
	}
	buf[num_chars] = '\0';
	for (size_t i = 0; i < num_bytes; i++) {
		char b = bytes[i];
		buf[2 * i] = hex_digits[b >> 4];
		buf[2 * i + 1] = hex_digits[b & 0xf];
	}
	return num_chars;
}

int base64_to_bytes(char* base64, char *buf, size_t buf_size) {
	size_t num_chars = strlen(base64);
	size_t max_num_bytes = num_bytes_from_base64(num_chars);
	if (max_num_bytes > buf_size) {
		return -1;
	}
	size_t num_bytes = 0;
	for (size_t i = 0; i < num_chars; i+= 4) {
		uint32_t group;  // the 24-bit group of bits we are handling now
		// Read in from input
		group = base64_digit_to_byte(base64[i]) << 18;
		if (i + 1 < num_chars) {
			group |= base64_digit_to_byte(base64[i + 1]) << 12;
		}
		if (i + 2 < num_chars && base64[i + 2] != '=') {
			group |= base64_digit_to_byte(base64[i + 2]) << 6;
		}
		if (i + 3 < num_chars && base64[i + 3] != '=') {
			group |= base64_digit_to_byte(base64[i + 3]);
		}
		// Write out to output
		if (i + 3 < num_chars && base64[i + 3] != '=') {
			++num_bytes;
			buf[i / 4 * 3 + 2] = group & 0xff;
		}
		group = group >> 8;
		if (i + 2 < num_chars && base64[i + 2] != '=') {
			++num_bytes;
			buf[i / 4 * 3 + 1] = group & 0xff;
		}
		group = group >> 8;
		buf[i / 4 * 3] = group & 0xff;
		++num_bytes;
	}
	return num_bytes;
}

int hex_to_bytes(char* hex, char *buf, size_t buf_size) {
	size_t num_chars = strlen(hex);
	size_t num_bytes = num_bytes_from_hex(num_chars);
	if (num_bytes > buf_size) {
		return -1;
	}
	for (size_t i = 0; i < num_chars; ++i) {
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

void print_bytes(char* bytes, size_t num_bytes) {
	printf("Printing %zd bytes:\n", num_bytes);
	for (size_t i = 0; i < num_bytes; i++) {
		printf(" %3d", (unsigned char)bytes[i]);
		if (i % 16 == 15) printf("\n");
	}
	printf("\n");
}

#endif /* ENCODING_C */