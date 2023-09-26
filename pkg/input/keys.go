package input

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"godev/pkg/config"

	raylib "github.com/gen2brain/raylib-go/raylib"
	json5 "github.com/zyedidia/json5"
)

var (
	Bindings map[string]map[string]string = make(map[string]map[string]string)

	prevControls map[string]map[string]bool = make(map[string]map[string]bool)
)

type Binding struct {
	Keys      map[string]bool
	Pressed   bool
	OnPress   *func()
	OnRelease *func()
}

// initialize bindings
func InitBindings() error {
	Bindings["Character"] = map[string]string{}
	for scope, _ := range Bindings {
		Bindings[scope] = DefaultBindings(scope)
	}
	raylib.SetExitKey(raylib.KeyF4)

	return LoadConfigBindings()
}

// load config bindings from file
func LoadConfigBindings() error {
	filename := filepath.Join(config.ConfigDir, "bindings.json")
	createBindingsIfNotExist(filename)

	var parsed map[string]map[string]interface{}
	if _, e := os.Stat(filename); e == nil {
		input, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		err = json5.Unmarshal(input, &parsed)
		if err != nil {
			return err
		}
	}

	for scope, binding := range parsed {
		for k, v := range binding {
			switch val := v.(type) {
			case string:
				bindKey(scope, k, val)
			default:
				return fmt.Errorf("Error reading bindings.json: non-string and non-map entry: %+v", k)
			}
		}
	}

	return nil
}

// attempt to create bindings file
func createBindingsIfNotExist(filename string) {
	if _, e := os.Stat(filename); os.IsNotExist(e) {
		ioutil.WriteFile(filename, []byte("{}"), 0644)
	}
}

// get currently pressed keys
func GetKeysPressed() map[string]bool {
	keys := map[string]bool{}

	//check keyboard
	for key, enum := range keyboard {
		if raylib.IsKeyDown(enum) {
			keys[key] = true
		}
	}

	//check mouse
	for key, enum := range mouse {
		if raylib.IsMouseButtonDown(enum) {
			keys[key] = true
		}
	}

	//check gamepads
	for name, controller := range controllers {
		if raylib.IsGamepadAvailable(controller) {
			for key, enum := range gamepad {
				if raylib.IsGamepadButtonDown(controller, enum) {
					keys[name+key] = true
				}
			}
		}
	}

	return keys
}

// get currently pressed modifiers
func GetModsPressed() map[string]bool {
	mods := map[string]bool{}

	//check keyboard
	for key, enum := range modifiers {
		if raylib.IsKeyDown(enum) {
			mods[key] = true
		}
	}

	_, l := mods["LeftShift"]
	_, r := mods["RightShift"]
	if l || r {
		mods["Shift"] = true
	}

	_, l = mods["LeftControl"]
	_, r = mods["RightControl"]
	if l || r {
		mods["Control"] = true
	}

	_, l = mods["LeftAlt"]
	_, r = mods["RightAlt"]
	if l || r {
		mods["Alt"] = true
	}

	_, l = mods["LeftSuper"]
	_, r = mods["RightSuper"]
	if l || r {
		mods["Super"] = true
	}

	return mods
}

// get pressed key combos
// TODO: make more efficent
func GetKeyCombos(key string, mods map[string]bool) map[string][]string {
	combos := map[string][]string{
		key: []string{key},
	}

	modlist := []string{}
	if _, ok := mods["LeftShift"]; ok {
		modlist = append(modlist, "LeftShift")
	}
	if _, ok := mods["RightShift"]; ok {
		modlist = append(modlist, "RightShift")
	}

	if _, ok := mods["LeftControl"]; ok {
		modlist = append(modlist, "LeftControl")
	}
	if _, ok := mods["RightControl"]; ok {
		modlist = append(modlist, "RightControl")
	}

	if _, ok := mods["LeftAlt"]; ok {
		modlist = append(modlist, "LeftAlt")
	}
	if _, ok := mods["RightAlt"]; ok {
		modlist = append(modlist, "RightAlt")
	}

	if _, ok := mods["LeftSuper"]; ok {
		modlist = append(modlist, "LeftSuper")
	}
	if _, ok := mods["RightSuper"]; ok {
		modlist = append(modlist, "RightSuper")
	}

	//combos specific to left/right
	for i, mod := range modlist {
		if key == mod {
			continue
		}
		combo := []string{key, mod}
		for _, mod := range modlist[i+1:] {
			if key == mod {
				continue
			}
			combo = append(combo, mod)
			combos[strings.Join(combo, " ")] = append(combo)
		}
	}

	if _, ok := mods["Shift"]; ok {
		modlist = append(modlist, "Shift")
	}
	if _, ok := mods["Control"]; ok {
		modlist = append(modlist, "Control")
	}
	if _, ok := mods["Alt"]; ok {
		modlist = append(modlist, "Alt")
	}
	if _, ok := mods["Super"]; ok {
		modlist = append(modlist, "Super")
	}

	//combos unspecific to left/right
	//do this last to prefer left/right over unsided
	for i, mod := range modlist {
		if strings.HasPrefix(mod, "Left") || strings.HasPrefix(mod, "Right") {
			break
		}
		combo := []string{key, mod}
		for _, mod := range modlist[i+1:] {
			if key == mod {
				continue
			}
			combo = append(combo, mod)
			combos[strings.Join(combo, " ")] = append(combo)
		}
	}
	return combos
}

// handle input keys
func HandleInput(scope string) {
	keys := GetKeysPressed()
	mods := GetModsPressed()
	nextControls := map[string]bool{}

	//handle each key pressed
	for key, _ := range keys {
		//search all key combos for possible controls
		combos := GetKeyCombos(key, mods)
		for combo, keys := range combos {
			if action, ok := Bindings[scope][combo]; ok {
				//handle action being pressed
				if err := handleBindingPress(action); err != nil {
					fmt.Printf("%+v\n", err)
				}
				//store that this action is pressed for tracking
				nextControls[action] = true
				delete(prevControls[scope], action)
				//remove keys from possible combo pool
				for _, key := range keys {
					delete(mods, key)
				}
				break
			}
		}
	}
	//search modifiers for actions as well
	//do this last to prefer key+mod rather than just mod
	for key, _ := range mods {
		//search all key combos for possible controls
		combos := GetKeyCombos(key, mods)
		for combo, keys := range combos {
			if action, ok := Bindings[scope][combo]; ok {
				//handle action being pressed
				if err := handleBindingPress(action); err != nil {
					fmt.Printf("%+v\n", err)
				}
				//store that this action is pressed for tracking
				nextControls[action] = true
				delete(prevControls[scope], action)
				//remove keys from possible combo pool
				for _, key := range keys {
					delete(mods, key)
				}
				break
			}
		}
	}

	//handle button releases
	for control, _ := range prevControls[scope] {
		if err := handleBindingRelease(control); err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
	//update previous pressed
	prevControls[scope] = nextControls

	return
}

// handle a binding action being pressed
func handleBindingPress(action string) error {
	if binding, ok := Actions[action]; ok {
		if !binding.Pressed {
			binding.Pressed = true
			if binding.OnPress != nil {
				(*binding.OnPress)()
			}
		}
	} else {
		return errors.New("Invalid action")
	}
	return nil
}

// handle a binding action being released
func handleBindingRelease(action string) error {
	if binding, ok := Actions[action]; ok {
		if binding.Pressed {
			binding.Pressed = false
			if binding.OnRelease != nil {
				(*binding.OnRelease)()
			}
		}
	} else {
		return errors.New("Invalid action")
	}
	return nil
}

// bind a key in a scope
func bindKey(scope string, k string, v string) {
	if _, ok := Bindings[scope]; !ok {
		return
	}
	Bindings[scope][k] = v
	Actions[v].Keys[k] = true
}

// unbind a key in a scope
func unbindKey(scope string, k string) {
	if _, ok := Bindings[scope][k]; ok {
		delete(Actions[Bindings[scope][k]].Keys, k)
		delete(Bindings[scope], k)
	}
}

// register an action
func RegisterAction(action string, press *func(), release *func(), overwrite bool) error {
	if _, ok := Actions[action]; ok && !overwrite {
		return errors.New("Already Existing Action")
	}
	Actions[action] = &Binding{
		Keys:      map[string]bool{},
		Pressed:   false,
		OnPress:   press,
		OnRelease: release,
	}

	return nil
}

// TryBindKey tries to bind a key by writing to config.ConfigDir/bindings.json
// Returns true if the keybinding already existed and a possible error
func TryBindKey(scope string, k string, v string, overwrite bool) (bool, error) {
	var e error
	var parsed map[string]map[string]string

	filename := filepath.Join(config.ConfigDir, "bindings.json")
	createBindingsIfNotExist(filename)

	if _, ok := Bindings[scope]; !ok {
		return false, errors.New("Scope does not exist for desired keybind")
	}

	if _, ok := Actions[v]; !ok {
		return false, errors.New("Action does not exist")
	}

	if _, e = os.Stat(filename); e == nil {
		input, err := ioutil.ReadFile(filename)
		if err != nil {
			return false, errors.New("Error reading bindings.json file: " + err.Error())
		}

		err = json5.Unmarshal(input, &parsed)
		if err != nil {
			return false, errors.New("Error reading bindings.json: " + err.Error())
		}

		found := false
		for key := range parsed[scope] {
			if key == k {
				if overwrite {
					parsed[scope][key] = v
				}
				found = true
				break
			}
		}

		if found && !overwrite {
			return true, nil
		} else if !found {
			parsed[scope][k] = v
		}
		if found && overwrite {
			unbindKey(scope, k)
		}

		bindKey(scope, k, v)

		txt, _ := json.MarshalIndent(parsed, "", "    ")
		return found, ioutil.WriteFile(filename, append(txt, '\n'), 0644)
	}
	return false, e
}

// UnbindKey removes the binding for a key from the bindings.json file
func UnbindKey(scope string, k string) error {
	var e error
	var parsed map[string]map[string]string

	if _, ok := Bindings[scope][k]; ok {
		return errors.New("Key does not exist")
	}
	unbindKey(scope, k)

	filename := filepath.Join(config.ConfigDir, "bindings.json")
	createBindingsIfNotExist(filename)
	if _, e = os.Stat(filename); e == nil {
		input, err := ioutil.ReadFile(filename)
		if err != nil {
			return errors.New("Error reading bindings.json file: " + err.Error())
		}

		err = json5.Unmarshal(input, &parsed)
		if err != nil {
			return errors.New("Error reading bindings.json: " + err.Error())
		}

		for key := range parsed[scope] {
			if key == k {
				delete(parsed[scope], key)
				break
			}
		}

		defaults := DefaultBindings(scope)
		if a, ok := defaults[k]; ok {
			bindKey(scope, k, a)
		}

		txt, _ := json.MarshalIndent(parsed, "", "    ")
		return ioutil.WriteFile(filename, append(txt, '\n'), 0644)
	}
	return e
}
