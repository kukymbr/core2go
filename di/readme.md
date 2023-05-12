## core2go di package

DI package with no-magic initializations.

### Usage

```go
package main

// Import the package
import "github.com/kukymbr/core2go/di"

func main() {
	
	// Create a builder instance
	builder := &di.Builder{}
	
	err := builder.Add(di.Def{
		// Define the dependency name, must be unique for the container
		Name: "dependency_name",
		
		// Define the build function
		Build: func(ctn *di.Container) (any, error) {
			return &MyObject{}, nil
		},
		
		// Validate something on add (optional)
		Validate: func(ctn *di.Container) (err error) {
			return nil
		},
		
		// Close dependency on the destruction
		Close: func(obj any) error {
			return obj.(*MyObject).Close()
		},
		
		// If true, dependency's build will be called on first get
		Lazy: false,
	    })
	
	if err != nil {
		panic(err)
    }
	
	// Build the container
	ctn, err := builder.Build()
	if err != nil {
		panic(err)
	}
	
	// Call the dependency
	myObj := ctn.Get("dependency_name").(*MyObject)
	// do something with myObj
}
```