package terminal

// #include "terminal.h"
import "C"

import "errors"

const (
	ctrue C.int = 1
	cfalse C.int = 0
)

type SizeInfo struct {
	Width int
	Height int
}

type Point struct {
	X int
	Y int
}

func WindowSize() (SizeInfo, error) {
	coord := new(C.COORDINATE)
	if C.GetTerminalWindowSize(coord) == cfalse {
		return SizeInfo{0, 0}, errors.New("Failed to get terminal size.")
	}

	return SizeInfo{int(coord.x), int(coord.y)}, nil
}

func SetWindowSize(x, y int) error {
	coord := C.COORDINATE { C.int(x), C.int(y) }
	if C.SetTerminalWindowSize(coord) == cfalse {
		return errors.New("Failed to set terminal size.")
	}

	return nil
}

func SetWindowSizeFromSizeInfo(size SizeInfo) error {
	return SetWindowSize(size.Width, size.Height)
}

func CursorPosition() (Point, error) {
	coord := new(C.COORDINATE)
	if C.GetCursorPosition(coord) == cfalse {
		return Point{0, 0}, errors.New("Failed to get cursor position.")
	}

	return Point{int(coord.x), int(coord.y)}, nil
}

func MoveCursor(x, y int) error {
	coord := C.COORDINATE { C.int(x), C.int(y) }
	if C.SetCursorPosition(coord) == cfalse {
		return errors.New("Failed to set cursor position.")
	}

	return nil
}

func MoveCursorToPoint(point Point) error {
	return MoveCursor(point.X, point.Y)
}

func ClearWindow() error {
	if C.ClearTerminalWindow() == cfalse {
		return errors.New("Failed to clear terminal window.")
	}

	return nil
}

func ClearBuffer() error {
	if C.ClearTerminalBuffer() == cfalse {
		return errors.New("Failed to clear terminal buffer.")
	}

	return nil
}

func IsTty() bool {
	return C.IsTty() == ctrue;
}
