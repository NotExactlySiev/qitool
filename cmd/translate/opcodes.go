package main

// OpCode represents the set of available instruction opcodes.
type OpCode uint32

// Enum values for all opcodes.
const (
	Text            OpCode = 100
	TextDif         OpCode = 101
	TextDifEX       OpCode = 1010
	TextDifEX2      OpCode = 1011
	TextEnd         OpCode = 102
	AutoPlay        OpCode = 103
	QuickPlay       OpCode = 104
	Notes           OpCode = 107
	TextChoice      OpCode = 108
	DisposeText     OpCode = 109
	DFloatButton    OpCode = 112
	UpdateUI        OpCode = 150
	BackGame        OpCode = 151
	If              OpCode = 200
	IfEnd           OpCode = 201
	Loop            OpCode = 202
	LoopAboveStart  OpCode = 203
	ButtonDif       OpCode = 204
	ButtonDifEnd    OpCode = 205
	JumpStory       OpCode = 206
	SkipStory       OpCode = 2062
	SCUIJumpStory   OpCode = 2063
	Var             OpCode = 207
	BackTitle       OpCode = 208
	LoopBreak       OpCode = 209
	Wait            OpCode = 210
	IfChoice        OpCode = 211
	ButtonDifChoose OpCode = 212
	VarEx           OpCode = 213
	CallMenu        OpCode = 214
	String          OpCode = 215
	AdvData         OpCode = 216
	IfEx            OpCode = 217
	MustSaveRead    OpCode = 218
	CallAD          OpCode = 230
	AdContent       OpCode = 231
	AdEnd           OpCode = 232
	CallSubStory    OpCode = 251
	CWeather        OpCode = 301
	Shake           OpCode = 302
	Flash           OpCode = 303
	BGMAdd          OpCode = 307
	CGAdd           OpCode = 308
	SimpleEffect    OpCode = 310
	ShowPic         OpCode = 400
	DisposePic      OpCode = 401
	MovePic         OpCode = 402
	RotatePic       OpCode = 404
	PreLoadPic      OpCode = 405
	ShowDynamicPic  OpCode = 406
	ShowCharacter   OpCode = 408
	StartBGM        OpCode = 501
	StartSE         OpCode = 502
	StartVoice      OpCode = 503
	StartBGS        OpCode = 504
	FadeBGM         OpCode = 505
	StopSE          OpCode = 506
	StopVoice       OpCode = 507
	FadeBGS         OpCode = 508
	StartVideo      OpCode = 600
	OperationVideo  OpCode = 601
	BubbleMode      OpCode = 220
	SendBubble      OpCode = 221
	HpLock          OpCode = 252
	StoryLock       OpCode = 253
	ShowDynamicPic2 OpCode = 4062 // duplicate class but different opcode
	ShowVideo       OpCode = 700
)

// String returns the name of the opcode for debugging or logging.
func (o OpCode) String() string {
	switch o {
	case Text:
		return "Text"
	case TextDif:
		return "TextDif"
	case TextDifEX:
		return "TextDifEX"
	case TextDifEX2:
		return "TextDifEX2"
	case TextEnd:
		return "TextEnd"
	case AutoPlay:
		return "AutoPlay"
	case QuickPlay:
		return "QuickPlay"
	case Notes:
		return "Notes"
	case TextChoice:
		return "TextChoice"
	case DisposeText:
		return "DisposeText"
	case DFloatButton:
		return "DFloatButton"
	case UpdateUI:
		return "UpdateUI"
	case BackGame:
		return "BackGame"
	case If:
		return "If"
	case IfEnd:
		return "IfEnd"
	case Loop:
		return "Loop"
	case LoopAboveStart:
		return "LoopAboveStart"
	case ButtonDif:
		return "ButtonDif"
	case ButtonDifEnd:
		return "ButtonDifEnd"
	case JumpStory:
		return "JumpStory"
	case SkipStory:
		return "SkipStory"
	case SCUIJumpStory:
		return "SCUIJumpStory"
	case Var:
		return "Var"
	case BackTitle:
		return "BackTitle"
	case LoopBreak:
		return "LoopBreak"
	case Wait:
		return "Wait"
	case IfChoice:
		return "IfChoice"
	case ButtonDifChoose:
		return "ButtonDifChoose"
	case VarEx:
		return "VarEx"
	case CallMenu:
		return "CallMenu"
	case String:
		return "String"
	case AdvData:
		return "AdvData"
	case IfEx:
		return "IfEx"
	case MustSaveRead:
		return "MustSaveRead"
	case CallAD:
		return "CallAD"
	case AdContent:
		return "AdContent"
	case AdEnd:
		return "AdEnd"
	case CallSubStory:
		return "CallSubStory"
	case CWeather:
		return "CWeather"
	case Shake:
		return "Shake"
	case Flash:
		return "Flash"
	case BGMAdd:
		return "BGMAdd"
	case CGAdd:
		return "CGAdd"
	case SimpleEffect:
		return "SimpleEffect"
	case ShowPic:
		return "ShowPic"
	case DisposePic:
		return "DisposePic"
	case MovePic:
		return "MovePic"
	case RotatePic:
		return "RotatePic"
	case PreLoadPic:
		return "PreLoadPic"
	case ShowDynamicPic:
		return "ShowDynamicPic"
	case ShowCharacter:
		return "ShowCharacter"
	case StartBGM:
		return "StartBGM"
	case StartSE:
		return "StartSE"
	case StartVoice:
		return "StartVoice"
	case StartBGS:
		return "StartBGS"
	case FadeBGM:
		return "FadeBGM"
	case StopSE:
		return "StopSE"
	case StopVoice:
		return "StopVoice"
	case FadeBGS:
		return "FadeBGS"
	case StartVideo:
		return "StartVideo"
	case OperationVideo:
		return "OperationVideo"
	case BubbleMode:
		return "BubbleMode"
	case SendBubble:
		return "SendBubble"
	case HpLock:
		return "HpLock"
	case StoryLock:
		return "StoryLock"
	case ShowDynamicPic2:
		return "ShowDynamicPic2"
	case ShowVideo:
		return "ShowVideo"
	default:
		return "Unknown"
	}
}
