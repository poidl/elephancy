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
var templateCacheFile = resourcedir + "/static_cache.json"
var templateCacheFileFingerprinted = fingerprintdir + "/static_cache_fingerprinted.json"

type siteData struct {
	Titlemobile  string
	Titledesktop string
}

type cacheData struct {
	Buttonpic  string
	Script     string
	Stylesheet string
}

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

// loadCacheData opens a json file and returns the contents as a struct
// TODO: handle errors
func loadCacheData(filename string) (cacheData, error) {

	bytearr, err := ioutil.ReadFile(filename)
	if err != nil {
		return cacheData{}, err
	}
	var m cacheData
	err = json.Unmarshal(bytearr, &m)
	if err != nil {
		return cacheData{}, err
	}
	return m, nil
}

func loadSiteData(filename string) (siteData, error) {

	bytearr, err := ioutil.ReadFile(filename)
	if err != nil {
		return siteData{}, err
	}
	var m siteData
	err = json.Unmarshal(bytearr, &m)
	if err != nil {
		return siteData{}, err
	}
	return m, nil
}

func writecacheData(cacheData cacheData, filename string) {

	data, err := json.Marshal(cacheData)
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
	tcf := templateCacheFile
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
	tcffp := templateCacheFileFingerprinted
	cp(tcf, tcffp)

	// load the resource data
	resource, err := loadCacheData(tcf)
	if err != nil {
		log.Fatal(err)
	}

	resource.Buttonpic = createFingerprintedResource(resource.Buttonpic)
	resource.Script = createFingerprintedResource(resource.Script)
	resource.Stylesheet = createFingerprintedResource(resource.Stylesheet)
	// write to file holding fingerprinted resources
	writecacheData(resource, tcffp)
}

func GenerateFingerprintedTemplate(ftmpl string, ftmplFingerpr string) {
	tmpl := template.New(path.Base(ftmpl))
	tmpl = tmpl.Delims("[[", "]]")
	tmpl, err := tmpl.ParseFiles(ftmpl)
	if err != nil {
		log.Fatal(err)
	}
	tcffp := templateCacheFileFingerprinted
	cachedat, err := loadCacheData(tcffp)
	if err != nil {
		log.Fatal(err)
	}
	// log.Fatal("debug")
	f, err := os.Create(ftmplFingerpr)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(f, &cachedat)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}

func FillTitle(ftemplFingerpr string) {
	// /////////////////////////////////
	templ := template.New(path.Base(ftemplFingerpr))
	templ = templ.Delims("||", "||")
	templ, err := templ.ParseFiles(ftemplFingerpr)
	if err != nil {
		log.Fatal(err)
	}
	tcffp := "../frontendserver/data/site.json"
	cachedat, err := loadSiteData(tcffp)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(ftemplFingerpr)
	if err != nil {
		log.Fatal(err)
	}

	err = templ.Execute(f, &cachedat)
	if err != nil {
		log.Fatal(err)
	}
}
