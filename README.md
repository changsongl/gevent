# gevent

GEvent is a golang package for observer pattern. 


The Observer pattern addresses the following problems:
- A one-to-many dependency between objects should be defined without making the objects tightly coupled.
- It should be ensured that when one object changes state, an open-ended number of dependent objects are updated automatically.
- It should be possible that one object can notify an open-ended number of other objects.


### Example:

```` go
// init.go
type test struct {
    data int
}

type testLog struct {
}

func (t testLog) Error(msg string) {
    myProgramLog.Error(msg)
}

var ge GEvent

func init() {
    ge = NewGEvent(NewLogOption(&testLog{}))
}

// account.go
func Register() {

    // some logic for register
	
    a := test{1}
    b := &test{2}
    ge.TriggerEvent("account-create", a, b)
}

// point.go
func init() {
    ge.AddObserver("account-create", "point", func(a test, b *test) {
       // do something for points
    }, true)
}

// user.go
func init() {
    ge.AddObserver("account-create", "user", func(a test, b *test) {
       // do something for user
    }, true)
}

````
