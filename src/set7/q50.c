#include <assert.h>
#include <stdio.h>

#include "crypt.c"

typedef struct transfer {
    int to, amount;
} transfer;

char key[16];

void compute_mac(char *msg, char iv[16], char mac[16]) {
    char *c;
    size_t len;

    len = strlen(msg);
    len = (1 + len / 16) * 16;  // make sure there is room for padding
    c = calloc(len, 1);
    strcpy(c, msg);
    assert(-1 != pkcs7(c, len, 16));
    cbc_encrypt(c, c, len, iv, key);
    memcpy(mac, c + len - 16, 16);
    free(c);
}

int main() {
    srand(1210);
    char iv[16], mac[16], buf[1024];
    char *orig = "alert('MZA who was that?');\n";
    size_t padded_len;
    //memset(mac, 0, sizeof(mac));

    memcpy(key, "YELLOW SUBMARINE", 16);
    memset(iv, 0, sizeof(iv));

    compute_mac(orig, iv, mac);
    printf("%s", orig);
    print_bytes(mac, sizeof(mac));
    printf("\n");

    strcpy(buf, "alert('Ayo, the Wu is back!'); //");
    compute_mac(buf, iv, mac);
    padded_len = pkcs7(buf, sizeof(buf), sizeof(mac));
    strcat(buf, orig);
    xor(buf + padded_len, mac, sizeof(mac), buf + padded_len);
    compute_mac(buf, iv, mac);
    printf("%s", buf);
    print_bytes(mac, sizeof(mac));
}
