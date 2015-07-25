#include <assert.h>
#include <stdio.h>

#include "crypt.c"
#include "encoding.c"
#include "md4.h"

// m should have room for the padding (up to 73 additional bytes I think)
// retuns the number of bytes after padding
size_t sha1_pad(char *m, size_t len) {
    size_t s = len;
    m[s++] = '\x80';
    while (s % 64 != 56) m[s++] = '\0';
    for (size_t i = 0; i < len * 8; ++i) {
        inc_counter_le(m + s, 8);
    }
    return s + 8;
}

int main() {
    srand(99);
    char key[16];
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }

    unsigned char mac[MD4_DIGEST_LENGTH];
    char *msg = "comment1=cooking%20MCs;userdata=foo;comment2=%20like%20a%20pound%20of%20bacon";
    md4_mac(key, sizeof(key), msg, strlen(msg), (char *)mac);
    printf("Original MAC:\n");
    print_bytes((char *)mac, MD4_DIGEST_LENGTH);

    MD4_CTX ctx;
    MD4Init(&ctx);
    memcpy(ctx.state, mac, MD4_DIGEST_LENGTH);
    ctx.count[0] = 1024;  // depends on size of message
    char *extension = calloc(64, 1);
    strcpy(extension, ";admin=true");
    MD4Update(&ctx, (uint8_t *)extension, strlen(extension));
    MD4Final(mac, &ctx);

    printf("Forged MAC:\n");
    print_bytes((char *)mac, MD4_DIGEST_LENGTH);

    printf("MAC of extended message:\n");
    char *s = calloc(64 * 3, 1);
    memcpy(s, key, sizeof(key));
    strcat(s + sizeof(key), msg);
    size_t n = sha1_pad(s, sizeof(key) + strlen(msg));
    strcat(s + n, extension);
    md4_mac("", 0, s, n + strlen(extension), (char *)mac);
    print_bytes((char *)mac, MD4_DIGEST_LENGTH);

    free(extension);
    free(s);
}
