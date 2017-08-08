#include <sys/socket.h>
#include <sys/stat.h>
#include <sys/types.h>

#include <netinet/in.h>

#include <fcntl.h>
#include <pwd.h>
#include <stdio.h>
#include <stdlib.h>
#include <syslog.h>
#include <unistd.h>

struct options {
	char *working_root;
	char *logfile;
	char *pidfile;
	int port;
	char *user;
};

void		set_defaults(struct options *);
int		check_opts(const struct options *);
void		daemonize(const struct options *);
void		forkerror(const pid_t);
void		prepare_logging(const char *, const char *);
void		close_fds(void);
void		deescalate(const char *);
int		init_server(int, int *);
void		score_loop(int, char *);
_Noreturn void	bail(const char *);

int
main(int argc, char **argv)
{
	int ch;
	static const char help[] = "Usage: scored [options]";
	struct options cli_opts;

	set_defaults(&cli_opts);
	while ((ch = getopt(argc, argv, "hd:l:p:")) != -1) {
		switch (ch) {
			case 'd':
				cli_opts.working_root = optarg;
				break;
			case 'h':
				fprintf(stderr, "%s: %s", argv[0], help);
				return 0;
			case 'p':
				cli_opts.port = strtol(optarg, NULL, 10);
				break;
				
		default:
			break;
		}
	}
	if (ch == '?')
		return 1;

	if (check_opts(&cli_opts))
		return 1;

	daemonize(&cli_opts);	

	return 0;
}

void
set_defaults(struct options *opts)
{
	opts->working_root = "/";
	opts->logfile = "/var/log/scored.log";
	opts->pidfile = "/var/run/scored.pid";
	opts->port = 30000;
	opts->user = "mark";
}

int check_opts(const struct options *opts)
{
	struct stat file_stat;

	if (!access(opts->pidfile, F_OK)) {
		fprintf(stderr, "scored: already running? %s\n", opts->pidfile);
		return 1;
	}
	/* working dir */
	if (stat(opts->working_root, &file_stat)) {
		perror(opts->working_root);
		return 1;
	} else if (!S_ISDIR(file_stat.st_mode)) {
		fprintf(stderr, "%s: Not a directory\n", opts->working_root);
		return 1;
	}

	if (opts->port == 0) {
		fprintf(stderr, "Invalid port\n");
		return 1;
	}

	return 0;
}

/*
 * run in the background
 * enter desired directory
 * prepare a socket
 * begin the _real_ program loop
 */
void
daemonize(const struct options *opts)
{
	int socket;

	forkerror(fork());
	if (setsid() < 0)
		bail(__func__);
	forkerror(fork());

	prepare_logging(opts->logfile, opts->user);
	close_fds();

	if (chdir(opts->working_root))
		exit(EXIT_FAILURE);
	if (init_server(opts->port, &socket))
		exit(EXIT_FAILURE);
	
	deescalate(opts->user);
	score_loop(socket, opts->logfile);
}

void
forkerror(const pid_t pid)
{
	if (pid < 0)
		bail(__func__);
	if (pid > 0)
		exit(EXIT_SUCCESS);
}

void
prepare_logging(const char *path, const char *username)
{
	struct passwd *userpw;

	if (access(path, F_OK))
		close(open(path, O_CREAT, 0644));

	userpw = getpwnam(username);
	if (userpw == NULL)
		bail(__func__);

	if (chown(path, userpw->pw_uid, userpw->pw_gid))
		bail(__func__);
	if (chmod(path, 0644))
		bail(__func__);
}

void
close_fds(void)
{
	int fd;

	for (fd = sysconf(_SC_OPEN_MAX); fd >= 0; fd--)
		close(fd);
}

void
deescalate(const char *username)
{
	struct passwd *userpw;

	userpw = getpwnam(username);
	if (getuid() == 0) {
		if (setgid(userpw->pw_gid) != 0)
			bail(__func__);
		if (setuid(userpw->pw_uid) != 0)
			bail(__func__);
	}
	if (seteuid(0) != -1)
		bail(__func__);
}

int
init_server(int port, int *fd)
{
	struct sockaddr_in addr = {0};

	if ((*fd = socket(AF_INET, SOCK_STREAM, 0)) == -1)
		return 1;
	setsockopt(*fd, SOL_SOCKET, SO_REUSEADDR, &(int){1}, sizeof(int));

	addr.sin_family = AF_INET;
	addr.sin_addr.s_addr = htonl(INADDR_ANY);
	addr.sin_port = htons(port);

	if (bind(*fd, (struct sockaddr *)&addr, sizeof(addr)))
		return 1;
	if (listen(*fd, 50))
		return 1;
	return 0;
}

_Noreturn void
bail(const char *err)
{
	perror(err);
	exit(EXIT_FAILURE);
}
