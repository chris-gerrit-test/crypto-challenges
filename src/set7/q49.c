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
    c = calloc(len + 1, 1);
    strcpy(c, msg);
    assert(-1 != pkcs7(c, len + 1, 16));
    cbc_encrypt(c, c, len, iv, key);
    memcpy(mac, c + len - 16, 16);
    free(c);
}

void sign_transfer(int from, int to, int amount, char iv[16], char mac[16]) {
    char buf[1024];
    sprintf(buf, "from=%d&to=%d&amount=%d", from, to, amount);
    compute_mac(buf, iv, mac);
}

void sign_transfers(int from, transfer *transfers, int num_transfers, char iv[16], char mac[16]) {
    char buf[1024];
    sprintf(buf, "from=%d&tx_list=", from);
    for (int i = 0; i < num_transfers; ++i) {
        sprintf(buf + strlen(buf), "%d:%d;", transfers[i].to, transfers[i].amount);
    }
    if (num_transfers) {
        // Strip trailing semicolon.
        buf[strlen(buf) - 1] = '\0';
    }
    compute_mac(buf, iv, mac);
}

int main() {
    srand(1210);
    char iv[16], new_iv[16], mac[16];
    //memset(mac, 0, sizeof(mac));

    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }
    for (size_t i = 0; i < sizeof(iv); ++i) {
        iv[i] = randn(256) - 128;
    }

    // Attack 1: can pick IV
    //
    // I am user 123 and I want to steal 1M from user 457.
    // Now intercept a MAC and message from the target.
    // (I don't see any way to do it if the target isn't already
    //  transferring $1M to someone.)
    // Message:
    //   from=457&to=968&
    //   amount=1000000..
    sign_transfer(457, 968, 1000000, iv, mac);
    print_bytes(mac, sizeof(mac));

    memcpy(new_iv, iv, sizeof(iv));
    // Swap 968 -> 123
    new_iv[12] ^= '9' ^ '1';
    new_iv[13] ^= '6' ^ '2';
    new_iv[14] ^= '8' ^ '3';
    // Same MAC should now be generated with a different "to" account.
    // Message:
    //   from=457&to=123&
    //   amount=1000000..
    sign_transfer(457, 123, 1000000, new_iv, mac);
    print_bytes(mac, sizeof(mac));

    // Attack 2: IV is fixed at zero
    //
    // Not sure about this one. Wikipedia describes a length extension
    // attack, but it only allows me to concatenate previously signed
    // messages (with one of them XORed by the other's MAC). I don't see
    // how to get the client to sign those. When I look at other people's
    // solutions on the web they just use the shared key, which seems to
    // miss the point--if I have that then I don't need to "break" anything.
}
