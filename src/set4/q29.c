#include <assert.h>
#include <stdio.h>

#include "crypt.c"
#include "encoding.c"
#include "sha1.h"

// m should have room for the padding (up to 73 additional bytes I think)
// retuns the number of bytes after padding
size_t sha1_pad(char *m, size_t len) {
    size_t s = len;
    m[s++] = '\x80';
    while (s % 64 != 56) m[s++] = '\0';
    for (size_t i = 0; i < len * 8; ++i) {
        inc_counter_be(m + s, 8);
    }
    return s + 8;
}

int main() {
    srand(99);
    char key[16];
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }

    unsigned char mac[SHA1HashSize];
    char *msg = "comment1=cooking%20MCs;userdata=foo;comment2=%20like%20a%20pound%20of%20bacon";
    sha1_mac(key, sizeof(key), msg, strlen(msg), (char *)mac);
    printf("Original MAC:\n");
    print_bytes((char *)mac, 20);

    SHA1Context ctx;
    assert(shaSuccess == SHA1Reset(&ctx));
    for (int i = 0; i < 5; ++i) {
        ctx.Intermediate_Hash[i] = (mac[4 * i] << 24) | (mac[4 * i + 1] << 16) | (mac[4 * i + 2] << 8) | mac[4 * i + 3];
    }
    ctx.Length_Low = 1024;  // depends on size of message
    char *extension = calloc(64, 1);
    strcpy(extension, ";admin=true");
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)extension, strlen(extension)));
    assert(shaSuccess == SHA1Result(&ctx, mac));

    printf("Forged MAC:\n");
    print_bytes((char *)mac, 20);

    printf("MAC of extended message:\n");
    char *s = calloc(64 * 3, 1);
    memcpy(s, key, sizeof(key));
    strcat(s + sizeof(key), msg);
    size_t n = sha1_pad(s, sizeof(key) + strlen(msg));
    strcat(s + n, extension);
    sha1_mac("", 0, s, n + strlen(extension), (char *)mac);
    print_bytes((char *)mac, 20);

    free(extension);
    free(s);
}
