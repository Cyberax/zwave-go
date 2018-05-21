# This makefile is used to re-generate the wrapper

MKFILE_PATH:=$(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR:=$(dir $(MKFILE_PATH))

openzwave.go: openzwave.i
	swig -go -cgo -intgosize 64 -c++ -I$(MKFILE_DIR)/open-zwave/cpp/src openzwave.i
