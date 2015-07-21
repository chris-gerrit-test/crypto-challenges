#include <assert.h>
#include <stdio.h>

#include "crypt.c"
#include "encoding.c"
#include "score.c"

char *inputs[] = {
   "SSBoYXZlIG1ldCB0aGVtIGF0IGNsb3NlIG9mIGRheQ==",
   "Q29taW5nIHdpdGggdml2aWQgZmFjZXM=",
   "RnJvbSBjb3VudGVyIG9yIGRlc2sgYW1vbmcgZ3JleQ==",
   "RWlnaHRlZW50aC1jZW50dXJ5IGhvdXNlcy4=",
   "SSBoYXZlIHBhc3NlZCB3aXRoIGEgbm9kIG9mIHRoZSBoZWFk",
   "T3IgcG9saXRlIG1lYW5pbmdsZXNzIHdvcmRzLA==",
   "T3IgaGF2ZSBsaW5nZXJlZCBhd2hpbGUgYW5kIHNhaWQ=",
   "UG9saXRlIG1lYW5pbmdsZXNzIHdvcmRzLA==",
   "QW5kIHRob3VnaHQgYmVmb3JlIEkgaGFkIGRvbmU=",
   "T2YgYSBtb2NraW5nIHRhbGUgb3IgYSBnaWJl",
   "VG8gcGxlYXNlIGEgY29tcGFuaW9u",
   "QXJvdW5kIHRoZSBmaXJlIGF0IHRoZSBjbHViLA==",
   "QmVpbmcgY2VydGFpbiB0aGF0IHRoZXkgYW5kIEk=",
   "QnV0IGxpdmVkIHdoZXJlIG1vdGxleSBpcyB3b3JuOg==",
   "QWxsIGNoYW5nZWQsIGNoYW5nZWQgdXR0ZXJseTo=",
   "QSB0ZXJyaWJsZSBiZWF1dHkgaXMgYm9ybi4=",
   "VGhhdCB3b21hbidzIGRheXMgd2VyZSBzcGVudA==",
   "SW4gaWdub3JhbnQgZ29vZCB3aWxsLA==",
   "SGVyIG5pZ2h0cyBpbiBhcmd1bWVudA==",
   "VW50aWwgaGVyIHZvaWNlIGdyZXcgc2hyaWxsLg==",
   "V2hhdCB2b2ljZSBtb3JlIHN3ZWV0IHRoYW4gaGVycw==",
   "V2hlbiB5b3VuZyBhbmQgYmVhdXRpZnVsLA==",
   "U2hlIHJvZGUgdG8gaGFycmllcnM/",
   "VGhpcyBtYW4gaGFkIGtlcHQgYSBzY2hvb2w=",
   "QW5kIHJvZGUgb3VyIHdpbmdlZCBob3JzZS4=",
   "VGhpcyBvdGhlciBoaXMgaGVscGVyIGFuZCBmcmllbmQ=",
   "V2FzIGNvbWluZyBpbnRvIGhpcyBmb3JjZTs=",
   "SGUgbWlnaHQgaGF2ZSB3b24gZmFtZSBpbiB0aGUgZW5kLA==",
   "U28gc2Vuc2l0aXZlIGhpcyBuYXR1cmUgc2VlbWVkLA==",
   "U28gZGFyaW5nIGFuZCBzd2VldCBoaXMgdGhvdWdodC4=",
   "VGhpcyBvdGhlciBtYW4gSSBoYWQgZHJlYW1lZA==",
   "QSBkcnVua2VuLCB2YWluLWdsb3Jpb3VzIGxvdXQu",
   "SGUgaGFkIGRvbmUgbW9zdCBiaXR0ZXIgd3Jvbmc=",
   "VG8gc29tZSB3aG8gYXJlIG5lYXIgbXkgaGVhcnQs",
   "WWV0IEkgbnVtYmVyIGhpbSBpbiB0aGUgc29uZzs=",
   "SGUsIHRvbywgaGFzIHJlc2lnbmVkIGhpcyBwYXJ0",
   "SW4gdGhlIGNhc3VhbCBjb21lZHk7",
   "SGUsIHRvbywgaGFzIGJlZW4gY2hhbmdlZCBpbiBoaXMgdHVybiw=",
   "VHJhbnNmb3JtZWQgdXR0ZXJseTo=",
   "QSB0ZXJyaWJsZSBiZWF1dHkgaXMgYm9ybi4=",
};

char key[16];
char nonce[8] = {0};

void init_keys() {
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }
}

size_t min(size_t x, size_t y) {
    return x < y ? x : y;
}

int main() {
    srand(3690);
    init_keys();

    char *key_guess = calloc(16, 1);

    char *buf = calloc(256, 1);
    for (int ci = 0; ci < 38; ++ci) {
        double best_score = 0.0;
        int best_char = 0;
        for (int c = -128; c < 128; ++c) {
            key_guess[ci] = c;
            int good = 1;
            double score = 0.0;
            for (int i = 0; i < 40; ++i) {
                char *s = inputs[i];
                strcpy(buf, s);
                int n = base64_to_bytes(buf, buf, strlen(buf));
                assert(n > 0);
                buf[n] = '\0';
                ctr_crypt(buf, buf, n, nonce, key);
                xor(buf, key_guess, n, buf);
                score += score_text(buf, n);
            }
            if (score > best_score) {
                best_score = score;
                best_char = c;

                //printf("char %d: %d (%f)\n", ci, c, best_score);
            } else {
                //printf("char %d: %d (%f)\n", ci, c, score);
            }
        }
        key_guess[ci] = best_char;
    }

    // Fixes
    key_guess[0] ^= 'i' ^ 'I';
    key_guess[30] ^= 'N' ^ 'n';
    key_guess[31] ^= 'd' ^ 'e';
    key_guess[34] ^= 'e' ^ 'a';
    key_guess[35] ^= 'E' ^ 'd';
    key_guess[36] ^= ' ' ^ 'n';

    for (int i = 0; i < 40; ++i) {
        char *s = inputs[i];
        strcpy(buf, s);
        int n = base64_to_bytes(buf, buf, strlen(buf));
        assert(n > 0);
        buf[n] = '\0';
        ctr_crypt(buf, buf, n, nonce, key);
        xor(buf, key_guess, n, buf);
        printf("%s\n", buf);
    }

    free(buf);

    free(key_guess);
}
