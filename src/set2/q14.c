#include <stdio.h>

#include "crypt.c"
#include "encoding.c"

char *suffix = "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK";
//char *suffix = "Um9sbG";
char *suffix_bytes;
size_t num_suffix_bytes;
size_t max_prefix_bytes = 100;

size_t pad_string(char *in, char *out) {
    size_t n = strlen(in);
    char *tmp = calloc(n, 1);
    memcpy(tmp, in, n);
    size_t prefix_len = randn(max_prefix_bytes);
    //prefix_len = 16;
    for (size_t i = 0; i < prefix_len; ++i) {
        out[i] = randn(256) - 128;
    }
    memcpy(out + prefix_len, tmp, n);
    memcpy(out + prefix_len + n, suffix_bytes, num_suffix_bytes);
    free(tmp);
    return n + prefix_len + num_suffix_bytes;
}

char key[16];

// s should have room for num_suffix_bytes extra, rounded up to 16 bytes.
size_t encryption_oracle(char *s) {
    //print_bytes(s, strlen(s));
    size_t n = pad_string(s, s);
    //print_bytes(s, n);
    aes_encrypt(s, s, n, key);
    return n;
}

size_t detect_block_size() {
    size_t n = num_suffix_bytes + max_prefix_bytes;
    if (n % 16 != 0) {
        n = n + 16 - (n % 16);
    }
    char *s = calloc(n, 1);
    size_t smallest_size = 65536;
    size_t second_smallest_size = 65536;
    for (int t = 0; t < 10000; ++t) {
        memset(s, 0, n);
        size_t size = encryption_oracle(s);
        if (size <= smallest_size) smallest_size = size;
        else if (size < second_smallest_size) second_smallest_size = size;
    }
    return second_smallest_size - smallest_size;
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
        size_t n = i + prefix_len + num_suffix_bytes + max_prefix_bytes;
        if (n % bs != 0) {
            n = n + bs - (n % bs);
        }
        char *s = calloc(n, 1);
        char *original_e = calloc(n, 1);
        char *original_d = calloc(n, 1);
        memset(s, 'k', prefix_len);
        memcpy(original_d, s, n);
        encryption_oracle(s);
        memcpy(original_e, s, n);
        // Now find the character at position i
        size_t max_matching_blocks = 0;
        for (int c = -127; c < 128; ++c) {
            if (c == 0) continue;  // necessary because encryption_oracle uses strlen
            while (1) {
                // Keep generating two ciphertexts until they are both
                // the right length to line things up.
                memcpy(s, original_d, n);
                size_t size1 = encryption_oracle(s) - num_suffix_bytes + i;
                //printf("%zu\n", size1);
                if (size1 % 16 != bs - 1) continue;
                //printf("%lu %lu\n", size1 - num_suffix_bytes - 16, (size1 - num_suffix_bytes + i) % bs);
                memcpy(original_e, s, n);
                memset(s, 0, n);
                memset(s, 'k', prefix_len);
                memcpy(s + prefix_len, unknown_string, i);
                s[i + prefix_len] = c;
                size_t size2 = encryption_oracle(s) - num_suffix_bytes + i;
                if (size1 + i + 1 != size2) {
                    //printf("%zu %zu %zu\n", i, size1 + i + 1, size2);
                    continue;
                }
                // 1 out of 1600 times, both will have the same amount of padding
                // and it will line us up so that this test works.
                size_t matching_blocks = 0;
                for (size_t k = 0; k <= size1; k += bs) {
                    if (!memcmp(s + k, original_e + k, bs)) {
                        //printf("size1=%zu size2=%zu\n", size1, size2);
                        ++matching_blocks;
                    }
                }
                if (matching_blocks > max_matching_blocks) {
                    //printf("%zu: %d: %zu blocks\n", i, c, matching_blocks);
                    max_matching_blocks = matching_blocks;
                    unknown_string[i] = c;
                }
                break;
                //printf("%d: %d / %d\n", c, matching_blocks, skip_blocks);
                // if (matching_blocks == skip_blocks + 17) {
                //     unknown_string[i] = c;
                //     goto found_byte;
                // }
            }
        }
        printf("Found byte %c\n", unknown_string[i]);
        free(s);
        free(original_e);
        free(original_d);
    }

    printf("%s\n", unknown_string);
    free(unknown_string);
}

int main() {
    srand(18);
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }

    num_suffix_bytes = num_bytes_from_base64(strlen(suffix));
    suffix_bytes = calloc(num_suffix_bytes, 1);
    num_suffix_bytes = base64_to_bytes(suffix, suffix_bytes, num_suffix_bytes);

    size_t bs = detect_block_size();
    bs = 16;
    printf("Detected block size %zu\n", bs);

    size_t n = bs * 100 + num_suffix_bytes + max_prefix_bytes;
    if (n % 16 != 0) {
        n = n + 16 - (n % 16);
    }
    char *s = calloc(n, 1);
    memset(s, 'k', bs * 100);
    encryption_oracle(s);

    if (!memcmp(s + 90 * bs, s + 91 * bs, bs)) {
        printf("Detected: ECB mode\n");
    } else {
        printf("Detected: CBC mode\n");
    }

    free(s);

    find_unknown_string(bs);
}
