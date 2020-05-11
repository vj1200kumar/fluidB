package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/tidwall/tile38/core"
	"github.com/tidwall/tile38/internal/hservice"
	"github.com/tidwall/tile38/internal/log"
	"github.com/tidwall/tile38/internal/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	dir         string
	port        int
	host        string
	verbose     bool
	veryVerbose bool
	devMode     bool
	quiet       bool
	pidfile     string
	cpuprofile  string
	memprofile  string
	pprofport   int
	nohup       bool
)

// TODO: Set to false in 2.*
var httpTransport = true

////////////////////////////////////////////////////////////////////////////////
//
// Fire up a webhook test server by using the --webhook-http-consumer-port
// for example
//   $ ./tile38-server --webhook-http-consumer-port 9999
//
// The create hooks like such...
//   SETHOOK myhook http://localhost:9999/myhook NEARBY mykey FENCE POINT 33.5 -115.5 1000
//
////////////////////////////////////////////////////////////////////////////////
//
// Memory profiling - start the server with the -pprofport flag
//
//   $ ./tile38-server -pprofport 6060
//
// Then, at any point, from a different terminal execute:
//   $ go tool pprof -svg http://localhost:6060/debug/pprof/heap > out.svg
//
// Load the SVG into a web browser to visualize the memory usage
//
////////////////////////////////////////////////////////////////////////////////

type hserver struct{}

func (s *hserver) Send(ctx context.Context, in *hservice.MessageRequest) (*hservice.MessageReply, error) {
	return &hservice.MessageReply{Ok: true}, nil
}

func main() {
	gitsha := " (" + core.GitSHA + ")"
	if gitsha == " (0000000)" {
		gitsha = ""
	}
	versionLine := `tile38-server version: ` + core.Version + gitsha

	output := os.Stderr
	flag.Usage = func() {
		fmt.Fprintf(output,
			versionLine+`

Usage: tile38-server [-p port]

Basic Options:
  -h hostname : listening host
  -p port     : listening port (default: 9470)
  -d path     : data directory (default: data)
  -q          : no logging. totally silent output
  -v          : enable verbose logging
  -vv         : enable very verbose logging

Advanced Options: 
  --pidfile path          : file that contains the pid
  --appendonly yes/no     : AOF persistence (default: yes)
  --appendfilename path   : AOF path (default: data/appendonly.aof)
  --queuefilename path    : Event queue path (default:data/queue.db)
  --http-transport yes/no : HTTP transport (default: yes)
  --protected-mode yes/no : protected mode (default: yes)
  --threads num           : number of network threads (default: num cores)
  --nohup                 : do not exist on SIGHUP

Developer Options:
  --dev                             : enable developer mode
  --webhook-http-consumer-port port : Start a test HTTP webhook server
  --webhook-grpc-consumer-port port : Start a test GRPC webhook server

`,
		)
	}

	if len(os.Args) == 3 && os.Args[1] == "--webhook-http-consumer-port" {
		log.SetOutput(os.Stderr)
		port, err := strconv.ParseUint(os.Args[2], 10, 16)
		if err != nil {
			log.Fatal(err)
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.HTTPf("http: %s : %s", r.URL.Path, string(data))
		})
		log.Infof("webhook server http://localhost:%d/", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			log.Fatal(err)
		}
		return
	}

	if len(os.Args) == 3 && os.Args[1] == "--webhook-grpc-consumer-port" {
		log.SetOutput(os.Stderr)
		port, err := strconv.ParseUint(os.Args[2], 10, 16)
		if err != nil {
			log.Fatal(err)
		}

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatal(err)
		}
		s := grpc.NewServer()
		hservice.RegisterHookServiceServer(s, &hserver{})
		log.Infof("webhook server grpc://localhost:%d/", port)
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
		return
	}

	var showEvioDisabled bool

	// parse non standard args.
	nargs := []string{os.Args[0]}
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--help":
			output = os.Stdout
			flag.Usage()
			return
		case "--version":
			fmt.Fprintf(os.Stdout, "%s\n", versionLine)
			return
		case "--protected-mode", "-protected-mode":
			i++
			if i < len(os.Args) {
				switch strings.ToLower(os.Args[i]) {
				case "no":
					core.ProtectedMode = "no"
					continue
				case "yes":
					core.ProtectedMode = "yes"
					continue
				}
			}
			fmt.Fprintf(os.Stderr, "protected-mode must be 'yes' or 'no'\n")
			os.Exit(1)
		case "--dev", "-dev":
			devMode = true
			continue
		case "--nohup", "-nohup":
			nohup = true
			continue
		case "--appendonly", "-appendonly":
			i++
			if i < len(os.Args) {
				switch strings.ToLower(os.Args[i]) {
				case "no":
					core.AppendOnly = false
					continue
				case "yes":
					core.AppendOnly = true
					continue
				}
			}
			fmt.Fprintf(os.Stderr, "appendonly must be 'yes' or 'no'\n")
			os.Exit(1)
		case "--appendfilename", "-appendfilename":
			i++
			if i == len(os.Args) || os.Args[i] == "" {
				fmt.Fprintf(os.Stderr, "appendfilename must have a value\n")
				os.Exit(1)
			}
			core.AppendFileName = os.Args[i]
		case "--queuefilename", "-queuefilename":
			i++
			if i == len(os.Args) || os.Args[i] == "" {
				fmt.Fprintf(os.Stderr, "queuefilename must have a value\n")
				os.Exit(1)
			}
			core.QueueFileName = os.Args[i]
		case "--http-transport", "-http-transport":
			i++
			if i < len(os.Args) {
				switch strings.ToLower(os.Args[i]) {
				case "1", "true", "yes":
					httpTransport = true
					continue
				case "0", "false", "no":
					httpTransport = false
					continue
				}
			}
			fmt.Fprintf(os.Stderr, "http-transport must be 'yes' or 'no'\n")
			os.Exit(1)
		case "--threads", "-threads":
			i++
			if i < len(os.Args) {
				n, err := strconv.ParseUint(os.Args[i], 10, 16)
				if err != nil {
					fmt.Fprintf(os.Stderr, "threads must be a valid number\n")
					os.Exit(1)
				}
				core.NumThreads = int(n)
				continue
			}
			fmt.Fprintf(os.Stderr, "http-transport must be 'yes' or 'no'\n")
			os.Exit(1)
		case "--evio", "-evio":
			i++
			if i < len(os.Args) {
				switch strings.ToLower(os.Args[i]) {
				case "no", "yes":
					showEvioDisabled = true
					continue
				}
			}
			fmt.Fprintf(os.Stderr, "evio must be 'yes' or 'no'\n")
			os.Exit(1)
		}
		nargs = append(nargs, os.Args[i])
	}
	os.Args = nargs

	flag.IntVar(&port, "p", 9470, "The listening port.")
	flag.StringVar(&pidfile, "pidfile", "", "A file that contains the pid")
	flag.StringVar(&host, "h", "", "The listening host.")
	flag.StringVar(&dir, "d", "data", "The data directory.")
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging.")
	flag.BoolVar(&quiet, "q", false, "Quiet logging. Totally silent.")
	flag.BoolVar(&veryVerbose, "vv", false, "Enable very verbose logging.")
	flag.IntVar(&pprofport, "pprofport", 0, "pprofport http at port")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to `file`")
	flag.StringVar(&memprofile, "memprofile", "", "write memory profile to `file`")
	flag.Parse()

	var logw io.Writer = os.Stderr
	if quiet {
		logw = ioutil.Discard
	}
	log.SetOutput(logw)
	if quiet {
		log.Level = 0
	} else if veryVerbose {
		log.Level = 3
	} else if verbose {
		log.Level = 2
	} else {
		log.Level = 1
	}
	core.DevMode = devMode
	core.ShowDebugMessages = veryVerbose

	hostd := ""
	if host != "" {
		hostd = "Addr: " + host + ", "
	}

	// pprof
	if cpuprofile != "" {
		log.Debugf("cpuprofile active")
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
	}
	if memprofile != "" {
		log.Debug("memprofile active")
	}

	var pprofcleanedup bool
	var pprofcleanupMu sync.Mutex
	pprofcleanup := func() {
		pprofcleanupMu.Lock()
		defer pprofcleanupMu.Unlock()
		if pprofcleanedup {
			return
		}
		// cleanup code
		if cpuprofile != "" {
			pprof.StopCPUProfile()
		}
		if memprofile != "" {
			f, err := os.Create(memprofile)
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
			f.Close()
		}
		pprofcleanedup = true
	}
	defer pprofcleanup()

	if pprofport != 0 {
		log.Debugf("pprof http at port %d", pprofport)
		go func() {
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", pprofport), nil))
		}()
	}

	// pid file
	var pidferr error
	var pidcleanedup bool
	var pidcleanupMu sync.Mutex
	pidcleanup := func() {
		if pidfile != "" {
			pidcleanupMu.Lock()
			defer pidcleanupMu.Unlock()
			if pidcleanedup {
				return
			}
			// cleanup code
			if pidfile != "" {
				os.Remove(pidfile)
			}
			pidcleanedup = true
		}
	}
	defer pidcleanup()
	if pidfile != "" {
		pidferr := ioutil.WriteFile(pidfile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0666)
		if pidferr == nil {
		}
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			if s == syscall.SIGHUP && nohup {
				continue
			}
			log.Warnf("signal: %v", s)
			pidcleanup()
			pprofcleanup()
			switch {
			default:
				os.Exit(-1)
			case s == syscall.SIGHUP:
				os.Exit(1)
			case s == syscall.SIGINT:
				os.Exit(2)
			case s == syscall.SIGQUIT:
				os.Exit(3)
			case s == syscall.SIGTERM:
				os.Exit(0xf)
			}
		}
	}()

	
	fmt.Fprintf(logw, `
   _______ _______
  |       |       |
  |____   |   _   |   FluidB %s%s %d bit (%s/%s)
  |       |       |   %sPort: %d, PID: %d
  |____   |   _   | 
  |       |       |   https://fluidb.icu
  |_______|_______| 
`+"\n", core.Version, gitsha, strconv.IntSize, runtime.GOARCH, runtime.GOOS, hostd, port, os.Getpid())
	if pidferr != nil {
		log.Warnf("pidfile: %v", pidferr)
	}
	if showEvioDisabled {
		// we don't currently support evio in Tile38
		log.Warnf("evio is not currently supported")
	}
	if err := server.Serve(host, port, dir, httpTransport); err != nil {
		log.Fatal(err)
	}
}
