// package sessions


// type Session interface {
// 	Set(key, value interface{}) error //set session value
// 	Get(key interface{}) interface{}  //get session value
// 	Delete(key interface{}) error     //delete session value
// 	SessionID() string                //back current sessionID
// }

// var sessionStorage = 

// // Register makes a session provider available by the provided name.
// // If a Register is called twice with the same name or if the driver is nil,
// // it panics.
// func Register(name string, provider Provider) {
// 	if provider == nil {
// 		panic("session: Register provider is nil")
// 	}
// 	if _, dup := provides[name]; dup {
// 		panic("session: Register called twice for provider " + name)
// 	}
// 	provides[name] = provider
// }

