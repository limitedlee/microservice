package micro

type Group struct {
	routes []Route
}

type Route struct {
	verb       HttpVerb
	prefix     string
	handleFunc HandleFunc
}

func (bs *Service) Group(prefix string) *Group {
	if bs.opts.serviceType != Http {
		panic("unsupported grpc defined route")
	}

	group := &Group{}

	if prefix == "/" {
		return bs.groups["/"]
	}

	if _, ok := bs.groups[prefix]; !ok {
		group.routes = make([]Route, 0, 999)
		bs.groups[prefix] = group
	}

	return group
}

func (g *Group) Get(prefix string, f HandleFunc) {
	route := Route{
		prefix: prefix,
		verb:   GET,
	}
	route.handleFunc = f
	g.routes = append(g.routes, route)
}

func (g *Group) POST(prefix string, f HandleFunc) {
	route := Route{
		prefix: prefix,
		verb:   POST,
	}
	route.handleFunc = f
	g.routes = append(g.routes, route)
}

func (g *Group) PUT(prefix string, f HandleFunc) {
	route := Route{
		prefix: prefix,
		verb:   PUT,
	}
	route.handleFunc = f
	g.routes = append(g.routes, route)
}

func (g *Group) DELETE(prefix string, f HandleFunc) {
	route := Route{
		prefix: prefix,
		verb:   DELETE,
	}
	route.handleFunc = f
	g.routes = append(g.routes, route)
}
