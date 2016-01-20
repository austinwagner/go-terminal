package terminal

import "testing"
import "time"

func TestWindowSize(t * testing.T) {
    origSize, err := WindowSize()
    if err != nil {
        t.Fatal(err)
    }

    err = SetWindowSize(100, 50);
    if err != nil {
        t.Fatal(err)
    }

    time.Sleep(100 * time.Millisecond)

    size, err := WindowSize()
    if err != nil {
        t.Fatal(err)
    }

    if size.Width != 100 || size.Height != 50 {
        t.Fatalf("bad size: got (%d, %d) expected (100, 50)", size.Width, size.Height)
    }

    err = SetWindowSizeFromSizeInfo(origSize)
    if err != nil {
        t.Fatal(err)
    }
}

func TestCursorPosition(t * testing.T) {
    origPos, err := CursorPosition()
    if err != nil {
        t.Fatal(err)
    }

    err = MoveCursor(5, 6)
    if err != nil {
        t.Fatal(err)
    }

    time.Sleep(100 * time.Millisecond)

    pos, err := CursorPosition()
    if err != nil {
        t.Fatal(err)
    }

    if pos.X != 5 || pos.Y != 6 {
        t.Fatalf("bad pos: got (%d, %d) expected (5, 6)", pos.X, pos.Y)
    }

    err = MoveCursorToPoint(origPos)
    if err != nil {
        t.Fatal(err)
    }
}

func TestIsTty(t * testing.T) {
    if !IsTty() {
        t.Fatal("IsTty returned false, expected true")
    }
}
