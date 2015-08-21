#include <assert.h>
#include <stdio.h>

#include <gmp.h>

#include "crypt.c"
#include "encoding.c"
#include "rsa.c"
#include "sha1.h"

int main() {
    char msg[512], buf[1025];
    rsa_key k1, k2, k3;
    mpz_t m, c1, c2, c3, m1, m2, m3, N;
    int n;

    mpz_init(m); mpz_init(c1); mpz_init(c2); mpz_init(c3);
    mpz_init(m1); mpz_init(m2); mpz_init(m3); mpz_init(N);

    k1 = new_rsa_key(1024);
    k2 = new_rsa_key(1024);
    k3 = new_rsa_key(1024);

    strcpy(msg, "This message has been encrypted thrice under different keys.");

    assert(-1 != bytes_to_hex(msg, strlen(msg), buf, sizeof(buf)));
    memset(msg, 0, sizeof(msg));
    mpz_set_str(m, buf, 16);
    mpz_powm(c1, m, *k1.e, *k1.n);
    mpz_powm(c2, m, *k2.e, *k2.n);
    mpz_powm(c3, m, *k3.e, *k3.n);
    assert(0 != mpz_cmp(c1, c2));
    assert(0 != mpz_cmp(c1, c3));
    assert(0 != mpz_cmp(c2, c3));
    mpz_set_ui(m, 0);

    mpz_mul(N, *k1.n, *k2.n);
    mpz_mul(N, N, *k3.n);

    mpz_mul(m1, *k2.n, *k3.n);
    assert(mpz_invert(m1, m1, *k1.n));
    mpz_mul(m1, m1, *k2.n);
    mpz_mul(m1, m1, *k3.n);
    mpz_mul(m1, m1, c1);

    mpz_mul(m2, *k1.n, *k3.n);
    assert(mpz_invert(m2, m2, *k2.n));
    mpz_mul(m2, m2, *k1.n);
    mpz_mul(m2, m2, *k3.n);
    mpz_mul(m2, m2, c2);

    mpz_mul(m3, *k1.n, *k2.n);
    assert(mpz_invert(m3, m3, *k3.n));
    mpz_mul(m3, m3, *k1.n);
    mpz_mul(m3, m3, *k2.n);
    mpz_mul(m3, m3, c3);

    mpz_add(m, m1, m2);
    mpz_add(m, m, m3);
    mpz_mod(m, m, N);
    mpz_root(m, m, 3);

    mpz_get_str(buf, 16, m);
    n = hex_to_bytes(buf, msg, sizeof(msg));
    assert(-1 != n);
    msg[n] = '\0';
    printf("%s\n", msg);
}
