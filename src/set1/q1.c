#include <stdio.h>

#include "encoding.c"

int main() {
	char* before = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d";
	//char* before = "ff";
	size_t max_bytes = (strlen(before) + 1)/2;
	byte* as_bytes = calloc(max_bytes, sizeof(byte));
	size_t num_bytes = hex_to_bytes(before, as_bytes, max_bytes);
	char* after = calloc(num_bytes * 2 + 1, sizeof(char));
	bytes_to_hex(as_bytes, num_bytes, after, num_bytes * 2 + 1);

    printf("Before: %s\n", before);
    print_bytes(as_bytes, num_bytes);
    printf("After: %s\n", after);

    free(as_bytes);
    free(after);
}
