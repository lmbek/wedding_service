The only rules we have is: 
* Always Strict interface implementation pattern
* that everything should be seperated in different packages and kept as minimal as possible but with interface first implementation to streamline everything as the same code style
* That all code inside each package should be of this coding style pattern with dependency injection:

type App interface {
}

type app struct {
}

func NewApp() App {
    return &app{}
}

THIS IS THE MOST IMPORTANT RULE!
We MUST ALWAYS have an interface and build out from it

And we must never expose structs to the outside world.

* This is not allowed, we must always not make our if statements like this, only err != nil allowed:
if err := g.Init("gateway_service", c.DebugLevel()); err != nil {
