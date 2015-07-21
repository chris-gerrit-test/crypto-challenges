#include <stdio.h>

#include "math.c"
#include "mersenne.c"

int main() {
    twister mt, mt_clone;
    mt_seed(&mt, 1324136);
    mt_seed(&mt_clone, 123);

    for (int i = 0; i < 624; ++i) {
        mt_clone.state[i] = untemper(mt_extract(&mt));
    }
    mt_clone.index = 624;

    printf("Original: %10u   Clone: %10u\n", mt_extract(&mt), mt_extract(&mt_clone));
    printf("Original: %10u   Clone: %10u\n", mt_extract(&mt), mt_extract(&mt_clone));
    printf("Original: %10u   Clone: %10u\n", mt_extract(&mt), mt_extract(&mt_clone));
    printf("Original: %10u   Clone: %10u\n", mt_extract(&mt), mt_extract(&mt_clone));
}
