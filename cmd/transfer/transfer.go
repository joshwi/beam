package main

import (
	"flag"
	"fmt"
	"log"

	"example.com/pkg/collect"
	"example.com/pkg/logger"
	"example.com/pkg/storage/system"
)

var (
	source    string
	target    string
	filetypes string
	logpath   string
	subfolder = "/tmp"
)

func init() {

	// Define flag arguments for the application
	flag.StringVar(&source, `s`, ``, `Source directory. Default: [SOURCE]`)
	flag.StringVar(&target, `t`, ``, `Source directory. Default: [TARGET]`)
	flag.StringVar(&filetypes, `f`, `cr2|jpg|png|mp4`, `Filetypes to consider. Default: cr2|arw|jpg|png|mp4`)
	flag.StringVar(&logpath, `l`, `/tmp/beam`, `Logpath directory. Default: /tmp/beam`)
	flag.Parse()

	logger.Init(logpath)

	logger.I.Info().Str("source", source).Str("target", target).Str("filetypes", filetypes).Str("status", "Running").Msg("Transfer")
}

func main() {

	// Backup source directory to target
	pass, total := system.Backup(source, target, filetypes, subfolder)
	if pass != total {
		log.Fatalf("Files failed to copy: %v/%v", pass, total)
	}

	// Scan copied files to get a filetree
	tmp_src := fmt.Sprintf("%v%v", target, subfolder)
	filetree, err := system.ScanDir(tmp_src, filetypes)
	if err != nil {
		log.Fatal(err)
	}

	// Move files into target sub directory sorted by metadata
	for _, path := range filetree {
		// info, err := os.Stat(fmt.Sprintf(`%v%v`, tmp_src, path))
		// if err != nil {
		// 	logger.I.Fatal().Err(err)
		// }
		// log.Println(info.Mode())
		// logger.I.Info().Interface("info", info).Msg("PERMISSIONS")
		metadata, err := collect.ReadImg(tmp_src, path)
		if err != nil {
			logger.I.Fatal().Err(err)
		}
		newpath := fmt.Sprintf("/%v/%v/%v/%v/%v", metadata.Make, metadata.FileType, metadata.Year, metadata.Date, metadata.FileName)
		logger.I.Debug().Str("path", newpath).Msg("DEBUG")
		s := tmp_src + path
		t := target + newpath
		system.Move(s, t)
	}

}
