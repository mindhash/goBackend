package rest

import (
	"fmt"
	"net/http"
	"flag"
	"github.com/mindhash/goBackend/base"
)

var config *ServerConfig


const(
	// Default value of ServerConfig.MaxIncomingConnections
	DefaultMaxIncomingConnections = 0
)

// JSON object that defines the server configuration.
type ServerConfig struct {
	Interface                      *string         // Interface to bind REST API to, default ":4984"
	SSLCert				   		   *string
	SSLKey						   *string
	ServerReadTimeout              *int            // maximum duration.Second before timing out read of the HTTP(S) request
	ServerWriteTimeout             *int            // maximum duration.Second before timing out write of the HTTP(S) response
	Pretty						   bool
	//Log                            []string        // Log keywords to enable
	//LogFilePath                    *string         // Path to log file, if missing write to stderr 
	Database                      *DbConfig     // Pre-configured databases, mapped by name
	MaxIncomingConnections 		   *int            // Max # of incoming HTTP connections to accept
}


type DbConfig struct {
	Name               string  
}

func ParseCommandLine() {
	dbname := flag.String("dbname","sampledb","Default Database Name")
	addr   := flag.String("addr","localhost:4984","HTTP Server Address") 
	pretty := flag.Bool("pretty", true, "Pretty-print JSON responses")
	flag.Parse()
	
	config = &ServerConfig { Interface: addr, Pretty: *pretty,Database: &DbConfig {Name: *dbname}}
}

func (config *ServerConfig) serve(addr string, handler http.Handler) {
	
	maxConns := DefaultMaxIncomingConnections
	if config.MaxIncomingConnections != nil {
		maxConns = *config.MaxIncomingConnections
	}

	err := base.ListenAndServeHTTP(addr, maxConns, config.SSLCert, config.SSLKey, handler, config.ServerReadTimeout, config.ServerWriteTimeout)
	if err != nil {
		base.LogFatal("Failed to start HTTP server on %s: %v", addr, err)
	}
}

func RunServer(config *ServerConfig) {
	
	//New Server Context
	sc := NewServerContext(config)
	
	//Open Database and add to server context 
	//if _, err := sc.AddDatabaseFromConfig(config.Database); err != nil {
	//		base.LogFatal("Error opening database: %v", err)
	//}
		
	//defer sc.CloseDatabase()
	
	base.Logf("Starting server on %s ...", *config.Interface)
	config.serve(*config.Interface, CreatePublicHandler(sc))

}



func ServerMain(){
	fmt.Println("Initiating Server..")
	ParseCommandLine()
	RunServer(config)
}
