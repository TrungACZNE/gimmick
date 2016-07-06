package parser

/* --- VM bytecode generation routines ---*/

func (node EOFNode) CodeGen(stream *CodeStream) {
}
func (node IntegerLiteralNode) CodeGen(stream *CodeStream) {
}
func (node FloatLiteralNode) CodeGen(stream *CodeStream) {
}
func (node KeywordNode) CodeGen(stream *CodeStream) {
}
func (node SymbolNode) CodeGen(stream *CodeStream) {
}
func (node IdentifierNode) CodeGen(stream *CodeStream) {
}
func (node ArgDeclNode) CodeGen(stream *CodeStream) {
}
func (node ArgListNode) CodeGen(stream *CodeStream) {
}
func (node ParamListNode) CodeGen(stream *CodeStream) {
}
func (node EmptyNode) CodeGen(stream *CodeStream) {
}
func (node FunctionDefNode) CodeGen(stream *CodeStream) {
}
func (node FunctionCallNode) CodeGen(stream *CodeStream) {
}
func (node BinaryOperatorNode) CodeGen(stream *CodeStream) {
}
func (node BlockNode) CodeGen(stream *CodeStream) {
}
