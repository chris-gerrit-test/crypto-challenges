#include <stdio.h>

#include "crypt.c"
#include "encoding.c"

char *suffix = "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK";
//char *suffix = "Um9sbG";
char *suffix_bytes;
size_t num_suffix_bytes;

size_t pad_string(char *in, char *out) {
    size_t n = strlen(in);
    memcpy(out, in, n);
    memcpy(out + n, suffix_bytes, num_suffix_bytes);
    return n + num_suffix_bytes;
}

char key[16];

// s should have room for num_suffix_bytes extra, rounded up to 16 bytes.
size_t encryption_oracle(char *s) {
    size_t n = pad_string(s, s);
    aes_encrypt(s, s, n, key);
    if (n % 16 != 0) {
        n = n + 16 - (n % 16);
    }
    return n;
}

size_t detect_block_size() {
    size_t n = num_suffix_bytes;
    if (n % 16 != 0) {
        n = n + 16 - (n % 16);
    }
    char *s = calloc(n, 1);
    size_t min_size = encryption_oracle(s);
    free(s);
    for (size_t i = 1; ; ++i) {
        n = i + num_suffix_bytes;
        if (n % 16 != 0) {
            n = n + 16 - (n % 16);
        }
        char *s = calloc(n, 1);
        for (size_t j = 0; j < i; ++j) {
            s[j] = 'k';
        }
        size_t encrypted_size = encryption_oracle(s);
        if (encrypted_size > min_size) {
            free(s);
            return encrypted_size - min_size;
        }
        free(s);
    }
}

void find_unknown_string(size_t bs) {
    // Could pretend we don't know the number of suffix bytes and find it
    // at the same time as we find the block size, but let's skip that.
    char *unknown_string = calloc(num_suffix_bytes + 1, 1);

    for (size_t i = 0; i < num_suffix_bytes; ++i) {
        // Need: (prefix_len + i) % bs == bs - 1
        // so prefix_len = (bs - i - 1) % bs
        // e.g. if i is 3 (4th byte) and bs is 16 then we need 12
        // bytes of kkkk prefix, plus the 3 bytes we already know.
        size_t prefix_len = (bs - i - 1) % bs;
        size_t skip_blocks = (prefix_len + i) / 16;
        size_t n = i + prefix_len + num_suffix_bytes;
        if (n % bs != 0) {
            n = n + bs - (n % bs);
        }
        char *s = calloc(n, 1);
        char *original_e = calloc(n, 1);
        memset(s, 'k', prefix_len);
        encryption_oracle(s);
        memcpy(original_e, s, n);
        // Now find the character at position i
        for (int c = -128; c < 128; ++c) {
            if (c == 0) continue;  // necessary because encryption_oracle uses strlen
            memset(s, 0, n);
            memset(s, 'k', prefix_len);
            memcpy(s + prefix_len, unknown_string, i);
            s[i + prefix_len] = c;
            encryption_oracle(s);
            if (!memcmp(s + skip_blocks * bs, original_e + skip_blocks * bs, bs)) {
                unknown_string[i] = c;
                break;
            }
        }
        free(s);
        free(original_e);
    }

    printf("%s\n", unknown_string);
    free(unknown_string);
}

int main() {
    srand(19);
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }

    num_suffix_bytes = num_bytes_from_base64(strlen(suffix));
    suffix_bytes = calloc(num_suffix_bytes, 1);
    num_suffix_bytes = base64_to_bytes(suffix, suffix_bytes, num_suffix_bytes);

    size_t bs = detect_block_size();
    printf("Detected block size %zu\n", bs);

    size_t n = bs * 4 + num_suffix_bytes;
    if (n % 16 != 0) {
        n = n + 16 - (n % 16);
    }
    char *s = calloc(n, 1);
    for (size_t i = 0; i < bs * 4; ++i) {
        s[i] = 'k';
    }
    encryption_oracle(s);

    if (!memcmp(s + bs, s + 2 * bs, bs)) {
        printf("Detected: ECB mode\n");
    } else {
        printf("Detected: CBC mode\n");
    }

    free(s);

    find_unknown_string(bs);
}
