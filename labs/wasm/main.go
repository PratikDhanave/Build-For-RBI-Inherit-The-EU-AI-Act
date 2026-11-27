//go:build js && wasm

package main

import "syscall/js"

func main() {
	js.Global().Set("goRunLab", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) < 2 {
			return map[string]any{"ok": false, "err": "expected (labName, goSource)"}
		}
		lab := args[0].String()
		src := args[1].String()

		var (
			out string
			err error
		)
		switch lab {
		case "drift":
			out, err = EvalApproval(src)
		case "chain":
			out, err = EvalAudit(src)
		default:
			return map[string]any{"ok": false, "err": "unknown lab: " + lab}
		}
		if err != nil {
			return map[string]any{"ok": false, "err": err.Error()}
		}
		return map[string]any{"ok": true, "out": out}
	}))
	if cb := js.Global().Get("__onGoWasmReady"); cb.Type() == js.TypeFunction {
		cb.Invoke()
	}
	select {}
}
