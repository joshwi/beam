module example.com

go 1.23.0

replace example.com/pkg => ./pkg

require (
	github.com/barasher/go-exiftool v1.10.0
	github.com/rs/zerolog v1.29.1
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)
