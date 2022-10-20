package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	webrtcvad "audio-denoise"
	"github.com/cryptix/wav"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("usage: example infile.wav")
	}

	filename := flag.Arg(0)

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
	rate := int(wavInfo.SampleRate)
	if wavInfo.Channels != 1 {
		log.Fatal("expected mono file")
	}
	if rate != 32000 {
		// log.Fatalf("expected 32kHz file:%v", rate)
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

	vad, err := webrtcvad.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := vad.SetMode(3); err != nil {
		log.Fatal(err)
	}

	checkFrameNums := rate / 1000 * 10
	frame := make([]byte, checkFrameNums)

	if ok := vad.ValidRateAndFrameLength(rate, checkFrameNums); !ok {
		log.Fatal("invalid rate or frame length")
	}

	var isActive bool
	var offset int

	report := func() {
		t := time.Duration(offset) * time.Second / time.Duration(rate) / 2
		fmt.Printf("isActive = %v, t = %v\n", isActive, t)
	}

	for {
		_, err := io.ReadFull(reader, frame)
		if err == io.EOF {
			// log.Printf("eof\n")
			break
		}
		if err == io.ErrUnexpectedEOF {
			// log.Printf("unexpected eof")
			break
		}
		// if err == io.EOF || err == io.ErrUnexpectedEOF {
		// 	break
		// }
		if err != nil {
			log.Fatal(err)
		}

		frameActive, err := vad.Process(rate, frame)
		if err != nil {
			log.Fatal(err)
		}

		if isActive != frameActive || offset == 0 {
			ms := time.Duration(int((float64(offset)/float64(rate))*1000000/2)/1000) * time.Millisecond
			// t := time.Duration(offset) * time.Second / time.Duration(rate) / 2
			if !isActive && frameActive {
				fmt.Printf("[%v-", ms)
			} else if isActive && !frameActive {
				fmt.Printf("%v]是人声\n", ms)
			}
			isActive = frameActive
			// report()
		}

		offset += len(frame)
	}

	report()
}
