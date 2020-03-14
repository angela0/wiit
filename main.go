package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anacrolix/torrent/metainfo"
)

type Torrent struct {
	Creator  string   `json:"creator"`
	Comment  string   `json:"comment"`
	Date     string   `json:"date"`
	Encoding string   `json:"encoding"`
	Hash     string   `json:"hash"`
	Name     string   `json:"name"`
	Files    []string `json:"files"`
	Size     int64    `json:"size"`
	Magnet   string   `json:"magnet"`
}

var (
	WithAll      bool
	WithCreator  bool
	WithComment  bool
	WithDate     bool
	WithEncoding bool
	WithHash     bool
	WithName     bool
	WithFile     bool
	WithSize     bool
	WithMagnet   bool

	UseJson bool
)

func main() {

	flag.BoolVar(&WithAll, "a", false, "print all, default behavior")
	flag.BoolVar(&WithCreator, "o", false, "print creator")
	flag.BoolVar(&WithComment, "c", false, "print comment")
	flag.BoolVar(&WithDate, "d", false, "print date")
	flag.BoolVar(&WithEncoding, "e", false, "print encoding")
	flag.BoolVar(&WithHash, "i", false, "print hash")
	flag.BoolVar(&WithName, "n", false, "print name")
	flag.BoolVar(&WithFile, "f", false, "print files")
	flag.BoolVar(&WithSize, "s", false, "print size")
	flag.BoolVar(&WithMagnet, "m", false, "print magnet")
	flag.BoolVar(&UseJson, "json", false, "print all using json format")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n    %s [options] <torrent file>\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	mi, err := metainfo.LoadFromFile(os.Args[len(os.Args)-1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading metainfo from stdin: %s", err)
		return
	}
	info, err := mi.UnmarshalInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error unmarshalling info: %s", err)
		return
	}

	var files []string
	for _, file := range info.Files {
		files = append(files, file.DisplayPath(&info))
	}
	torrentInfo := Torrent{
		Creator:  mi.CreatedBy,
		Comment:  mi.Comment,
		Date:     time.Unix(mi.CreationDate, 0).String(),
		Encoding: mi.Encoding,
		Hash:     mi.HashInfoBytes().String(),
		Name:     info.Name,
		Files:    files,
		Size:     info.TotalLength(),
		Magnet:   mi.Magnet(info.Name, mi.HashInfoBytes()).String(),
	}

	if WithCreator {
		fmt.Println(torrentInfo.Creator)
	}
	if WithComment {
		fmt.Println(torrentInfo.Comment)
	}
	if WithDate {
		fmt.Println(torrentInfo.Date)
	}
	if WithEncoding {
		fmt.Println(torrentInfo.Encoding)
	}
	if WithHash {
		fmt.Println(torrentInfo.Hash)
	}
	if WithName {
		fmt.Println(torrentInfo.Name)
	}
	if WithFile {
		fmt.Println(strings.Join(torrentInfo.Files, "\n"))
	}
	if WithSize {
		fmt.Println(torrentInfo.Size)
	}
	if WithMagnet {
		fmt.Println(torrentInfo.Magnet)
	}

	if WithCreator || WithComment || WithDate || WithEncoding || WithHash || WithName || WithFile || WithSize || WithMagnet {
		return
	}

	if UseJson {
		v, err := json.Marshal(torrentInfo)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(v))
		return
	}

	fmt.Printf("Creator\t\t%s\n", torrentInfo.Creator)
	fmt.Printf("Comment\t\t%s\n", torrentInfo.Comment)
	fmt.Printf("Date\t\t%s\n", torrentInfo.Date)
	fmt.Printf("Encode\t\t%s\n", torrentInfo.Encoding)
	fmt.Printf("Hash\t\t%s\n", torrentInfo.Hash)
	fmt.Printf("Name\t\t%s\n", torrentInfo.Name)
	fmt.Printf("Files\t\t%s\n", torrentInfo.Files)
	fmt.Printf("Size\t\t%d\n", torrentInfo.Size)
	fmt.Printf("Magnet\t\t%s\n", torrentInfo.Magnet)
}
