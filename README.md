ASSIMP-GO
=========

This is a [Go][golang] library that wraps the use of the Open Asset Import Library
known as [assimp][assimp-link].


Requirements
------------

This does require `cgo` which means that `gcc` should be in your path
when trying to build using this library.

Software requirements:

* [Assimp][assimp-link] version 3.1 - tested with this version
* [Mathgl][mgl] - for 3d math
* [Gombz][gombz-link] - used as a file format and data structure for
  the information pulled from assimp.

Usage
-----

The module can be used to load files supported by [Assimp][assimp-link] and converts
them into the [Gombz][gombz-link] meshes. Once imported, you can load an file
(e.g. .OBJ or .FBX file) by using the following call:


```
srcMeshes, err := assimp.ParseFile(srcFilepath)
```



LICENSE
=======

Assimp-go is released under the BSD license. See the [LICENSE][license-link] file for more details.


[golang]: https://golang.org/
[license-link]: https://raw.githubusercontent.com/tbogdala/assimp-go/master/LICENSE
[assimp-link]: http://assimp.sourceforge.net/
[mgl]: https://github.com/go-gl/mathgl
[gombz-link]: https://github.com/tbogdala/gombz
