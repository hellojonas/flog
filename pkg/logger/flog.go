package logger

import (
	"bufio"
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

const (
	LEVEL_INFO = iota
	LEVEL_WARN
	LEVEL_ERROR
)

type logger struct {
	rootDir  string
	baseDir  string
	fullpath string
	groups   map[string][]Entry
}

type Entry struct {
	Group   string    `json:"group"`
	Time    time.Time `json:"time"`
	Level   int       `json:"level"`
	Message string    `json:"message"`
}

func New(dir string) (*logger, error) {
	rootDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	basedir := time.Now().Format("2006_01_02")
	fullpath := path.Join(rootDir, basedir)
	if err = os.MkdirAll(fullpath, os.ModePerm); err != nil {
		return nil, err
	}
	return &logger{
		rootDir:  rootDir,
		baseDir:  basedir,
		fullpath: fullpath,
		groups:   make(map[string][]Entry),
	}, err
}

func (l *logger) flush() {
	l.createDir()
	wg := sync.WaitGroup{}
	for group, entries := range l.groups {
		wg.Add(1)
		go func() {
			if len(entries) == 0 {
				wg.Done()
				return
			}
			gf, err := l.groupFile(group)
			if err != nil {
				// TODO: log this on application level
				wg.Done()
				log.Println(err.Error())
				return
			}
			out := bufio.NewWriter(gf)
			for _, entry := range entries {
				enc := json.NewEncoder(out)
				// b, err := json.Marshal(entry)
				if err = enc.Encode(entry); err != nil {
					// TODO: log this on application level
					wg.Done()
					log.Println(err.Error())
					return
				}
			}
			if err = out.Flush(); err != nil {
				// TODO: log this on application level
				wg.Done()
				log.Println(err.Error())
			}
			if err = gf.Close(); err != nil {
				// TODO: log this on application level
				wg.Done()
				log.Println(err.Error())
			}
			l.groups[group] = make([]Entry, 0)
			wg.Done()
		}()
	}
	wg.Wait()
}

func (l *logger) Log(entries chan Entry, poll int) {
	ticker := time.NewTicker(time.Duration(poll) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			l.flush()
		case entry := <-entries:
			l.groups[entry.Group] = append(l.groups[entry.Group], entry)
		}
	}
}

func (l *logger) createDir() error {
	dir := time.Now().Format("2006_01_02")
	if l.baseDir != dir {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return err
		}
		l.baseDir = dir
		l.fullpath = path.Join(l.rootDir, dir)
	}
	return nil
}

func (l *logger) groupFile(group string) (*os.File, error) {
	fpath := path.Join(l.fullpath, group+".log")
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, fs.ModePerm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (l *logger) log(level int, msg string, args ...any) {}
