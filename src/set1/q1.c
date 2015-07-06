#include <stdio.h>

#include "encoding.c"

int main() {
	char* before = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d";
	//char* before = "ff";
	byte* as_bytes;
	size_t num_bytes;
	bytes_from_hex(before, &as_bytes, &num_bytes);
	char* after = hex_from_bytes(as_bytes, num_bytes);

    printf("Before: %s\n", before);
    print_bytes(as_bytes, num_bytes);
    printf("After: %s\n", after);
}
