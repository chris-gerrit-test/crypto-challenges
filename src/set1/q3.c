#include <stdio.h>

#include "encoding.c"
#include "math.c"
#include "score.c"

int main() {
	char encrypted[] = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736";
	size_t max_bytes = num_bytes_from_hex(strlen(encrypted));
	byte* c = calloc(max_bytes + 1, sizeof(byte));
	size_t num_bytes = hex_to_bytes(encrypted, c, max_bytes);
	c[num_bytes] = '\0';  // turn into a cstring

	double best_score = -1.0;
	byte best_key = 0;
    for (int i = 0; i < 256; i++) {
    	byte k = i;
    	repeated_xor(c, num_bytes, &k, 1, c);
    	double score = score_text((char*)c, num_bytes);
    	if (score > best_score) {
    		best_key = k;
    		best_score = score;
    	}
    	// Undo the "decryption"
    	repeated_xor(c, num_bytes, &k, 1, c);
    }

    // Decrypt.
	repeated_xor(c, num_bytes, &best_key, 1, c);
	printf("%s\n", c);

    free(c);
}
