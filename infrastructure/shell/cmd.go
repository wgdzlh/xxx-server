package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "xxx-server/application/logger"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"

	"go.uber.org/zap"
)

type SubCmd struct {
	sync.Mutex
	cmd       string
	script    string
	workers   []*subWorker
	workerIdx int
	workerSem chan struct{}
	timeout   time.Duration
	logTag    string
}

type subWorker struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	// stderr *bufio.Reader
}

const (
	Sub_CMD_LOG_TAG    = "SubCmd::"
	PYTHON_CMD_LOG_TAG = "PythonCmd::"
	ErrorSigChar       = '!'

	CMD_SUBMIT_TIMEOUT = time.Second * 5
)

var (
	enterBs = []byte{'\n'}
)

func init() {
	signal.Ignore(syscall.SIGCHLD) // 避免子进程异常退出时变zombie（父进程没有wait）
}

func NewCmd(cmd string) repo.SubCmdRepo {
	p := &SubCmd{
		cmd:    cmd,
		logTag: Sub_CMD_LOG_TAG,
	}
	return p
}

func NewPythonCmd(script string, qSize ...int) repo.PythonCmdRepo {
	p := &SubCmd{
		cmd:    "python",
		script: script,
		logTag: PYTHON_CMD_LOG_TAG,
	}
	if len(qSize) > 0 && qSize[0] > 0 {
		if err := p.initWorkers(qSize[0]); err != nil {
			log.Fatal(p.logTag+"init workers err", zap.Error(err))
		}
	}
	return p
}

func (p *SubCmd) newSubWorker(openIn bool, args ...string) (w *subWorker, err error) {
	if p.script != "" {
		args = append([]string{p.script}, args...)
	}
	w = &subWorker{
		cmd: exec.Command(p.cmd, args...),
	}
	if openIn {
		if w.stdin, err = w.cmd.StdinPipe(); err != nil {
			return
		}
	}
	stdout, err := w.cmd.StdoutPipe()
	if err != nil {
		return
	}
	w.cmd.Stderr = os.Stdout
	w.stdout = bufio.NewReader(stdout)

	err = w.cmd.Start()
	return
}

func (p *SubCmd) initWorkers(size int) (err error) {
	p.workers = make([]*subWorker, size)
	p.workerSem = make(chan struct{}, size)
	// if config.C.Py.SubmitTimeout > 0 {
	// 	p.timeout = time.Second * time.Duration(config.C.Py.SubmitTimeout)
	// } else {
	p.timeout = CMD_SUBMIT_TIMEOUT
	// }
	var w *subWorker
	for i := 0; i < size; i++ {
		if w, err = p.newSubWorker(true); err != nil {
			return
		}
		p.workers[i] = w
	}
	p.workerIdx = size - 1
	return
}

func (p *SubCmd) cleanInput(in []byte) []byte {
	if bytes.Contains(in, enterBs) {
		return bytes.ReplaceAll(in, enterBs, nil)
	}
	return in
}

func (p *SubCmd) Exec(input entity.AnyJson, args ...string) (out []byte, err error) {
	openIn := len(input) > 0
	w, err := p.newSubWorker(openIn, args...)
	if err != nil {
		log.Error(p.logTag+"subprocess start error", zap.Error(err))
		return
	}
	defer w.cmd.Wait()
	log.Info(p.logTag+"subprocess start succeed", zap.String("script", p.script), zap.Any("args", args))
	if openIn {
		w.stdin.Write(p.cleanInput(input))
		w.stdin.Write(enterBs)
		w.stdin.Close()
	}
	out, err = io.ReadAll(w.stdout)
	if err != nil {
		return
	}
	if len(out) == 0 || out[0] == ErrorSigChar {
		err = fmt.Errorf("err in %s: %s", p.cmd, out)
	}
	return
}

func (p *SubCmd) Submit(input entity.AnyJson) (out entity.AnyJson, err error) {
	if p.workers == nil { // 当未初始化备用worker队列时，调用单次命令逻辑
		return p.Exec(input)
	}
	var w *subWorker
	select {
	case p.workerSem <- struct{}{}:
		w = p.popWorker()
		defer func() {
			p.pushWorker(w)
			<-p.workerSem
		}()
	case <-time.After(p.timeout):
		err = repo.ErrPyCmdSubmitTimeout
		return
	}
	if err = w.cmd.Process.Signal(syscall.Signal(0)); err != nil { // 检查子进程是存在
		log.Warn(p.logTag+"subprocess invalid", zap.Error(err))
		w.cmd.Wait() // 子进程已退出，此处只为关闭fd资源
		if w, err = p.newSubWorker(true); err != nil {
			return
		}
	}
	if _, err = w.stdin.Write(p.cleanInput(input)); err != nil {
		return
	}
	if _, err = w.stdin.Write(enterBs); err != nil {
		return
	}
	if out, err = readLine(w.stdout); err != nil {
		return
	}
	if len(out) == 0 || out[0] == ErrorSigChar {
		err = fmt.Errorf("err in %s: %s", p.cmd, out)
		out = nil
		log.Error("operation failed", zap.Error(err), zap.Any("input", input))
	}
	return
}

func (p *SubCmd) popWorker() (w *subWorker) {
	if len(p.workers) == 1 {
		return p.workers[0]
	}
	p.Lock()
	w = p.workers[p.workerIdx]
	p.workerIdx--
	p.Unlock()
	return
}

func (p *SubCmd) pushWorker(w *subWorker) {
	if len(p.workers) == 1 {
		return
	}
	p.Lock()
	p.workerIdx++
	p.workers[p.workerIdx] = w
	p.Unlock()
}

func readLine(br *bufio.Reader) (out []byte, err error) {
	var (
		line     []byte
		isPrefix bool
	)
	for {
		line, isPrefix, err = br.ReadLine()
		if err != nil {
			return
		}
		out = append(out, line...)
		if !isPrefix {
			return
		}
	}
}
