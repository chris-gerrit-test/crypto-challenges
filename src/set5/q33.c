#include <assert.h>
#include <stdio.h>

#include <gmp.h>

char *p_str = "ffffffffffffffffc90fdaa22168c234c4c6628b80dc1cd129024"
              "e088a67cc74020bbea63b139b22514a08798e3404ddef9519b3cd"
              "3a431b302b0a6df25f14374fe1356d6d51c245e485b576625e7ec"
              "6f44c42e9a637ed6b0bff5cb6f406b7edee386bfb5a899fa5ae9f"
              "24117c4b1fe649286651ece45b3dc2007cb8a163bf0598da48361"
              "c55d39a69163fa8fd24cf5f83655d23dca3ad961c62f356208552"
              "bb9ed529077096966d670c354e4abc9804f1746c08ca237327fff"
              "fffffffffffff";

char *g_str = "2";

gmp_randstate_t rand_state;

int main() {
    mpz_t p, g, a, A, b, B, s;

    gmp_randinit_default(rand_state);
    gmp_randseed_ui(rand_state, 0);

    // Constants
    mpz_init_set_str(p, p_str, 16);
    mpz_init_set_str(g, g_str, 16);

    // DH
    mpz_init(a);
    mpz_urandomm(a, rand_state, p);
    mpz_init(b);
    mpz_urandomm(b, rand_state, p);

    mpz_init(A);
    mpz_powm(A, g, a, p);

    mpz_init(B);
    mpz_powm(B, g, b, p);

    mpz_init(s);
    mpz_powm(s, B, a, p);

    printf("(B**a)%%p: %s\n\n", mpz_get_str(NULL, 16, s));

    mpz_powm(s, A, b, p);
    printf("(A**b)%%p: %s\n", mpz_get_str(NULL, 16, s));
}
