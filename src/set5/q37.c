#include <assert.h>
#include <stdio.h>

#include <gmp.h>

#include "crypt.c"
#include "encoding.c"
#include "dh.c"
#include "sha1.h"

typedef struct credential {
    char salt[16];
    mpz_t *v;
} credential;

void fill_cred(char *password, mpz_t *N, mpz_t *g, credential *cred) {
    SHA1Context ctx;
    char hash[20], hashx[41];
    size_t i;

    cred->v = (mpz_t*)calloc(1, sizeof(mpz_t));
    for (i = 0; i < 16; ++i) {
        cred->salt[i] = randn(256) - 128;
    }
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)cred->salt, 16));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)password, strlen(password)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)hash));
    assert(-1 != bytes_to_hex(hash, 20, hashx, 41));
    mpz_init_set_str(*cred->v, hashx, 16);
    mpz_powm(*cred->v, *g, *cred->v, *N);
}

void create_server_key(mpz_t *b, mpz_t *B, mpz_t *N, mpz_t* g, mpz_t *k, credential *cred) {
    mpz_t z;

    mpz_init(*b);
    mpz_init(*B);
    mpz_init(z);
    mpz_urandomm(*b, rand_state, *N);
    mpz_powm(*B, *g, *b, *N);
    mpz_mul(z, *k, *cred->v);
    mpz_add(*B, *B, z);
}

void free_ints(dh_key k) {
    free(k.p);
    free(k.g);
    free(k.public);
    free(k.private);
}

int main() {
    mpz_t N, g, k, b, B, u, z1;
    credential cred;
    dh_key key_a;
    char hash[20], hashx[41], buf[1024], cK[20], sK[20];
    SHA1Context ctx;

    init_dh();
    mpz_init_set_str(N, nist_p, 16);
    mpz_init_set_str(g, nist_g, 16);
    mpz_init_set_ui(k, 3);

    // Save password
    fill_cred("$secret", &N, &g, &cred);

    // Create session keys
    key_a = new_NIST_dh_key();
    create_server_key(&b, &B, &N, &g, &k, &cred);


    // u (known to server and client)
    assert(shaSuccess == SHA1Reset(&ctx));
    mpz_get_str(buf, 16, *key_a.public);
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)buf, strlen(buf)));
    mpz_get_str(buf, 16, B);
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)buf, strlen(buf)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)hash));
    assert(-1 != bytes_to_hex(hash, 20, hashx, 41));
    mpz_init_set_str(u, hashx, 16);

    // client key derivation -- A = 0
    mpz_init_set_ui(z1, 0);
    mpz_set_ui(*key_a.public, 0);
    mpz_get_str(buf, 16, z1);
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)buf, strlen(buf)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)cK));

    // server key derivation
    mpz_powm(z1, *cred.v, u, N);
    mpz_mul(z1, *key_a.public, z1);
    mpz_powm(z1, z1, b, N);
    mpz_get_str(buf, 16, z1);
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)buf, strlen(buf)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)sK));

    printf("Client's key:\n");
    print_bytes(cK, 20);
    printf("Server's key:\n");
    print_bytes(sK, 20);

    // client key derivation -- A = 3 * N
    mpz_init_set_ui(z1, 0);
    mpz_set(*key_a.public, N);
    mpz_mul_ui(*key_a.public, *key_a.public, 3);
    mpz_get_str(buf, 16, z1);
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)buf, strlen(buf)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)cK));

    // server key derivation
    mpz_powm(z1, *cred.v, u, N);
    mpz_mul(z1, *key_a.public, z1);
    mpz_powm(z1, z1, b, N);
    mpz_get_str(buf, 16, z1);
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)buf, strlen(buf)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)sK));

    printf("Client's key:\n");
    print_bytes(cK, 20);
    printf("Server's key:\n");
    print_bytes(sK, 20);
}
