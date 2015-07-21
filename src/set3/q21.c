#include <stdio.h>

#include "mersenne.c"

int main() {
    twister mt;
    mt_seed(&mt, 10);
    for (int i = 0; i < 10; ++i) {
        printf("%u\n", mt_extract(&mt));
    } 
}
