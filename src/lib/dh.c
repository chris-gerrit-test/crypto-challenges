#ifndef DH_C
#define DH_C

#include <stdlib.h>

#include <gmp.h>

typedef struct dh_key {
	mpz_t *p, *g, *private, *public;
} dh_key;

char *nist_p = "ffffffffffffffffc90fdaa22168c234c4c6628b80dc1cd129024"
               "e088a67cc74020bbea63b139b22514a08798e3404ddef9519b3cd"
               "3a431b302b0a6df25f14374fe1356d6d51c245e485b576625e7ec"
               "6f44c42e9a637ed6b0bff5cb6f406b7edee386bfb5a899fa5ae9f"
               "24117c4b1fe649286651ece45b3dc2007cb8a163bf0598da48361"
               "c55d39a69163fa8fd24cf5f83655d23dca3ad961c62f356208552"
               "bb9ed529077096966d670c354e4abc9804f1746c08ca237327fff"
               "fffffffffffff";
char *nist_g = "2";

int dh_inited = 0;
gmp_randstate_t rand_state;

void init_dh() {
	if (dh_inited) return;
	dh_inited = 1;

    gmp_randinit_default(rand_state);
    gmp_randseed_ui(rand_state, 0);
}

struct dh_key new_NIST_dh_key() {
	mpz_t* p = (mpz_t*)calloc(1, sizeof(mpz_t));
	mpz_t* g = (mpz_t*)calloc(1, sizeof(mpz_t));
	mpz_t* a = (mpz_t*)calloc(1, sizeof(mpz_t));
	mpz_t* A = (mpz_t*)calloc(1, sizeof(mpz_t));

	init_dh();
    mpz_init_set_str(*p, nist_p, 16);
    mpz_init_set_str(*g, nist_g, 16);
    mpz_init(*a);
    mpz_urandomm(*a, rand_state, *p);
    mpz_init(*A);
    mpz_powm(*A, *g, *a, *p);

	return (dh_key) {.p = p, .g = g, .private = a, .public = A};
}

void dh_session_key(dh_key key, mpz_t *sess, mpz_t *other_pubkey) {
    mpz_init(*sess);
    mpz_powm(*sess, *other_pubkey, *key.private, *key.p);
}

#endif /* DH_C */