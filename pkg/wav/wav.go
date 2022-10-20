package wav

import (
	"io"
	"log"
	"os"

	"github.com/cryptix/wav"
)

type Wav struct {
	*wav.File
	Reader io.Reader
}

func ReadWav(filename string) *Wav {
	info, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	wavReader, err := wav.NewReader(file, info.Size())
	if err != nil {
		log.Fatal(err)
	}
	reader, err := wavReader.GetDumbReader()
	if err != nil {
		log.Fatal(err)
	}

	wavInfo := wavReader.GetFile()

	wavReader.ReadSample()

	return &Wav{
		&wavInfo,
		reader,
	}

	// samples := make([]int32, 0)
	// framesNum := 0
	// for {
	// 	sample, err := wavReader.ReadSample()
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		log.Fatal(err)
	// 	}
	// 	framesNum++
	// 	samples = append(samples, sample)
	// }
	//
	// fmt.Printf("时长:%vms\n", wavInfo.Duration.Milliseconds())
	// fmt.Printf("帧率:%v\n", wavInfo.SampleRate)
	// fmt.Printf("帧数:%v\n", len(samples))
}
