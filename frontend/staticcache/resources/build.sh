errexit() {
    echo "*MYERROR* "$1 1>&2
    exit 1
}
builddir="./jsbuild"

echo "tsc..."
tsc || errexit "tsc failed"

echo "babel..."
babel --no-babelrc $builddir"/script.js" > $builddir"/scriptbabel.js" || errexit "babel failed" 

echo "browserify..."
browserify  $builddir"/scriptbabel.js" > $builddir"/scriptbabelbrowser.js" || errexit "browserify failed" 
