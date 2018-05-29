# This makefile is used to re-generate the wrapper

MKFILE_PATH:=$(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR:=$(dir $(MKFILE_PATH))

openzwave.go: openzwave.i CallbackAdapter.h
	swig -go -cgo -intgosize 64 -c++ -I$(MKFILE_DIR)/open-zwave/cpp/src openzwave.i

# perl -pi -e 's|#define intgo swig_intgo|#cgo CXXFLAGS: -std=c++0x -fwrapv -O4 -DNDEBUG -I/usr/include/ortools -DARCH_K8 -Wno-deprecated -DUSE_CBC -DUSE_CLP -DUSE_GLOP -DUSE_BOP\n#cgo LDFLAGS: -L/usr/lib -lortools -lrt -lpthread\n#define intgo swig_intgo|' ortools.go

