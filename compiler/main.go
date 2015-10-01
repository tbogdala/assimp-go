// Copyright 2015, Timothy Bogdala <tdb@animal-machine.com>
// See the LICENSE file for more details.

package main

/*

	Compiler is a module for assimp-go that creates a command-line program to take
	input files that ASSIMP supports and compile them into the GOMBZ binary format.

	Known limitations:
		* Only one mesh per file is supported.

*/

import (
	"flag"
	"fmt"
	"github.com/tbogdala/assimp-go"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	srcFileFlag := flag.String("src", "", "the source file to compile")
	quietFlag := flag.Bool("q", false, "silences the normal output of the compiler")
	flag.Parse()

	// sanity checks
	if len(*srcFileFlag) < 1 {
		fmt.Printf("ERROR: at least one -src file must be specified.\n")
		os.Exit(1)
	}

	// what are we compiling
	var srcFile string = *srcFileFlag

	if *quietFlag == false {
		fmt.Printf("Compiling file: %s\n", srcFile)
	}

	srcMeshes, err := assimp.ParseFile(srcFile)
	if err != nil {
		fmt.Printf("ERROR: failed to load the source file (%s) for the ComponentMesh.\n%v\n", srcFile, err)
		os.Exit(1)
	}

	// send out a warning if we have more than one mesh returned from the assimp parsing
	// TODO: support more than one mesh
	numOfSrcMeshes := len(srcMeshes)
	if numOfSrcMeshes > 1 {
		if *quietFlag == false {
			fmt.Printf("WARNING: source file has %d meshes. Only one mesh is supported!\n", numOfSrcMeshes)
		}
	}

	// set which mesh to use
	srcMesh := srcMeshes[0]

	// construct the output path
	ext := path.Ext(srcFile)
	outfile := srcFile[0:len(srcFile)-len(ext)] + ".gombz"

	// encode the file
	meshBytes, err := srcMesh.Encode()
	if err != nil {
		fmt.Printf("ERROR: failed to encode the mesh.")
		os.Exit(1)
	}

	// we've encoded, now write the file out
	err = ioutil.WriteFile(outfile, meshBytes, os.ModePerm)
	if err != nil {
		fmt.Printf("ERROR: failed to write mesh to file: %s\n", outfile)
		os.Exit(1)
	}

	if *quietFlag == false {
		fmt.Printf("Wrote mesh: %s\n", outfile)
	}
}
