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

void create_server_key(mpz_t *b, mpz_t *B, mpz_t *u, mpz_t *N, mpz_t* g) {
    char us[16], usx[33];
    size_t i;

    mpz_init(*b);
    mpz_init(*B);
    mpz_init(*u);
    mpz_urandomm(*b, rand_state, *N);
    mpz_powm(*B, *g, *b, *N);
    for (i = 0; i < 16; ++i) {
        us[i] = randn(256) - 128;
    }
    assert(-1 != bytes_to_hex(us, 16, usx, 33));
    mpz_init_set_str(*u, usx, 16);
}

void free_ints(dh_key k) {
    free(k.p);
    free(k.g);
    free(k.public);
    free(k.private);
}

int main() {
    mpz_t N, g, b, B, u, x, z1;
    credential cred;
    dh_key key_a;
    char hash[20], hashx[41], buf[1024], cK[20], sK[20];
    SHA1Context ctx;
    char *line;
    int n;
    size_t s;
    FILE *f;
    char* password = "earthenhearted";

    init_dh();
    mpz_init_set_str(N, nist_p, 16);
    mpz_init_set_str(g, nist_g, 16);

    // Save password
    fill_cred(password, &N, &g, &cred);

    // Create session keys
    key_a = new_NIST_dh_key();
    create_server_key(&b, &B, &u, &N, &g);

    printf("Verify protocol:\n");

    // x (know to client only)
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)cred.salt, 16));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)password, strlen(password)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)hash));
    assert(-1 != bytes_to_hex(hash, 20, hashx, 41));
    mpz_init_set_str(x, hashx, 16);

    // client key derivation
    mpz_init(z1);
    mpz_mul(z1, x, u);
    mpz_add(z1, *key_a.private, z1);
    mpz_powm(z1, B, z1, N);
    mpz_get_str(buf, 16, z1);
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)buf, strlen(buf)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)cK));
    sha1_mac(cK, 20, cred.salt, 16, cK);

    // server key derivation
    mpz_powm(z1, *cred.v, u, N);
    mpz_mul(z1, *key_a.public, z1);
    mpz_powm(z1, z1, b, N);
    mpz_get_str(buf, 16, z1);
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)buf, strlen(buf)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)sK));
    sha1_mac(sK, 20, cred.salt, 16, sK);

    printf("Client's HMAC:\n");
    print_bytes(cK, 20);
    printf("Server's HMAC:\n");
    print_bytes(sK, 20);

    printf("\nGuess password with no salt, u=1, b=1:\n");
    mpz_set_ui(u, 1);
    mpz_set_ui(b, 1);
    mpz_set(B, g);

    // client key derivation
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)password, strlen(password)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)hash));
    assert(-1 != bytes_to_hex(hash, 20, hashx, 41));
    mpz_init_set_str(x, hashx, 16);
    mpz_init(z1);
    mpz_mul(z1, x, u);
    mpz_add(z1, *key_a.private, z1);
    mpz_powm(z1, B, z1, N);
    mpz_get_str(buf, 16, z1);
    assert(shaSuccess == SHA1Reset(&ctx));
    assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)buf, strlen(buf)));
    assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)cK));
    sha1_mac(cK, 20, cred.salt, 16, cK);
    printf("Client's HMAC:\n");
    print_bytes(cK, 20);

    assert(f = fopen("/usr/share/dict/words", "r"));
    // Takes 35s on my machine
    while ((n = getline(&line, &s, f)) != -1) {
        line[n - 1] = '\0';
        assert(shaSuccess == SHA1Reset(&ctx));
        assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)line, strlen(line)));
        assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)hash));
        assert(-1 != bytes_to_hex(hash, 20, hashx, 41));
        mpz_init_set_str(z1, hashx, 16);
        mpz_powm(z1, g, z1, N);
        mpz_powm(z1, z1, u, N);
        mpz_mul(z1, *key_a.public, z1);
        mpz_powm(z1, z1, b, N);
        mpz_get_str(buf, 16, z1);
        assert(shaSuccess == SHA1Reset(&ctx));
        assert(shaSuccess == SHA1Input(&ctx, (uint8_t *)buf, strlen(buf)));
        assert(shaSuccess == SHA1Result(&ctx, (uint8_t *)sK));
        sha1_mac(sK, 20, cred.salt, 16, sK);
        if (!memcmp(cK, sK, 20)) {
            printf("Found password: %s\n", line);
            return 0;
        }
    }
    printf("Did not find password\n");
    return 1;
}
