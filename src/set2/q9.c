#include <stdio.h>

#include "crypt.c"
#include "encoding.c"

int main() {
    char str[21] = "YELLOW SUBMARINE";

    print_bytes(str, sizeof(str));
    pkcs7(str, sizeof(str), 20);
    print_bytes(str, sizeof(str));
}
