package main

//  tox save data decrypt/encrypt

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/TokTok/go-toxcore-c"
	"github.com/TokTok/go-toxcore-c/toxenums"
)

func init() {
	log.SetFlags(log.Flags() ^ log.Ldate ^ log.Ltime | log.Lshortfile)
}

var tsfile string = "tox_save.tox"
var pass string
var tofile string = "./tsdec.bin"

var crypt_mode string // enc/dec
var decrypt_mode = false
var encrypt_mode = false

func printHelp() {
	log.Println("For help: /path/to/tsdec [options] <tsfile> -h")
}

func main() {
	flag.StringVar(&tsfile, "tsfile", tsfile, "tox save data file")
	flag.StringVar(&pass, "pass", pass, "tox save data password")
	flag.StringVar(&tofile, "tofile", tofile, "result file")

	flag.Parse()

	log.Println(tsfile)
	log.Println(pass)

	data, err := ioutil.ReadFile(tsfile)
	if err != nil {
		log.Println(err)
		return
	}

	isencrypt := tox.IsDataEncrypted(data)
	log.Println("Is encrypt: ", isencrypt)
	if isencrypt && pass == "" {
		log.Println("Need -pass argument.")
		flag.Usage()
		return
	}

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
	log.Println(t)
	log.Println("Self Name:", t.SelfGetName_l())
	log.Printf("Self ID: %X\n", t.Address)
	mystmsg := t.SelfGetStatusMessage_l()
	log.Println("Status:", mystmsg)
	log.Println("------------------------------------------")
	log.Println("Friend Count:", len(flist))

	if _, err = os.Stat(tofile); err != nil {
		log.Println(err)
		// return
	}

	if isencrypt { // do decrypt
		log.Println("Decrypting...")
		err := tox.WriteSavedata(tofile, t.GetSavedata_l())
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Save decrypted data OK: ", tofile)
		}
	} else { // do encrypt
		log.Println("Encrypting...")
		pakey, err := tox.Derive([]byte(pass))
		defer pakey.Free()
		if err != nil {
			log.Println(err)
			return
		}
		encdata, err := pakey.Encrypt(data)
		if err != nil {
			log.Println(err)
			return
		}

		err = ioutil.WriteFile(tofile, encdata, 0755)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Save encrypted data OK: ", tofile)
		}
	}
}

func encrypt() {

}

// @param pt plain tox
func decrypt(pt *tox.Tox) {

}
