package llvmgen

import (
	"compiler_project/tac"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"strconv"
)

type LLVMBuilder struct {
	mod    *ir.Module
	fnMain *ir.Func
	block  *ir.Block
	vars   map[string]*ir.InstAlloca
}

func NewLLVMBuilder() *LLVMBuilder {
	mod := ir.NewModule()
	mainFn := mod.NewFunc("main", types.I32)
	entry := mainFn.NewBlock("entry")

	return &LLVMBuilder{
		mod:    mod,
		fnMain: mainFn,
		block:  entry,
		vars:   map[string]*ir.InstAlloca{},
	}
}

func (b *LLVMBuilder) GenerateFromTAC(instructions []tac.TACInstruction) {
	for _, instr := range instructions {
		switch instr.Op {
		case "=":
			val := b.getValue(instr.Arg1)
			ptr := b.ensureVar(instr.Res)
			b.block.NewStore(val, ptr)

		case "+", "-", "*", "/":
			l := b.getValue(instr.Arg1)
			r := b.getValue(instr.Arg2)
			var result value.Value
			switch instr.Op {
			case "+":
				result = b.block.NewAdd(l, r)
			case "-":
				result = b.block.NewSub(l, r)
			case "*":
				result = b.block.NewMul(l, r)
			case "/":
				result = b.block.NewSDiv(l, r)
			}
			ptr := b.ensureVar(instr.Res)
			b.block.NewStore(result, ptr)
		}
	}
	// Возврат из main
	b.block.NewRet(constant.NewInt(types.I32, 0))
}

func (b *LLVMBuilder) ensureVar(name string) *ir.InstAlloca {
	if v, ok := b.vars[name]; ok {
		return v
	}
	ptr := b.block.NewAlloca(types.I32)
	b.vars[name] = ptr
	return ptr
}

func (b *LLVMBuilder) getValue(name string) value.Value {
	if v, ok := b.vars[name]; ok {
		return b.block.NewLoad(types.I32, v)
	}
	// Попробуем парсить как число
	if i, err := strconv.Atoi(name); err == nil {
		return constant.NewInt(types.I32, int64(i))
	}
	panic("неизвестное значение: " + name)
}

func (b *LLVMBuilder) IR() *ir.Module {
	return b.mod
}
