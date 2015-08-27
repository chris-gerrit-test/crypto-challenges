#include <assert.h>
#include <limits.h>
#include <stdio.h>
#include <string.h>
#include <zlib.h>

#include "math.c"

char key[16];
char nonce[16];

char *template = "POST / HTTP/1.1\n"
                 "Host: hapless.com\n"
                 "Cookie: sessionid=TmV2ZXIgcmV2ZWFsIHRoZSBXdS1UYW5nIFNlY3JldCE=\n"
                 "Content-Length: %zu\n"
                 "%s";

size_t oracle_ctr(char *msg) {
    char formatted[8192], compressed[8192];
    size_t formatted_len, compressed_len = sizeof(compressed);

    sprintf(formatted, template, strlen(msg), msg);
    formatted_len = strlen(formatted);
    assert(Z_OK == compress((Bytef *)compressed, &compressed_len, (Bytef *)formatted, formatted_len));
    return compressed_len;  // Encrypted length will be the same.
}

size_t oracle_cbc(char *msg) {
    char formatted[8192], compressed[8192];
    size_t formatted_len, compressed_len = sizeof(compressed);

    sprintf(formatted, template, strlen(msg), msg);
    formatted_len = strlen(formatted);
    assert(Z_OK == compress((Bytef *)compressed, &compressed_len, (Bytef *)formatted, formatted_len));
    return (1 + compressed_len / 16) * 16;
}

void show(char *msg) {
    printf("%5zu %s\n", oracle_ctr(msg), msg);
}

void guess_id(size_t (*oracle)(char *)) {
    char *id_start, best_char, *id_prefix="sessionid=", buf[8192], id[100];
    size_t prefix_len = strlen(id_prefix), len, best_len, pos;

    id_start = strstr(template, id_prefix);
    assert(NULL != id_start);
    id_start += prefix_len;

    memset(id, 0, sizeof(id));
    for (pos = 0; pos < sizeof(id); ++pos) {
        best_len = INT_MAX;
        for (int i = 1; i < 256; ++i) {
            memset(buf, 0, sizeof(buf));
            char guess = i;
            for (int k = 0; k < 8; ++k) {
                memcpy(buf + strlen(buf), id_start - prefix_len - k, prefix_len + k);
                strcat(buf, id);
                buf[strlen(buf)] = guess;
            }
            len = oracle(buf);
            if (len < best_len) {
                best_len = len;
                best_char = guess;
            }
        }
        if (best_char == '\n') break;
        id[pos] = best_char;
    }

    printf("Guessed:  %s\n", id);
    assert(!strcmp(id, "TmV2ZXIgcmV2ZWFsIHRoZSBXdS1UYW5nIFNlY3JldCE="));
}

int main() {
    srand(1019);
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }
    for (size_t i = 0; i < sizeof(nonce); ++i) {
        nonce[i] = randn(256) - 128;
    }

    printf("CTR:\n");
    guess_id(oracle_ctr);
    printf("\nCBC:\n");
    guess_id(oracle_cbc);
}
