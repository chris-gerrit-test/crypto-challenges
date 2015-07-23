#include <stdio.h>

#include "crypt.c"
#include "encoding.c"

char key[16];
char nonce[16];

char *encode_user_data(char *user_data, size_t *n) {
    char *prefix = "comment1=cooking%20MCs;userdata=";
    char *suffix = ";comment2=%20like%20a%20pound%20of%20bacon";
    *n = strlen(user_data) + strlen(prefix) + strlen(suffix) + 100;
    char *encoded = calloc(*n, 1);
    strcat(encoded, prefix);
    char *q = encoded + strlen(encoded);
    for (char *p = user_data; *p; ++p) {
        // Get rid of characters we don't like.
        if (*p != '=' && *p != ';') {
            *q++ = *p;
        }
    }
    strcat(q, suffix);
    *n = pkcs7(encoded, *n, 16);
    ctr_crypt(encoded, encoded, *n, nonce, key);
    return encoded;
}

int main() {
    srand(90);

    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }
    for (size_t i = 0; i < sizeof(nonce); ++i) {
        nonce[i] = randn(256) - 128;
    }

    size_t n;
    char *e = encode_user_data("", &n);
    xor(e, "comment1=cooking", 16, e);
    xor(e, ";admin=true;", 12, e);

    ctr_crypt(e, e, 12, nonce, key);
    printf("%.12s\n", e);

    free(e);
}
