## Worker pool
For testing purpose there is a file `top-1m.csv` for looping domain and making GET request.
Additionally it logs response time and data size in kb.

You can configure `workerCount` and do `go run main.go`. 
Job count depends on the size of domain list.
