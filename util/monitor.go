package util

import (
    "log"
    "os"
    "os/exec"
    "strings"
    "time"
)

func (w *GoMonitor) Monitor() {
    c := time.Tick(time.Second * time.Duration(w.Interval))
    for {
        select {
        case <-w.change: // file has been change
            w.BuildAndRun()
        case <-c: // walk every 5 seconds
            w.WalkFile()
        }
    }
}

func (w *GoMonitor) BuildAndRun() {
    // if the process is running ,kill it before build
    if w.cmd != nil && w.cmd.Process != nil {
        log.Printf("the process :%d\n", w.cmd.Process.Pid)
        if err := w.cmd.Process.Kill(); err != nil {
            log.Printf("%s %s\n", w.cmd.ProcessState.String(), err)
        } else {
            log.Printf("%s\n", w.cmd.ProcessState.String())
        }
    }
    if err := w.Build(); err != nil {
        log.Printf("build fail %s\n", err)
    } else {
        log.Printf("start process :%s\n", w.RunCmd)
        go w.Run()
    }
}

func (w *GoMonitor) Run() {
    args := strings.Split(w.RunCmd, " ")
    w.cmd = exec.Command(args[0], args[1:]...)
    w.cmd.Stdin = os.Stdin
    w.cmd.Stdout = os.Stdout
    w.cmd.Stderr = os.Stderr
    err := w.cmd.Run()
    if err != nil {
        log.Printf("%s\n", err)
    }
    w.cmd = nil
}

func (w *GoMonitor) Build() (err error) {
    args := strings.Split(w.BuildCmd, " ")
    w.cmd = exec.Command(args[0], args[1:]...)
    log.Printf("[build cmd] %s", w.BuildCmd)
    w.cmd.Stdin = os.Stdin
    w.cmd.Stdout = os.Stdout
    w.cmd.Stderr = os.Stderr
    err = w.cmd.Run()
    if err != nil {
        log.Printf("%s\n", err)
        return
    }
    log.Printf("build sucess\n")

    return
}
func (w *GoMonitor) WalkFile() {
    var change bool
    for file, modtime := range w.FileStatus {
        info, err := os.Stat(file)
        if err != nil {
            log.Printf("filename err :%s", err.Error())
        }
        newModTime := info.ModTime().Unix()
        if modtime != newModTime {
            log.Printf("File :%s has been changed", file)
            w.FileStatus[file] = newModTime
            change = true
        }
    }

    if change {
        w.change <- true
    }

}
