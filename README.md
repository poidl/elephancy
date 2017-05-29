# elephancy

A beginner's exercise in full-stack web development.

Elephancy is a tiny content display system. The objective of this project is
to learn

* Http
* Server-side programming (currently Go)
* Client-side programming (Typescript/Javascript/CSS/HTML)
* Possibly data bases in the future, currently it's "flat-file" only.

### Installation

Note that the dependencies listed below are *only* necessary on your local 
machine. Once everything is compiled on the local machine, a sinlge Go 
binary and all the resources can be uploaded to the server. No further 
dependencies need to be installed on the server.

##### Client-side frontend

The Client-side frontend is written in Typescript. To compile this into javascript that runs in a web-browser, you will need:

- typescript
- webpack
- awesome-typescript-loader
- source-map-loader

To install these, go to the directory `frontendclient/typescript` and run

```
npm install -g webpack  
npm install --save-dev typescript awesome-typescript-loader source-map-loader
```

Furthermore, Typescript needs to know about the standard node.js modules:

```
npm install --save @types/node
```

Then run `webpack` to generate the Javascript code

```
./webpack
```

##### Backend and server-side frontend

Both are written in Go. In the folder `/elephancy`, run

```
go build
./elephancy
```

Open http://localhost:8080/ in a browser.

