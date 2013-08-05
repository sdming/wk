// +build appengine

package fsw

import (
	"fmt"
	"log"
	"os"
	"path"
)

type Mode int

const (
	_           = iota
	Create Mode = 1 << (2 * iota)
	Modify
	Delete
	Rename
)

var Logger *log.Logger
var Debug bool = false

func debug(v ...interface{}) {
	if Debug && Logger != nil {
		Logger.Println(v)
	}
}

func fatal(v ...interface{}) {
	if Logger != nil {
		Logger.Println(v)
	}
}

// 
type Event struct {
	Name string
	Mode Mode
}

// String
func (e Event) String() string {
	s := e.Name
	if e.Mode|Create == Create {
		s = s + " Create"
	}
	if e.Mode|Modify == Modify {
		s = s + " Modify"
	}
	if e.Mode|Delete == Delete {
		s = s + " Delete"
	}
	if e.Mode|Rename == Rename {
		s = s + " Rename"
	}
	return s
}

// 
type NotifyFunc func(e Event)

// FsWatcher is file system watcher
type FsWatcher struct {
	Roots     []string
	Listeners []NotifyFunc
}

// String
func (w *FsWatcher) String() string {
	return fmt.Sprint(w.Roots)
}

// Listen register a event listener
func (w *FsWatcher) Listen(fn NotifyFunc) {
	w.Listeners = append(w.Listeners, fn)
}

// Watch add file or dir to watch list
func (w *FsWatcher) Watch(name string) error {
	return nil
}

// Close stop watch 
func (w *FsWatcher) Close() {

}

func (w *FsWatcher) notifyCallback() {

}

func NewFsWatcher(roots ...string) (*FsWatcher, error) {
	if len(roots) == 0 {
		return nil, fmt.Errorf("length of roots is 0")
	}

	roots = roots[:]
	for i := 0; i < len(roots); i++ {
		root := path.Clean(roots[i])
		_, err := os.Stat(root)
		if err != nil {
			return nil, err
		}
		roots[i] = root
	}

	w := &FsWatcher{
		Roots:     roots,
		Listeners: make([]NotifyFunc, 0),
	}

	return w, nil
}
