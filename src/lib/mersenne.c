#ifndef MERSENNE_C
#define MERSENNE_C

#include <stdint.h>

#define MERSENNE_C_N 624
#define MERSENNE_C_R 31
#define MERSENNE_C_LOWER_MASK ((1 << MERSENNE_C_R) - 1)

uint32_t w = 32;
uint32_t n = MERSENNE_C_N;
uint32_t m = 397;
uint32_t r = MERSENNE_C_R;
uint32_t a = 0x9908B0DF;
uint32_t u = 11;
uint32_t s = 7;
uint32_t b = 0x9D2C5680;
uint32_t t = 15;
uint32_t c = 0xEFC60000;
uint32_t l = 18;
uint32_t f = 1812433253;
uint32_t lower_mask = MERSENNE_C_LOWER_MASK;
uint32_t upper_mask = ~MERSENNE_C_LOWER_MASK;

typedef struct twister {
    uint32_t state[MERSENNE_C_N];
    uint32_t index;
} twister;

void mt_seed(twister *mt, uint32_t seed) {
    mt->index = n;
    mt->state[0] = seed;
    for (uint32_t i = 1; i < n; ++i) {
        mt->state[i] = f * (mt->state[i - 1] ^ (mt->state[i - 1] >> (w - 2))) + i;
    }
}

void mt_twist(twister *mt) {
    for (uint32_t i = 0; i < n; ++i) {
        uint32_t x = (mt->state[i] & upper_mask) + (mt->state[(i + 1) % n] & lower_mask);
        uint32_t xA = x >> 1;
        if (x & 1) {
            xA ^= a;
        }
        mt->state[i] = mt->state[(i + m) % n] ^ xA;
    }
    mt->index = 0;
}

uint32_t mt_extract(twister *mt) {
    if (mt->index >= n) {
        mt_twist(mt);
    }

    uint32_t y = mt->state[mt->index];
    y ^= (y >> u);
    y ^= ((y << s) & b);
    y ^= ((y << t) & c);
    y ^= (y >> l);
    mt->index++;
    return y;
}

uint32_t untemper(uint32_t y) {
    y ^= (y >> l);
    y ^= ((y << t) & c);
    y = y ^ ((y ^ ((y ^ ((y ^ (y << s) & b) << s) & b) << s) & b) << s) & b;
    y = y ^ ((y ^ (y >> u)) >> u);
    return y;
}

#endif /* MERSENNE_C */
