#include <assert.h>
#include <stdio.h>

#include "crypt.c"
#include "encoding.c"

char *inputs[] = {
    "MDAwMDAwTm93IHRoYXQgdGhlIHBhcnR5IGlzIGp1bXBpbmc=",
    "MDAwMDAxV2l0aCB0aGUgYmFzcyBraWNrZWQgaW4gYW5kIHRoZSBWZWdhJ3MgYXJlIHB1bXBpbic=",
    "MDAwMDAyUXVpY2sgdG8gdGhlIHBvaW50LCB0byB0aGUgcG9pbnQsIG5vIGZha2luZw==",
    "MDAwMDAzQ29va2luZyBNQydzIGxpa2UgYSBwb3VuZCBvZiBiYWNvbg==",
    "MDAwMDA0QnVybmluZyAnZW0sIGlmIHlvdSBhaW4ndCBxdWljayBhbmQgbmltYmxl",
    "MDAwMDA1SSBnbyBjcmF6eSB3aGVuIEkgaGVhciBhIGN5bWJhbA==",
    "MDAwMDA2QW5kIGEgaGlnaCBoYXQgd2l0aCBhIHNvdXBlZCB1cCB0ZW1wbw==",
    "MDAwMDA3SSdtIG9uIGEgcm9sbCwgaXQncyB0aW1lIHRvIGdvIHNvbG8=",
    "MDAwMDA4b2xsaW4nIGluIG15IGZpdmUgcG9pbnQgb2g=",
    "MDAwMDA5aXRoIG15IHJhZy10b3AgZG93biBzbyBteSBoYWlyIGNhbiBibG93"
};

char key[16];
char iv[16];

void init_keys() {
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }
}

// enc should fit at least 128 bytes; iv should fit 16.
// returns the size of the ciphertext
int get_ciphertext(char *enc, char *iv) {
    for (size_t i = 0; i < sizeof(iv); ++i) {
        iv[i] = randn(256) - 128;
    }
    int n = base64_to_bytes(inputs[randn(10)], enc, 128);
    assert(n > 0);
    // Assumes the plaintexts have no zero bytes.
    n = pkcs7(enc, 128, 16);
    assert(n > 0);
    cbc_encrypt(enc, enc, n, iv, key);
    return n;
}

int has_valid_padding(char *enc, size_t n, char *iv) {
    cbc_decrypt(enc, enc, n, iv, key);
    int ret = NULL != strip_pkcs7(enc, n, 16);

    // re-encrypt so buffer comes back the same
    cbc_encrypt(enc, enc, n, iv, key);

    return ret;
}

void decrypt_block(char *enc, size_t block_pos, char *iv, char *dec) {
    assert(block_pos % 16 == 0);
    char *block_start = enc + block_pos;
    assert(block_start >= enc);

    char *prev_block = calloc(16, 1);
    char *prev_block_start;
    if (block_pos == 0) {
        prev_block_start = iv;
    } else {
        prev_block_start = block_start - 16;
    }

    dec = dec + block_pos;

    for (int i = 0; i < 16; ++i) {
        char padding_byte = i + 1;
        memcpy(prev_block, prev_block_start, 16);
        for (int j = 0; j < i; ++j) {
            prev_block[16 - j - 1] = prev_block[16 - j - 1] ^ dec[16 - j - 1] ^ padding_byte;
        }
        for (int c = -128; c < 128; ++c) {
            prev_block[16 - i - 1] = c;
            int is_valid = has_valid_padding(block_start, 16, prev_block);
            if (i == 0 && is_valid) {
                // Make sure we didn't get lucky on the first byte and get
                // \2\2: change the second byte and try again.
                prev_block[16 - 2] += 1;
                is_valid = has_valid_padding(block_start, 16, prev_block);
            }
            if (is_valid) {
                dec[16 - i - 1] = (prev_block_start[16 - i - 1]) ^ c ^ padding_byte;
                //printf("Found byte %d: %d\n", i, dec[16 - i - 1]);
                break;
            }
        }
        assert(dec[16 - i - 1] != '\0');
    }

    free(prev_block);
}

int main() {
    srand(99);
    init_keys();

    char *enc = calloc(128, 1);
    char *iv = calloc(16, 1);
    int n = get_ciphertext(enc, iv);
    // Sanity checks.
    assert(n % 16 == 0);
    assert(has_valid_padding(enc, n, iv));
    // Validate assumption: can feed previous block as IV
    assert(has_valid_padding(enc + n - 16, 16, enc + n - 32));

    char *dec = calloc(n, 1);
    for (int i = 0; i < n; i += 16) {
        decrypt_block(enc, i, iv, dec);
    }

    char *end = strip_pkcs7(dec, n, 16);
    assert(end != 0);
    *end = '\0';
    printf("%s\n", dec);

    free(enc);
    free(iv);
    free(dec);
}
