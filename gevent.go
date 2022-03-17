package gevent

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// gevent implementation
type gevent struct {
	observers observers
	l         Log
}

// optFunc is a function implements the Option interface
type optFunc func(*gevent)

// Option is a option interface
type Option interface {
	apply(*gevent)
}

// apply applies option to gevent
func (f optFunc) apply(g *gevent) {
	f(g)
}

// NewLogOption creates and returns a log Option
func NewLogOption(l Log) Option {
	return optFunc(func(g *gevent) {
		g.l = l
	})
}

// observers implementation
type observers struct {
	selectorMap map[string]map[string]observer
	sync.RWMutex
}

// observer implementation
type observer struct {
	selectorValue reflect.Value
	async         bool
}

// GEvent is observer pattern interface
type GEvent interface {
	TriggerEvent(eventName string, params ...interface{})
	AddObserver(eventName, observerName string, selector interface{}, async bool) error
	RemoveObserver(eventName, observerName string) bool
}

// TriggerEvent trigger a event
func (g *gevent) TriggerEvent(eventName string, params ...interface{}) {
	if len(g.observers.selectorMap[eventName]) == 0 {
		return
	}

	var values []reflect.Value
	for _, param := range params {
		values = append(values, reflect.ValueOf(param))
	}

	for observerName, obs := range g.observers.selectorMap[eventName] {
		if obs.async {
			go obs.Call(observerName, values, g.l)
		} else {
			obs.Call(observerName, values, g.l)
		}
	}
}

// AddObserver adds observer
func (g *gevent) AddObserver(eventName, observerName string, selector interface{}, async bool) error {
	selectorValue := reflect.ValueOf(selector)
	if selectorValue.Kind() != reflect.Func {
		return errors.New("selector is not a function")
	}

	g.observers.Lock()
	defer g.observers.Unlock()

	m := g.observers.selectorMap[eventName]
	if len(m) == 0 {
		g.observers.selectorMap[eventName] = map[string]observer{observerName: {async: async, selectorValue: selectorValue}}
	} else {
		g.observers.selectorMap[eventName][observerName] = observer{async: async, selectorValue: selectorValue}
	}

	return nil
}

// RemoveObserver removes observer
func (g *gevent) RemoveObserver(eventName, observerName string) bool {
	g.observers.Lock()
	defer g.observers.Unlock()

	m := g.observers.selectorMap[eventName]
	if len(m) == 0 {
		return false
	}

	if _, ok := m[observerName]; !ok {
		return false
	}

	delete(m, observerName)
	return true
}

// AddObserver adds observer
func (observers *observers) AddObserver(eventName, observerName string, selector reflect.Value, async bool) {
	observers.Lock()
	defer observers.Unlock()

	m := observers.selectorMap[eventName]
	if len(m) == 0 {
		observers.selectorMap[eventName] = map[string]observer{observerName: {async: async, selectorValue: selector}}
	} else {
		observers.selectorMap[eventName][observerName] = observer{async: async, selectorValue: selector}
	}
}

// Call calls function
func (o observer) Call(observerName string, values []reflect.Value, log Log) {
	defer func() {
		if r := recover(); r != nil {
			if log != nil {
				log.Error(fmt.Sprintf("call failed by observer %s, %s", observerName, r))
			}
		}
	}()

	o.selectorValue.Call(values)
}

// NewGEvent creates a new gevent
func NewGEvent(options ...Option) GEvent {
	g := &gevent{
		observers: observers{
			selectorMap: make(map[string]map[string]observer),
		},
		l: newConsoleLogger(),
	}

	for _, opt := range options {
		opt.apply(g)
	}

	return g
}
