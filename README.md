ASSIMP-GO
=========

This is a [Go][golang] library that wraps the use of the Open Asset Import Library
known as [assimp][assimp-link].

UNDER CONSTRUCTION
==================

The library is currently in an alpha state, but you are welcome to see how
it's progressing.

Requirements
------------

This does require `cgo` which means that `gcc` should be in your path
when trying to build using this library.

Software requirements:

* [Assimp][assimp-link] version 3.1 - tested with this version
* [Mathgl][mgl] - for 3d math
* [Gombz][gombz-link] - used as a file format and data structure for
  the information pulled from assimp.

TODO
----

The following need to be addressed in order to start releases:

* documentation
* api comments
* samples
* more data from assimp like animations


LICENSE
=======

Assimp-go is released under the BSD license. See the [LICENSE][license-link] file for more details.


[golang]: https://golang.org/
[license-link]: https://raw.githubusercontent.com/tbogdala/assimp-go/master/LICENSE
[assimp-link]: http://assimp.sourceforge.net/
[mgl]: https://github.com/go-gl/mathgl
[gombz-link]: https://github.com/tbogdala/gombz
