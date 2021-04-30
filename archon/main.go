package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	port  string = "8080"
	token string
)

func init() {

	// overwrite port
	if value, ok := os.LookupEnv("PORT"); ok {
		port = value
	}

	// set token
	if value, ok := os.LookupEnv("TOKEN"); ok {
		token = value
	} else {
		log.Fatal("request token undefined")
	}

}

func sourceScript(name string) (script []Task, err error) {

	// open script
	scriptYAML, err := os.Open(name)
	if err != nil {
		return
	}
	defer scriptYAML.Close()

	// parse script
	err = yaml.NewDecoder(scriptYAML).Decode(&script)
	if err != nil {
		return
	}

	return

}

func sourceFleet(name string) (fleet []string, err error) {

	// open fleet
	fleetYAML, err := os.Open(name)
	if err != nil {
		return
	}

	// parse fleet
	err = yaml.NewDecoder(fleetYAML).Decode(&fleet)
	if err != nil {
		return
	}

	return

}

func execute(errChan chan<- error, script []Task, server string) {

	// run through tasks in script
	for index, task := range script {

		// log start
		log.Printf("[%s] started #%d %s", server, index+1, task)

		// run task; return on error
		err := task.Run(server, token)
		if err != nil {
			log.Printf("[%s] failed #%d %s; %s", server, index+1, task, err)
			errChan <- err
			return
		}

		// log complete
		log.Printf("[%s] compleated #%d %s", server, index+1, task)

	}

	// write nil to error channel
	errChan <- nil

}

func main() {

	// validate number of arguments to command
	if len(os.Args) != 3 {
		log.Fatal("invalid number of arguments; see documentation for help")
	}

	// source script from yaml; cmd arg 1
	script, err := sourceScript(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// source fleet from yaml; cmd arg 2
	fleet, err := sourceFleet(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	// error channel for routines to write to
	errChan := make(chan error)

	// concurrently execute script on each server
	for _, server := range fleet {
		go execute(errChan, script, server)
	}

	// block until each routine writes to channel; count errors
	var errCount int
	for range fleet {
		if err := <-errChan; err != nil {
			errCount++
		}
	}

	// exit with status 1 if errors were encountered
	if errCount > 0 {
		log.Fatalf("script failed on %d out of %d servers", errCount, len(fleet))
	}

	// success
	log.Print("script ran successfully on fleet")

}
