#include <arpa/inet.h>
#include <assert.h>
#include <err.h>
#include <errno.h>
#include <netdb.h>
#include <netinet/in.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/time.h>
#include <sys/types.h> 
#include <sys/socket.h>
#include <unistd.h>

#include "crypt.c"
#include "encoding.c"

uint64_t GetTimeStamp() {
    struct timeval tv;
    gettimeofday(&tv,NULL);
    return tv.tv_sec*(uint64_t)1000000+tv.tv_usec;
}

uint64_t check_sig(FILE *fp, char *file, char *sig) {
    char buf[1024];
    buf[0] = '\0';
    strcat(buf, "GET /test?file=");
    strcat(buf, file);
    strcat(buf, "&signature=");
    strcat(buf, sig);
    strcat(buf, " HTTP/1.1\r\n\r\n");
    uint64_t before = GetTimeStamp();
    fputs(buf, fp);
    fflush(fp);
    char *line = NULL;
    size_t linecap = 0;
    ssize_t n;
    bool success = false;
    while ((n = getline(&line, &linecap, fp)) > 0) {
        if (strstr(line, "HTTP/1.1 200") != 0) {
            success = true;
        }
        if (n == 2 && memcmp(line, "\r\n", 2) == 0) break;
    }
    if (n == -1) err(1, "client couldn't read");
    uint64_t after = GetTimeStamp();
    if (success) return 0;
    return after - before;
}

int cmp(const void *e1, const void *e2) {
    uint64_t u1 = *((uint64_t *)e1);
    uint64_t u2 = *((uint64_t *)e2);
    if (u1 < u2) return -1;
    if (u1 > u2) return 1;
    return 0;
}

void guess_byte(FILE *fp, char *file, char *mac, size_t byte_pos) {
    uint64_t best_time = 0;
    char best_byte = 0;
    for (int b = -128; b < 128; ++b) {
        char c = b;
        assert(-1 != bytes_to_hex(&c, 1, mac + 2 * byte_pos, 3));
        int num_trials = 1;
        uint64_t trials[num_trials];
        for (int i = 0; i < num_trials; ++i) {
            trials[i] = check_sig(fp, file, mac);
        }
        qsort(trials, num_trials, sizeof(uint64_t), cmp);
        uint64_t t = trials[num_trials / 2];  // median
        // if (byte_pos == 11) {
        //     printf("%.2x: %zd\n", (unsigned char)c, t);
        // }
        if (best_time == 0 || t == 0 || (best_time != 0 && t > best_time)) {
            for (int i = 0; i < num_trials; ++i) {
                printf("%zd ", trials[i]);
            }
            printf("\n");
            best_byte = c;
            best_time = t;
        }
    }
    assert(-1 != bytes_to_hex(&best_byte, 1, mac + 2 * byte_pos, 3));
    printf("%.2x has best time: %zd\n", (unsigned char)best_byte, best_time);
    mac[2 * (byte_pos + 1) + 1] = '\0';
    printf("%s\n", mac);
}

int main() {
    srand(333);
    char key[20];
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }
    char ipad[20];
    memset(ipad, '6', 20);
    xor(ipad, key, 20, ipad);
    char opad[20];
    memset(opad, '\\', 20);
    xor(opad, key, 20, opad);



    int one = 1;
    struct sockaddr_in svr_addr;

    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock < 0) err(1, "can't open socket");

    setsockopt(sock, SOL_SOCKET, SO_REUSEADDR, &one, sizeof(int));

    int port = 7878;
    svr_addr.sin_family = AF_INET;
    svr_addr.sin_addr.s_addr = INADDR_ANY;
    svr_addr.sin_port = htons(port);

    if (connect(sock, (struct sockaddr *) &svr_addr, sizeof(svr_addr)) == -1) { 
        close(sock);
        err(1, "Can't connect");
    }
    FILE *fp = fdopen(sock, "w+");

    // Correct MAC is 50c5e5c9cc84bdbe69750da241822b6c9657ae89
    char *mac = calloc(21, 1);
    for (int i = 0; i < 20; ++i) {
        guess_byte(fp, "foo", mac, i);
    }

    printf("MAC: %s\n", mac);

    fclose(fp);

    free(mac);
}