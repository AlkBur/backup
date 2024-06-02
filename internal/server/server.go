package server

import (
	"backup/internal/config"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc defines the request handler used by server
type HandlerFunc func(*Context)

// Server implement the interface of ServeHTTP
type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // support middleware
		parent      *RouterGroup  // support nesting
		server      *Server       // all groups share a Server instance
	}

	Server struct {
		*RouterGroup
		router        *router
		groups        []*RouterGroup     // store all groups
		htmlTemplates *template.Template // for html render
		funcMap       template.FuncMap   // for html render
		cfgInfo       config.Notification
	}
)

// New is the constructor of Server
func New() *Server {
	server := &Server{router: newRouter()}
	server.RouterGroup = &RouterGroup{server: server}
	server.groups = []*RouterGroup{server.RouterGroup}
	return server
}

// Default use Logger() & Recovery middlewares
func Default() *Server {
	server := New()
	server.Use(Logger(), Recovery())
	return server
}

func (server *Server) Notification(info config.Notification) {
	server.cfgInfo = info
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Server instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	server := group.server
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		server: server,
	}
	server.groups = append(server.groups, newGroup)
	return newGroup
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.server.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// for custom render function
func (server *Server) SetFuncMap(funcMap template.FuncMap) {
	server.funcMap = funcMap
}

func (server *Server) LoadHTMLGlob(pattern string) {
	server.htmlTemplates = template.Must(template.New("").Funcs(server.funcMap).ParseGlob(pattern))
}

// Run defines the method to start a http server
func (server *Server) Run(addr string) (err error) {
	return http.ListenAndServe(addr, server)
}

func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range server.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.server = server
	server.router.handle(c)
}
