#ifndef DH_C
#define DH_C

#include <stdlib.h>

#include <gmp.h>

typedef struct rsa_key {
	mpz_t *e, *d, *n;
} rsa_key;

int rsa_inited = 0;
gmp_randstate_t rsa_rand_state;

void init_rsa() {
	if (rsa_inited) return;
	rsa_inited = 1;

    gmp_randinit_default(rsa_rand_state);
    gmp_randseed_ui(rsa_rand_state, 5128);
}


struct rsa_key new_rsa_key(int bits) {
	init_rsa();

	mpz_t* e = (mpz_t*)calloc(1, sizeof(mpz_t));
	mpz_t* d = (mpz_t*)calloc(1, sizeof(mpz_t));
	mpz_t* n = (mpz_t*)calloc(1, sizeof(mpz_t));
	mpz_t p, q, et;
	int r = 0;

	mpz_init(*e); mpz_init(*d); mpz_init(*n);
	mpz_init(p); mpz_init(q); mpz_init(et);

	while (r == 0) {
		mpz_urandomb(p, rsa_rand_state, bits/2);
		mpz_nextprime(p, p);
		mpz_urandomb(q, rsa_rand_state, bits/2);
		mpz_nextprime(q, q);

		mpz_mul(*n, p, q);
		mpz_sub_ui(p, p, 1);
		mpz_sub_ui(q, q, 1);
		mpz_mul(et, p, q);
		mpz_set_ui(*e, 3);
		r = mpz_invert(*d, *e, et);
	}

	return (rsa_key) {.e = e, .d = d, .n = n};
}

#endif /* DH_C */