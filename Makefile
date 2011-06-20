include $(GOROOT)/src/Make.inc

TARG=anagramizer
GOFILES=status.go\
	anagramizer.go\
	wordsorter.go\

include $(GOROOT)/src/Make.cmd
