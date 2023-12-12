# DI package

DI package with no-magic initializations.

## Usage

1. Import the package:
   ```
    import "github.com/kukymbr/core2go/di"`
   ``` 
   
2. Create a Builder and add dependencies definitions:
   ```go
   builder := &di.Builder{}
   
   err := builder.Add(di.Def{
        // Define the dependency name, 
        // must be unique for the container.
        Name: "dependency_name",
      
        // Define the build function, 
        // this is a mandatory field.
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
      
        // If true, dependency's build will be called
        // on the first dependency call, not on the build.
        Lazy: false,
   })
   if err != nil {
       panic(err)
   }

   ```
   
3. Build the Container:
   ```go
   ctn, err := builder.Build()
   if err != nil {
        panic(err)
   }
   ```

4. Call the dependency:
   ```go
       myObj := ctn.Get("dependency_name").(*MyObject)
       // do something with myObj
   ```