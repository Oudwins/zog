---
sidebar_position: 3
---

# zenv

`zenv` helps validate environment variables. Since `os.Getenv` does not perform any validation or type coercion this is a great way to ensure you didn't forget to set an environment variable. Which we have all done at some point....

```go
import (
  z "github.com/Oudwins/zog"
  "github.com/Oudwins/zog/zenv"
)

var envSchema = z.Struct(z.Schema{
	"PORT": z.Int().GT(1000).LT(65535).Default(3000),
	"DB": z.Struct(z.Schema{
		"Host": z.String().Default("localhost"),
		"User": z.String().Default("root"),
		"Pass": z.String().Default("root"),
	}),
})
var Env = struct {
	PORT int // zog will automatically coerce the PORT env to an int
	DB   struct {
		Host string `zog:"DB_HOST"` // we specify the zog tag to tell zog to parse the field from the DB_HOST environment variable
		User string `zog:"DB_USER"`
		Pass string `zog:"DB_PASS"`
	}
}{}

// Init our typesafe env vars, panic if any envs are missing
func Init() {
  errs := envSchema.Parse(zenv.NewDataProvider(), &Env)
  if errs != nil {
    log.Fatal(errs)
  }
}

// if you want to always panic on error
var Env = parse()
func Parse() env {
	var e env
	errs := envSchema.Parse(zenv.NewDataProvider(), &e)
	if errs != nil {
		fmt.Println("FAILURE TO PARSE ENV VARIABLES")
		log.Fatal(z.Issues.SanitizeMap(errs))
	}
	return e
}
```
