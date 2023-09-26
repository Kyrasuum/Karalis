package input

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

var (
	controllers map[string]int32 = map[string]int32{
		// Gamepad Number
		"Controller1": raylib.GamepadPlayer1,
		"Controller2": raylib.GamepadPlayer2,
		"Controller3": raylib.GamepadPlayer3,
		"Controller4": raylib.GamepadPlayer4,
	}

	modifiers map[string]int32 = map[string]int32{
		// Keyboard Function Keys
		"LeftShift":    raylib.KeyLeftShift,
		"LeftControl":  raylib.KeyLeftControl,
		"LeftAlt":      raylib.KeyLeftAlt,
		"LeftSuper":    raylib.KeyLeftSuper,
		"RightShift":   raylib.KeyRightShift,
		"RightControl": raylib.KeyRightControl,
		"RightAlt":     raylib.KeyRightAlt,
		"RightSuper":   raylib.KeyRightSuper,
	}

	keyboard map[string]int32 = map[string]int32{
		// Keyboard Function Keys
		"Space":        raylib.KeySpace,
		"Escape":       raylib.KeyEscape,
		"Enter":        raylib.KeyEnter,
		"Tab":          raylib.KeyTab,
		"Backspace":    raylib.KeyBackspace,
		"Insert":       raylib.KeyInsert,
		"Delete":       raylib.KeyDelete,
		"Right":        raylib.KeyRight,
		"Left":         raylib.KeyLeft,
		"Down":         raylib.KeyDown,
		"Up":           raylib.KeyUp,
		"PageUp":       raylib.KeyPageUp,
		"PageDown":     raylib.KeyPageDown,
		"Home":         raylib.KeyHome,
		"End":          raylib.KeyEnd,
		"CapsLock":     raylib.KeyCapsLock,
		"ScrollLock":   raylib.KeyScrollLock,
		"NumLock":      raylib.KeyNumLock,
		"PrintScreen":  raylib.KeyPrintScreen,
		"Pause":        raylib.KeyPause,
		"F1":           raylib.KeyF1,
		"F2":           raylib.KeyF2,
		"F3":           raylib.KeyF3,
		"F4":           raylib.KeyF4,
		"F5":           raylib.KeyF5,
		"F6":           raylib.KeyF6,
		"F7":           raylib.KeyF7,
		"F8":           raylib.KeyF8,
		"F9":           raylib.KeyF9,
		"F10":          raylib.KeyF10,
		"F11":          raylib.KeyF11,
		"F12":          raylib.KeyF12,
		"KbMenu":       raylib.KeyKbMenu,
		"LeftBracket":  raylib.KeyLeftBracket,
		"BackSlash":    raylib.KeyBackSlash,
		"RightBracket": raylib.KeyRightBracket,
		"Grave":        raylib.KeyGrave,

		// Keyboard Number Pad Keys
		"Kp0":        raylib.KeyKp0,
		"Kp1":        raylib.KeyKp1,
		"Kp2":        raylib.KeyKp2,
		"Kp3":        raylib.KeyKp3,
		"Kp4":        raylib.KeyKp4,
		"Kp5":        raylib.KeyKp5,
		"Kp6":        raylib.KeyKp6,
		"Kp7":        raylib.KeyKp7,
		"Kp8":        raylib.KeyKp8,
		"Kp9":        raylib.KeyKp9,
		"KpDecimal":  raylib.KeyKpDecimal,
		"KpDivide":   raylib.KeyKpDivide,
		"KpMultiply": raylib.KeyKpMultiply,
		"KpSubtract": raylib.KeyKpSubtract,
		"KpAdd":      raylib.KeyKpAdd,
		"KpEnter":    raylib.KeyKpEnter,
		"KpEqual":    raylib.KeyKpEqual,

		// Keyboard Alpha Numeric Keys
		"Apostrophe": raylib.KeyApostrophe,
		"Comma":      raylib.KeyComma,
		"Minus":      raylib.KeyMinus,
		"Period":     raylib.KeyPeriod,
		"Slash":      raylib.KeySlash,
		"Zero":       raylib.KeyZero,
		"One":        raylib.KeyOne,
		"Two":        raylib.KeyTwo,
		"Three":      raylib.KeyThree,
		"Four":       raylib.KeyFour,
		"Five":       raylib.KeyFive,
		"Six":        raylib.KeySix,
		"Seven":      raylib.KeySeven,
		"Eight":      raylib.KeyEight,
		"Nine":       raylib.KeyNine,
		"Semicolon":  raylib.KeySemicolon,
		"Equal":      raylib.KeyEqual,
		"A":          raylib.KeyA,
		"B":          raylib.KeyB,
		"C":          raylib.KeyC,
		"D":          raylib.KeyD,
		"E":          raylib.KeyE,
		"F":          raylib.KeyF,
		"G":          raylib.KeyG,
		"H":          raylib.KeyH,
		"I":          raylib.KeyI,
		"J":          raylib.KeyJ,
		"K":          raylib.KeyK,
		"L":          raylib.KeyL,
		"M":          raylib.KeyM,
		"N":          raylib.KeyN,
		"O":          raylib.KeyO,
		"P":          raylib.KeyP,
		"Q":          raylib.KeyQ,
		"R":          raylib.KeyR,
		"S":          raylib.KeyS,
		"T":          raylib.KeyT,
		"U":          raylib.KeyU,
		"V":          raylib.KeyV,
		"W":          raylib.KeyW,
		"X":          raylib.KeyX,
		"Y":          raylib.KeyY,
		"Z":          raylib.KeyZ,

		// Android keys
		"Back":       raylib.KeyBack,
		"Menu":       raylib.KeyMenu,
		"VolumeUp":   raylib.KeyVolumeUp,
		"VolumeDown": raylib.KeyVolumeDown,
	}

	mouse map[string]int32 = map[string]int32{
		// Mouse Buttons
		"LeftButton":    raylib.MouseLeftButton,
		"RightButton":   raylib.MouseRightButton,
		"MiddleButton":  raylib.MouseMiddleButton,
		"SideButton":    raylib.MouseSideButton,
		"ExtraButton":   raylib.MouseExtraButton,
		"ForwardButton": raylib.MouseForwardButton,
		"BackButton":    raylib.MouseBackButton,
	}

	gamepad map[string]int32 = map[string]int32{
		// PS3 USB Controller Buttons
		"Ps3Triangle": raylib.GamepadPs3ButtonTriangle,
		"Ps3Circle":   raylib.GamepadPs3ButtonCircle,
		"Ps3Cross":    raylib.GamepadPs3ButtonCross,
		"Ps3Square":   raylib.GamepadPs3ButtonSquare,
		"Ps3L1":       raylib.GamepadPs3ButtonL1,
		"Ps3R1":       raylib.GamepadPs3ButtonR1,
		"Ps3L2":       raylib.GamepadPs3ButtonL2,
		"Ps3R2":       raylib.GamepadPs3ButtonR2,
		"Ps3Start":    raylib.GamepadPs3ButtonStart,
		"Ps3Select":   raylib.GamepadPs3ButtonSelect,
		"Ps3Up":       raylib.GamepadPs3ButtonUp,
		"Ps3Right":    raylib.GamepadPs3ButtonRight,
		"Ps3Down":     raylib.GamepadPs3ButtonDown,
		"Ps3Left":     raylib.GamepadPs3ButtonLeft,
		"Ps3Ps":       raylib.GamepadPs3ButtonPs,

		// PS3 USB Controller Axis
		"Ps3AxisLeftX":  raylib.GamepadPs3AxisLeftX,
		"Ps3AxisLeftY":  raylib.GamepadPs3AxisLeftY,
		"Ps3AxisRightX": raylib.GamepadPs3AxisRightX,
		"Ps3AxisRightY": raylib.GamepadPs3AxisRightY,
		// [1..-1] (pressure-level)
		"Ps3AxisL2": raylib.GamepadPs3AxisL2,
		// [1..-1] (pressure-level)
		"Ps3AxisR2": raylib.GamepadPs3AxisR2,

		// Xbox360 USB Controller Buttons
		"XboxA":      raylib.GamepadXboxButtonA,
		"XboxB":      raylib.GamepadXboxButtonB,
		"XboxX":      raylib.GamepadXboxButtonX,
		"XboxY":      raylib.GamepadXboxButtonY,
		"XboxLb":     raylib.GamepadXboxButtonLb,
		"XboxRb":     raylib.GamepadXboxButtonRb,
		"XboxSelect": raylib.GamepadXboxButtonSelect,
		"XboxStart":  raylib.GamepadXboxButtonStart,
		"XboxUp":     raylib.GamepadXboxButtonUp,
		"XboxRight":  raylib.GamepadXboxButtonRight,
		"XboxDown":   raylib.GamepadXboxButtonDown,
		"XboxLeft":   raylib.GamepadXboxButtonLeft,
		"XboxHome":   raylib.GamepadXboxButtonHome,

		// Xbox360 USB Controller Axis
		// [-1..1] (left->right)
		"XboxAxisLeftX": raylib.GamepadXboxAxisLeftX,
		// [1..-1] (up->down)
		"XboxAxisLeftY": raylib.GamepadXboxAxisLeftY,
		// [-1..1] (left->right)
		"XboxAxisRightX": raylib.GamepadXboxAxisRightX,
		// [1..-1] (up->down)
		"XboxAxisRightY": raylib.GamepadXboxAxisRightY,
		// [-1..1] (pressure-level)
		"XboxAxisLt": raylib.GamepadXboxAxisLt,
		// [-1..1] (pressure-level)
		"XboxAxisRt": raylib.GamepadXboxAxisRt,

		// Android Gamepad Controller (SNES CLASSIC)
		"AndroidDpadUp":     raylib.GamepadAndroidDpadUp,
		"AndroidDpadDown":   raylib.GamepadAndroidDpadDown,
		"AndroidDpadLeft":   raylib.GamepadAndroidDpadLeft,
		"AndroidDpadRight":  raylib.GamepadAndroidDpadRight,
		"AndroidDpadCenter": raylib.GamepadAndroidDpadCenter,

		"AndroidButtonA":  raylib.GamepadAndroidButtonA,
		"AndroidButtonB":  raylib.GamepadAndroidButtonB,
		"AndroidButtonC":  raylib.GamepadAndroidButtonC,
		"AndroidButtonX":  raylib.GamepadAndroidButtonX,
		"AndroidButtonY":  raylib.GamepadAndroidButtonY,
		"AndroidButtonZ":  raylib.GamepadAndroidButtonZ,
		"AndroidButtonL1": raylib.GamepadAndroidButtonL1,
		"AndroidButtonR1": raylib.GamepadAndroidButtonR1,
		"AndroidButtonL2": raylib.GamepadAndroidButtonL2,
		"AndroidButtonR2": raylib.GamepadAndroidButtonR2,
	}
)
