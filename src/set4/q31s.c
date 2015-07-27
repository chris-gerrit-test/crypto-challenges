#include <arpa/inet.h>
#include <assert.h>
#include <err.h>
#include <netdb.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h> 
#include <sys/socket.h>
#include <unistd.h>

#include "crypt.c"
#include "encoding.c"

char *pass = "HTTP/1.1 200 OK\r\n\r\n";
char *fail = "HTTP/1.1 403 Forbidden\r\n\r\n";

int insecure_compare(char *mac1, char *mac2) {
    for (int i = 0; i < 20; ++i) {
        if (mac1[i] != mac2[i]) return 1;
        usleep(50000);
    }
    return 0;
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

    int one = 1, client_fd;
    struct sockaddr_in svr_addr, cli_addr;
    socklen_t sin_len = sizeof(cli_addr);

    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock < 0) err(1, "can't open socket");

    setsockopt(sock, SOL_SOCKET, SO_REUSEADDR, &one, sizeof(int));

    int port = 7878;
    svr_addr.sin_family = AF_INET;
    svr_addr.sin_addr.s_addr = INADDR_ANY;
    svr_addr.sin_port = htons(port);

    if (bind(sock, (struct sockaddr *) &svr_addr, sizeof(svr_addr)) == -1) { 
        close(sock);
        err(1, "Can't bind");
        return 1;
    }

    listen(sock, 5);
    while (1) {
        client_fd = accept(sock, (struct sockaddr *) &cli_addr, &sin_len);

        if (client_fd == -1) {
            perror("Can't accept");
            continue;
        }
        FILE *fp = fdopen(client_fd, "w+");

        while(1) {
            char *line = NULL;
            size_t linecap = 0;
            ssize_t n;
            if (!fp) err(1, "Can't fdopen");
            char *file = NULL;
            char *sig = NULL;
            while ((n = getline(&line, &linecap, fp)) > 0) {
                if (file == NULL) {
                    // crappy parsing of GET line
                    file = calloc(n, 1);
                    sig = calloc(n, 1);
                    if (strstr(line, "GET /test?") != line) break;
                    char *start = strstr(line, "file=");
                    if (start == NULL) break;
                    start = start + strlen("file=");
                    char *end = strstr(start, "&");
                    if (end == NULL) end = strstr(start, " ");
                    if (end == NULL) end = line + n - 2;
                    memcpy(file, start, end - start);
                    start = strstr(line, "signature=");
                    if (start == NULL) break;
                    start = start + strlen("signature=");
                    end = strstr(start, "&");
                    if (end == NULL) end = strstr(start, " ");
                    if (end == NULL) end = line + n - 2;
                    memcpy(sig, start, end - start);
                }
                if (n == 2 && memcmp(line, "\r\n", 2) == 0) break;
            }
            if (n == -1) err(1, "Server couldn't read");
            char *as_bytes = calloc(20, 1);
            if (sig) hex_to_bytes(sig, as_bytes, 20);

            unsigned char mac[20];
            if (file) {
                unsigned char buf[20];
                sha1_mac(ipad, sizeof(ipad), file, strlen(file), (char *)buf);
                sha1_mac(opad, sizeof(opad), (char *)buf, 20, (char *)mac);
            }

            if (insecure_compare(as_bytes, (char *)mac)) {
                if (fputs(fail, fp) <= 0) {
                    err(1, "Couldn't write");
                }
            } else {
                if (fputs(pass, fp) <= 0)  {
                    err(1, "Couldn't write");
                }
            }

            free(file);
            free(sig);
        }

        fclose(fp);
        close(client_fd);
    }
}