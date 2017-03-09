package module

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/bohler/lib/dlog"
)

func Run(mods ...Module) {
	for i := 0; i < len(mods); i++ {
		Register(mods[i])
	}
	Init()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGKILL, syscall.SIGINT)
	log.Log.Debug("running...")
	sig := <-c
	log.Log.Warningf("App shutdown (siggal: %v)", sig)
	Destroy()
}
