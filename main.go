package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/pflag"
)

// MVC(Model, View, Controller)
// The model has your state and can be initialized with defaults
// The view is the way you show your current model state
// The controller changes your model state

var Release string

type ContentBuffer struct {
	Filename string
	Content  []rune
}

type Pager struct {
	screen           tcell.Screen
	buffers          []ContentBuffer
	cursorX, cursorY int
	viewOffset       int
	screenWidth      int
	screenHeight     int
	currentFileIndex int
	showNumbers      bool
	relativeNumbers  bool
	quitAtEOF        bool
}

const tabWidth = 4

func main() {
	// Flags using pflag
	initialViewOffset := pflag.IntP("offset", "o", 0, "Initial view offset")
	showVersion := pflag.BoolP("version", "v", false, "Show version number")
	showNumbers := pflag.BoolP("numbers", "n", false, "Show line numbers")
	relativeNumbers := pflag.BoolP("relative", "r", false, "Show relative line numbers")
	quitAtEOF := pflag.BoolP("quit", "q", false, "Quits when EOF is reached")
	pflag.Parse()

	// Show version and exit if the version flag is set
	if *showVersion {
		fmt.Println("pager version", Release)
		return
	}

	// Get file names from command-line arguments
	fileNames := pflag.Args()

	// Initialize buffers
	buffers := []ContentBuffer{}

	// Read from stdin if available
	stdinStat, _ := os.Stdin.Stat()
	if stdinStat.Mode()&os.ModeCharDevice == 0 {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("Failed to read input", err)
		}
		buffers = append(buffers, ContentBuffer{
			Filename: "stdin",
			Content:  []rune(string(content)),
		})
	}

	// Read from files if file arguments are provided
	if len(fileNames) > 0 {
		buffers = append(buffers, readFilesOrStdin(fileNames)...)
	}

	// Check if buffers are empty
	if len(buffers) == 0 {
		fmt.Println("Nothing to display.")
		pflag.Usage()
		return
	}

	// Ensure initial offset does not exceed the file length
	for i, buffer := range buffers {
		lines := len(strings.Split(string(buffer.Content), "\n"))
		if *initialViewOffset >= lines {
			fmt.Printf("Warning: Initial offset %d exceeds file length %d for file %s. Setting to last line.\n", *initialViewOffset, lines, buffer.Filename)
			*initialViewOffset = lines - 1
		}
		buffers[i] = buffer
	}

	// Initialize screen
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal("Failed to create new screen", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatal("Failed to initialize new screen", err)
	}
	defer screen.Fini()

	// Get screen size
	screenWidth, screenHeight := screen.Size()

	// Create Pager instance
	pager := Pager{
		screen:           screen,
		buffers:          buffers,
		cursorX:          0,
		cursorY:          *initialViewOffset,
		viewOffset:       *initialViewOffset,
		screenWidth:      screenWidth,
		screenHeight:     screenHeight,
		currentFileIndex: 0,
		showNumbers:      *showNumbers,
		relativeNumbers:  *relativeNumbers,
		quitAtEOF:        *quitAtEOF,
	}

	// Initial draw
	pager.drawContent()

	// Handle key events
	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
			pager.screenWidth, pager.screenHeight = screen.Size()
			pager.drawContent()
		case *tcell.EventKey:
			switch {
			case ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' || ev.Rune() == 'Q' || ev.Rune() == 'Z':
				return
			case ev.Key() == tcell.KeyLeft || ev.Rune() == 'h':
				if pager.cursorX > 0 {
					pager.cursorX--
					pager.drawContent()
				}
			case ev.Key() == tcell.KeyRight || ev.Rune() == 'l':
				lineLength := len([]rune(strings.Split(string(pager.buffers[pager.currentFileIndex].Content), "\n")[pager.cursorY]))
				if pager.cursorX < lineLength {
					pager.cursorX++
					pager.drawContent()
				}
			case ev.Key() == tcell.KeyDown || ev.Rune() == 'j':
				if pager.cursorY < len(strings.Split(string(pager.buffers[pager.currentFileIndex].Content), "\n"))-1 {
					pager.cursorY++
					if pager.cursorY-pager.viewOffset >= pager.screenHeight {
						pager.viewOffset++
					}
					lineLength := len([]rune(strings.Split(string(pager.buffers[pager.currentFileIndex].Content), "\n")[pager.cursorY]))
					if pager.cursorX > lineLength {
						pager.cursorX = lineLength
					}
					pager.drawContent()
					// Handle EOF
					if pager.quitAtEOF && pager.cursorY == len(strings.Split(string(pager.buffers[pager.currentFileIndex].Content), "\n"))-1 {
						return
					}
				}
			case ev.Key() == tcell.KeyUp || ev.Rune() == 'k':
				if pager.cursorY > 0 {
					pager.cursorY--
					if pager.cursorY < pager.viewOffset {
						pager.viewOffset--
					}
					lineLength := len([]rune(strings.Split(string(pager.buffers[pager.currentFileIndex].Content), "\n")[pager.cursorY]))
					if pager.cursorX > lineLength {
						pager.cursorX = lineLength
					}
					pager.drawContent()
				}
			case ev.Key() == tcell.KeyCtrlU || ev.Key() == tcell.KeyPgUp:
				pager.cursorY -= 10
				if pager.cursorY < 0 {
					pager.cursorY = 0
				}
				if pager.cursorY < pager.viewOffset {
					pager.viewOffset -= 10
					if pager.viewOffset < 0 {
						pager.viewOffset = 0
					}
				}
				pager.drawContent()
			case ev.Key() == tcell.KeyCtrlD || ev.Key() == tcell.KeyPgDn:
				pager.cursorY += 10
				if pager.cursorY >= len(strings.Split(string(pager.buffers[pager.currentFileIndex].Content), "\n")) {
					pager.cursorY = len(strings.Split(string(pager.buffers[pager.currentFileIndex].Content), "\n")) - 1
				}
				if pager.cursorY-pager.viewOffset >= pager.screenHeight {
					pager.viewOffset += 10
				}
				pager.drawContent()
			case ev.Rune() == 'g':
				pager.cursorY = 0
				pager.viewOffset = 0
				pager.drawContent()
			case ev.Rune() == 'G':
				pager.cursorY = len(strings.Split(string(pager.buffers[pager.currentFileIndex].Content), "\n")) - 1
				pager.viewOffset = pager.cursorY - pager.screenHeight + 1
				if pager.viewOffset < 0 {
					pager.viewOffset = 0
				}
				pager.drawContent()
			case ev.Rune() == 'n':
				if pager.currentFileIndex < len(pager.buffers)-1 {
					pager.currentFileIndex++
					pager.cursorX, pager.cursorY = 0, *initialViewOffset
					pager.viewOffset = *initialViewOffset
					pager.drawContent()
				}
			case ev.Rune() == 'p' || ev.Rune() == 'b':
				if pager.currentFileIndex > 0 {
					pager.currentFileIndex--
					pager.cursorX, pager.cursorY = 0, *initialViewOffset
					pager.viewOffset = *initialViewOffset
					pager.drawContent()
				}
			case ev.Rune() == 'r':
				pager.drawContent()
			}
		}
	}
}

// Function to read files
func readFilesOrStdin(fileNames []string) []ContentBuffer {
	var buffers []ContentBuffer

	// Read from files
	for _, filename := range fileNames {
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filename, err)
			continue
		}
		buffers = append(buffers, ContentBuffer{
			Filename: filename,
			Content:  []rune(string(content)),
		})
	}

	return buffers
}

// Pager method to draw content on the screen
func (pager *Pager) drawContent() {
	pager.screen.Clear()
	y := 0
	screenCursorY := pager.cursorY - pager.viewOffset // Adjust cursor position
	content := pager.buffers[pager.currentFileIndex].Content

	// Calculate the line number width once
	lineNumberWidth := len(fmt.Sprintf("%d", len(strings.Split(string(content), "\n"))))

	// Start rendering from viewOffset
	for lineNum := pager.viewOffset; lineNum < len(strings.Split(string(content), "\n")); lineNum++ {
		if y >= pager.screenHeight {
			break
		}

		line := strings.Split(string(content), "\n")[lineNum]
		x := 0

		// Display line numbers if enabled
		if pager.showNumbers {
			lineNumber := fmt.Sprintf("%*d ", lineNumberWidth, lineNum+1)
			for _, r := range lineNumber {
				if x >= pager.screenWidth {
					break
				}
				pager.screen.SetContent(x, y, r, nil, tcell.StyleDefault)
				x++
			}
		}

		// Utility function to calculate the absolute difference
		abs := func(x int) int {
			if x < 0 {
				return -x
			}
			return x
		}

		// Display relative line numbers if enabled
		if pager.relativeNumbers && !pager.showNumbers {
			relativeNumber := fmt.Sprintf("%*d ", lineNumberWidth, abs(pager.cursorY-lineNum))
			for _, r := range relativeNumber {
				if x >= pager.screenWidth {
					break
				}
				pager.screen.SetContent(x, y, r, nil, tcell.StyleDefault)
				x++
			}
		}

		// Set the view content with line wrapping
		for _, r := range line {
			if r == '\t' {
				// Handle tab character by converting it to spaces
				for i := 0; i < tabWidth; i++ {
					if x >= pager.screenWidth {
						x = lineNumberWidth + 1 // Add padding for wrapped lines
						y++
						if y >= pager.screenHeight {
							break
						}
					}
					pager.screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
					x++
				}
			} else {
				if x >= pager.screenWidth {
					x = lineNumberWidth + 1 // Add padding for wrapped lines
					y++
					if y >= pager.screenHeight {
						break
					}
				}
				pager.screen.SetContent(x, y, r, nil, tcell.StyleDefault)
				x++
			}
		}
		y++
	}

	// Ensure cursor does not move over the line numbers
	if pager.showNumbers || pager.relativeNumbers {
		if pager.cursorX < lineNumberWidth+1 {
			pager.cursorX = lineNumberWidth + 1
		}
	}

	// Show cursor if it's within the visible screen area
	if screenCursorY >= 0 && screenCursorY < pager.screenHeight {
		pager.screen.ShowCursor(pager.cursorX, screenCursorY)
	} else {
		pager.screen.HideCursor()
	}

	pager.screen.Show()
}
