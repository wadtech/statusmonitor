# Go Status Monitor

Playing with high-level networking and concurrency, a nice break from the single-threaded.

## Usage
provide a services.json file that looks like this

```json
{
  "services": [
    {
      "description": "port 8080 on localhost"
      "port": "8080",
      "host": "localhost"
    }
  ]
}
```

then run `statusmonitor`

cli options

```bash
-port    numeric port value to serve the status page from.
```

## What?

statusmonitor will start serving a status page (by default) at localhost:8080 attempt to connect to the ports defined and report on success for each.

Statusmonitor speaks JSON, just send the correct request header `Accept: application/json`
