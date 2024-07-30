# dohc
dohc is a DoH server scanner. With this tool you can check wich DoH servers are working on your network.


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


## Source of DoH servers data

### Scraping cURL DoH wiki
You can use [scrape_doh_providers.py](misc/scrape_doh_providers.py) python3 script to scrape [cURL DoH wiki](https://github.com/curl/curl/wiki/DNS-over-HTTPS)
to make a list of DoH resolvers. This method is recommended and most up-to-date.
```bash
python3 ./misc/scrape_doh_providers.py '"{}".format(o["url"])' > ./doh_servers.txt
```
I took the python script from [HERE](https://gist.github.com/kimbo/dd65d539970e3a28a10628f15398247b).


### Using resolvers data
> [!NOTE]
> This method also uses cURL DoH wiki but it may be outdated or old.

I took DoH resolvers data ([This file](misc/doh_resolvers_data_20240119.json)) from [encrypted-dns-resolvers](https://github.com/cslev/encrypted_dns_resolvers) repository.
Then I extracted URIs with [converter.sh](misc/converter.sh) shell script. The final file (`doh_servers.txt`) is a list of DoH resolvers each on a seperate line.
```bash
./misc/converter.sh ./misc/doh_resolvers_data_20240119.json ./doh_servers.txt
```


## Contribution
If you can improve the source code or make this software better, feel free to send a PR :)
