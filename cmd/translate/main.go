package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/ghostiam/binstruct"
)

// TODO: The bug.
// TODO: Print function, just print every story in full.

// Extract all scripts for all scenes to separate files.
// Extract one scene and print it.
// Extract only dialogue strings from a scene.
// Inject replacement strings to a scene.

var b = binary.LittleEndian

type StoryLoc struct {
	Offset int
	Id     uint32
	Name   string
	Length int
}

type Tool struct {
	data []byte
	ptr  int
}

func (t *Tool) Skip(count int) {
	t.ptr += count
}

func (t *Tool) GetU32() uint32 {
	ret := b.Uint32(t.data[t.ptr : t.ptr+4])
	t.ptr += 4
	return ret
}

func (t *Tool) GetString() string {
	strLen := int(t.GetU32())
	ret := string(t.data[t.ptr : t.ptr+strLen])
	t.ptr += strLen
	return ret
}

// The size parameter is the size of the data being replaced. The size of the
// new data is just len(data).
// Returns the difference in size
func (t *Tool) Patch(data []byte, offset int, size int) int {
	if len(data) == size {
		// Do it in place
		copy(t.data[offset:offset+size], data[:])
	} else {
		// Do whatever the fuck this is. Split it up and stitch back together.
		// Idk if this is the best way to do it but it works.
		before := t.data[0:offset]
		after := t.data[offset+size:]

		// Is this okay??
		newBuffer := make([]byte, len(before)+len(data)+len(after))
		copy(newBuffer[0:offset], before)
		copy(newBuffer[offset:offset+len(data)], data)
		copy(newBuffer[offset+len(data):], after)
		t.data = newBuffer
	}
	return len(data) - size
}

func (t *Tool) PatchU32(val uint32, offset int) {
	b.PutUint32(t.data[offset:offset+4], val)
}

// Add an amount to a u32
func (t *Tool) DiffU32(amount int, offset int) {
	sav := t.ptr
	t.ptr = offset
	old := t.GetU32()
	new := int(old) + amount
	t.PatchU32(uint32(new), offset)
	t.ptr = sav
}

func (t *Tool) Pad(amount int, offset int) {
	// before := t.data[:offset]
	// after := t.data[offset+amount:] // Strange

	zeros := []byte{}
	for range amount {
		zeros = append(zeros, 0)
	}
	slices.Insert(t.data, offset, zeros...)
	// t.data = append(append(before, zeros...), after...)
}

// Returns the difference in size
func (t *Tool) PatchString(s string, offset int) int {
	oldSize := b.Uint32(t.data[offset : offset+4])
	newBytes := []byte(s)
	// fmt.Println("the old size is", int(oldSize))
	// fmt.Println("the new size is", len(newBytes))
	if oldSize != uint32(len(newBytes)) {
		// fmt.Printf("AT %08x: %d -> %d: %s\n", offset, int(oldSize), len(newBytes), s)
	}

	t.PatchU32(uint32(len(newBytes)), offset)
	return t.Patch(newBytes, offset+4, int(oldSize))
}

func (t *Tool) GetStoryHeader() StoryLoc {
	offset := t.ptr
	byteLen := t.GetU32() + 4
	name := t.GetString()
	id := t.GetU32()

	return StoryLoc{
		Offset: offset,
		Length: int(byteLen),
		Name:   name,
		Id:     id,
	}
}

// func (t *Tool) GetStory(loc StoryLoc) {
func (t *Tool) GetStory(w io.Writer) {
	// Assume these are done.
	// t.ptr = loc.Offset

	storyOffset := t.ptr

	hdr := t.GetStoryHeader()
	fmt.Fprintf(w, "-------- %d\t%08x  %s\n", hdr.Id, storyOffset, hdr.Name)

	eventsCount := t.GetU32()

	// There are strings that we want to translate
	// There are string that we don't want to translate but want to include in
	// the extracted file for translation context
	// There are strings we don't want to translate or extract

	storyHasAny := false

	for range eventsCount {
		eventLenOffset := t.ptr
		_ = t.GetU32()
		opcode := OpCode(t.GetU32())
		_ = t.GetU32()
		argvCount := t.GetU32()

		// fmt.Println(opcode.String())

		// Store the argv strings and their offsets here, since we might need
		// to patch some of them.
		argv := []struct {
			str    string
			offset int
		}{}
		for range argvCount {
			offset := t.ptr
			arg := t.GetString()
			argv = append(argv, struct {
				str    string
				offset int
			}{arg, offset})
		}

		// If this event needs to be patched at all.
		hasAny := false

		addLine := func(offset int, s string) {
			if strings.TrimSpace(s) == "" {
				return
			}

			// Include but don't translate file names.
			if strings.HasSuffix(s, ".jpg") || strings.HasSuffix(s, ".png") {
				offset = 0
			}

			if !hasAny {
				if !storyHasAny {
					fmt.Fprintf(w, "STORY %08x\n", storyOffset)
					storyHasAny = true
				}

				fmt.Fprintf(w, "EVENT %08x\n", eventLenOffset)
				hasAny = true
			}

			if offset != 0 {
				s = strings.ReplaceAll(s, "\n", "\\\n")
				fmt.Fprintf(w, "%08x %s\n", offset, s)
			} else {
				fmt.Fprintf(w, "-------- %s\n", s)
			}
		}

		switch opcode {
		case SendBubble:
			// sender is number 4, we don't care about that right now
			sender := argv[4]
			text := argv[6]
			addLine(0, sender.str+":")
			// addLine(sender.offset, sender.str+":")
			addLine(text.offset, text.str)

		case Text:
			text := argv[2]
			// s := text.str

			// Don't include empty strings.
			// if strings.TrimSpace(s) == "" {
			// 	continue
			// }

			// // Include but don't translate file names.
			// if strings.HasSuffix(s, ".jpg") || strings.HasSuffix(s, ".png") {
			// 	text.offset = 0
			// }

			addLine(text.offset, text.str)
		}
	}
}

func (t *Tool) Save(path string) {
	err := os.WriteFile(path, t.data, 0644)
	if err != nil {
		panic(err)
	}
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// TODO: This shouldn't be hardcoded. Can change with updates.
var storiesOffset = 0x2E0E2F

func extractStories(binPath, mapPath, outPath string) {
	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	binMd5, err := FileMd5(binPath)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(f, "BINMD5", binMd5)

	mapMd5, err := FileMd5(mapPath)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(f, "MAPMD5", mapMd5)

	fileData, err := os.ReadFile(binPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// First let's find the offset of all the stories.
	t := Tool{
		data: fileData,
		ptr:  storiesOffset,
	}

	t.GetU32()    // ByteLen
	t.GetString() // ProjectName

	storiesCnt := t.GetU32()

	fmt.Fprintln(f, "-------- NUMBER OF STORIES:", storiesCnt)

	// This is a shallow read loop. Notes where the stories are located but does
	// not parse each one.
	//
	// stories := map[uint32]StoryLoc{}
	// for i := range storiesCnt {
	// 	savPtr := t.ptr
	// 	story := t.GetStory()
	// 	// fmt.Println(story.Offset)
	// 	// stories[story.Id] = story
	// 	stories[i] = story
	// 	t.ptr = savPtr
	// 	t.Skip(story.Length)
	// }

	fmt.Fprintf(f, "TOTAL %08x\n", 0x50)
	fmt.Fprintf(f, "TOTAL %08x\n", storiesOffset)
	// for id, story := range stories {
	for range storiesCnt {
		t.GetStory(f)
	}
}

func FileMd5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func main() {

	// TODO: MD5 check. Note the MD5 of the bin file we extracted the script
	// from, and then refuse to inject any file that doesn't match that. We can
	// even store the MD5 of the filemap.bin file in our script, even tho we
	// don't need that file for the extraction phase.

	cmd := os.Args[1]
	switch cmd {
	case "print":
		// fmt.Println("Usage: toolqi print <data04.bin>")
		printStories(os.Args[2])
		return

	case "extract":
		// fmt.Println("Usage: toolqi extract <data04.bin> <filemap.bin> <script.txt>")
		extractStories(os.Args[2], os.Args[3], os.Args[4])
		return

	case "inject":
		// fmt.Println("Usage: toolqi inject <data04.bin> <filemap.bin> <translation.txt>")

	default:
		fmt.Println("Unknown command:", cmd, " use extract or inject.")
		return
	}

	fileData, err := os.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	// First let's find the offset of all the stories.
	t := Tool{
		data: fileData,
		ptr:  storiesOffset,
	}

	var totalDiff int
	{
		patchFile, err := os.Open(os.Args[4])
		if err != nil {
			panic(err)
		}

		s := bufio.NewScanner(patchFile)
		s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}

			unescape := func(s []byte) []byte {
				return bytes.ReplaceAll(s, []byte("\\\n"), []byte("\n"))
			}

			searchFrom := 0
			for i := bytes.IndexByte(data[searchFrom:], '\n'); i >= 0; i = bytes.IndexByte(data[searchFrom:], '\n') {
				// Is there a backslash before it?
				if data[searchFrom+i-1] == '\\' {
					// Skip the backslash and continue
					searchFrom += i + 1
					continue
				}
				return searchFrom + i + 1, unescape(dropCR(data[0 : searchFrom+i])), nil
			}
			// If we're at EOF, we have a final, non-terminated line. Return it.
			if atEOF {
				return len(data), unescape(dropCR(data)), nil
			}
			// Request more data.
			return 0, nil, nil
		})

		// First read the first two lines, which should contain the checksums
		// for the original files this script was generated from. The modified
		// script can only apply to those files.

		if !s.Scan() {
			panic("no line 1")
		}
		line1 := s.Text()
		if !s.Scan() {
			panic("no line 2")
		}
		line2 := s.Text()

		binMd5 := ""
		mapMd5 := ""
		if _, err := fmt.Sscanf(string(line1), "BINMD5 %s", &binMd5); err != nil {
			panic("invalid line 1. expected BINMD5 <md5>")
		}
		if _, err := fmt.Sscanf(string(line2), "MAPMD5 %s", &mapMd5); err != nil {
			panic("invalid line 2. expected MAPMD5 <md5>")
		}

		if binMd5 == "" || mapMd5 == "" {
			panic("checksums not provided")
		}

		// Now check the checksums against the actual files
		binExpected, err := FileMd5(os.Args[2])
		if err != nil {
			panic(err)
		}
		mapExpected, err := FileMd5(os.Args[3])
		if err != nil {
			panic(err)
		}
		if binMd5 != binExpected {
			panic("bin md5 mismatch")
		}
		if mapMd5 != mapExpected {
			panic("map md5 mismatch")
		}

		// var totalDiff int
		var totalPtrs []uint64
		var storyDiff int
		var storyPtr uint64
		var eventDiff int
		var eventPtr uint64
	scan_loop:
		for s.Scan() {
			line := s.Text()
			switch {
			case strings.HasPrefix(line, "------"):
				continue

			case strings.HasPrefix(line, "STORY "):
				// Apply previous story diff
				if storyPtr != 0 && storyDiff != 0 {
					t.DiffU32(storyDiff, int(storyPtr))
					storyDiff = 0
				}

				storyPtr, err = strconv.ParseUint(line[6:], 16, 64)
				if err != nil {
					panic(err)
				}

				storyPtr += uint64(totalDiff)

				// When a story ends, its final event also ends.
				fallthrough

			case strings.HasPrefix(line, "EVENT "):
				// Apply previous event diff
				if eventPtr != 0 && eventDiff != 0 {
					t.DiffU32(eventDiff, int(eventPtr))
					eventDiff = 0
				}

				eventPtr, err = strconv.ParseUint(line[6:], 16, 64)
				if err != nil {
					panic(err)
				}

				fmt.Printf("new event at %06x (%06x)\n", eventPtr, eventPtr+uint64(totalDiff))
				eventPtr += uint64(totalDiff)

			case strings.HasPrefix(line, "TOTAL "):
				// We apply these at the very end.
				totalPtr, err := strconv.ParseUint(line[6:], 16, 64)
				if err != nil {
					panic(err)
				}

				if totalDiff != 0 {
					panic("expected totalDiff to be 0 here")
				}

				totalPtrs = append(totalPtrs, totalPtr+uint64(totalDiff))

			case line == "\n":
				break scan_loop

			default:
				parts := strings.SplitN(line, " ", 2)
				strPtr, err := strconv.ParseUint(parts[0], 16, 64)
				if err != nil {
					panic(err)
				}

				newStr := ""
				if len(parts) == 2 {
					newStr = parts[1]
				}

				diff := t.PatchString(newStr, int(strPtr)+totalDiff)
				eventDiff += diff
				storyDiff += diff
				totalDiff += diff
			}
		}

		if s.Err() != nil {
			panic(s.Err())
		}

		// Apply all the remaining diffs.
		if eventPtr != 0 && eventDiff != 0 {
			t.DiffU32(eventDiff, int(eventPtr))
		}

		if storyPtr != 0 && storyDiff != 0 {
			t.DiffU32(storyDiff, int(storyPtr))
		}

		// Now apply all the total diffs.
		for _, totalPtr := range totalPtrs {
			t.DiffU32(totalDiff, int(totalPtr))
		}

		fmt.Println(totalDiff)
		t.Save("new_data04.bin")
	}

	{
		// Now the file map.
		fileData, err := os.ReadFile(os.Args[3])
		if err != nil {
			fmt.Println(err)
			return
		}

		t := Tool{
			data: fileData,
			ptr:  0,
		}

		// String url;
		// String suffix;
		// u16 unk;
		// u32 offset;
		// u32 size;
		// u8 flags;

		getOneEntry := func() (size uint32, offset uint32, ptr int) {
			t.GetString()
			t.GetString()
			t.Skip(2)
			ptr = t.ptr
			offset = t.GetU32()
			size = t.GetU32()
			t.Skip(1)
			return
		}

		fileIndex := -1
		for {
			size, offset, ptr := getOneEntry()
			if offset == 0 {
				fileIndex += 1
			}

			if offset == 0 {
				fmt.Printf("data%02d.bin starts with entry of size %d\n", fileIndex, size)
			}

			if fileIndex == 4 {
				// Here we need to modify the entries.
				if offset == 0 {
					// First entry, only modify the size.
					t.DiffU32(totalDiff, ptr+4)
				} else {
					// Everything after that, only modify the offset.
					t.DiffU32(totalDiff, ptr)
				}
			}

			if fileIndex == 5 {
				// We're done.
				break
			}
		}

		t.Save("new_filemap.bin")
	}

}

func printStories(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	rdr := binstruct.NewReaderFromBytes(data, binary.LittleEndian, false)

	var df DataFile
	err = rdr.Unmarshal(&df)
	if err != nil {
		fmt.Println(err)
		return
	}

	// for _, button := range df.System.Title.Buttons {
	// 	fmt.Println("BUTTON:", button.Reserved)
	// }

	// for i, button := range df.System.Buttons {
	// 	fmt.Println(i, button.Name.String(), button.Image1.FileName.String(), button.Image2.FileName.String())
	// }

	for i, story := range df.ScriptData.Stories {
		f := os.Stdout

		fmt.Printf("[%d]\t%d\t%s\n", i, story.ID, story.Name)

		for _, ev_ := range story.Events {
			ev := ev_.Event

			for range ev.Indent + 1 {
				fmt.Fprintf(f, "    ")
			}

			argv := []string{}
			for _, arg := range ev.Argv {
				argv = append(argv, arg.String())
			}

			fmt.Fprintf(f, "%s(%s)\n", ev.Code.String(), strings.Join(argv, ", "))
		}

		fmt.Println()
		fmt.Println()
		fmt.Println()
	}
}
