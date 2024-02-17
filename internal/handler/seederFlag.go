package handler

import (
	"github.com/alecthomas/kingpin"
	"os"
)

func HandleArgs() (string, int) {
	application := kingpin.New("Seeds", "For Seeds")
	seedsMethod := application.Flag("seeds", "Name Seeds Method if greater than 1 use coma ','").Required().String()
	fill := application.Flag("fill", "For Total Seeds to fill database").Default("100").Short('f').Int()
	if *fill == 0 {
		*fill = 100
	}
	kingpin.MustParse(application.Parse(os.Args[1:]))
	return *seedsMethod, *fill
}
