// +build !appengine

package fsw

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path"
	"path/filepath"
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

	quit   chan string
	notify *fsnotify.Watcher
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
	return w.notify.Watch(name)
}

// Close stop watch 
func (w *FsWatcher) Close() {
	w.quit <- "close"
	w.notify.Close()
}

func (w *FsWatcher) notifyCallback() {
	for {
		select {
		case e := <-w.notify.Event:
			debug("watcher notify", e)

			var fi os.FileInfo
			var err error
			var name = e.Name
			var mode Mode

			if e.IsCreate() {
				fi, err = os.Stat(name)
				if err == nil && fi.IsDir() {
					err = w.notify.Watch(name)
					debug("watcher add", name, err)
				}
				mode = mode | Create
			}
			if e.IsModify() {
				mode = mode | Modify
			}
			if e.IsDelete() {
				//remove ??
				//err = w.notify.Watch(name)
				//debug("watcher remove", name, err)
				mode = mode | Delete
			}
			if e.IsRename() {
				fi, err = os.Stat(name)
				if err == nil && fi.IsDir() {
					err = w.notify.Watch(name)
					debug("watcher add", name, err)
				}

				mode = mode | Rename
			}

			if err != nil {
				fatal("watcher error", err, name)
			}

			for i := 0; i < len(w.Listeners); i++ {
				w.Listeners[i](Event{Name: name, Mode: mode})
			}

		case err := <-w.notify.Error:
			fatal(err)
		case q := <-w.quit:
			debug("watcher quit", q)
			return
		}
	}
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

	notify, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &FsWatcher{
		Roots:     roots,
		Listeners: make([]NotifyFunc, 0),
		notify:    notify,
		quit:      make(chan string, 0),
	}

	go w.notifyCallback()

	for i := 0; i < len(roots); i++ {
		root := roots[i]

		if f, err := os.Stat(root); err != nil {
			return nil, fmt.Errorf("Stat %s fail: %v", root, err)
		} else if !f.IsDir() {
			err = w.notify.Watch(root)
			if err != nil {
				return nil, fmt.Errorf("watch %s fail: %v", root, err)
			}
			continue
		}

		err := filepath.Walk(root, func(p string, fi os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("walk %s fail: %v", p, err)
			}
			if fi.IsDir() {
				err = w.notify.Watch(p)
				if err != nil {
					return fmt.Errorf("watch %s fail: %v", p, err)
				}
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return w, nil

}
