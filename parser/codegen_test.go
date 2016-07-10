package parser

import "testing"

func TestCodeGen(t *testing.T) {
	code := `
def main() {
	do_something(100, 5)
}

def do_something(x : int, y:int) {
	x * y
}
`
	token, _, err := Module(NewParser(code), 0)
	if err != nil {
		t.Error(err)
	}

	module, ok := token.(ModuleNode)
	if !ok {
		t.Error(module)
	}
}
