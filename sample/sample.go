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
	roggy.Flush()

	log.Sync = true
	log.Error("Sync: This should print This as well")
	log.Notice("Sync: This should print")
	log.Info("Sync: This should print")
	log.Verbose("Sync: This should print")
	log.Debug("Sync: Lorem ipsum dolor et whatever what teh fuck did you just say to me you little shit")
	log.Verbose("Sync: This should print")

	log.NoTrace = true
	log.Error("NoTrace: This should print This as well")
	log.Notice("NoTrace: This should print")
	log.Info("NoTrace: This should print")
	log.Verbose("NoTrace: This should print")
	log.Debug("NoTrace: Lorem ipsum dolor et whatever what teh fuck did you just say to me you little shit")
	log.Verbose("NoTrace: This should print")

	roggy.Wait()
}
