#include <stdio.h>

#include "crypt.c"
#include "encoding.c"
#include "sha1.h"

int main() {
    char mac[SHA1HashSize];
    sha1_mac("christopher", 11, "hundt", 5, mac);
    print_bytes(mac, SHA1HashSize);
}
