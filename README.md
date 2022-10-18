## Worker pool
For testing purpose there is a file `top-1m.csv` with list of domains.
It makes requests using working pool and logs response time with data size in kb.

You can configure `workerCount` and do `go run main.go`. It will produce logs in console with required data.
Job count depends on the size of domain list.
