package main

/*
TODO:

[ ] Secure the connection between client and a server using TLS. Issue the certificate. And do the authentication.
[ ] Add commands:
   [x] mv -> move file to a different folder.
   [x] pwd -> displays current working directory.
   [ ] send -> (sends file from server to client). Example: :send myfile.py <server_destination>
   [ ] :ftp command should except two arguments, filepath on the server, and destrination on the client.
       Example: :ftp <filepath> <destpath>

[ ] Implement :tree command (our own using ascii).

[x] Improve :du command.

[ ] Polish :rm command (maybe redesign).

[ ] Go through all the commands and their functions and make sure they work as expected.
	Apply modifications if needed.
	Things to take into consideration: Correctness of error handling.

[ ] On connection every client should be authenticated.

[ ] :ftp and :send commands should use workers in order to speed up file transfer, analogous to P2P connection.
    Each worker should upload chunk of the file and notify client that file transfer is finished.
	All the work between workers should be synchronized using channels.
	When the last chunk have been delivered to client, it should receive a notification
	that file has been transferred successfully.

[ ] All the operations on a server should be protected against race conditions.
    Make the server read-only?

[ ] Implement logging system. All invalid operations should be logged into a file.
    Explore zap-logger: https://github.com/uber-go/zap

[ ] Implement command capturing. When the connection established (session),
	all the commands executed on a server in this particular session should be captured
	into a command registry, so they can be reproduced later.

	For example: :history command might display all the commands that were issues during this session.
	Each entry in a history should contain (name of the command, its arguments, time when it was executed on a server).

	One design decision might be (after establising the connection):

	beingCommandCapture(...)
	// All the commands executed in between will be captured.
	// :ls
	// :
	endCommandCapture(...)

	All the commands should be written into database using client's credentials and session's id.
	We could consider in memory database (Redis), or even go for more heavier solution like mysql
	(would be a great learning opportunity.)

[ ] Implement propper command parsing.

[ ] List all connected clients. Maybe with :clients/peers command.
    Investigate whether it will be possible to exchange the data between clients, aka P2P.
	We would probably have to come with some command for that.
	Something like:
	:connect <client_name> or
	:p2p <client> -> should establish TCP connection with <client>

[ ] Write benchmarks for :ftp command (with workers disabled/enabled).
	That will solidify our knowledge on how to write benchmarks in GO.

[ ] Support both TCP and UDP connections.

[ ] Make sure that the data is encrypted, to avoid man-in-the-middle attacks.
*/
