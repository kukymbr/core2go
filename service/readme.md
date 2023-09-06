## Service package

`Service` is a structure, combining parts of the `core2go` kit into single application.

### Service initialization

There are two ways to initialize the service: with default or custom container.

**Default container initialization:**

```go
package main

import (
	"github.com/kukymbr/core2go/service"
)

func main() {
	srv, err := service.NewWithDefaultContainer()
	if err != nil {
		panic(err)
    }
	
	if err := srv.Run(); err != nil {
		panic(err)
    }
}
```

Default service contains base context, config raw, config, logger and gin router runner.
It receives connections on the port defined in the config.

**Custom container initialization:**

```go
package main

import (
	"github.com/kukymbr/core2go/di"
	"github.com/kukymbr/core2go/service"
)

func main() {
	builder := &di.Builder{}
	if err := builder.Add(/* add your dependencies */); err != nil {
		panic(err)  
	}
	
	/* or:
	
	builder, err := service.GetDefaultDIBuilder()
	if err != nil {
	    panic(err)
	}

	runner := GetSomeRunner()
	
	if err := builder.Add(runner); err != nil {
		panic(err)
	}
	
	 */

	ctn, err := builder.Build()
	if err != nil {
		panic(err)
	}

	srv := service.New(ctn)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
```

Don't forget to register required service's dependencies: 
* `service.DIKeyBaseContext` of `service.ContextWithCancel`
* `service.DIKeyLogger` of `*zap.Logger`
* `service.DIKeyConfig` of `*service.Config`
* `service.DIKeyRunner` of `service.Runner`

### DI definitions

`Service` uses the dependency definitions with the keys, listed in the [di.go](di.go) file constants.
To redefine some of them, register the other instance (with a matched type) 
in your container with the corresponding key.

### Configuration

When using the default dependencies definitions, service config should be provided. 
Copy the [config.example.yml](config.example.yml) to the `config.yml` inside your project's root
or redefine the `service.DIKeyConfigRAW` DI item (returning the `*viper.Viper` instance)
to use another way to read the configuration values.