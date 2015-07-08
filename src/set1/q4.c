#include <stdio.h>

#include "encoding.c"
#include "math.c"
#include "score.c"

int main() {
    double best_score = -1.0;
    char* best_line = NULL;
    char* line = NULL;
    size_t n = 0;
    while (getline(&line, &n, stdin) != -1) {
        size_t len = strlen(line);
        if (line[len - 1] == '\n') {
            line[len - 1] = '\0';
            --len;
        }
        size_t max_bytes = num_bytes_from_hex(len);
        byte *c = calloc(max_bytes + 1, sizeof(byte));
        size_t num_bytes = hex_to_bytes(line, c, max_bytes);
        c[num_bytes] = '\0';  // turn into a cstring
        int found_best_key = 0;
        byte best_key = -1;
        for (int i = 0; i < 256; i++) {
            byte k = i;
            repeated_xor(c, num_bytes, &k, 1, c);
            double score = score_text((char*)c, num_bytes);
            if (score > best_score) {
                found_best_key = 1;
                best_key = k;
                best_score = score;
            }
            // Undo the "decryption"
            repeated_xor(c, num_bytes, &k, 1, c);
        }
        if (found_best_key) {
            repeated_xor(c, num_bytes, &best_key, 1, c);
            free(best_line);
            best_line = (char*)c;
        } else {
            free(c);
        }
        line = NULL;
        n = 0;
    }

	printf("%s\n", best_line);
    free(best_line);
}
