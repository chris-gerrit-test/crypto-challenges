#include <stdio.h>

#include "crypt.c"
#include "encoding.c"

char key[16];

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
    cbc_encrypt(encoded, encoded, *n, key, key);
    return encoded;
}

typedef struct error {
    char msg[1024];
} error;

int verify(char *data, size_t len, error *err) {
    cbc_decrypt(data, data, len, key, key);
    int valid = 1;
    for (size_t i = 0; i < len; ++i) {
        if (data[i] < 0) {
            valid = 0;
            err->msg[0] = '\0';
            strcpy(err->msg, "Invalid message: ");
            memcpy(err->msg + strlen(err->msg), data, len);  // overrun
            break;
        }
    }
    // Undo changes.
    cbc_encrypt(data, data, len, key, key);
    return valid;
}

int main() {
    srand(9);

    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }
    printf("Original key:\n");
    print_bytes(key, 16);

    size_t n;
    char *e = encode_user_data("", &n);
    error err;
    memset(e + 16, '\0', 16);
    memcpy(e + 32, e, 16);
    verify(e, n, &err);
    char *p = err.msg + strlen("Invalid message: ");
    xor(p, p + 32, 16, p);
    printf("Reversed key:\n");
    print_bytes(p, 16);

    free(e);
}
