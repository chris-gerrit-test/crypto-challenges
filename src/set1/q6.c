#include <stdio.h>

#include "encoding.c"
#include "math.c"
#include "score.c"

int main() {
    char* doc = calloc(1, 1);
    size_t doc_len = 0;
    char* line = NULL;
    size_t n = 0;
    while (getline(&line, &n, stdin) != -1) {
    	size_t line_len = strlen(line);
    	doc = realloc(doc, doc_len + line_len);
    	strncpy(doc + doc_len, line, line_len - 1);
    	doc_len += line_len - 1;
    	doc[doc_len] = '\0';
    }
    
    size_t num_bytes = num_bytes_from_base64(doc_len);
    char* as_bytes = calloc(num_bytes + 1, 1);
    num_bytes = base64_to_bytes(doc, as_bytes, num_bytes);
    as_bytes[num_bytes] = '\0';

    double best_hamming = 1e10;
    int best_keysize = -1;
    // Start at 10 because my score gets false positives at 5 and 9.
    for (int k = 10; k <= 40; ++k) {
    	double hamming = 1.0 * (
    		hamming_distance_n(as_bytes, as_bytes + k, k)
    		+ hamming_distance_n(as_bytes, as_bytes + 2 * k, k)
    		+ hamming_distance_n(as_bytes + 3 * k, as_bytes + 4 * k, k)
    		+ hamming_distance_n(as_bytes + 3 * k, as_bytes + 5 * k, k)) / k;
    	if (hamming < best_hamming) {
    		best_hamming = hamming;
    		best_keysize = k;
    	}
    }

    printf("Best keysize: %d: (score %.4f)\n", best_keysize, best_hamming);

    size_t col_size = num_bytes / best_keysize;
    char* col = calloc(col_size, 1);
    char* key = calloc(best_keysize, 1);
    for (int i = 0; i < best_keysize; ++i) {
    	for (size_t j = 0; j < col_size; ++j) {
    		col[j] = as_bytes[j * best_keysize + i];
    	}
		double best_score = -1.0;
		char best_key = 0;
	    for (int h = 0; h < 256; h++) {
	    	char k = h;
	    	repeated_xor(col, col_size, &k, 1, col);
	    	double score = score_text(col, col_size);
	    	if (score > best_score) {
	    		best_key = k;
	    		best_score = score;
	    	}
	    	// Undo the "decryption"
	    	repeated_xor(col, col_size, &k, 1, col);
	    }
	    key[i] = best_key;
    }

    printf("Key: %s\n", key);

    repeated_xor(as_bytes, num_bytes, key, best_keysize, as_bytes);
    printf("Decoded: %s\n", as_bytes);

    free(as_bytes);
    free(col);
    free(key);
}
