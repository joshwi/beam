package collect

import (
	"fmt"
	"strings"
	"time"

	"example.com/pkg/utils"
	"github.com/barasher/go-exiftool"
)

type Image struct {
	Make         string
	Model        string
	Lens         string
	LensModel    string
	FocalLength  string
	FileName     string
	FilePath     string
	FileNumber   string
	FileType     string
	FileSize     string
	ISO          string
	ShutterSpeed string
	Aperture     string
	Megapixels   string
	CreateDate   string
	ModifyDate   string
}

func ReadImg(directory string, filepath string) (utils.Image, error) {

	path := fmt.Sprintf(`%v%v`, directory, filepath)

	metadata := utils.Image{}

	output := map[string]string{}

	exif, err := exiftool.NewExiftool()
	if err != nil {
		fmt.Printf("Error when intializing: %v\n", err)
		return metadata, err
	}
	defer exif.Close()

	data := exif.ExtractMetadata(path)

	// logger.I.Debug().Interface("data", data).Msg("DEBUG")

	for _, info := range data {

		if info.Err != nil {
			return metadata, info.Err

		}

		for k, v := range info.Fields {
			output[k] = fmt.Sprintf("%v", v)
		}
	}

	raw_date, _ := time.Parse("2006:01:02 15:04:05", output["CreateDate"])
	year := raw_date.Year()
	date := raw_date.Format("2006-01-02")

	metadata = utils.Image{
		Make:         output["Make"],
		Model:        output["Model"],
		Lens:         output["Lens"],
		LensModel:    output["LensModel"],
		FocalLength:  output["FocalLength"],
		FilePath:     filepath,
		FileName:     output["FileName"],
		FileNumber:   output["FileNumber"],
		FileType:     strings.ToUpper(output["FileTypeExtension"]),
		FileSize:     output["FileSize"],
		ISO:          output["ISO"],
		ShutterSpeed: output["ShutterSpeed"],
		Aperture:     output["Aperture"],
		Megapixels:   output["Megapixels"],
		Date:         date,
		Year:         year,
		CreateDate:   output["CreateDate"],
		ModifyDate:   output["ModifyDate"],
	}

	// logger.I.Debug().Interface("metadata", metadata).Msg("DEBUG")

	return metadata, nil

}
