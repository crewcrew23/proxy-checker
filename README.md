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
| Flag       | Description                                                                                                  | Required |
|------------|--------------------------------------------------------------------------------------------------------------|----------|
| `--input`  | Path to file with proxy list (one per line)                                                                  | ✅        |
| `--type`   | Proxy type: `http` or `socks5` Default http:(all proxies in the file must be of the same type)                            | ❌        |
| `--target` | URL of the resource through which the availability of the proxy is checked                                   | ✅        |
| `--timeout`| Connection timeout in seconds (default 5)                                                                    | ❌        |
| `--save`   | File for saving working proxies (in CSV format)                                             | ❌        |
| `--threshold`   | threshold of the number of proxies in the list, upon reaching which the worker pool will be used for processing (default 100)                                            | ❌        |

```bash
./proxy-checker --input <file_with_proxy> --type <proxy_type> --target <target_URL> --timeout <second> --save <output_file>
```

### Example
```bash
./proxy-checker --input proxies-socks5.txt --type socks5 --target https://www.google.com --timeout 5 --save good-socks5.csv
```


## Input file format:
proxy without auth <br>
``` host:port ```
``` protocol://host:port ```

proxy with auth <br>
``` username:password@host:port ```
``` protocol://username:password@host:port ```

### Example
```
127.0.0.1:8888
user:pass@127.0.0.1:8889
http://user:pass@127.0.0.1:8889
```
## Testing
```make test-unit```<br>
run unit tests <br>
``` make test-e2e ``` <br>
will run the docker-compose with proxy <br>
if all good, you will see like that
```
🔍 Checking 3 proxies...
✅ http://user:pass@127.0.0.1:8889 [618.2278ms]
✅ user:pass@127.0.0.1:8889 [618.2278ms]
✅ 127.0.0.1:8888 [618.2278ms]

Total: 3, Alive: 3
✅ Saved good proxies to test_bin\good-http.csv
🔍 Checking 3 proxies...
✅ socks5://test:secret@127.0.0.1:1081 [643.6372ms]
✅ 127.0.0.1:1080 [643.636ms]
✅ test:secret@127.0.0.1:1081 [642.626ms]

Total: 3, Alive: 3
✅ Saved good proxies to test_bin\good-socks5.csv
```