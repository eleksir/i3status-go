package lib

import (
	"log"
	"reflect"
	"regexp"

	"go.i3wm.org/i3"
)

type WinList struct {
	Instance Collection
	Class    Collection
}

// WL contains collections of all windows currently running, collection of their classes and collection of their
// instances.
var WL WinList

// UpdateI3WinList subsribes to all window-related events and send them to I3EventParser() for more detailed parsing.
func (c *MyConfig) UpdateI3WinList() {
	// TODO: Здесь получить начальный список окон, управляемых i3wm (через GetTree()) и сложить их атрибутику в переменную WL
	Tree, err := i3.GetTree()

	if err != nil {
		log.Printf("Unable to get list of windows: %s", err)
	}

	c.ExtractProps(Tree.Root)

	i3wm := i3.Subscribe(i3.WindowEventType)

	for i3wm.Next() {
		e := i3wm.Event().(*i3.WindowEvent)
		go c.I3EventParser(e)
	}

	log.Fatal(i3wm.Close())
}

// I3EventParser parses i3wm events and updates WL (aka Windows List).
func (c *MyConfig) I3EventParser(e *i3.WindowEvent) {
	switch e.Change {
	case "new":
		var (
			count uint64
		)

		countI, exist := WL.Class.Get(e.Container.WindowProperties.Class)

		if exist {
			count = interfaceToUint64(countI)
		} else {
			count = 0
		}

		count++
		WL.Class.Set(e.Container.WindowProperties.Class, count)

		countI, exist = WL.Instance.Get(e.Container.WindowProperties.Instance)

		if exist {
			count = interfaceToUint64(countI)
		} else {
			count = 0
		}

		count++
		WL.Instance.Set(e.Container.WindowProperties.Instance, count)

		c.Channels.UpdateReady <- true

	case "close":
		countI, exist := WL.Class.Get(e.Container.WindowProperties.Class)

		if exist {
			count := interfaceToUint64(countI)

			if count > 1 {
				count--

				WL.Class.Set(e.Container.WindowProperties.Class, count)
			} else {
				WL.Class.Delete(e.Container.WindowProperties.Class)
			}
		}

		countI, exist = WL.Instance.Get(e.Container.WindowProperties.Instance)

		if exist {
			count := interfaceToUint64(countI)

			if count > 1 {
				count--

				WL.Instance.Set(e.Container.WindowProperties.Instance, count)
			} else {
				WL.Instance.Delete(e.Container.WindowProperties.Instance)
			}
		}

		c.Channels.UpdateReady <- true
	}
}

// interfaceToUint превращает данный интерфейс в Uint64.
// Если может, конечно :) .
func interfaceToUint64(iface interface{}) uint64 {
	// А теперь мы начинаем дурдом, нам надо превратить ёбанный interface{} в []string
	// Поскольку interface{} может быть чем угодно, перестрахуемся
	if reflect.TypeOf(iface).Kind() == reflect.Uint64 {
		return reflect.ValueOf(iface).Uint()
	}

	return 0
}

// HasWindows returns true if given Window Class and/or Instance has more than 0 windows according our observations.
func HasWindows(className string, instanceName string) bool {
	var (
		ret        bool
		myclass    bool
		myinstance bool
	)

	switch {
	case className != "" && instanceName != "":
		// Since we agree that we use regexp for searching window class and instance, but not yet decide what to do
		// with exact match, comment it out and leave it here.
		/*
			_, myclass = WL.Class.Get(className)
			_, myinstance = WL.Instance.Get(instanceName)
		*/
		myclass = FindWindowClass(className)
		myinstance = FindWindowInstance(instanceName)
		ret = myclass && myinstance

	case className != "":
		// _, ret = WL.Class.Get(className)
		ret = FindWindowClass(className)

	case instanceName != "":
		// _, ret = WL.Instance.Get(instanceName)
		ret = FindWindowInstance(instanceName)
	}

	return ret
}

// FindWindowInstance returns true if given regex matces at least one pattern in WL.Class collection.
// Otherwise, returns false.
func FindWindowInstance(regexpStr string) bool {
	var found = false

	f := func(key interface{}, value interface{}) bool {
		keyStr := reflect.ValueOf(key).String()

		re, err := regexp.Compile(regexpStr) //nolint:nolintlint, gocritic, we do not want shit our pants and panic on
		//        erroneous regexp, we want log it to stderr, that it.

		if err != nil {
			log.Printf("Unable to perform regexp match: %s", err)

			return true
		}

		if re.Match([]byte(keyStr)) { //nolint:mirror
			found = true

			return false
		}

		return true
	}

	WL.Instance.Range(f)

	return found
}

// FindWindowClass returns true if given regex matces at least one pattern in WL.Class collection.
// Otherwise, returns false.
func FindWindowClass(regexpStr string) bool {
	var found = false

	f := func(key interface{}, value interface{}) bool {
		keyStr := reflect.ValueOf(key).String()

		re, err := regexp.Compile(regexpStr) //nolint:nolintlint, gocritic, we do not want shit our pants and panic on
		//        erroneous regexp, we want log it to stderr, that it.

		if err != nil {
			log.Printf("Unable to perform regexp match: %s", err)

			return true
		}

		if re.Match([]byte(keyStr)) { //nolint:mirror
			found = true

			return false
		}

		return true
	}

	WL.Class.Range(f)

	return found
}

// ExtractProps fills WL stricture collections with instance name and class name of given node (x11, that is).
func (c *MyConfig) ExtractProps(n *i3.Node) {
	var e i3.WindowEvent
	e.Change = "new"
	e.Container = *n

	c.I3EventParser(&e)

	for _, node := range n.Nodes {
		c.ExtractProps(node)
	}

	for _, node := range n.FloatingNodes {
		c.ExtractProps(node)
	}
}
