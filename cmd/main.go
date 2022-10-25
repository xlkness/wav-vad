package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	webrtcvad "audio-denoise"
	"github.com/cryptix/wav"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func vad(wavInfo *wav.File, samples []byte, totalDuraMs float64) {
	vad, err := webrtcvad.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := vad.SetMode(1); err != nil {
		log.Fatal(err)
	}

	// checkFrameNums := rate / 1000 * 20
	// frame := make([]byte, checkFrameNums)
	//
	// if ok := vad.ValidRateAndFrameLength(rate, checkFrameNums); !ok {
	// 	log.Fatal("invalid rate or frame length")
	// }

	rate := int(wavInfo.SampleRate)
	sectionMs := 30
	samplesNum := rate / 1000 * sectionMs
	sectionNum := int(float64(wavInfo.NumberOfSamples) / float64(samplesNum))

	isActive := false
	for i := 0; i < sectionNum; i++ {
		realStart := i * samplesNum * 2
		realEnd := realStart + samplesNum*2
		frameActive, err := vad.Process(rate, samples[realStart:realEnd])
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("[%v-%v] %v\n", realStart/2, realEnd/2, frameActive)

		if isActive != frameActive || i == 0 {
			startMs := int(float64(realStart/2) / float64(wavInfo.NumberOfSamples) * float64(totalDuraMs))
			// endMs := float64(realEnd/2) / float64(wavInfo.NumberOfSamples) * float64(totalDuraMs)
			// ms := time.Duration(int((float64(realStart/2)/float64(rate))*1000000/2)/1000) * time.Millisecond
			// t := time.Duration(offset) * time.Second / time.Duration(rate) / 2
			if !isActive && frameActive {
				fmt.Printf("[%vms-", startMs)
			} else if isActive && !frameActive {
				fmt.Printf("%vms]是人声\n", startMs)
			}
			isActive = frameActive
			// report()
		}
	}
}

func ans(wavInfo *wav.File, samples []byte, totalDuraMs float64) {
	ans, err := webrtcvad.NewAns(int(wavInfo.SampleRate))
	if err != nil {
		panic(err)
	}
	ans.SetMode(3)

	rate := int(wavInfo.SampleRate)
	sectionMs := 10
	samplesNum := rate / 1000 * sectionMs
	sectionNum := int(float64(wavInfo.NumberOfSamples) / float64(samplesNum))

	for i := 0; i < sectionNum; i++ {
		realStart := i * samplesNum * 2
		realEnd := realStart + samplesNum*2
		_, err := ans.Process(rate, samples[realStart:realEnd])
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("[%v-%v] %v\n", realStart/2, realEnd/2, frameActive)
	}
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
	// reader, err := wavReader.GetDumbReader()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	wavInfo := wavReader.GetFile()
	rate := int(wavInfo.SampleRate)
	if wavInfo.Channels != 1 {
		log.Fatal("expected mono file")
	}
	if rate != 32000 {
		// log.Fatalf("expected 32kHz file:%v", rate)
	}

	fmt.Printf("采样数量:%v\n", wavInfo.NumberOfSamples)

	samples := make([]byte, 0)
	framesNum := 0
	for {
		sample, err := wavReader.ReadRawSample()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		framesNum++
		samples = append(samples, sample...)
	}

	fmt.Printf("时长:%vms\n", wavInfo.Duration.Milliseconds())
	fmt.Printf("深度:%v\n", wavInfo.SignificantBits)
	fmt.Printf("帧率:%v\n", wavInfo.SampleRate)
	fmt.Printf("帧数:%v\n", len(samples)/int(wavInfo.SignificantBits))
	fmt.Printf("通道:%v\n", wavInfo.Channels)
	totalDuraMs := float64(wavInfo.NumberOfSamples) / float64(wavInfo.SampleRate/1000)
	fmt.Printf("时长:%vms\n", totalDuraMs)
	if wavInfo.Channels != 1 {
		panic(wavInfo.Channels)
	}

	ans(&wavInfo, samples, totalDuraMs)
	// vad(&wavInfo, samples, totalDuraMs)
}
