#include <assert.h>
#include <stdio.h>

#include <gmp.h>

#include "crypt.c"
#include "dh.c"
#include "sha1.h"

dh_key A_get_key() {
    return new_NIST_dh_key();
}

dh_key B_get_key(mpz_t *p, mpz_t *g) {
    mpz_t* new_p = (mpz_t*)calloc(1, sizeof(mpz_t));
    mpz_t* new_g = (mpz_t*)calloc(1, sizeof(mpz_t));
    mpz_t* a = (mpz_t*)calloc(1, sizeof(mpz_t));
    mpz_t* A = (mpz_t*)calloc(1, sizeof(mpz_t));

    init_dh();
    mpz_init_set(*new_p, *p);
    mpz_init_set(*new_g, *g);
    mpz_init(*a);
    mpz_urandomm(*a, rand_state, *p);
    mpz_init(*A);
    mpz_powm(*A, *g, *a, *p);

    return (dh_key) {.p = new_p, .g = new_g, .private = a, .public = A};
}

void A_send_msg(char *msg, dh_key key_a, mpz_t *B, char *out, char iv[16]) {
    char key[16];
    mpz_t s;
    char *s_str;
    SHA1Context ctx;

    // Random IV
    for (size_t i = 0; i < 16; ++i) {
        iv[i] = randn(256) - 128;
    }

    // Get the session key and hash it to make the encryption key
    dh_session_key(key_a, &s, B);
    s_str = mpz_get_str(NULL, 16, s);
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)s_str, strlen(s_str)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)out));
    memcpy(key, out, 16);

    // Encrypt
    cbc_encrypt(msg, out, strlen(msg), iv, key);

    free(s_str);
}

void B_respond(char *msg, size_t msg_size, dh_key key_b, mpz_t *A, char *out, char iv[16]) {
    char key[16];
    mpz_t s;
    char *s_str;
    SHA1Context ctx;

    // Get the session key and hash it to make the encryption key
    dh_session_key(key_b, &s, A);
    s_str = mpz_get_str(NULL, 16, s);
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)s_str, strlen(s_str)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)out));
    memcpy(key, out, 16);

    // Decrypt
    cbc_decrypt(msg, out, msg_size, iv, key);

    printf("  B received: %s\n", out);

    // Re-encrypt
    memcpy(out, msg, msg_size);

    free(s_str);
}

void M_decrypt(char *msg, size_t msg_size, char *out, char iv[16]) {
    char key[16];
    mpz_t s;
    char *s_str;
    SHA1Context ctx;

    // Get the session key and hash it to make the encryption key
    mpz_init_set_ui(s, 0);
    s_str = mpz_get_str(NULL, 16, s);
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)s_str, strlen(s_str)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)out));
    memcpy(key, out, 16);

    // Decrypt
    cbc_decrypt(msg, out, msg_size, iv, key);

    printf("  M decrypted: %s\n", out);

    free(s_str);
}

void free_ints(dh_key k) {
    free(k.p);
    free(k.g);
    free(k.public);
    free(k.private);
}

gmp_randstate_t rand_state;

int main() {
    char *msg, *buf, *buf2, iv[16];
    dh_key key_a;
    dh_key key_b;

    key_a = A_get_key();
    key_b = B_get_key(key_a.p, key_a.g);

    msg = "Hi, there!";

    printf("Checking basic protocol...\n");
    buf = calloc(1024, 1), buf2 = calloc(1024, 1);
    A_send_msg(msg, key_a, key_b.public, buf, iv);
    B_respond(buf, strlen(msg), key_b, key_a.public, buf2, iv);
    assert(!memcmp(buf, buf2, strlen(msg)));

    printf("Checking with interception...\n");
    A_send_msg(msg, key_a, key_b.p, buf, iv);
    M_decrypt(buf, strlen(msg), buf2, iv);
    B_respond(buf, strlen(msg), key_b, key_a.p, buf2, iv);
    assert(!memcmp(buf, buf2, strlen(msg)));

    free_ints(key_a);
    free_ints(key_b);
}
