#include <stdio.h>
#include <time.h>
#include <unistd.h>

#include "crypt.c"
#include "encoding.c"
#include "math.c"
#include "mersenne.c"

char *get_token(char *username, time_t expiration, uint32_t seed) {
    size_t n = strlen(username) + sizeof(time_t) * 2 + 2;
    char *tok = calloc(n, 1);
    sprintf(tok, "%x|%s", (unsigned)expiration, username);
    mersenne_crypt(tok, tok, n, seed);
    return tok;
}

int is_time_based(char *tok, size_t len, char *known_substring) {
    uint32_t now = time(NULL) * 1000;
    size_t k = strlen(known_substring);
    char buf[len];
    for (int i = 0; i < 5000; ++i) {
        mersenne_crypt(tok, buf, len, now - i);
        for (char *c = buf; c + k <= buf + len; ++c) {
            if (!memcmp(c, known_substring, k)) {
                return 1;
            }
        }
    }
    return 0;
}

int main() {
    srand(309);
    uint16_t seed = randn(0xffff);

    int prefix_len = randn(500) + 500;
    int n = prefix_len + 15;
    char *message = calloc(n, 1);
    for (int i = 0; i < prefix_len; ++i) {
        message[i] = randn(256) - 128;
    }
    memset(message + prefix_len, 'a', 14);
    mersenne_crypt(message, message, n, seed);
    //print_bytes(message, n);

    char *decrypted = calloc(n, 1);
    for (int i = 0; i < 0xffff; ++i) {
        mersenne_crypt(message, decrypted, n, i);
        if (!memcmp(decrypted - 15 + n, "aaaaaaaaaaaaaa", 14)) {
            printf("Guess: %d   Seed: %d\n", i, seed);
        }
    }
    free(decrypted);
    free(message);

    char *email = "chris@charlesbearllc.com";
    char *tok = get_token(email, time(NULL) + 60 * 60, time(NULL) * 1000 + randn(1000));
    sleep(randn(3) + 1);
    printf("actually time-based; detected: %d\n", is_time_based(tok, n, email));
    free(tok);
    tok = get_token(email, time(NULL) + 60 * 60, randn(0xffff));
    sleep(randn(3) + 1);
    printf("not actually time-based; detected: %d\n", is_time_based(tok, n, email));
    free(tok);
}
