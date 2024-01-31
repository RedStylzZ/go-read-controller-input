package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

type JSType uint8

const (
	JSEventButton = 0x01
	JSEventAxis   = 0x02
	JSEventInit   = 0x80
)

type JSEvent struct {
	Time   uint32
	Value  int16
	Type   JSType
	Number uint8
}

func readInput(r io.Reader) (JSEvent, error) {
	var event JSEvent
	order := binary.LittleEndian

	err := binary.Read(r, order, &event.Time)
	if err != nil {
		return JSEvent{}, err
	}

	err = binary.Read(r, order, &event.Value)
	if err != nil {
		return JSEvent{}, err
	}

	err = binary.Read(r, order, &event.Type)
	if err != nil {
		return JSEvent{}, err
	}

	err = binary.Read(r, order, &event.Number)
	if err != nil {
		return JSEvent{}, err
	}

	return event, nil
}

var buttons = map[uint8]string{
	0:  "A",
	1:  "B",
	3:  "X",
	4:  "Y",
	6:  "LB",
	7:  "RB",
	10: "Share",
	11: "Start",
	13: "L3",
	14: "R3",
}

var axis = map[uint8]string{
	0: "LeftAxis",
	1: "LeftAxis",
	2: "RightAxis",
	3: "RightAxis",
	4: "RT",
	5: "LT",
	6: "DPad",
	7: "DPad",
}

func printEvent(event JSEvent) {
	// fmt.Printf("Time: %d\nValue: %d\nType: 0x%X\nNumber: %d\n", event.Time, event.Value, event.Type, event.Number)
	// fmt.Printf("IsInit: %t\n", event.Type == (JSEventButton|JSEventInit))

	switch event.Type {
	case JSEventButton:
		state := "pressed"
		if event.Value == 0 {
			state = "released"
		}
		fmt.Printf("Button '%s' %s\n", buttons[event.Number], state)
	case JSEventAxis:
		switch event.Number {
		case 4, 5:
			state := "pressed"
			if event.Value < 0 {
				state = "released"
			}
			fmt.Printf("Axis '%s' %s: %d\n", axis[event.Number], state, event.Value)
		case 6, 7:
			if event.Value == 0 {
				fmt.Printf("Axis '%s' released\n", axis[event.Number])
				break
			}

			switch event.Number {
			case 6:
				if event.Value > 0 {
					fmt.Printf("Axis '%s' pressed right: %d\n", axis[event.Number], event.Value)
				} else {
					fmt.Printf("Axis '%s' pressed left: %d\n", axis[event.Number], event.Value)
				}
			case 7:
				if event.Value > 0 {
					fmt.Printf("Axis '%s' pressed down: %d\n", axis[event.Number], event.Value)
				} else {
					fmt.Printf("Axis '%s' pressed up: %d\n", axis[event.Number], event.Value)
				}
			}

		default:
			fmt.Printf("Pressed '%s': %d\n", axis[event.Number], event.Value)
		}
	}
}

func main() {
	f, err := os.Open("/dev/input/js0")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for {
		m := make([]byte, 64)
		_, err := f.Read(m)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		r := bytes.NewReader(m)

		event, err := readInput(r)
		if err != nil {
			break
		}

		printEvent(event)
	}
}
