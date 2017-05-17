package frontendserver

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// update and fingerprint each file in the staticcache. TODO: make it optional
// which files to update
var resourcedir = "../frontendclient/resources"
var fingerprintdir = "../frontendclient/resources_fingerprinted"
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

func FingerprintFile(fname string) string {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// check if the resource file exists
	fp, err := os.Stat(fname)
	if os.IsNotExist(err) {
		log.Fatal(err)
	}
	// get the basename
	basename := fp.Name()
	// check if the name already contains the tag
	tag := "hashstart"
	if strings.Contains(basename, tag) {
		log.Fatal("Resource name contains the tag \"" + tag + "\"")
	}
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))[:15]
	suffix := filepath.Ext(basename)
	basename_nosuffix := basename[0 : len(basename)-len(suffix)]
	basename_new := basename_nosuffix + "_" + tag + hash + suffix
	fname_new := fname[0:len(fname)-len(basename)] + basename_new
	cp(fname, fname_new)
	return fname_new
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
	fullname := resourcedir + name
	fpf := FingerprintFile(fullname)
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
	return "/staticcache" + relpath + fpfBasename
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