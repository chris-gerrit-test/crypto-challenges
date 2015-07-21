#include <stdio.h>
#include <time.h>

#include "math.c"
#include "mersenne.c"

int main() {
    srand(100);
    int now = time(NULL);
    twister mt;

    int wait_time = 40 + randn(960);
    mt_seed(&mt, now + wait_time);

    uint32_t n = mt_extract(&mt);

    for (int guess = 40; guess < 1000; ++guess) {
        mt_seed(&mt, now + guess);
        if (mt_extract(&mt) == n) {
            printf("Guessed: %d\n", guess);
            // Keep going to see if there are multiple guesses.
        }
    }

    printf("Actual: %d\n", wait_time);
}
