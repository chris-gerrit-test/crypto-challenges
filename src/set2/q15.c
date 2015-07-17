#include <stdio.h>

#include "crypt.c"
#include "encoding.c"

void check_padding(char *s) {
    char *stripped = strip_pkcs7(s, strlen(s), 16);
    print_bytes(s, strlen(s));
    if (stripped == NULL) {
        printf(" - BAD PADDING\n");
    } else {
        printf(" -> ");
        print_bytes(s, stripped - s);
    }
}

int main() {
    check_padding("ICE ICE BABY\x04\x04\x04\x04");
    check_padding("ICE ICE BABY\x05\x05\x05\x05");
    check_padding("ICE ICE BABY\x01\x02\x03\x04");
    check_padding("ICE ICE BABY!\x03\x03\x03");
    check_padding("ICE ICE BABY!!!!\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10");

    char *s = calloc(32, 1);
    memcpy(s, "YELLOW SUBMARINE", 16);
    pkcs7(s, 32, 16);
    check_padding(s);
}
