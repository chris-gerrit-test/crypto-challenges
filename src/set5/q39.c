#include <assert.h>
#include <stdio.h>

#include <gmp.h>

#include "crypt.c"
#include "encoding.c"
#include "rsa.c"
#include "sha1.h"

int main() {
    char msg[512], buf[1025];
    rsa_key k;
    mpz_t m, c;
    int n;

    mpz_init(m); mpz_init(c);

    k = new_rsa_key(1024);

    strcpy(msg, "This message is very secret so don't try to read it.");

    assert(-1 != bytes_to_hex(msg, strlen(msg), buf, sizeof(buf)));
    mpz_set_str(m, buf, 16);
    mpz_powm(c, m, *k.e, *k.n);
    mpz_powm(m, c, *k.d, *k.n);
    mpz_get_str(buf, 16, m);
    n = hex_to_bytes(buf, msg, sizeof(msg));
    assert(-1 != n);
    msg[n] = '\0';
    printf("%s\n", msg);
}
