package llvmgen

import (
	"compiler_project/tac"
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"strconv"
)

type LLVMBuilder struct {
	mod    *ir.Module
	fnMain *ir.Func
	block  *ir.Block
	vars   map[string]*ir.InstAlloca
	printf *ir.Func
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
	labelBlocks := make(map[string]*ir.Block)

	for idx := 0; idx < len(instructions); idx++ {
		instr := instructions[idx]

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

		case "show":
			val := b.getValue(instr.Arg1)
			printf := b.ensurePrintf()
			formatStr := b.ensureGlobalString("%d\n", "fmt")
			b.block.NewCall(printf, formatStr, val)

		case "showstr":
			formatStr := b.ensureGlobalString("%s\n", "fmt_str")
			str := b.ensureGlobalString(instr.Arg1, b.uniqueGlobalName("str"))
			printf := b.ensurePrintf()
			b.block.NewCall(printf, formatStr, str)

		case "iffalse":
			cond := b.getValue(instr.Arg1)
			zero := constant.NewInt(types.I32, 0)
			condVal := b.block.NewICmp(enum.IPredNE, cond, zero)

			trueBlock := b.fnMain.NewBlock(b.uniqueLabel("true"))
			falseBlock := labelBlocks[instr.Res]
			if falseBlock == nil {
				falseBlock = b.fnMain.NewBlock(instr.Res)
				labelBlocks[instr.Res] = falseBlock
			}
			b.block.NewCondBr(condVal, trueBlock, falseBlock)
			b.block = trueBlock

		case "goto":
			targetBlock := labelBlocks[instr.Res]
			if targetBlock == nil {
				targetBlock = b.fnMain.NewBlock(instr.Res)
				labelBlocks[instr.Res] = targetBlock
			}
			b.block.NewBr(targetBlock)

		case "label":
			block := labelBlocks[instr.Res]
			if block == nil {
				block = b.fnMain.NewBlock(instr.Res)
				labelBlocks[instr.Res] = block
			}
			b.block = block

		case "call":
			fn := b.mod.NewFunc(instr.Arg1, types.I32)
			argsCount, _ := strconv.Atoi(instr.Arg2)
			args := make([]value.Value, argsCount)
			for i := 0; i < argsCount; i++ {
				idx++
				argInstr := instructions[idx]
				args[i] = b.getValue(argInstr.Arg1)
			}
			result := b.block.NewCall(fn, args...)
			ptr := b.ensureVar(instr.Res)
			b.block.NewStore(result, ptr)
		}
	}

	pauseStr := b.ensurePauseString()
	pausePtr := b.block.NewGetElementPtr(
		pauseStr.ContentType,
		pauseStr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	systemFn := b.ensureSystemFunc()
	b.block.NewCall(systemFn, pausePtr)

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

func (b *LLVMBuilder) ensureGlobalString(s, name string) value.Value {
	g := b.mod.NewGlobalDef(name, constant.NewCharArrayFromString(s+"\x00"))
	g.Linkage = enum.LinkagePrivate
	return b.block.NewGetElementPtr(g.ContentType, g, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
}

func (b *LLVMBuilder) ensurePrintf() *ir.Func {
	if b.printf != nil {
		return b.printf
	}
	b.printf = b.mod.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	b.printf.Sig.Variadic = true
	return b.printf
}

func (b *LLVMBuilder) ensurePauseString() *ir.Global {
	for _, g := range b.mod.Globals {
		if g.Name() == ".pause_str" {
			return g
		}
	}

	strVal := constant.NewCharArrayFromString("pause\x00")
	global := b.mod.NewGlobalDef(".pause_str", strVal)
	global.Immutable = true
	// b.mod.Globals уже содержит global, не нужно явно добавлять
	return global
}

func (b *LLVMBuilder) ensureSystemFunc() *ir.Func {
	for _, fn := range b.mod.Funcs {
		if fn.Name() == "system" {
			return fn
		}
	}
	// Объявляем extern int system(i8*)

	//systemType := types.NewFunc(types.I32, types.NewPointer(types.I8))
	fn := b.mod.NewFunc("system", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	fn.Linkage = enum.LinkageExternal
	return fn
}

func (b *LLVMBuilder) uniqueGlobalName(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, len(b.mod.Globals))
}

func (b *LLVMBuilder) uniqueLabel(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, len(b.fnMain.Blocks))
}
