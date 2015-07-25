#include <assert.h>
#include <stdio.h>

#include "crypt.c"
#include "encoding.c"
#include "sha1.h"

// m should have room for 64 bytes
// and len should be less than 56 bytes
void sha1_pad(char *m, size_t len) {
    assert(len < 56);
    m[len] = '\x80';
    memset(m + len + 1, '\0', 64 - len - 1);
    for (size_t i = 0; i < len; ++i) {
        inc_counter_be(m + 60, 4);
        inc_counter_be(m + 60, 4);
        inc_counter_be(m + 60, 4);
        inc_counter_be(m + 60, 4);
        inc_counter_be(m + 60, 4);
        inc_counter_be(m + 60, 4);
        inc_counter_be(m + 60, 4);
        inc_counter_be(m + 60, 4);
    }
}

int main() {
    srand(99);
    char key[16];
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }

    unsigned char mac[SHA1HashSize];
    char *msg = "test";
    sha1_mac(key, sizeof(key), msg, strlen(msg), (char *)mac);
    printf("Original MAC:\n");
    print_bytes((char *)mac, 20);

    SHA1Context ctx;
    assert(shaSuccess == SHA1Reset(&ctx));
    for (int i = 0; i < 5; ++i) {
        ctx.Intermediate_Hash[i] = (mac[4 * i] << 24) | (mac[4 * i + 1] << 16) | (mac[4 * i + 2] << 8) | mac[4 * i + 3];
    }
    ctx.Length_Low = 512;
    //printf("Copied:\n");
    for (int i = 0; i < 5; ++i) {
        //printf("%x\n", ctx.Intermediate_Hash[i]);
    }
    char *extension = calloc(64, 1);
    strcpy(extension, "...plus more!");
    //sha1_pad(extension, strlen(extension));
    //char *extension = "...plus more!";
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)extension, strlen(extension)));
    assert(shaSuccess == SHA1Result(&ctx, mac));

    printf("Forged MAC:\n");
    print_bytes((char *)mac, 20);

    char *s = calloc(64 + strlen(extension) + 1, 1);
    memcpy(s, key, sizeof(key));
    strcat(s + sizeof(key), msg);
    sha1_pad(s, sizeof(key) + strlen(msg));
    strcat(s + 64, extension);
    //print_bytes(s, 64 + strlen(extension));
    sha1_mac("", 0, s, 64 + strlen(extension), (char *)mac);
    printf("MAC of extended message:\n");
    print_bytes((char *)mac, 20);

    free(extension);
    free(s);
}
