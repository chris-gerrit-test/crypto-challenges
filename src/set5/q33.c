#include <assert.h>
#include <stdio.h>

#include <gmp.h>

#include "dh.c"

gmp_randstate_t rand_state;

int main() {
    dh_key key_a, key_b;
    mpz_t s;

    key_a = new_NIST_dh_key();
    key_b = new_NIST_dh_key();

    dh_session_key(key_a, &s, key_b.public);
    printf("(B**a)%%p: %s\n\n", mpz_get_str(NULL, 16, s));

    dh_session_key(key_b, &s, key_a.public);
    printf("(A**b)%%p: %s\n", mpz_get_str(NULL, 16, s));
}
