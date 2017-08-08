# CyberPatriot Scoring Engine
This will be a daemon that records points for changes in configuration and
state on another system. It is not usable yet, but should be soon. Clients ran
on the testing machine will connect to the server, authenticate, send
information on scoring changes, then disconnect.

After this is working I'll work on documentation.

# Known issues
* **INSECURE** - Passwords are hashed using SHA512 without salts right now.
Later it will use scrypt.
