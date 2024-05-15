package shell

import (
	log "xxx-server/application/logger"
)

type Repositories struct {
	// PySearchModule   repo.PythonCmdRepo
	// SiteShpProc repo.SubCmdRepo
}

var (
	R *Repositories
)

func SetupRepositories() (err error) {
	R = &Repositories{
		// PySearchModule:   NewPythonCmd("search_module.py", config.C.Py.SearchPs),
		// SiteShpProc: NewCmd("ogr2ogr"),
	}
	log.Info("setup sub cmds succeed")
	return
}
