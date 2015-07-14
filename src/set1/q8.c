#include <stdio.h>

#include "crypt.c"
#include "encoding.c"
#include "math.c"
#include "score.c"

int main() {
    char* line = NULL;
    size_t n = 0;
    while (getline(&line, &n, stdin) != -1) {
        size_t len = strlen(line);
        if (line[len - 1] == '\n') {
            line[len - 1] = '\0';
            --len;
        }
        size_t max_bytes = num_bytes_from_hex(len);
        char *c = calloc(max_bytes + 1, 1);
        size_t num_bytes = hex_to_bytes(line, c, max_bytes);
        c[num_bytes] = '\0';  // turn into a cstring
        
        // Look for 4-byte strings that are repeated.
        int num_repetitions = 0;
        for (size_t i = 0; i < num_bytes - 4; ++i) {
            for (size_t j = i + 4; j < num_bytes; ++j) {
                if (!memcmp(c + i, c + j, 4)) ++num_repetitions;
            }
        }

        if (num_repetitions > 10) {
            printf("Found %d repetitions in line: %s\n", num_repetitions, line);
        }
    }
}
