package vm

/* --- Generic code builder, hopefully extensible --- */

type NameType struct {
	Name string
	Type string
}

type ScopedBuilder func(scopedBuilder CodeBuilder)

type CodeBuilder interface {
	Push(instructions ...Instruction)
	DefineFunc(signature []NameType, builder ScopedBuilder)
	Resolve(symbol string) int64
	ResolveOrDefine(symbol string) int64
}

/* --- Default code builder --- */

const (
	SYM_FUN int64 = iota
	SYM_VAR
)

type SymbolType int64

type Symbol struct {
	ID   int64
	Type SymbolType
}

type Scope struct {
	// doesn't seem like anything else goes inside a scope, does there?
	SymbolTable map[string]Symbol
}

func NewScope() Scope {
	return Scope{make(map[string]Symbol)}
}

type ScopeStack struct {
	Stack []Scope
}

// PushScope creates a new scope
func (stack ScopeStack) PushScope() {
	stack.Stack = append(stack.Stack, NewScope())
}

type GimmickBuilder struct {
	Scopes ScopeStack
}

// Interface methods

func (builder GimmickBuilder) Push(instructions ...Instruction) {
}

func (builder GimmickBuilder) DefineFunc(signature []NameType, scopedBuilder ScopedBuilder) {
}

func (builder GimmickBuilder) Resolve(symbol string) int64 {
	return -1
}

func (builder GimmickBuilder) ResolveOrDefine(symbol string) int64 {
	return -1
}

// private methods
