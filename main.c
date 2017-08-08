#include <sys/types.h>
#include <sys/socket.h>

#include <arpa/inet.h>
#include <netinet/in.h>

#include <openssl/evp.h>
#include <openssl/sha.h>

#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#define MAXPASSWDSIZ 2048
#define SHA512_HASH_HEX_LEN 129

void	score_loop(int, char *);
int	get_passwd_hash(char *);
void	test_request_validity(const char *, int *, const char *);
int	authenticate(const char *, const char *);
void	generate_hd(unsigned char *, char *, size_t );
void	handle_request(const int, int, char *);
void	execute_instructions(const char *);
void	logrequest(const char *, const char *, const int, const char *,
		const char *);

static const char *RESPONSE[] = {
	"0:SUCCESS:EOM:\n",
	"1:AUTH_FAIL:EOM:\n",
	"2:REQ_FAIL:EOM:\n"
};

/* program starts here */
void
score_loop(int socket, char *logpath)
{
	char client_name[INET_ADDRSTRLEN];
	char passwd_hash[MAXPASSWDSIZ];
	char requestbuf[8192];
	int connection;
	int status_code;
	socklen_t destlen;
	struct sockaddr_in dest;

	destlen = sizeof(dest);

	for (;;) {
		connection = accept(socket, (struct sockaddr *)&dest, &destlen);
		inet_ntop(AF_INET, &dest.sin_addr.s_addr, client_name, INET_ADDRSTRLEN);
		recv(connection, requestbuf, 8192, 0);

		if (get_passwd_hash(passwd_hash) != 0)
			status_code = 0;
		else
			test_request_validity(requestbuf, &status_code,
					passwd_hash);
		handle_request(status_code, connection, requestbuf);
		
		logrequest(client_name, logpath, status_code, "", "");

		close(connection);
	}
}

int
get_passwd_hash(char *buf)
{
	FILE *fp;

	if ((fp = fopen("/etc/scored-passwd", "r")) == NULL)
		return 1;
	
	if (fgets(buf, SHA512_HASH_HEX_LEN, fp) == NULL)
		return 1;

	return 0;
}

void
test_request_validity(const char *request, int *status,
		const char real_pw_hash[SHA512_DIGEST_LENGTH])
{
	char passwd[MAXPASSWDSIZ];
	char *tmp;

	if (strstr(request, ":EOM:") == NULL) {
		*status = 2;
		return;
	}
	tmp = strdup(request);
	strncpy(passwd, strtok(tmp, ":"), MAXPASSWDSIZ);
	free(tmp);

	*status = authenticate(real_pw_hash, passwd);

	return;
}


int
authenticate(const char real_pw_hash[SHA512_DIGEST_LENGTH],
		const char *client_pw)
{
	char client_pw_hash[SHA512_HASH_HEX_LEN];
	char *tmp;
	EVP_MD_CTX *mdctx;
	int i;
	unsigned char digest[SHA512_DIGEST_LENGTH];

	/* some OpenSSL magic to get the raw digest */
	mdctx = EVP_MD_CTX_create();
	EVP_DigestInit_ex(mdctx, EVP_sha512(), NULL);
	EVP_DigestUpdate(mdctx, client_pw, strlen(client_pw));
	OPENSSL_malloc(EVP_MD_size(EVP_sha512()));
	EVP_DigestFinal_ex(mdctx, digest,
			&(unsigned int){SHA512_DIGEST_LENGTH});

	for (i = 0; i < SHA512_DIGEST_LENGTH; i++) {
		tmp = strdup(client_pw_hash);
		snprintf(client_pw_hash, SHA512_HASH_HEX_LEN, "%s%02x", tmp,
				digest[i]);
		free(tmp);
	}

	if (strcmp(real_pw_hash, client_pw_hash))
		return 1;
	else
		return 0;
}

void
handle_request(const int status, int consocket, char *request)
{
	switch (status) {
	case 0:
		execute_instructions(request);
		send(consocket, RESPONSE[0], strlen(RESPONSE[0]), 0);
		break;
	case 1:
		send(consocket, RESPONSE[1], strlen(RESPONSE[1]), 0);
		break;
	case 2:
		send(consocket, RESPONSE[2], strlen(RESPONSE[2]), 0);
		break;
	}
}

void
execute_instructions(const char *static_instructions)
{
	char *instructions;

	instructions = strdup(static_instructions);
	free(instructions);
}

void
logrequest(const char *client, const char *log, const int status,
		const char *pt_change, const char *awarded_for)
{
	FILE *fp;
	time_t since_epoch;

	fp = fopen(log, "a");

	since_epoch = time(NULL);
	fprintf(fp, "%li:\t[%s]\t%s\t%s\t%s", since_epoch, client,
			RESPONSE[status], pt_change, awarded_for);

	fclose(fp);
}
