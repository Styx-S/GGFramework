package Define

type Server struct {
	ModuleMap map[string]Module
}

func (svr *Server) GetModule(moduleName string) Module {
	return svr.ModuleMap[moduleName]
}

func (svr *Server) AddModules(moduleName string, module Module) {
	svr.ModuleMap[moduleName] = module
	module.ConnectToSvr(svr)
}

func (svr *Server) Start() {
	for _, m := range svr.ModuleMap {
		go m.Start()
	}
}
