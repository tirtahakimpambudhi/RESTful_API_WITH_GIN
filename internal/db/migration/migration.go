package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"go_gin/internal/domain/model"
	"log"
	"os"
)

var (
	Migration     = kingpin.New("migration", "Testing For Migration")
	Driver        = Migration.Flag("driver", "Driver Database Like e mysql,postgres,sqlite,etc...").Short('d').Required().String()
	ConnectString = Migration.Flag("connect", "Database Connection String").Short('s').Required().String()
	Command       = Migration.Flag("command", "Command Goose Like e create,up,down,status,etc").Short('c').Required().String()
	Directory     = Migration.Flag("dir", "Directory Database Migration").Required().String()
	File          = Migration.Flag("file", "File Database Migration").Short('f').Required().String()
	Version       = Migration.Flag("version", "Version File Required if Command 'up-to' and 'down-to'").Short('v').Int64()
)

func main() {
	kingpin.MustParse(Migration.Parse(os.Args[1:]))
	fmt.Printf("%v %v %v", *Directory, *File, *ConnectString)
	conn := model.NewConnection(*Version, *Driver, *ConnectString, *File, *Directory)
	err := conn.Run(*Command)
	if err != nil {
		log.Fatal(err.Error())
	}
}
