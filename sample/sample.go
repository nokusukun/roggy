package main

import (
    "github.com/nokusukun/roggy"
)

func main() {
    log := roggy.Printer("main-service")
    roggy.LogLevel = roggy.TypeDebug
    log.Error("This should print This as well")
    log.Notice("This should print")
    log.Info("This should print")
    log.Verbose("This should print")
    log.Debug("Lorem ipsum dolor et whatever what teh fuck did you just say to me you little shit")
    log.Verbose("This should print")


    roggy.Wait()
}