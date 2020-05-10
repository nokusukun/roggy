package roggy

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "runtime"
    "strings"
    "time"

    "github.com/nokusukun/stemp"
)

const (
    TypeRoggy = iota - 1
    TypeError
    TypeNotice
    TypeInfo
    TypeVerbose
    TypeDebug
)

var (
    Enable           = true
    LogQueue         = make(chan LogShard, 100)
    LogLevel         = 1
    Simple           = false
    running          = false
    CurrentSupMinute = ""
    Filter           = ""
    Hook             = func(shard LogShard) {}

    StdOut = os.Stdout
    StdErr = os.Stderr

    LevelText = map[int]string{
        TypeRoggy:   "ROGY",
        TypeError:   "ERRO",
        TypeNotice:  "NOTI",
        TypeInfo:    "INFO",
        TypeVerbose: "VERB",
        TypeDebug:   "DEBU",
    }

    wait = make(chan interface{})
)

type LogShard struct {
    Type        string        `json:"type"`
    Service     string        `json:"service"`
    Trace       string        `json:"source"`
    Message     []interface{} `json:"message"`
    LogLevel    int           `json:"loglevel"`
    Color       string        `json:"-"`
    Destination *os.File      `json:"-"`
}

func getFrame(skipFrames int) runtime.Frame {
    targetFrameIndex := skipFrames + 2

    programCounters := make([]uintptr, targetFrameIndex+2)
    n := runtime.Callers(0, programCounters)

    frame := runtime.Frame{Function: "unknown"}
    if n > 0 {
        frames := runtime.CallersFrames(programCounters[:n])
        for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
            var frameCandidate runtime.Frame
            frameCandidate, more = frames.Next()
            if frameIndex == targetFrameIndex {
                frame = frameCandidate
            }
        }
    }

    return frame
}

func Printer(service string) *LogPrinter {
    if !running {
        go start()
        running = !running
    }
    return &LogPrinter{service, 0, "", StdOut}
}

type LogPrinter struct {
    Service     string
    level       int
    logFile     string
    Destination *os.File
}

func (p *LogPrinter) Sub(service string) *LogPrinter {
    return &LogPrinter{
        fmt.Sprintf("%v/%v", p.Service, Clr(service, p.level+1)),
        p.level + 1,
        p.logFile,
        p.Destination,
    }
}

func (p *LogPrinter) Notice(message ...interface{}) {
    if LogLevel < TypeNotice {
        return
    }
    rawSource := fmt.Sprintf("%v:%v:%v",  getFrame(1).File, getFrame(1).Function, getFrame(1).Line)
    source := strings.Split(rawSource, "/")
    if !running {
        return
    }
    LogQueue <- LogShard{
        Type:     "Notice",
        Service:  p.Service,
        Trace:    source[len(source)-1],
        Message:  message,
        LogLevel: TypeNotice,
        Color:    "\u001b[32;1m"}
}

func (p *LogPrinter) Noticef(f string, message ...interface{}) {
    if LogLevel < TypeNotice {
        return
    }
    msg := fmt.Sprintf(f, message...)
    rawSource := fmt.Sprintf("%v:%v:%v",  getFrame(1).File, getFrame(1).Function, getFrame(1).Line)
    source := strings.Split(rawSource, "/")
    if !running {
        return
    }
    LogQueue <- LogShard{
        Type:     "Notice",
        Service:  p.Service,
        Trace:    source[len(source)-1],
        Message:  []interface{}{msg},
        LogLevel: TypeNotice,
        Color: "\u001b[32;1m",
        Destination: p.Destination,
    }
}

func (p *LogPrinter) Info(message ...interface{}) {
    if LogLevel < TypeInfo {
        return
    }
    rawSource := fmt.Sprintf("%v:%v:%v",  getFrame(1).File, getFrame(1).Function, getFrame(1).Line)
    source := strings.Split(rawSource, "/")
    if !running {
        return
    }
    LogQueue <- LogShard{
        Type:        "Info",
        Service:     p.Service,
        Trace:       source[len(source)-1],
        Message:     message,
        LogLevel:    TypeInfo,
        Color:       "\u001b[34;1m",
        Destination: p.Destination,
    }
}

func (p *LogPrinter) Infof(f string, message ...interface{}) {
    if LogLevel < TypeInfo {
        return
    }
    msg := fmt.Sprintf(f, message...)
    rawSource := fmt.Sprintf("%v:%v:%v",  getFrame(1).File, getFrame(1).Function, getFrame(1).Line)
    source := strings.Split(rawSource, "/")
    if !running {
        return
    }
    LogQueue <- LogShard{
        Type:        "Info",
        Service:     p.Service,
        Trace:       source[len(source)-1],
        Message:     []interface{}{msg},
        LogLevel:    TypeInfo,
        Color:       "\u001b[34;1m",
        Destination: p.Destination,
    }
}

func (p *LogPrinter) Error(message ...interface{}) {
    if LogLevel < TypeError {
        return
    }
    rawSource := fmt.Sprintf("%v:%v:%v",  getFrame(1).File, getFrame(1).Function, getFrame(1).Line)
    source := strings.Split(rawSource, "/")
    if !running {
        return
    }
    LogQueue <- LogShard{
        Type:        "Error",
        Service:     p.Service,
        Trace:       source[len(source)-1],
        Message:     message,
        LogLevel:    TypeError,
        Color:       "\u001b[31;1m",
        Destination: p.Destination,
    }
}

func (p *LogPrinter) Errorf(f string, message ...interface{}) {
    if LogLevel < TypeError {
        return
    }
    msg := fmt.Sprintf(f, message...)
    rawSource := fmt.Sprintf("%v:%v:%v",  getFrame(1).File, getFrame(1).Function, getFrame(1).Line)
    source := strings.Split(rawSource, "/")
    if !running {
        return
    }
    LogQueue <- LogShard{
        Type:        "Error",
        Service:     p.Service,
        Trace:       source[len(source)-1],
        Message:     []interface{}{msg},
        LogLevel:    TypeError,
        Color:       "\u001b[31;1m",
        Destination: p.Destination,
    }
}

func (p *LogPrinter) Verbose(message ...interface{}) {
    if LogLevel < TypeVerbose {
        return
    }
    rawSource := fmt.Sprintf("%v:%v:%v",  getFrame(1).File, getFrame(1).Function, getFrame(1).Line)
    source := strings.Split(rawSource, "/")
    if !running {
        return
    }
    LogQueue <- LogShard{
        Type:        "Verbose",
        Service:     p.Service,
        Trace:       source[len(source)-1],
        Message:     message,
        LogLevel:    TypeVerbose,
        Color:       "\u001b[33;1m",
        Destination: p.Destination,
    }
}

func (p *LogPrinter) Verbosef(f string, message ...interface{}) {
    if LogLevel < TypeVerbose {
        return
    }
    msg := fmt.Sprintf(f, message...)
    rawSource := fmt.Sprintf("%v:%v:%v",  getFrame(1).File, getFrame(1).Function, getFrame(1).Line)
    source := strings.Split(rawSource, "/")
    if !running {
        return
    }
    LogQueue <- LogShard{
        Type:        "Verbose",
        Service:     p.Service,
        Trace:       source[len(source)-1],
        Message:     []interface{}{msg},
        LogLevel:    TypeVerbose,
        Color:       "\u001b[33;1m",
        Destination: p.Destination,
    }
}

func (p *LogPrinter) Debug(message ...interface{}) {
    if LogLevel < TypeDebug {
        return
    }
    rawSource := fmt.Sprintf("%v:%v:%v",  getFrame(1).File, getFrame(1).Function, getFrame(1).Line)
    source := strings.Split(rawSource, "/")
    if !running {
        return
    }
    LogQueue <- LogShard{
        Type:        "Debug",
        Service:     p.Service,
        Trace:       source[len(source)-1],
        Message:     message,
        LogLevel:    TypeDebug,
        Color:       "\u001b[36;1m",
        Destination: p.Destination,
    }
}

func (p *LogPrinter) Debugf(f string, message ...interface{}) {
    if LogLevel < TypeDebug {
        return
    }
    msg := fmt.Sprintf(f, message...)
    rawSource := getFrame(1)
    source := strings.Split(rawSource.Function, "/")
    if !running {
        return
    }
    LogQueue <- LogShard{
        Type:        "Debug",
        Service:     p.Service,
        Trace:       source[len(source)-1],
        Message:     []interface{}{msg},
        LogLevel:    TypeDebug,
        Color:       "\u001b[36;1m",
        Destination: p.Destination,
    }
}

func Wait() {
    if !running {
        return
    }

    running = false
    close(LogQueue)
    <-wait
}

func commandListener() {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        txt := scanner.Text()
        if strings.HasPrefix(txt, ":l") {
            _, _ = fmt.Sscanf(txt, ":l%d", &LogLevel)
            LogQueue <- LogShard{
                Type:     "ROGGY",
                Service:  "ROGGY",
                Trace:    "INTERNAL",
                Message:  []interface{}{"Changed log level to ", LogLevel},
                LogLevel: -1,
                Color:    "\u001b[36;1m"}
        }

        if strings.HasPrefix(txt, ":f") {
            _, _ = fmt.Sscanf(txt, ":f%v", &Filter)
            msg := []interface{}{"Setting Filter to : ", Filter}
            if Filter == ":c" {
                Filter = ""
                msg = []interface{}{"Clearing filter"}
            }
            LogQueue <- LogShard{
                Type:     "ROGGY",
                Service:  "ROGGY",
                Trace:    "INTERNAL",
                Message:  msg,
                LogLevel: -1,
                Color:    "\u001b[36;1m"}
        }
    }
}

func start() {
    go commandListener()
    for log := range LogQueue {
        go Hook(log)
        csm := time.Now().Format("02/01/06 03:04PM")
        if csm != CurrentSupMinute {
            CurrentSupMinute = csm
            //                   BG Yellow     FG Black        FG Yellow    BG Black
            fmt.Printf("\u001b[43;30m [ %v ] \u001b[33;40m\u001b[0m\n", csm)
        }

        now := time.Now().Format("05.999")
        if log.LogLevel <= LogLevel && Enable {
            //fmt.Println(log)
            cl := "\u001b[0m"+log.Color
            re := "\u001b[30;1m"
            lbr := Clr("[", 2)
            rbr := Clr("]", 2)

            if Simple {
                cl = ""
                re = ""
                lbr = ""
                rbr = ""
            }

            msg := stemp.Compile(
                `{now:j=l,w=7} {col}[{type}]{lbr}{service}{rbr} {col}`,
                map[string]interface{}{
                    "now":     fmt.Sprintf("\u001b[33m%v", now),
                    "type":    LevelText[log.LogLevel],
                    "service": log.Service,
                    "col":     cl,
                    "lbr":     lbr,
                    "rbr":     rbr,
                    "reset":   re,
                })
            msg += fmt.Sprint(log.Message...)
            msg += fmt.Sprint(" \u001b[30;1m>", log.Trace, "\u001b[0m\n")
            if Filter == "" || (Filter != "" && strings.Contains(msg, Filter)) || log.LogLevel == -1 {
                if Filter != "" {
                    msg = fmt.Sprintf("F:%v|", Filter) + msg
                }
                _, _ = fmt.Fprint(StdOut, msg)
            }
        }
    }
    wait <- 1
}

func Rainbowize(text string) string {
    offset := rand.Intn(7)
    ret := ""

    for i, s := range strings.Split(text, "") {
        col := (i + offset) % 7
        ret += fmt.Sprintf("\u001b[3%v;1m%v", col, s)
    }

    ret += "\u001b[0m"
    return ret
}

func Clr(text string, offset int) string {
    col := offset % 7
    return fmt.Sprintf("\u001b[3%v;1m%v\u001b[0m", col, text)
}
