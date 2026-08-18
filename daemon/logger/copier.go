package logger

import (
	"bufio"
	"io"
	"sync"
	"time"
)

// Message is the structure for log messages
type Message struct {
	Source    string
	Line      []byte
	Timestamp time.Time
}

// Logger is the interface for loggers
type Logger interface {
	Log(*Message) error
}

// Copier can copy logs from stdout/stderr to a logger
type Copier struct {
	src      map[string]io.Reader
	dst      Logger
	copyJobs sync.WaitGroup
}

// NewCopier creates a new Copier
func NewCopier(src map[string]io.Reader, dst Logger) *Copier {
	return &Copier{
		src: src,
		dst: dst,
	}
}

// Run starts the copier and blocks until all sources are copied
func (c *Copier) Run() {
	for name, src := range c.src {
		c.copyJobs.Add(1)
		go c.copySrc(name, src)
	}
	c.copyJobs.Wait()
}

// copySrc copies a single source to the destination logger
func (c *Copier) copySrc(name string, src io.Reader) {
	defer c.copyJobs.Done()
	reader := bufio.NewReader(src)
	var line []byte
	for {
		l, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				// log error or handle it
			}
			if len(l) > 0 {
				line = append(line, l...)
				c.dst.Log(&Message{
					Source:    name,
					Line:      line,
					Timestamp: time.Now().UTC(),
				})
			}
			break
		}
		line = append(line, l...)
		if !isPrefix {
			c.dst.Log(&Message{
				Source:    name,
				Line:      line,
				Timestamp: time.Now().UTC(),
			})
			line = nil
		}
	}
}
