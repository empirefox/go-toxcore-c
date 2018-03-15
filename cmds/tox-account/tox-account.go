package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	tox "github.com/TokTok/go-toxcore-c"
)

var secret = flag.String("secret", "", "hex encoded secret, auto generate if not set")
var nospam = flag.Uint("nospam", 0, "nospam, auto generate if not set")

func main() {
	flag.Parse()

	if *nospam == 0 {
		*nospam = uint(time.Now().Nanosecond())
	}

	var a *tox.Account
	var err error
	if *secret == "" {
		a, err = tox.GenerateAccount(rand.Reader, uint32(*nospam))
	} else {
		a, err = tox.NewAccountFrom(*secret, uint32(*nospam))
	}
	if err != nil {
		log.Fatalln(err)
	}
	ha := a.HumanReadable()
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(ha)
	if err != nil {
		log.Fatalln(err)
	}
}
