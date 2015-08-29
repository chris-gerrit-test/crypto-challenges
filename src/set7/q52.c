#include <assert.h>
#include <stdio.h>

#include "crypt.c"

#define MAX_HASH_SIZE 6 // in bytes

// hash should already be inited
void aes_hash(char *msg, char hash[MAX_HASH_SIZE], size_t hash_size) {
    char buf[8192], key[16];
    int padded_size, offset;
    AES_KEY aes_key;

    assert(hash_size <= 16);
    memset(hash, 0, MAX_HASH_SIZE);

    memcpy(key, hash, hash_size);
    assert(0 == AES_set_encrypt_key((unsigned char*)key, 128, &aes_key));

    strcpy(buf, msg);
    padded_size = pkcs7(buf, sizeof(buf), 16);
    assert(-1 != padded_size);

    for (offset = 0; offset < padded_size; offset += 16) {
        memset(key + hash_size, 0, sizeof(key) - hash_size);
        //print_bytes(key, 16);
        assert(0 == AES_set_encrypt_key((unsigned char*)key, 128, &aes_key));
        AES_encrypt((unsigned char*)(buf + offset), (unsigned char*)key, &aes_key);
    }
    memcpy(hash, key, hash_size);
}

size_t calls[MAX_HASH_SIZE];

void aes_hash_nopad(char *msg, char hash[MAX_HASH_SIZE], size_t hash_size) {
    calls[hash_size - 1] += 1;
    char key[16];
    int padded_size, offset;
    AES_KEY aes_key;

    assert(hash_size <= 16);

    memcpy(key, hash, hash_size);
    assert(0 == AES_set_encrypt_key((unsigned char*)key, 128, &aes_key));

    padded_size = strlen(msg);
    assert(padded_size % 16 == 0);

    for (offset = 0; offset < padded_size; offset += 16) {
        memset(key + hash_size, 0, sizeof(key) - hash_size);
        assert(0 == AES_set_encrypt_key((unsigned char*)key, 128, &aes_key));
        AES_encrypt((unsigned char*)(msg + offset), (unsigned char*)key, &aes_key);
    }
    memcpy(hash, key, hash_size);
}

void print_hash(char *str, size_t hash_size) {
    char hash[MAX_HASH_SIZE];

    printf("Text:\n");
    print_bytes(str, strlen(str));
    aes_hash(str, hash, hash_size);
    printf("Hash:\n");
    print_bytes(hash, hash_size);
}

typedef struct bucket {
    size_t num_entries, *entries;
} bucket;

char counter[MAX_HASH_SIZE + 1];

void init_counter_input() {
    memset(counter, 1, sizeof(counter));
    counter[sizeof(counter) - 1] = '\0';
}

void next_counter_input() {
    inc_counter_be(counter, MAX_HASH_SIZE);
}

char *perm1, *perm2;
char *perm_input = NULL;
size_t perm_n, perm_counter;

void init_permutation_input(char *s1, char *s2, size_t n) {
    perm1 = s1;
    perm2 = s2;
    perm_n = n;
    perm_counter = 0;
    free(perm_input);
    perm_input = calloc(strlen(s1) + 1, 1);
    strcpy(perm_input, s1);
}

void next_permutation_input() {
    long long unsigned k;

    ++perm_counter;
    assert(perm_counter < (1 << perm_n));  // don't roll over
    strcpy(perm_input, perm1);
    for (k = 0; k < perm_n; ++k) {
        if (perm_counter & (1 << k)) {
            memcpy(perm_input + 16 * k, perm2 + 16 * k, 16);
        }
    }
}

// s1 and s2 should have room for at least input_size + 1 bytes
// input_size should be a multiple of 16
void find_collision(
    char state[MAX_HASH_SIZE],
    char *s1, char *s2,
    size_t input_size,
    size_t hash_size,
    char *input, void (*next_input)()) {
    char h1[MAX_HASH_SIZE], h2[MAX_HASH_SIZE];
    bucket buckets[131072];
    char *arena = NULL;
    size_t arena_size = 0;
    long long unsigned n;
    bucket *b;

    for (size_t i = 0; i < sizeof(buckets) / sizeof(bucket); ++i) {
        buckets[i].num_entries = 0;
        buckets[i].entries = NULL;
    }

    memset(s1, 0, input_size + 1);
    strcpy(s1, input);
    assert(-1 != pkcs7(s1, input_size + 1, 16));
    memset(s2, 0, input_size + 1);
    next_input();
    strcpy(s2, input);
    assert(-1 != pkcs7(s2, input_size + 1, 16));
    //print_bytes(s1, input_size);
    //print_bytes(s2, input_size);

    while (1) {
        memcpy(h2, state, MAX_HASH_SIZE);
        aes_hash_nopad(s2, h2, hash_size);
        n = be_counter_val(h2, hash_size) % (sizeof(buckets) / sizeof(bucket));
        //printf("%llu\n", n);
        b = &buckets[n];
        for (size_t k = 0; k < b->num_entries; ++k) {
            if (!memcmp(s2, arena + b->entries[k], input_size)) {
                // Already saw this input
                continue;
            }
            memcpy(h1, state, MAX_HASH_SIZE);
            aes_hash_nopad(arena + b->entries[k], h1, hash_size);
            if (!memcmp(h1, h2, MAX_HASH_SIZE)) {
                // Found a match!
                memcpy(s1, arena + b->entries[k], input_size + 1);
                return;
            }
        }
        b->entries = realloc(b->entries, (b->num_entries + 1) * sizeof(size_t));
        assert(b->entries != NULL);
        b->entries[b->num_entries] = arena_size;
        b->num_entries += 1;
        arena = realloc(arena, arena_size + input_size + 1);
        assert(arena != NULL);
        memcpy(arena + arena_size, s2, input_size + 1);
        arena_size += input_size + 1;
        do {
            next_input();
            strcpy(s2, input);
            assert(-1 != pkcs7(s2, input_size + 1, 16));
        } while (memchr(s2, 0, input_size));
    }
}

// s1 and s2 should have size at least 16 * n + 1
void find_many_collisions(int n, char *s1, char *s2, size_t hash_size) {
    int i;
    char state[MAX_HASH_SIZE];

    memset(state, 0, MAX_HASH_SIZE);
    init_counter_input();

    for (i = 0; i < n; ++i) {
        find_collision(state, s1 + 16 * i, s2 + 16 * i, 16, hash_size, counter, next_counter_input);
        memset(state, 0, MAX_HASH_SIZE);
        aes_hash_nopad(s1, state, hash_size);
    }

}

int main() {
    size_t small_hash_size = 4;
    size_t big_hash_size = 5;
    size_t n = big_hash_size * 4;
    char s1[16 * (n + 1) + 1], s2[sizeof(s1)], c1[sizeof(s1)], c2[sizeof(s1)];
    char h1[MAX_HASH_SIZE], h2[MAX_HASH_SIZE];

    find_many_collisions(n, s1, s2, small_hash_size);
    // print_bytes(s1, sizeof(s1));
    // print_bytes(s2, sizeof(s2));
    s1[sizeof(s1) - 1] = '\0';
    s2[sizeof(s2) - 1] = '\0';

    aes_hash(s1, h1, small_hash_size);
    aes_hash(s2, h2, small_hash_size);
    assert(memcmp(s1, s2, sizeof(s1)));
    assert(!memcmp(h1, h2, small_hash_size));
    printf("Small hash:\n");
    print_bytes(h1, small_hash_size);

    init_permutation_input(s1, s2, n);
    memset(h1, 0, MAX_HASH_SIZE);
    find_collision(h1, c1, c2, sizeof(c1) - 1, big_hash_size, perm_input, next_permutation_input);

    print_bytes(c1, sizeof(c1));
    print_bytes(c2, sizeof(c2));
    assert(memcmp(c1, c2, sizeof(c1)));
    aes_hash(c1, h1, small_hash_size);
    aes_hash(c2, h2, small_hash_size);
    assert(!memcmp(h1, h2, MAX_HASH_SIZE));
    printf("Small hash:\n");
    print_bytes(h1, small_hash_size);
    aes_hash(c1, h1, big_hash_size);
    aes_hash(c2, h2, big_hash_size);
    printf("Big hash:\n");
    print_bytes(h1, big_hash_size);
    assert(!memcmp(h1, h2, MAX_HASH_SIZE));

    printf("Small hash calls: %zu\n", calls[small_hash_size - 1]);
    printf("Big hash calls:   %zu\n", calls[big_hash_size - 1]);
}
