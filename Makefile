CFLAGS=-g -Wall -Wextra -pedantic -O2
LDFLAGS=-lcrypto

OBJ = \
      init.o main.o

all: scored

scored: $(OBJ)
	$(CC) $(OBJ) -o scored $(LDFLAGS)

%.o: %.c
	$(CC) -c $(CFLAGS) $< -o $@

clean:
	rm $(OBJ) scored
