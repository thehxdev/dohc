# dohc
dohc is a DoH server scanner. With this tool you can check wich DoH servers are working on your network.


## Download
Check [releases page](https://github.com/thehxdev/dohc/releases/latest) to download binary packages.


## Build

### Linux / macOS
```bash
CGO_ENABLED=0 go build -ldflags='-s -buildid=' .
```

### Windows
```powershell
$env:CGO_ENABLED=0
go build -ldflags='-s -buildid=' .
```

### Makefile
```bash
make
```

### Cross-Platform build
```bash
make cross-plat
```


## Usage
To print a help message:
```bash
./dohc -help
```

To start testing:
```bash
./dohc -f doh_servers.txt
```


## Contribution
If you can improve the source code or make this software better, feel free to send a PR :)
