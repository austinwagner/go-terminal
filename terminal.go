package terminal

/*
#define bool int
#define true 1
#define false 0

typedef struct _COORDINATE {
	int x;
	int y;
} COORDINATE;


#if defined(__unix__) || defined(__CYGWIN__) || (defined(__APPLE__) && defined(__MACH__))

#include <termios.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>
#include <sys/ioctl.h>
#include <stdio.h>

int getch(int fd) {
	char buffer[1];
	if (read(fd, buffer, 1) <= 0) {
		return -1;
	}

	return buffer[0];
}

bool putch(int fd, char c) {
	char buffer[1] = {c};
	if (write(fd, buffer, 1) != 1) {
		return false;
	}

	return true;
}

bool putstr(int fd, const char * str) {
	size_t len = strlen(str);

	const char * head = str;
	const char * tail = str + len;

	while (head < tail) {
		ssize_t bytesRead = write(fd, head, tail - head);

		if (bytesRead > 0) {
			head += bytesRead;
		}
		else if (bytesRead != (ssize_t)-1 && errno != EINTR && errno != EAGAIN) {
			return false;
		}
	}

	return true;
}

bool GetTerminalWindowSize(COORDINATE * coord) {
	struct winsize size;
	if (ioctl(STDOUT_FILENO, TIOCGWINSZ, &size) != 0) {
		return false;
	}

	coord->x = size.ws_col;
	coord->y = size.ws_row;
	return true;
}

// http://www.linuxquestions.org/questions/programming-9/get-cursor-position-in-c-947833/
bool GetCursorPosition(COORDINATE * coord) {
	bool result = false;
	struct termios ti_orig;
	struct termios ti_new;

	coord->x = 0;
	coord->y = 0;

	if (tcgetattr(STDOUT_FILENO, &ti_orig) != 0) {
		return false;
	}

	memcpy(&ti_new, &ti_orig, sizeof(ti_orig));

	ti_new.c_lflag &= ~ICANON;
	ti_new.c_lflag &= ~ECHO;
	ti_new.c_cflag &= ~CREAD;

	if (tcsetattr(STDOUT_FILENO, TCSADRAIN, &ti_new) != 0) {
		goto cleanup;
	}

	if (putstr(STDOUT_FILENO, "\033[6n")) {
		goto cleanup;
	}

	int c = getch(STDOUT_FILENO);
	if (c != 27) {
		goto cleanup;
	}

	c = getch(STDOUT_FILENO);
	if (c != '[') {
		goto cleanup;
	}

	c = getch(STDOUT_FILENO);
	while ('0' <= c && c <= '9') {
		coord->y = (10 * coord->y) + (c - '0');
		c = getch(STDOUT_FILENO);
	}

	if (c != ';') {
		goto cleanup;
	}

	c = getch(STDOUT_FILENO);
	while ('0' <= c && c <= '9') {
		coord->x = (10 * coord->x) + (c - '0');
		c = getch(STDOUT_FILENO);
	}

	if (c != 'R') {
		goto cleanup;
	}

	result = true;

cleanup:
	tcsetattr(STDOUT_FILENO, TCSADRAIN, &ti_orig);

	return result;
}

// http://stackoverflow.com/a/101613/370556
int ipow(int base, int exp) {
	int result = 1;
	while (exp) {
		if (exp & 1) {
			result *= base;
		}
		exp >>= 1;
		base *= base;
	}

	return result;
}

bool SetCursorPosition(COORDINATE coord) {
	bool result = false;

	COORDINATE adjustedCoord;
	adjustedCoord.x = coord.x + 1;
	adjustedCoord.y = coord.y + 1;

	struct termios ti_orig;
	struct termios ti_new;

	if (tcgetattr(STDOUT_FILENO, &ti_orig) != 0) {
		return false;
	}

	memcpy(&ti_new, &ti_orig, sizeof(ti_orig));

	ti_new.c_lflag &= ~ICANON;
	ti_new.c_lflag &= ~ECHO;
	ti_new.c_cflag &= ~CREAD;

	if (tcsetattr(STDOUT_FILENO, TCSADRAIN, &ti_new) != 0) {
		goto cleanup;
	}

	char output[64];
	sprintf(output, "\x1B[%d;%d;", adjustedCoord.y, adjustedCoord.x);

	if (!putstr(STDOUT_FILENO, output)) {
		goto cleanup;
	}

	result = true;

cleanup:
	tcsetattr(STDOUT_FILENO, TCSADRAIN, &ti_orig);

	return result;
}

bool ClearTerminalWindow() {
	bool result = false;

	struct termios ti_orig;
	struct termios ti_new;

	if (tcgetattr(STDOUT_FILENO, &ti_orig) != 0) {
		return false;
	}

	memcpy(&ti_new, &ti_orig, sizeof(ti_orig));

	ti_new.c_lflag &= ~ICANON;
	ti_new.c_lflag &= ~ECHO;
	ti_new.c_cflag &= ~CREAD;

	if (tcsetattr(STDOUT_FILENO, TCSADRAIN, &ti_new) != 0) {
		goto cleanup;
	}

	if (!putstr(STDOUT_FILENO, "\x1B[2J")) {
		goto cleanup;
	}

	result = true;

cleanup:
	tcsetattr(STDOUT_FILENO, TCSADRAIN, &ti_orig);

	return result;
}

bool ClearTerminalBuffer() {
	bool result = false;

	struct termios ti_orig;
	struct termios ti_new;

	if (tcgetattr(STDOUT_FILENO, &ti_orig) != 0) {
		return false;
	}

	memcpy(&ti_new, &ti_orig, sizeof(ti_orig));

	ti_new.c_lflag &= ~ICANON;
	ti_new.c_lflag &= ~ECHO;
	ti_new.c_cflag &= ~CREAD;

	if (tcsetattr(STDOUT_FILENO, TCSADRAIN, &ti_new) != 0) {
		goto cleanup;
	}

	if (!putstr(STDOUT_FILENO, "\x1B[3J")) {
		goto cleanup;
	}

	result = true;

cleanup:
	tcsetattr(STDOUT_FILENO, TCSADRAIN, &ti_orig);

	return result;
}

bool IsTty() {
	return isatty(STDOUT_FILENO);
}

bool SetTerminalWindowSize(COORDINATE coord) {
	bool result = false;

	struct termios ti_orig;
	struct termios ti_new;

	if (tcgetattr(STDOUT_FILENO, &ti_orig) != 0) {
		return false;
	}

	memcpy(&ti_new, &ti_orig, sizeof(ti_orig));

	ti_new.c_lflag &= ~ICANON;
	ti_new.c_lflag &= ~ECHO;
	ti_new.c_cflag &= ~CREAD;

	if (tcsetattr(STDOUT_FILENO, TCSADRAIN, &ti_new) != 0) {
		goto cleanup;
	}

	char output[64];
	sprintf(output, "\x1B[8;%d;%dt", coord.y, coord.x);

	if (!putstr(STDOUT_FILENO, output)) {
		goto cleanup;
	}

	result = true;

cleanup:
	tcsetattr(STDOUT_FILENO, TCSADRAIN, &ti_orig);

	return result;
}

#elif defined(WIN32) || defined(_WIN32) || defined(__WIN32)

#include <Windows.h>

bool GetTerminalWindowSize(COORDINATE * coord) {
	HANDLE stdout = GetStdHandle(STD_OUTPUT_HANDLE);
	if (stdout == INVALID_HANDLE_VALUE) {
		return false;
	}

	CONSOLE_SCREEN_BUFFER_INFO csbi;
	if (!GetConsoleScreenBufferInfo(stdout, &csbi)) {
		return false;
	}

	coord->x = csbi.srWindow.Right - csbi.srWindow.Left + 1;
	coord->y = csbi.srWindow.Bottom - csbi.srWindow.Top + 1;
	return true;
}

bool GetCursorPosition(COORDINATE * coord) {
	HANDLE stdout = GetStdHandle(STD_OUTPUT_HANDLE);
	if (stdout == INVALID_HANDLE_VALUE) {
		return false;
	}

	CONSOLE_SCREEN_BUFFER_INFO csbi;
	if (!GetConsoleScreenBufferInfo(stdout, &csbi)) {
		return false;
	}

	coord->x = csbi.dwCursorPosition.X - csbi.srWindow.Left;
	coord->y = csbi.dwCursorPosition.Y - csbi.srWindow.Top;

	return true;
}

bool SetCursorPosition(COORDINATE coord) {
	HANDLE stdout = GetStdHandle(STD_OUTPUT_HANDLE);
	if (stdout == INVALID_HANDLE_VALUE) {
		return false;
	}

	CONSOLE_SCREEN_BUFFER_INFO csbi;
	if (!GetConsoleScreenBufferInfo(stdout, &csbi)) {
		return false;
	}

	COORD c;
	c.X = csbi.srWindow.Left + coord.x;
	c.Y = csbi.srWindow.Top + coord.y;

	if (c.X > csbi.srWindow.Right) {
		c.X = csbi.srWindow.Right;
	}

	if (c.Y > csbi.srWindow.Bottom) {
		c.Y = csbi.srWindow.Bottom;
	}

	return SetConsoleCursorPosition(stdout, c);
}

// https://support.microsoft.com/en-us/kb/99261
bool ClearTerminalWindow() {
	HANDLE stdout = GetStdHandle(STD_OUTPUT_HANDLE);
	if (stdout == INVALID_HANDLE_VALUE) {
		return false;
	}

	DWORD cCharsWritten;
	CONSOLE_SCREEN_BUFFER_INFO csbi;
	DWORD dwConSize;

	if (!GetConsoleScreenBufferInfo(stdout, &csbi)) {
		return false;
	}

	COORD coordWindow;
	coordWindow.X = csbi.srWindow.Left;
	coordWindow.Y = csbi.srWindow.Top;

	dwConSize =
		(csbi.srWindow.Right - csbi.srWindow.Left + 1) *
		(csbi.srWindow.Bottom - csbi.srWindow.Top + 1);

	if (!FillConsoleOutputCharacterW(stdout, L' ', dwConSize, coordWindow, &cCharsWritten)) {
		return false;
	}

	if (!GetConsoleScreenBufferInfo(stdout, &csbi)) {
		return false;
	}

	if (!FillConsoleOutputAttribute(stdout, csbi.wAttributes, dwConSize, coordWindow, &cCharsWritten)) {
		return false;
	}

	return true;
}

bool ClearTerminalBuffer() {
	HANDLE stdout = GetStdHandle(STD_OUTPUT_HANDLE);
	if (stdout == INVALID_HANDLE_VALUE) {
		return false;
	}

	COORD coordScreen = { 0, 0 };
	DWORD cCharsWritten;
	CONSOLE_SCREEN_BUFFER_INFO csbi;
	DWORD dwConSize;

	if (!GetConsoleScreenBufferInfo(stdout, &csbi)) {
		return false;
	}

	dwConSize = csbi.dwSize.X * csbi.dwSize.Y;

	if (!FillConsoleOutputCharacterW(stdout, L' ', dwConSize, coordScreen, &cCharsWritten)) {
		return false;
	}

	if (!GetConsoleScreenBufferInfo(stdout, &csbi)) {
		return false;
	}

	if (!FillConsoleOutputAttribute(stdout, csbi.wAttributes, dwConSize, coordScreen, &cCharsWritten)) {
		return false;
	}

	if (!SetConsoleCursorPosition(stdout, coordScreen)) {
		return false;
	}

	return true;
}

bool IsTty() {
	HANDLE stdout = GetStdHandle(STD_OUTPUT_HANDLE);
	if (stdout == INVALID_HANDLE_VALUE) {
		return false;
	}

	DWORD mode;
	return GetConsoleMode(stdout, &mode);
}

bool SetTerminalWindowSize(COORDINATE coord) {
	HANDLE stdout = GetStdHandle(STD_OUTPUT_HANDLE);
	if (stdout == INVALID_HANDLE_VALUE) {
		return false;
	}

	CONSOLE_SCREEN_BUFFER_INFO csbi;
	if (!GetConsoleScreenBufferInfo(stdout, &csbi)) {
		return false;
	}

	SMALL_RECT rect;
	rect.Top = csbi.srWindow.Top;
	rect.Left = csbi.srWindow.Left;
	rect.Bottom = rect.Top + coord.y - 1;
	rect.Right = rect.Left + coord.x - 1;

	return SetConsoleWindowInfo(stdout, TRUE, &rect);
}

#else

#error Unsupported OS

#endif
*/
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
