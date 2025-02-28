package system

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"example.com/pkg/logger"
	"example.com/pkg/utils"
)

// Scan a directory for files and subfolders
func Scan(directory string) ([]string, error) {

	output := []string{}

	err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		rel_path := strings.ReplaceAll(path, directory, "")
		output = append(output, rel_path)
		return err
	})

	if err != nil {
		logger.E.Error().Str("directory", directory).Str("status", "Failed").Err(err).Msg("Scan")
		return nil, err
	} else {
		logger.I.Info().Str("directory", directory).Str("status", "Success").Msg("Scan")
	}

	return output, nil
}

// Scan a directory for files and subfolders
func ScanDir(directory string, filetypes string) ([]string, error) {

	output := []string{}

	var REG_FILE_MATCH = regexp.MustCompile(fmt.Sprintf(`(?i)\.(%v)$`, filetypes))

	err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		file_match := REG_FILE_MATCH.FindAllString(path, -1)
		if !f.IsDir() && len(file_match) > 0 {
			rel_path := strings.ReplaceAll(path, directory, "")
			output = append(output, rel_path)
		}
		return err
	})

	if err != nil {
		logger.E.Error().Str("directory", directory).Str("status", "Failed").Err(err).Msg("Scan")
		return nil, err
	} else {
		logger.I.Info().Str("directory", directory).Str("status", "Success").Msg("Scan")
	}

	return output, nil
}

func Copy(source string, target string) error {
	srcFile, _ := os.Open(source)
	defer srcFile.Close()

	_, err := os.Stat(target)
	if os.IsNotExist(err) {
		// Creates any directories in the path that don't exist
		err = os.MkdirAll(path.Dir(target), 0755)
		if err != nil {
			logger.E.Error().Str("source", source).Str("target", target).Str("status", "Failed").Err(err).Msg("Copy")
			return err
		}
	}

	destFile, _ := os.Create(target) // creates if file doesn't exist
	defer destFile.Close()

	// Move the file to new location
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		logger.E.Error().Str("source", source).Str("target", target).Str("status", "Failed").Err(err).Msg("Copy")
		return err
	}

	logger.I.Info().Str("source", source).Str("target", target).Str("status", "Success").Msg("Copy")

	return nil
}

// Move a file to a new directory
func Move(source string, destination string) error {

	// Check if the file path exists
	_, err := os.Stat(destination)
	if os.IsNotExist(err) {
		// Creates any directories in the path that don't exist
		err = os.MkdirAll(path.Dir(destination), 0755)
		if err != nil {
			logger.E.Error().Str("source", source).Str("target", destination).Str("status", "Failed").Err(err).Msg("Move")
			return err
		}
	}

	// Move the file to new location
	err = os.Rename(source, destination)
	if err != nil {
		logger.E.Error().Str("source", source).Str("target", destination).Str("status", "Failed").Err(err).Msg("Move")
		return err
	}

	logger.I.Info().Str("source", source).Str("target", destination).Str("status", "Success").Msg("Move")

	return nil

}

// Read contents of a file
func Read(filename string) ([]byte, error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		logger.E.Error().Str("file", filename).Str("status", "Failed").Err(err).Msg("Read")
		return nil, err
	} else {
		logger.I.Info().Str("file", filename).Str("status", "Success").Msg("Read")
	}

	return data, nil

}

// Write contents to a file
func Write(filename string, data []byte, mode int) error {

	// Check if file already exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// Creates any directories that don't exist
		err = os.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			logger.E.Error().Str("file", filename).Str("status", "Failed").Err(err).Msg("Write")
			return err
		}
		// Creates file
		_, err = os.Create(filename)
		if err != nil {
			logger.E.Error().Str("file", filename).Str("status", "Failed").Err(err).Msg("Write")
			return err
		}
	}

	// Writes byte data to the file
	err = os.WriteFile(filename, data, os.FileMode(mode))
	if err != nil {
		logger.E.Error().Str("file", filename).Str("status", "Failed").Err(err).Msg("Write")
		return err
	} else {
		logger.I.Info().Str("file", filename).Str("status", "Success").Msg("Write")
	}

	return nil

}

func Backup(source string, target string, filetypes string, subfolder string) (int, int) {

	VALIDATE := regexp.MustCompile(fmt.Sprintf(`(?i)\.(%v)$`, filetypes))

	selected_files := []utils.Tag{}

	start := time.Now()

	filetree, err := Scan(source)
	if err != nil {
		log.Println(err)
	}

	// Look through results of directory scan
	for _, item := range filetree {
		// Does the filetype match specified types?
		match := VALIDATE.FindString(item)
		if len(match) > 0 {
			_, err = os.Stat(target + subfolder + item)
			if os.IsNotExist(err) {
				selected_files = append(selected_files, utils.Tag{Name: source + item, Value: target + subfolder + item})
			}
		}
	}

	// Create channels for data flow and error reporting
	files := make(chan utils.Tag, 10)
	errs := make(chan error)

	// Input selected files into channel
	go func() {
		for _, entry := range selected_files {
			files <- entry
		}
	}()

	// Run worker to copy files from source to target
	for i := 0; i < cap(files); i++ {
		go func(files chan utils.Tag, errs chan error) {
			for item := range files {
				err := Copy(item.Name, item.Value)
				if err != nil {
					errs <- err
				}
				errs <- nil
			}
		}(files, errs)
	}

	err_list := []error{}
	counter := 0

	// Count up errors
	for range selected_files {
		entry := <-errs
		err_list = append(err_list, entry)
		if entry == nil {
			counter++
		}
	}

	// Quick mafs
	end := time.Now()
	elapsed := end.Sub(start)
	duration := fmt.Sprintf("%v", elapsed.Round(time.Second/1000))
	percent := 100.0
	if len(selected_files) > 0 {
		percent = (float64(counter) / float64(len(selected_files))) * 100
	}

	success := fmt.Sprintf("%v%%", percent)

	logger.I.Info().Str("source", source).Str("target", target).Str("types", filetypes).Str("types", filetypes).Str("duration", duration).Str("success", success).Int("files", counter).Int("total", len(selected_files)).Msg("Backup")

	// Close channels
	close(files)
	close(errs)

	return counter, len(err_list)

}

func Transfer(source string, target string, filetypes string, subfolder string) (int, int) {

	VALIDATE := regexp.MustCompile(fmt.Sprintf(`(?i)\.(%v)$`, filetypes))

	selected_files := []utils.Tag{}

	start := time.Now()

	filetree, err := Scan(source)
	if err != nil {
		log.Println(err)
	}

	// Look through results of directory scan
	for _, item := range filetree {
		// Does the filetype match specified types?
		match := VALIDATE.FindString(item)
		if len(match) > 0 {
			_, err = os.Stat(target + subfolder + item)
			if os.IsNotExist(err) {
				selected_files = append(selected_files, utils.Tag{Name: source + item, Value: target + subfolder + item})
			}
		}
	}

	// Create channels for data flow and error reporting
	files := make(chan utils.Tag, 10)
	errs := make(chan error)

	// Input selected files into channel
	go func() {
		for _, entry := range selected_files {
			files <- entry
		}
	}()

	// Run worker to copy files from source to target
	for i := 0; i < cap(files); i++ {
		go func(files chan utils.Tag, errs chan error) {
			for item := range files {
				err := Move(item.Name, item.Value)
				if err != nil {
					errs <- err
				}
				errs <- nil
			}
		}(files, errs)
	}

	err_list := []error{}
	counter := 0

	// Count up errors
	for range selected_files {
		entry := <-errs
		err_list = append(err_list, entry)
		if entry == nil {
			counter++
		}
	}

	// Quick mafs
	end := time.Now()
	elapsed := end.Sub(start)
	duration := fmt.Sprintf("%v", elapsed.Round(time.Second/1000))
	percent := 100.0
	if len(selected_files) > 0 {
		percent = (float64(counter) / float64(len(selected_files))) * 100
	}

	success := fmt.Sprintf("%v%%", percent)

	logger.I.Info().Str("source", source).Str("target", target).Str("types", filetypes).Str("types", filetypes).Str("duration", duration).Str("success", success).Int("files", counter).Int("total", len(selected_files)).Msg("Transfer")

	// Close channels
	close(files)
	close(errs)

	return counter, len(err_list)

}
