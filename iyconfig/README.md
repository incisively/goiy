Config Package
=====================

Parse JSON encoded configurations using either the `config.Unmarshal` function or a `config.Decoder`.

Given the following JSON data.
```json
{
	"test": {
		"a": "some kind of string",
		"b": 100
	},
	"production": {
		"a": "production worthy string",
		"b": 9001
	}
}
```

The following can be done to parse this configuration file.
```go
type Config struct {
    A string `json:"a"`
    B int    `json:"b"`
}

var conf Config
if err := config.Unmarshal(jsondata, &conf, "test"); err != nil {
    panic(err)
}
```

This will result in the following state of Config, given `jsondata` is populated with the example data above.

```go
Config{
    A: "some kind of string",
    B: 100,
}
```

Alternatively the `config.Decoder` can be used when parsing configuration from a file. If the example data were stored in the file `./config.json`. It could be decoded in the following way.
```go
type Config struct {
    A string `json:"a"`
    B int    `json:"b"`
}

fi, err := os.Open("./config.json")
if err != nil {
    panic(err)
}

var conf Config
if err := config.NewDecoder(fi).Decode(&conf, "production"); err != nil {
    panic(err)
}
```

This will result in the following state of the Config struct.
```go
Config{
    A: "production worthy string",
    B: 9001,
}
```

Both `config.Unmarshal()` and the `config.Decoder` will parse the entire json stream/byte slice, without aggregating by an environment if the `env` string is not provided. 
