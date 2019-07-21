package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type (
	Options struct {
		ShowHelp   bool
		RemovePort bool
		Port       int
		HttpCode   int
		ServerAddr string
	}
)

func parseflags() (*Options, *flag.FlagSet, error) {

	fs := flag.NewFlagSet("", flag.ContinueOnError)

	options := new(Options)

	fs.BoolVar(&options.ShowHelp, "help", false, "show help")
	fs.BoolVar(&options.RemovePort, "n", false, "remove port from URI")
	fs.IntVar(&options.Port, "p", 443, "redefine port")
	fs.IntVar(&options.HttpCode, "c", 308, "redirection code")
	fs.StringVar(&options.ServerAddr, "s", "0.0.0.0:80", "set address this server listening on")

	for k, v := range os.Args {
		fmt.Println(k, v)
	}

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, nil, err
	}

	return options, fs, nil
}

func getaddrs(opts *Options, fs *flag.FlagSet) (*net.TCPAddr, error) {

	serve_addr := (opts.ServerAddr)

	saddr, err := net.ResolveTCPAddr("tcp", serve_addr)
	if err != nil {
		return nil, errors.New("error resolving s parameter: " + err.Error())
	}

	return saddr, nil
}

func main() {

	log.Print("Welcome to Simple HTTPS redirect")
	log.Print("https://github.com/AnimusPEXUS/simplehttpsredirect")

	opts, fs, err := parseflags()
	if err != nil {
		log.Fatal("flags parse error: " + err.Error())
	}

	{
		b, _ := json.MarshalIndent(opts, "  ", "  ")
		fmt.Println(string(b))
	}

	fs.PrintDefaults()

	if opts.ShowHelp {
		fs.PrintDefaults()
		log.Fatal("showing help")
	}

	saddr, err := getaddrs(opts, fs)
	if err != nil {
		log.Fatal("error working with supplied address: " + err.Error())
	}

	// fmt.Println("saddr", saddr)

	log.Fatal(
		http.ListenAndServe(
			saddr.String(),
			&Server{opts, fs},
		),
	)

	return

}

type Server struct {
	opts *Options
	fs   *flag.FlagSet
}

func (self *Server) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {

	var c *url.URL

	log.Print("request url ", rq.URL.String())

	{
		t := *rq.URL
		tt := t
		c = &tt
	}

	c.Scheme = "https"

	if self.opts.RemovePort || self.opts.Port == 443 {
		c.Host = c.Hostname()
	} else {
		c.Host = net.JoinHostPort(c.Hostname(), strconv.Itoa(self.opts.Port))
	}

	h := rw.Header()
	h.Set("Location", c.String())

	rw.WriteHeader(self.opts.HttpCode)

}
