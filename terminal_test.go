package terminal

import "testing"
import "time"

func TestWindowSize(t * testing.T) {
    err := SetWindowSize(100, 50);
    if err != nil {
        t.Fatal(err)
    }

    time.Sleep(time.Second)

    size, err := WindowSize()
    if err != nil {
        t.Fatal(err)
    }

    if size.Width != 100 || size.Height != 50 {
        t.Fatal("Mismatched size")
    }
}

func TestCursorPosition(t * testing.T) {
    err := MoveCursor(5, 6)
    if err != nil {
        t.Fatal(err)
    }

    time.Sleep(time.Second)

    pos, err := CursorPosition()
    if err != nil {
        t.Fatal(err)
    }

    if pos.X != 5 || pos.Y != 6 {
        t.Fatal("Mismatched position")
    }
}

func TestIsTty(t * testing.T) {
    if !IsTty() {
        t.Fatal("IsTty returned false, expected true")
    }
}
