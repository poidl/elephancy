package frontend

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// update and fingerprint each file in the staticcache. TODO: make it optional
// which files to update
var resourcedir = "./frontend/staticcache/resources"
var fingerprintdir = "./frontend/staticcache/fingerprinted"
var templateCacheFile = "simple_cache.json"
var templateCacheFileFingerprinted = "simple_cache_fingerprinted.json"

func cp(src string, dest string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd := "cp"
	args := []string{src, dest}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.Fatal("Copying files failed")
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

func fingerprintFile(fname string) string {
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

type TemplateData struct {
	Buttonpic  string
	Script     string
	Stylesheet string
}

// loadTemplateData opens a json file and returns the contents as a struct
// TODO: handle errors
func loadTemplateData(filename string) (TemplateData, error) {

	bytearr, err := ioutil.ReadFile(filename)
	if err != nil {
		return TemplateData{}, err
	}
	var m TemplateData
	err = json.Unmarshal(bytearr, &m)
	if err != nil {
		return TemplateData{}, err
	}
	return m, nil
}

func writeTemplateData(TemplateData TemplateData, filename string) {

	data, err := json.Marshal(TemplateData)
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal("Writing " + filename + " failed.")
	}
}

func createFingerprintedResource(name string) string {
	// fingerprint resource
	fpf := fingerprintFile(name)
	fpffd, err := os.Stat(fpf)
	if err != nil {
		log.Fatal(err)
	}

	// Copy to fingerprint directory
	fpfBasename := fpffd.Name()
	src := fpf
	relpath := fpf[len(resourcedir) : len(fpf)-len(fpfBasename)]
	err = os.MkdirAll(fingerprintdir+relpath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	dest := fingerprintdir + relpath + fpfBasename
	mv(src, dest)
	return dest
}

func SetupcacheNew() {
	// check if file defining cache resources exists
	tcf := resourcedir + "/" + templateCacheFile
	_, err := os.Stat(tcf)
	if err != nil {
		log.Fatal(err)
	}
	// create directory for fingerprinted resources if it doesn't exist
	err = os.MkdirAll(fingerprintdir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	// create file holding fingerprinted resources, overwriting it if it exists
	tcffp := fingerprintdir + "/" + templateCacheFileFingerprinted
	cp(tcf, tcffp)

	// load the resource data
	resource, err := loadTemplateData(tcf)
	if err != nil {
		log.Fatal(err)
	}

	resource.Buttonpic = createFingerprintedResource(resource.Buttonpic)
	resource.Script = createFingerprintedResource(resource.Script)
	resource.Stylesheet = createFingerprintedResource(resource.Stylesheet)
	// write to file holding fingerprinted resources
	writeTemplateData(resource, tcffp)
}

func GenerateFingerprintedTemplate(ftmpl string, ftmplFingerpr string) {
	tmpl := template.New(path.Base(ftmpl))
	tmpl = tmpl.Delims("[[", "]]")
	tmpl, err := tmpl.ParseFiles(ftmpl)
	if err != nil {
		log.Fatal(err)
	}
	tcffp := fingerprintdir + "/" + templateCacheFileFingerprinted
	cachedat, err := loadTemplateData(tcffp)
	if err != nil {
		log.Fatal(err)
	}
	// log.Fatal("debug")
	f, err := os.Create(ftmplFingerpr)
	err = tmpl.Execute(f, &cachedat)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}
