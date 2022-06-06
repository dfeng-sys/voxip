# voxip
AVOXI code challenge

Given an IP address and a country whitelist, returns the IP address's country and whether or not it is whitelisted. Uses a .mmdb for lookup, downloaded from GeoLite2.

Run with `go run main.go`.

Send a `GET` request to `localhost:9001/api/v1/ip` with the request body format:

    {
        "ip": "8.8.8.8",
        "whitelist": [
            "US",
            "GB",
            "AU",
             ...
            "DE",
            "JP",
            "CN"
        ]
    }

The API expects the country whitelist to be a string array of ISO 3166-1 alpha-2 codes.

#### Keeping mapping data updated

Not fully implemented, but one way to do this would be to create a second API endpoint. This endpoint would accept a URL from which to download an updated .mmdb file. The new file can then be verified by computing a checksum before replacing the existing file.

#### Features to add

1. Add validation of countries in 'whitelist' field.
2. Allow other formats of country codes/names.
3. Automated updating of mapping data.
