package main

// SizedString represents a length‑prefixed UTF‑8 string.
type SizedString struct {
	Len  uint32
	Data string `bin:"len:Len"`
}

// -------------------- Basic Types --------------------

type Viewport struct {
	Reserved uint32
	X        int32
	Y        int32
	Width    int32
	Height   int32
}

type FileName struct {
	Reserved uint32
	From     uint32
	FileName SizedString
}

type MusicItem struct {
	Reserved uint32
	FileName FileName
	Volume   uint32
}

type ButtonIndex struct {
	Reserved uint32
	Index    int32
	X        int32
	Y        int32
}

type Point struct {
	X uint32
	Y uint32
}

// -------------------- Header --------------------

type Header struct {
	Magic        [6]byte
	DataVer      uint32
	GameWidth    uint32
	GameHeight   uint32
	MaxPicWidth  uint32
	MaxPicHeight uint32
	GUID         SizedString // 26 -> skip
	Title        SizedString // -> skip
	PVer         uint32
	Reserved1    uint32
	Reserved2    uint32 // + 12
}

// -------------------- Title --------------------

type Title struct {
	Reserved   uint32
	ShowLog    uint32
	LogoImage  FileName
	TitleImage FileName
	DrawTitle  uint32
	BGM        MusicItem
	ButtonCnt  uint32
	Buttons    []ButtonIndex `bin:"len:ButtonCnt"`
}

// -------------------- GameMenu --------------------

type GameMenu struct {
	Reserved  uint32
	BackImage FileName
	ButtonCnt uint32
	Buttons   []ButtonIndex `bin:"len:ButtonCnt"`
}

// -------------------- CG --------------------

type CGItem struct {
	Reserved uint32
	Name     SizedString
	CgPath   FileName
	Message  SizedString
}

type CG struct {
	Reserved    uint32
	BackImage   FileName
	Column      uint32
	SpanRow     uint32
	SpanCol     uint32
	ShowMessage uint32
	MegX        uint32
	MegY        uint32
	Zoom        uint32
	CgX         uint32
	CgY         uint32
	NoPic       FileName
	CgListCnt   uint32
	CgList      []CGItem `bin:"len:CgListCnt"`
	Viewport    Viewport
	BackButton  ButtonIndex
	CloseButton ButtonIndex
}

// -------------------- BGM --------------------

type BGMItem struct {
	Reserved uint32
	Name     SizedString
	BgmPath  FileName
	PicPath  FileName
	Message  SizedString
}

type BGM struct {
	Reserved     uint32
	BackImage    FileName
	Column       uint32
	SpanRow      uint32
	SpanCol      uint32
	ShowPic      uint32
	ShowMsg      uint32
	Px           uint32
	Py           uint32
	Mx           uint32
	My           uint32
	Nx           uint32
	Ny           uint32
	NoName       SizedString
	NoPic        FileName
	BgmListCnt   uint32
	BgmList      []BGMItem `bin:"len:BgmListCnt"`
	Viewport     Viewport
	SelectButton ButtonIndex
	CloseButton  ButtonIndex
}

// -------------------- SaveData (conditional) --------------------

type SaveData struct {
	Reserved1   uint32
	Reserved2   uint32
	Reserved3   uint32
	BackImage   FileName
	Reserved4   uint32
	Reserved5   uint32
	Reserved6   uint32
	Reserved7   uint32
	Reserved8   uint32
	Reserved9   uint32
	Reserved10  uint32
	Reserved11  uint32
	Reserved12  uint32
	Reserved13  uint32
	Reserved14  uint32
	Reserved15  uint32
	Viewport    Viewport
	BackButton  ButtonIndex
	CloseButton ButtonIndex

	// Conditional fields (if DataVer >= 106)
	TurnPageSwitch    uint32      // `bin:"if=DataVer >= 106"`
	PrevPage          ButtonIndex // `bin:"if=DataVer >= 106"`
	NextPage          ButtonIndex // `bin:"if=DataVer >= 106"`
	PageNumberTextPos Point       // `bin:"if=DataVer >= 106"`
	FontName          SizedString // `bin:"if=DataVer >= 106"`
	FontSize          uint32      // `bin:"if=DataVer >= 106"`
}

// -------------------- MessageBox --------------------

type TalkWin struct {
	Reserved        uint32
	BackX           uint32
	BackY           uint32
	BackImage       FileName
	FaceBorderImage FileName
	FaceBorderX     uint32
	FaceBorderY     uint32
	TextX           uint32
	TextWidth       uint32
	TextY           uint32
	Reserved2       uint32
	ButtonsCnt      uint32
	Buttons         []ButtonIndex `bin:"len:ButtonsCnt"`
}

type NameWin struct {
	Reserved  uint32
	BackX     uint32
	BackY     uint32
	BackImage FileName
	IsCenter  uint32
	TextX     uint32
	TextY     uint32
}

type MessageBox struct {
	Reserved          uint32
	FaceStyle         uint32
	ChoiceButtonIndex uint32
	TalkWin           TalkWin
	NameWin           NameWin
}

// -------------------- Replay --------------------

type Replay struct {
	Reserved    uint32
	BackImage   FileName
	CloseButton ButtonIndex
	Viewport    Viewport
}

// -------------------- Setting (conditional) --------------------

type Setting struct {
	Reserved  uint32
	BackImage FileName
	BarNone   FileName
	BarMove   FileName
	BgmX      uint32
	BgmY      uint32
	SeX       uint32
	SeY       uint32
	VoiceX    uint32
	VoiceY    uint32
	ShowFull  uint32
	ShowAuto  uint32
	ShowBGM   uint32
	ShowSE    uint32
	Reserved1 uint32
	Reserved2 uint32
	Button1   ButtonIndex
	Button2   ButtonIndex
	Button3   ButtonIndex
	Button4   ButtonIndex
	Button5   ButtonIndex
	Button6   ButtonIndex

	// Conditional fields (if DataVer >= 108)
	Reserved3 uint32 //`bin:"if=DataVer >= 108"`
	Reserved4 uint32 //`bin:"if=DataVer >= 108"`
	Reserved5 uint32 //`bin:"if=DataVer >= 108"`
}

// -------------------- Button --------------------

type Button struct {
	Reserved uint32
	Name     SizedString
	Image1   FileName
	Image2   FileName
	X        uint32
	Y        uint32
}

// -------------------- Event --------------------

type Event struct {
	Code      OpCode
	Indent    uint32
	ArgvCount uint32
	Argv      []SizedString `bin:"len:ArgvCount"`
}

// -------------------- CustomUI --------------------

type CustomUIData struct {
	Len  uint32
	Data []byte `bin:"len:Len"`
}

// (CUIInnerData is commented out in the original; we keep it as placeholder)
type CUIInnerData struct {
	LoadEventCount uint32
	// LoadEvent       []Event `bin:"len:LoadEventCount"`
	// AfterEventCount uint32
	// AfterEvent      []Event `bin:"len:AfterEventCount"`
}

// -------------------- System (conditional) --------------------

type System struct {
	Reserved      uint32
	FontName      SizedString
	FontSize      uint32
	FontTalkColor SizedString
	FontUiColor   SizedString

	// if DataVer >= 101
	EffectKind uint32 //`bin:"if=DataVer >= 101"`

	SkipTitle         uint32
	StartStoryId      uint32
	Reserved2         uint32
	IconName          FileName
	ShowAD            uint32
	AuthorWords       uint32
	AuthorWordsTiming uint32
	AutoRun           uint32
	ShowSystemMenu    uint32
	SEClick           MusicItem
	SEMove            MusicItem
	SECancel          MusicItem
	SEError           MusicItem
	Title             Title
	GameMenu          GameMenu
	CG                CG
	BGM               BGM
	SaveData          SaveData
	MsgBox            MessageBox
	Replay            Replay
	Setting           Setting
	ButtonsCnt        uint32
	Buttons           []Button `bin:"len:ButtonsCnt"`
	UIInitSave        uint32

	// if DataVer >= 103
	CuisCount uint32         //`bin:"if=DataVer >= 103"`
	Cuis      []CustomUIData `bin:"len:CuisCount"` //;if=DataVer >= 103"`
	MenuIndex uint32         //`bin:"if=DataVer >= 103"`
}

// -------------------- Character --------------------

type ChatCharacter struct {
	Reserved    uint32
	ID          uint32
	Name        SizedString
	FacePath    SizedString
	Reserved1   uint32
	Reserved2   uint32
	Reserved3   uint32
	BgPath      SizedString
	Reserved4   uint32
	CgPath      SizedString
	UiIndex     SizedString
	Reserved5   uint32
	Reserved6   uint32
	BgLeftW     uint32
	BgRightW    uint32
	Reserved7   uint32
	Reserved8   uint32
	Reserved9   uint32
	Reserved10  uint32
	Reserved11  uint32
	Reserved12  uint32
	Reserved13  uint32
	MsgFontPath SizedString
	MsgFontSize uint32
	Reserved14  uint32
	Reserved15  uint32
	Reserved16  uint32
	Reserved17  uint32
	A           SizedString
	B           SizedString
}

type Something struct {
	Str       SizedString
	Reserved  uint32
	Reserved2 uint32
}

type Character struct {
	ChatCharactersCnt uint32
	ChatCharacters    []ChatCharacter `bin:"len:ChatCharactersCnt"`
	SthCnt            uint32
	Sth               []Something `bin:"len:SthCnt"`
}

// -------------------- Event_ & Story --------------------

type Event_ struct {
	ByteLen uint32
	Event   Event
}

type Story struct {
	ByteLen   uint32
	Name      SizedString
	ID        uint32
	EventsCnt uint32
	Events    []Event_ `bin:"len:EventsCnt"`
}

func (s SizedString) String() string {
	return s.Data
}

// -------------------- DataFile (root) --------------------

type ScriptData struct {
	ByteLen     uint32
	ProjectName SizedString
	StoriesCnt  uint32
	Stories     []Story `bin:"len:StoriesCnt"`
}

type FloatButtonData struct {
	EventsCount uint32
	Events      []Event_ `bin:"len:EventsCount"`

	X             uint32
	Y             uint32
	Name          SizedString
	ElementsCount uint32
	// Elements      []FloatElement `bin:"len:ElementsCount"`
}

type DataFile struct {
	Header    Header
	System    System
	Character Character

	// We want a pointer to here.
	ScriptData ScriptData

	// FloatCount   uint32
	// FloatButtons []FloatButtonData `bin:"len:FloatCount"`
}

// 返回
