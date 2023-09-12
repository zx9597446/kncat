# kncat

crypted netcat in golang. a network tunneling software working as an encryption wrapper between clients and servers (remote/local). It can pipe connections(remote/local) or cli app's stdin/stdout. It works interactive(shell or reverse shell) or batch mode.

# examples

1. send messages to server stdout.

   server:

   ```kncat -s secret_key```

   client:

   ```cat file.log | kncat -s secret_key -c svr:9597```

2. get a shell from server.

      server:

      ```kncat -s secret_key -e "/bin/bash"```

      client:

      ```kncat -s secret_key -c svr:9597```

3. get a reverse shell from client on server.

      server:

      ```kncat -s secret_key -r```

      client:

      ```kncat -s secret_key -r -e "cmd.exe" -c svr:9597```

4. pipe server's redis-cli to local.

      server:

      ```kncat -s secret_key -e "redis-cli"```

      client:

      ```kncat -s secret_key -c svr:9597```

5. pipe server's redis port 127.0.0.1:6379 to local 127.0.0.1:6379.

      server:

      ```kncat -s secret_key -f 127.0.0.1:6379```

      client:

      ```kncat -s secret_key -c svr:9597 -f 127.0.0.1:6379```


# installtion

 ``` go install github.com/zx9597446/kncat@latest ```

      or download from releases


# usage

```
Usage of kncat:
  -c string
        connect to address
  -e string
        program to execute (cmd.exe or /bin/bash or with args: cat -- some.log, use -- split args)
  -f string
        forward address(server: connect to this address. client: accept on this address)
  -l string
        listen on address (default ":9597")
  -m string
        crypto method (rc4 or aes256cfb) (default "aes256cfb")
  -n string
        network protocol: tcp tcp4 tcp6 (default "tcp")
  -r    reverse mode: connect and execute program on client side, to get a reverse shell
  -s string
        your secret key
  -v    verbose output
```

# reference & thanks:
1. https://github.com/jiguangsdf/netcat
2. https://github.com/getqujing/qtunnel