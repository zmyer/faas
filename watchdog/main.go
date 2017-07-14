// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/alexellis/faas/watchdog/types"
)

func main() {
	osEnv := types.OsEnv{}
	readConfig := ReadConfig{}
	config := readConfig.Read(osEnv)

	if len(config.faasProcess) == 0 {
		log.Panicln("Provide a valid process via fprocess environmental variable.")
		return
	}

	readTimeout := config.readTimeout
	writeTimeout := config.writeTimeout

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
	}

	http.HandleFunc("/", makeRequestHandler(&config))

	ticker := time.NewTicker(time.Millisecond * 20)

	// ListenAndServe is blocking, but we need the lock file to be written only after the server is listening.
	go func() {
		<-ticker.C

		if config.suppressLock == false {
			path := "/tmp/.lock"
			log.Printf("Writing lock-file to: %s\n", path)
			writeErr := ioutil.WriteFile(path, []byte{}, 0660)
			if writeErr != nil {
				log.Panicf("Cannot write %s. To disable lock-file set env suppress_lock=true.\n Error: %s.\n", path, writeErr.Error())
			}
		}
	}()

	log.Fatal(s.ListenAndServe())
}
