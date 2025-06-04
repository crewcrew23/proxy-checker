# Proxy Checker
![GitHub](https://img.shields.io/badge/Go-1.24.2+-blue)

Go utility for checking proxy server availability.

## Description
The program receives a list of proxy servers and checks their availability by trying to connect to the specified target resource through each proxy.

##  Install

```bash
  git clone https://github.com/crewcrew23/proxy-checker.git
  cd proxy-checker
  make build # or go build -o bin/proxy-checker cmd/app/main.go

  cd bin #or set proxy-checker/bin in Env Var
```

## Usage

### Parameters
| Флаг      | Описание                                                                 |
|-----------|--------------------------------------------------------------------------|
| `-input`  | Path to file with proxy list (one per line)                   |
| `-type`   | Proxy type: `http` or `socks5` (all proxies in the file must be of the same type) |
| `-target` | URL of the resource through which the availability of the proxy is checked              |
| `-timeout`| Connection timeout in seconds (recomends 5)                         |
| `-save`   | File for saving working proxies (in CSV format)                      |

```bash
./proxy-checker -input <file_with_proxy> -type <proxy_type> -target <target_URL> -timeout <second> -save <output_file>
```

### Example
```bash
./proxy-checker -input proxies-socks5.txt -type socks5 -target https://www.google.com -timeout 5 -save good-socks5.csv
```


## Input file format:
proxy without auth <br>
``` host:port ```

proxy with auth <br>
``` host:port:username:password ```

### Example
```
127.0.0.1:1080
127.0.0.1:1081:user:pass
proxy.example.com:8888
```

## Testing
Run <br>
``` make test ``` <br>
will run the docker-compose with proxy <br>
if all good, you will see like that
```
🔍 Checking 2 proxies...
✅ 127.0.0.1:8888 [148.1436ms]
✅ 127.0.0.1:8889:user:pass [174.8869ms]

Total: 2, Alive: 2
✅ Saved good proxies to test_bin\good-http.csv
🔍 Checking 2 proxies...
✅ 127.0.0.1:1080 [127.5531ms]
✅ 127.0.0.1:1081:test:secret [129.812ms]

Total: 2, Alive: 2
✅ Saved good proxies to test_bin\good-socks5.csv
```