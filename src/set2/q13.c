#include <stddef.h>
#include <stdio.h>

#include "crypt.c"
#include "encoding.c"

typedef struct entry {
    char *key;
    char *val;
} entry;

typedef struct profile {
    entry *entries;
    size_t num_entries;
} profile;

char key[16];

char *encode_profile(profile *pr, size_t *n) {
    *n = 16;
    char *out = calloc(*n, 1);
    char *p = out;
    for (size_t i = 0; i < pr->num_entries; ++i) {
        entry entry = pr->entries[i];
        for (size_t k = 0; k < strlen(entry.key); ++k) {
            char c = entry.key[k];
            if (c == '&' || c == '=') continue;
            if (p >= out + *n - 2) {
                ptrdiff_t offs = p - out;
                *n *= 2;
                out = realloc(out, *n);
                p = out + offs;
            }
            *p++ = c;
        }
        *p++ = '=';
        for (size_t v = 0; v < strlen(entry.val); ++v) {
            char c = entry.val[v];
            if (c == '&' || c == '=') continue;
            if (p >= out + *n - 2) {
                ptrdiff_t offs = p - out;
                *n *= 2;
                out = realloc(out, *n);
                p = out + offs;
            }
            *p++ = c;
        }
        if (i < pr->num_entries - 1) {
            *p++ = '&';
        }
    }
    *p = '\0';
    //printf("profile:%s\n", out);
    aes_encrypt(out, out, *n, key);
    return out;
}

char *profile_for(char *email, size_t *n) {
    profile pr = {
        (entry[]){{"email", email}, {"uid", "10"}, {"role", "user"}},
        3
    };
    return encode_profile(&pr, n);
}

int main() {
    // TODO: go back and make a better solution with padding
    srand(19);
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }

    size_t n1 = 0;
    char *s1 = profile_for(
        "0123456789"        // end current block (starts with email=)
        "admin", &n1);

    size_t n2 = 0;
    char *s2 = profile_for("0123456789012", &n2);

    char *attack = calloc(48, 1);
    memcpy(attack, s2, 32);
    memcpy(attack + 32, s1 + 16, 16);

    aes_decrypt(attack, attack, 48, key);
    printf("%s\n", attack);

    free(s1);
    free(s2);
    free(attack);
}
