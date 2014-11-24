monitor
=======

a small package to help with monitoring services.

### Statter

Interface for counting / measuring statistics

#### Implementations

##### Stat Hat

###### configure
`SH_KEY`:(optional) Stat Hat key to log either counts / values against;

###### example usage
```Go
	// Using key defined by the environment variable 'SH_KEY'
	stats.Count("endpoint /ping hit", 10)
	stats.Measure("endpoint /ping responded within (ms)", 0.12)

	// Using a key defined via construction of a new StatHat
	sth := stats.NewStatHat("new-key")
	sth.Count("endpoint /pong hit", 5)
	sth.Measure("endpoint /pong responded within (ms)", 0.01)
```
