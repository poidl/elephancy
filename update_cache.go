package main

import (
	"html/template"
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

func getCacheResources(fname string) TemplateMap {
	m, err := loadJSONStruct(templateCacheFile)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func setupcacheNew() {
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
	resource, err := loadJSONmsi(tcf)
	if err != nil {
		log.Fatal(err)
	}

	// load the fingerprinted data
	resourceFP, err := loadJSONmsi(tcffp)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range resource {

		// fingerprint resource
		fpf := fingerprint(v.(string))
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
		// update map
		resourceFP[k] = dest
	}
	// write to file holding fingerprinted resources
	writeJson(tcffp, resourceFP)
}

func generateFingerprintedTemplate() {
	tmpl := template.New(path.Base(ftempl))
	tmpl = tmpl.Delims("[[", "]]")
	tmpl, err := tmpl.ParseFiles(ftempl)
	if err != nil {
		log.Fatal(err)
	}
	tcffp := fingerprintdir + "/" + templateCacheFileFingerprinted
	cachedat, err := loadJSONStruct(tcffp)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(ftemplFingerpr)
	err = tmpl.Execute(f, &cachedat)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}
