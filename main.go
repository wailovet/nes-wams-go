package main

import (
	"github.com/fogleman/nes/nes"
	"log"
	"syscall/js"
)

func main() {

	c := make(chan struct{}, 0)
	println("Go WebAssembly Initialized!")

	console := start()
	ta := js.Global().Get("Uint8Array").New(256 * 240 * 4)

	js.Global().Set("nesButton", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		s := [8]bool{
			args[0].Bool(),
			args[1].Bool(),
			args[2].Bool(),
			args[3].Bool(),
			args[4].Bool(),
			args[5].Bool(),
			args[6].Bool(),
			args[7].Bool(),
		}

		console.SetButtons1(s)
		return false
	}))
	js.Global().Set("nes", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		console.StepFrame()
		img := console.Buffer()

		js.CopyBytesToJS(ta, img.Pix)
		return ta
	}))
	<-c
}

var fps = 0

func start() nes.Console {
	file, err := assetFS().Open("test.nes")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer file.Close()
	cartridge, err := LoadNES(file)
	ram := make([]byte, 2048)
	controller1 := nes.NewController()
	controller2 := nes.NewController()
	console := nes.Console{
		nil, nil, nil, cartridge, controller1, controller2, nil, ram}
	mapper, err := nes.NewMapper(&console)
	if err != nil {
		log.Fatalln(err.Error())
	}
	console.Mapper = mapper
	console.CPU = nes.NewCPU(&console)
	console.APU = nes.NewAPU(&console)
	console.PPU = nes.NewPPU(&console)
	return console
}
