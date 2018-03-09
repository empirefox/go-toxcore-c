package main

//  tox save data explorer

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/TokTok/go-toxcore-c"
	"github.com/TokTok/go-toxcore-c/toxenums"
)

func init() {
	log.SetFlags(log.Flags() ^ log.Ldate ^ log.Ltime)
}

var tsfile string
var pass string

func printHelp() {
	log.Println("For help: /path/to/rsexp -h")
}

func main() {
	// flag.StringVar(&tsfile, "tsfile", "", "tox save data file")
	flag.StringVar(&pass, "pass", pass, "tox save data password")
	flag.Parse()
	// log.Println(flag.Args())
	if len(flag.Args()) < 1 {
		printHelp()
		flag.Usage()
		return
	}
	tsfile = flag.Arg(0)

	data, err := ioutil.ReadFile(tsfile)
	if err != nil {
		log.Println(err)
		return
	}
	isencrypt := tox.IsDataEncrypted(data)
	log.Println("Is encrypt: ", isencrypt)
	if isencrypt {
		salt, err := tox.GetSalt(data)
		if err != nil {
			log.Println(err, len(salt), salt)
		}
		pkey, err := tox.DeriveWithSalt([]byte(pass), salt)
		defer pkey.Free()
		if err != nil {
			log.Println(err)
		}
		datad, err := pkey.Decrypt(data)
		if err != nil {
			// log.Println(ok, err, len(datad), datad[0:32])
			log.Println("Decrypt error, check your -pass:", err)
			return
		}

		data = datad
	}

	opts := tox.NewToxOptions()
	opts.Savedata_type = toxenums.TOX_SAVEDATA_TYPE_TOX_SAVE
	opts.Savedata_data = data
	t, err := tox.NewTox(opts)
	if err != nil {
		log.Println(err)
		return
	}
	fnums := t.SelfGetFriendList()
	log.Println("Self Name:", t.SelfGetName())
	log.Printf("Self ID: %X\n", t.SelfGetAddress()[:])
	mystmsg := t.SelfGetStatusMessage()
	log.Println("Status:", mystmsg)
	log.Println("------------------------------------------")
	log.Println("Friend Count:", len(fnums))

	if len(fnums) > 0 {
		log.Println("num\tname\tID\tseen\tstatus\tstmsg")
	}
	for i := 0; i < len(fnums); i++ {
		pubkey, err := t.FriendGetPublicKey(fnums[i])
		tm, err := t.FriendGetLastOnline(fnums[i])
		fname, err := t.FriendGetName(fnums[i])
		stmsg, err := t.FriendGetStatusMessage(fnums[i])
		status, err := t.FriendGetConnectionStatus(fnums[i])
		if err != nil {
			log.Println("wtf", i)
		} else {
			otm := time.Unix(int64(tm), 0)
			log.Println(fmt.Sprintf("Friend %d: ", fnums[i]),
				fname, fmt.Sprintf("%X", pubkey[:]), otm, status, stmsg)
		}
	}
	if len(fnums) > 20 {
		log.Println("Friend Count:", len(fnums))
	}
	if len(fnums) > 0 {
		log.Println()
	}
}
