package main

//  tox save data explorer

import (
	"flag"
	"io/ioutil"
	"log"

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
	flist := t.SelfGetFriendList()
	log.Println("Self Name:", t.SelfGetName_l())
	log.Printf("Self ID: %X\n", t.Address)
	mystmsg := t.SelfGetStatusMessage_l()
	log.Println("Status:", mystmsg)
	log.Println("------------------------------------------")
	log.Println("Friend Count:", len(flist))

	if len(flist) > 0 {
		log.Println("num\tpublic_key")
	}
	for fn, pk := range flist {
		log.Printf("%d\t%X\n", fn, pk[:])
	}
	if len(flist) > 20 {
		log.Println("Friend Count:", len(flist))
	}
	if len(flist) > 0 {
		log.Println()
	}
}
