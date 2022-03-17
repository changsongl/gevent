package gevent

import (
	"errors"
	"reflect"
	"sync"
)

type gevent struct {
	observers observers
}

type observers struct {
	selectorMap map[string]map[string]observer
	lock        *sync.RWMutex
}

type observer struct {
	selectorValue reflect.Value
	async         bool
}

type GEvent interface {
	TriggerEvent(eventName string, params ...interface{})
	AddObserver(eventName, observerName string, selector interface{}, async bool) error
	RemoveObserver(eventName, observerName string) bool
}

func (g *gevent) TriggerEvent(eventName string, params ...interface{}) {
	for _, obs := range g.observers.selectorMap[eventName] {
		var values []reflect.Value
		for _, param := range params {
			values = append(values, reflect.ValueOf(param))
		}

		obs.selectorValue.Call(values)
	}
}

func (g *gevent) AddObserver(eventName, observerName string, selector interface{}, async bool) error {
	selectorValue := reflect.ValueOf(selector)
	if selectorValue.Kind() != reflect.Func {
		return errors.New("selector is not a function")
	}

	g.observers.lock.Lock()
	defer g.observers.lock.Unlock()

	m := g.observers.selectorMap[eventName]
	if len(m) == 0 {
		g.observers.selectorMap[eventName] = map[string]observer{observerName: {async: async, selectorValue: selectorValue}}
	} else {
		g.observers.selectorMap[eventName][observerName] = observer{async: async, selectorValue: selectorValue}
	}

	return nil
}

func (observers observers) AddObserver(eventName, observerName string, selector reflect.Value, async bool) {
	observers.lock.Lock()
	defer observers.lock.Unlock()

	m := observers.selectorMap[eventName]
	if len(m) == 0 {
		observers.selectorMap[eventName] = map[string]observer{observerName: {async: async, selectorValue: selector}}
	} else {
		observers.selectorMap[eventName][observerName] = observer{async: async, selectorValue: selector}
	}
}

func (g *gevent) RemoveObserver(eventName, observerName string) bool {
	g.observers.lock.Lock()
	defer g.observers.lock.Unlock()

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

// NewGEvent creates a new gevent
func NewGEvent() GEvent {
	return &gevent{}
}
