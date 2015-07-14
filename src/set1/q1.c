#include <stdio.h>

#include "encoding.c"

int main() {
	char before[] = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d";
	size_t max_bytes = num_bytes_from_hex(strlen(before));
	char* as_bytes = calloc(max_bytes, 1);
	size_t num_bytes = hex_to_bytes(before, as_bytes, max_bytes);
	char* as_base64 = calloc(num_base64_digits(num_bytes) + 1, 1);
	bytes_to_base64(as_bytes, num_bytes, as_base64, num_base64_digits(num_bytes) + 1);

    printf("Hex: %s\n", before);
    printf("Base64: %s\n", as_base64);

    base64_to_bytes(as_base64, as_bytes, num_bytes);
    bytes_to_hex(as_bytes, num_bytes, before, strlen(before) + 1);
    printf("Back to hex: %s\n", before);

    free(as_bytes);
    free(as_base64);
}
