### Mayonaise

Using websockets to pass linux commands between a server and a client


- Using websockets to open a shell instead of using ssh


- The client sends unix commands via the sockets and the server runs the commands via os.Exec. The stdout is returned to the clinet and print as stdout
