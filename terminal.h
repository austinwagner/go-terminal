#define bool int
#define true 1
#define false 0

typedef struct _COORDINATE {
	int x;
	int y;
} COORDINATE;

bool GetTerminalWindowSize(COORDINATE * coord);
bool GetCursorPosition(COORDINATE * coord);
bool SetCursorPosition(COORDINATE coord);
bool ClearTerminalWindow();
bool ClearTerminalBuffer();
bool SetTerminalWindowSize(COORDINATE coord);
bool IsTty();