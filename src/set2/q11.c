#include <stdio.h>

#include "crypt.c"
#include "encoding.c"

int pad_string(char *in, char *out) {
    size_t prepend_count = randn(6) + 5;
    size_t append_count = randn(6) + 5;
    char *p = out;
    for (size_t i = 0; i < prepend_count; ++i) {
        *p++ = randn(256) - 128;
    }
    memcpy(p, in, strlen(in));
    p += strlen(in);
    for (size_t i = 0; i < append_count; ++i) {
        *p++ = randn(256) - 128;
    }
    return strlen(in) + prepend_count + append_count;
}

char key[16];
char iv[16];

int cbc_mode;

void encryption_oracle(char *s) {
    if (cbc_mode) {
        printf("CBC mode\n");
        cbc_encrypt(s, s, strlen(s), iv, key);
    } else {
        printf("ECB mode\n");
        aes_encrypt(s, s, strlen(s), key);
    }
}

int main() {
    srand(9);

    cbc_mode = randn(2);
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }
    for (size_t i = 0; i < sizeof(iv); ++i) {
        iv[i] = randn(256) - 128;
    }

    char *in = "kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk";
    size_t n = strlen(in) + 20;
    if (n % 32 != 0) {
        n = n + 32 - (n % 32);
    }
    char *s = calloc(n, 1);
    pad_string(in, s);
    encryption_oracle(s);

    if (!memcmp(s + 16, s + 32, 16)) {
        printf("Detected: ECB mode\n");
    } else {
        printf("Detected: CBC mode\n");
    }

    free(s);
}
