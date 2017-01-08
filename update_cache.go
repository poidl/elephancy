package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// update and fingerprint each file in the staticcache. TODO: make it optional
// which files to update
var resourcepath = "./frontend/staticcache/resources"
var fingerprintpath = "./frontend/staticcache/fingerprinted"

const chunkSize = 64000

func deepCompare(file1, file2 string) bool {
	// by Pith on SO: http: //stackoverflow.com/questions/29505089/how-can-i-compare-two-files-in-golang
	// Check file size ...
	chunkSize := 64000
	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}
		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

// func filesDiffer(file1 string, file2 string) bool {
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)
// 	f1, err := os.Open(file1)
// 	defer f1.Close()
// 	_, err = f1.Stat()
// 	if err != nil {
// 		log.Fatal("Error opening " + file1)
// 	}
// 	f2, err := os.Open(file2)
// 	defer f2.Close()
// 	_, err = f2.Stat()
// 	println(file2)
// 	if err != nil {
// 		log.Fatal("Error opening " + file2)
// 	}
// 	cmd := "diff"
// 	args := []string{file1, file2}
// 	var bout []byte
// 	println(file1)
// 	println(file2)
// 	bout, err = exec.Command(cmd, args...).Output()
// 	// if bout, err = exec.Command(cmd, args...).Output(); err != nil {
// 	// println(err)
// 	// log.Fatal("Error executing diff in bash")
// 	// }
// 	out := string(bout)
// 	println(out)
// 	log.Fatal(err)
// 	if out == "" {
// 		return false
// 	}
// 	return true
// }

func mv(src string, dest string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd := "mv"
	args := []string{src, dest}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.Fatal("Moving files failed")
	}
}
func fingerprint(fname string) string {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd := "./scripts/fingerprint.sh"
	args := []string{fname}
	var bout []byte
	var err error
	if bout, err = exec.Command(cmd, args...).Output(); err != nil {
		log.Fatal("Error executing fingerprint in bash")
	}
	fout := strings.TrimSpace(string(bout))
	return fout
}

func setupcache() {
	filename := "./frontend/staticcache/cache.json"
	stylesheet := resourcepath + "/mystyle.css"
	ccoll, err := loadJson(filename)
	if err != nil {
		log.Fatal(err)
	}
	pag := ccoll["Frame.html"].(map[string]interface{})
	stylesheetOld := pag["Stylesheet"].(string)
	stylesheetOld = fingerprintpath + "/" + stylesheetOld
	// log.Println("Reading " + filename ". Found entry \\\"Stylesheet)
	_, err = os.Stat(stylesheetOld)
	if err == nil {
		if deepCompare(stylesheet, stylesheetOld) {
			log.Println("Stylesheet " + stylesheetOld + " exists and is up to date.")
		} else {
			log.Println("Stylesheet " + stylesheetOld + " exists but is not up to date. Removing " + stylesheetOld + ".")
			os.Remove(stylesheetOld)
		}
	}
	_, err = os.Stat(stylesheetOld)
	if os.IsNotExist(err) {
		stylesheetNew := fingerprint(resourcepath + "/mystyle.css")
		f1, err := os.Stat(stylesheetNew)
		if err != nil {
			log.Fatal(err)
		}
		basename := f1.Name()
		mv(stylesheetNew, fingerprintpath)
		log.Print("Stylesheet named " + stylesheetOld + " does not exists. Created " + fingerprintpath + "/" + basename)
		pag["Stylesheet"] = basename
		writeJson(filename, ccoll)
	}

	// switch _, err := os.Stat(stylesheetOld); err {
	// case nil:
	// 	println("Nil")
	// case os.IsNotExist:
	// 	println("not exist")
	// }
}
