# Go Status Monitor

Playing with high-level networking and concurrency, a nice break from the single-threaded stuff you get to do with most web applications.

[![baby-gopher](https://raw2.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org)

## Usage
provide a config.json file that looks like this

```json
{
  "port": "8080",
  "workers": 10,
  "services": [
    {
      "description": "localhost on port 8080",
      "host": "127.0.0.1",
      "port": "8080"
    }
  ]
}
```

then run `statusmonitor`. It will look in the current working directory for the config file, though it's recommended to supply the path to the file yourself (see below)
You can

cli options

```bash
-c    config file location.
```

statusmonitor will start serving a status page (by default) at localhost:8080 attempt to connect to the ports defined and report on success for each.

Statusmonitor speaks JSON, just send the correct request header `Accept: application/json`
