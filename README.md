# kncat
crypted netcat in golang

```
usage: netcat [-l] [-v] [-p port] [-n tcp]
options:
  -e    shell mode
  -h string
        host addr to connect or listen (default "0.0.0.0")
  -help
        print this help
  -k string
        secret key to crypt
  -l    listen mode
  -m string
        crypto method: rc4 aes256cfb (default "rc4")
  -n string
        network protocol (default "tcp")
  -p int
        host port to connect or listen (default 4000)
  -v    verbose mode (default true)
```

# reference:
1. https://github.com/jiguangsdf/netcat
2. https://github.com/getqujing/qtunnel