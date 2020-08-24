package spinner

import (
	"fmt"
	"io"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

type Spinner struct {
	sync.Mutex
	Title     string
	Charset   []string
	Framerate time.Duration
	runchan   chan struct{}
	stoponce  sync.Once
	Output    io.Writer
	NoTty     bool
	IsDone    chan bool
}

var DefaultCharset = []string{"|", "/", "-", "\\"}

const (
	// 150ms per frame
	DEFAULT_FRAME_RATE = time.Millisecond * 150
)

func NewSpinner(title string) *Spinner {
	sp := &Spinner{
		Title:     title,
		Charset:   DefaultCharset,
		Framerate: DEFAULT_FRAME_RATE,
		runchan:   make(chan struct{}),
		IsDone:    make(chan bool),
	}
	if !terminal.IsTerminal(syscall.Stdout) {
		sp.NoTty = true
	}
	return sp
}

func (sp *Spinner) animate() {
	var out string
	for i := 0; i < len(sp.Charset); i++ {
		out = sp.Charset[i] + " " + sp.Title
		switch {
		case sp.Output != nil:
			fmt.Fprint(sp.Output, out)
		case !sp.NoTty:
			fmt.Print(out)
		}
		time.Sleep(sp.Framerate)
		sp.clearLine()
	}
}

func (sp *Spinner) clearLine() {
	if !sp.NoTty {
		fmt.Printf("\033[2K")
		fmt.Println()
		fmt.Printf("\033[1A")
	}
}

func (sp *Spinner) writer() {
	sp.animate()
	for {
		select {
		case <-sp.runchan:
			sp.IsDone <- true
			return
		default:
			sp.animate()
		}
	}
}

func (sp *Spinner) Stop() {
	//prevent multiple calls
	sp.stoponce.Do(func() {
		close(sp.runchan)
		sp.clearLine()
	})
}

func (sp *Spinner) Start() *Spinner {
	go sp.writer()
	return sp
}

func StartNew(title string) *Spinner {
	return NewSpinner(title).Start()
}

func (sp *Spinner) SetSpeed(rate time.Duration) *Spinner {
	sp.Lock()
	sp.Framerate = rate
	sp.Unlock()
	return sp
}

func (sp *Spinner) SetCharset(charset []string) *Spinner {
	sp.Lock()
	sp.Charset = charset
	sp.Unlock()
	return sp
}

// var rcv make(chan string)
// var snd make(chan string)

//Goroutine
// // rcv <- "sending message" in other thread
// snd <- rcv
// var recievedInteger <- snd

// Normal
// // rcv = "sending message" in other thread
// snd = rcv
// var recievedInteger = snd

// rcv = "sending message" //in other thread
// snd = rcv

// snd = rcv
// rcv = "sending message" //in other thread
