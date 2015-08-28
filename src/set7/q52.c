#include <assert.h>
#include <stdio.h>

#include "crypt.c"

#define MAX_HASH_SIZE 4 // in bytes

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

void aes_hash_nopad(char *msg, char hash[MAX_HASH_SIZE], size_t hash_size) {
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

// s1 and s2 should have room for at least 17 bytes
void find_collision(char state[MAX_HASH_SIZE], char *s1, char *s2, size_t hash_size) {
    char h1[MAX_HASH_SIZE], h2[MAX_HASH_SIZE];
    bucket buckets[131072];
    char *arena = NULL;
    size_t arena_size = 0;
    int search_bytes = hash_size;
    long long unsigned n;
    bucket *b;

    for (size_t i = 0; i < sizeof(buckets) / sizeof(bucket); ++i) {
        buckets[i].num_entries = 0;
        buckets[i].entries = NULL;
    }

    memset(s1, 0, 17);
    memset(s1, 1, search_bytes);
    assert(-1 != pkcs7(s1, 17, 16));
    memset(s2, 0, 17);
    memset(s2, 1, search_bytes);
    assert(-1 != pkcs7(s2, 17, 16));
    inc_counter_be(s2, search_bytes);

    while (1) {
        memcpy(h2, state, MAX_HASH_SIZE);
        aes_hash_nopad(s2, h2, hash_size);
        n = be_counter_val(h2, hash_size) % (sizeof(buckets) / sizeof(bucket));
        b = &buckets[n];
        for (size_t k = 0; k < b->num_entries; ++k) {
            memcpy(h1, state, MAX_HASH_SIZE);
            aes_hash_nopad(arena + b->entries[k], h1, hash_size);
            if (!memcmp(h1, h2, MAX_HASH_SIZE)) {
                memcpy(s1, arena + b->entries[k], 17);
                return;
            }
        }
        b->entries = realloc(b->entries, (b->num_entries + 1) * sizeof(size_t));
        assert(b->entries != NULL);
        b->entries[b->num_entries] = arena_size;
        b->num_entries += 1;
        arena = realloc(arena, arena_size + 17);
        assert(arena != NULL);
        memcpy(arena + arena_size, s2, 17);
        arena_size += 17;
        do {
            inc_counter_be(s2, search_bytes);
        } while (memchr(s2, 0, search_bytes));
    }
}

// s1 and s2 should have size at least 16 * n + 1
void find_many_collisions(int n, char *s1, char *s2, size_t hash_size) {
    int i;
    char state[MAX_HASH_SIZE];

    memset(state, 0, MAX_HASH_SIZE);

    for (i = 0; i < n; ++i) {
        find_collision(state, s1 + 16 * i, s2 + 16 * i, hash_size);
        memset(state, 0, MAX_HASH_SIZE);
        aes_hash_nopad(s1, state, hash_size);
    }

}

int main() {
    size_t small_hash_size = 2;
    size_t big_hash_size = 4;
    size_t n = big_hash_size;
    char s1[16 * n + 1], s2[sizeof(s1)], s[sizeof(s1)];
    char h1[MAX_HASH_SIZE], h2[MAX_HASH_SIZE];
    unsigned i, k;

    find_many_collisions(n, s1, s2, small_hash_size);

    aes_hash(s1, h1, small_hash_size);
    aes_hash(s2, h2, small_hash_size);
    assert(!memcmp(h1, h2, small_hash_size));
    // Make a string out of pieces of s1 and s2
    for (i = 0; i < (unsigned)2 << (4 * n); ++i) {
        memcpy(s, s1, sizeof(s1));
        for (k = 0; k < n; ++k) {
            if (i & (1 << k)) {
                memcpy(s + 16 * k, s2 + 16 * k, 16);
            }
        }
        aes_hash(s, h2, small_hash_size);
        assert(!memcmp(h1, h2, small_hash_size));
    }
}
