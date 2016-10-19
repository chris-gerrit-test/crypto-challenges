#include <assert.h>
#include <stdio.h>

#include "crypt.c"

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

    printf("*** ATTACK 1 ***\n");
    // Attack 1: can pick IV
    //
    // I am user 123 and I want to steal 1M from user 457.
    // First generate a message sending the money to myself.
    // Message:
    //   from=123&to=123&
    //   amount=1000000..
    sign_transfer(123, 123, 1000000, iv, mac);
    printf("Legit message MAC:\n");
    print_bytes(mac, sizeof(mac));

    memcpy(new_iv, iv, sizeof(iv));
    // Swap 123 -> 457
    new_iv[5] ^= '1' ^ '4';
    new_iv[6] ^= '2' ^ '5';
    new_iv[7] ^= '3' ^ '7';
    // Same MAC should now be generated with a different "to" account.
    // Message:
    //   from=457&to=123&
    //   amount=1000000..
    sign_transfer(457, 123, 1000000, new_iv, mac);
    printf("Modified message MAC w/ new IV:\n");
    print_bytes(mac, sizeof(mac));

    printf("\n*** ATTACK 2 ***\n");
    // Attack 2: IV is fixed at zero
    //
    // Not sure about this one. Wikipedia describes a length extension
    // attack, but it only allows me to concatenate previously signed
    // messages (with one of them XORed by the other's MAC). I don't see
    // how to get the client to sign those. When I look at other people's
    // solutions on the web they just use the shared key, which seems to
    // miss the point--if I have that then I don't need to "break" anything.

    memset(iv, 0, sizeof(iv));

    // Message 1: user 47 sends $1 to user 48
    char msg1[1024];
    sprintf(msg1, "from=%d&tx_list=%d:%d", 47, 48, 1);
    compute_mac(msg1, iv, mac);
    printf("Target message:\n %s\nTarget message MAC:\n", msg1);
    print_bytes(mac, sizeof(mac));

    // Message 2: attacker (123) sends $1M to herself.
    char msg2[1024];
    sprintf(msg2, "from=%d&tx_list=123:0;123:0;0:0;%d:%d", 123, 123, 1000000);
    compute_mac(msg2, iv, mac);
    printf("\nAttacker message:\n %s\nAttacker message MAC:\n", msg2);
    print_bytes(mac, sizeof(mac));

    // Message 3: 1 + 2 (with first block of 2 mangled)
    // First add padding to msg 1.
    size_t len = strlen(msg1);
    len = (1 + len / 16) * 16;  // make sure there is room for padding
    char *padded_msg1 = calloc(1024, 1);
    strcpy(padded_msg1, msg1);
    int padded_len1 = pkcs7(padded_msg1, len + 1, 16);
    assert((int)len == padded_len1);
    // Now add message 2 and mangle first block.
    memcpy(padded_msg1 + padded_len1, msg2, strlen(msg2));
    xor(padded_msg1 + padded_len1, mac, 16, padded_msg1 + padded_len1);
    // Verify MAC.
    compute_mac(msg2, iv, mac);
    printf("\nCombined message:\n %s\nCombined message MAC:\n", padded_msg1);
    print_bytes(mac, sizeof(mac));

}
