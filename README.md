# Daemon
This is an extension of the ![](github.com/sevlyar/go-daemon) library.
It is an example and a template of how to create a simple daemon that
behaves how you would expect a Unix daemon to behave.

first you run the daemon and it forks itself into the background.
Then, you can send signals to the daemon to make it quit, reload or stop.
But we also want to be able to send commands to a daemon... so I've added
a simple Client Server that communicates over a Unix socket. After forking
the daemon off into the background it will start a listener Unix socket,
you can then re-run the same binary with different flags/args and the
daemon will process those commands.

Fork it:
```bash
./daemon
```

Send a command:
```bash
./daemon -e "notify-send hello world"
```
you can process these commands and flags like you would with any other golang program.
The limits are your imagination here. An example would be asking a wall-paper daemon
to change the wallpaper. This is all handled in `server.go`

Reload the daemon configuration:
```bash
./daemon -s reload
```

Stop the daemon cleanly:
```bash
./daemon -s stop
```
