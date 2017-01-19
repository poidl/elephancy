package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// update and fingerprint each file in the staticcache. TODO: make it optional
// which files to update
var templateCacheFile = "./frontend/staticcache/simple_cache.json"
var cachepath = "./frontend/staticcache/cache.json"
var resourcedir = "./frontend/staticcache/resources"
var fingerprintdir = "./frontend/staticcache/fingerprinted"

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

type TemplateMap struct {
	Buttonpic  string
	Script     string
	Stylesheet string
}

func getCacheTemplateData() TemplateMap {
	bla, err := loadJson(cachepath)
	if err != nil {
		log.Fatal(err)
	}

	tmp := TemplateMap{bla["pics/menu.png"].(string), bla["script.js"].(string), bla["mystyle.css"].(string)}
	return tmp
}

func getCacheResources() {
	cachedat, err := loadJson(templateCacheFile)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range cachedat {
		println(k)
		println(v.(string))
	}

	// bla := cachedat["frame.html"].(map[string]interface{})

	// tmp := TemplateMap{bla["pics/menu.png"].(string), bla["script.js"].(string), bla["mystyle.css"].(string)}
	// return tmp
}

func setupcache() {
	err := os.MkdirAll(fingerprintdir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	filename := cachepath
	checkResource("mystyle.css", filename)
	checkResource("script.js", filename)
	checkResource("pics/menu.png", filename)
}

func checkResource(resource string, filename string) {
	template, err := loadJson(filename)
	if err != nil {
		log.Fatal(err)
	}
	resourceFilepath := resourcedir + "/" + resource
	old := template[resource].(string)
	fold, err := os.Stat(old)
	if err == nil {
		log.SetFlags(log.LstdFlags)
		if deepCompare(resourceFilepath, old) {
			log.Println("Resource " + fold.Name() + " is up to date.")
		} else {
			log.Println("Resource " + fold.Name() + " is not up to date, deleting.")
			os.Remove(old)
		}
	}
	_, err = os.Stat(old)
	if os.IsNotExist(err) {
		new := fingerprint(resourceFilepath)
		f1, err := os.Stat(new)
		if err != nil {
			log.Fatal(err)
		}
		basename := f1.Name()
		relpath := new[len(resourcedir) : len(new)-len(basename)]
		err = os.MkdirAll(fingerprintdir+relpath, 0755)
		if err != nil {
			log.Fatal(err)
		}
		mv(new, fingerprintdir+new[len(resourcedir):])
		log.SetFlags(log.LstdFlags)
		log.Println("Created " + basename + ".")
		template[resource] = fingerprintdir + new[len(resourcedir):]
		writeJson(filename, template)
	} else if err != nil {
		log.Fatal(err)
	}

}

func generateFingerprintedTemplate() {
	tmpl := template.New("frame_new.html")
	tmpl = tmpl.Delims("[[", "]]")
	tmpl, err := tmpl.ParseFiles(ftempl)
	if err != nil {
		log.Fatal(err)
	}
	cachedat := getCacheTemplateData()
	f, err := os.Create(ftemplFingerpr)
	err = tmpl.Execute(f, &cachedat)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}
