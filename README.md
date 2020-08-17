# ipInfo

An ipinfo.io clone, without the rate limiting. And some other goodies.

## Preamble

ipinfo.io sets a rate limiting of 1000 requests per day. I understand it,
although that is a bit of a bummer.

## Usage

This is really similiar to the way ipinfo.io does it.

### Basic usage

Make a GET request to `http://localhost/[ip]`.

```json
$ curl "http://localhost/8.8.8.8"
{
  "ip": "8.8.8.8",
  "city": "Mountain View",
  "region": "California",
  "country": {
    "code": "US",
    "name": "United States"
  },
  "continent": {
    "code": "NA",
    "name": "North America"
  },
  "location": {
    "latitude": 37.386,
    "longitude": -122.0838
  },
  "postal":"94040",
  "asn":15169,
  "organization": "Google LLC"
}
```

Much better! :smile:

You can also get returned your information by just calling `/`.

```sh
$ curl "http://localhost/"
```

We're are not done yet! You want to use JSONP. You guessed it, just provide a
`callback` parameter to your GET request.

```javascript
$ curl "localhost/8.8.8.8?pretty=1&callback=myFancyFunction"
/**/ typeof myFancyFunction === 'function' && myFancyFunction({
  "ip": "8.8.8.8",
  "city": "Mountain View",
  "region": "California",
  "country": {
    "code": "US",
    "name": "United States"
  },
  "continent": {
    "code": "NA",
    "name": "North America"
  },
  "location": {
    "latitude": 37.386,
    "longitude": -122.0838
  },
  "postal":"94040",
  "asn": 15169,
  "organization": "Google LLC"
});
```

```html
<script>
var myFancyFunction = function(data) {
  alert("The city of the IP address 8.8.8.8 is: " + data.city);
}
</script>
<script src="http://localhost/8.8.8.8?callback=myFancyFunction"></script>
```

## Differences from ipinfo.io

### Features we have, that ipinfo.io does not

* JSON minified, so it gets to your server quicker.
* Full name for the country!
* We also got continent info, with the full name too.

### Some advantages over ipinfo.io

* We are using Go and not nodejs like them. Go is a compiled language, and
  therefore is [amazingly fast](docs/benchmarks.md). A response can be generated
  in a very short time.
* We get data only from one data source. Which means no lookups on other
  databases, which results in being faster overall.
* We are open source. Which means you can compile and put it on your own
  server!

### Features that are not here (and not going to be implemented)

* Hostname.  We would have to pick that data from another data source, out of
  scope.

## Contributing

Feel free to open an issue or pull request for anything! If you want to run it
locally for whatever reason, you can do so this way if you don't need to touch
the code:

```sh
go get -d http://github.com/jnovack/ipinfo
cd $GOPATH/src/github.com/jnovack/ipinfo
go build
./ipinfo # .exe if you're on windows
```

You can not just do `go get` and then execute it.  The software requires
`GeoLite2-City.mmdb`, please see the [Data Source](#data-source) section.

If you want to hack in the future, this is a better way:

```sh
cd $GOPATH
mkdir -p src/github.com/jnovack
cd src/github.com/jnovack
git clone git@github.com:jnovack/ipinfo.git
cd ipinfo
go build
./ipinfo
# Or if you don't want to create the binary in the folder
go run cmd/ipinfo/main.go
```

## Data Source

This product uses the GeoLite2 data created by MaxMind, available from
http://www.maxmind.com.

## Credits

Inspired by [ip.zxq.co](http://ip.zxq.co/)
([source](https://github.com/thehowl/ip.zxq.co)).