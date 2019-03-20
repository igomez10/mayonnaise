### Mayonaise

#### UML


- [class diagram](https://drive.google.com/file/d/1wuLOeAyHpikEeDS1eUIts6C0kMiApF31/view?usp=sharing)

- [client server interaction](https://drive.google.com/file/d/1Rp7O8UJIbLV_HGRlFTIDxu93Qr4aytoo/view?usp=sharing)


Using websockets to pass linux commands between a server and a client


- Using websockets to open a shell instead of using ssh


- The client sends unix commands via the sockets and the server runs the commands via os.Exec. The stdout is returned to the clinet and print as stdout

- The server should input the commands to be run in multiple clients. Follow a master -> slave relationship where the slave returns the output of the command sent by the master

- Read commands from stdinput in master and broadcast commands to slaves
