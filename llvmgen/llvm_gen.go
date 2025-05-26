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
	mod           *ir.Module
	fnMain        *ir.Func
	block         *ir.Block
	vars          map[string]*ir.InstAlloca
	printf        *ir.Func
	namedValues   map[string]value.Value
	globalStrings map[string]*ir.Global
	varsTypes     map[string]types.Type
}

func NewLLVMBuilder() *LLVMBuilder {
	mod := ir.NewModule()
	mainFn := mod.NewFunc("main", types.I32)
	entry := mainFn.NewBlock("entry")

	return &LLVMBuilder{
		mod:           mod,
		fnMain:        mainFn,
		block:         entry,
		vars:          map[string]*ir.InstAlloca{},
		namedValues:   make(map[string]value.Value),
		globalStrings: make(map[string]*ir.Global),
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

			// Если типы не совпадают — возможно, нужно привести val
			// или просто сделать store с правильным типом
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

		case "less", "more", "equal":
			l := b.getValue(instr.Arg1)
			r := b.getValue(instr.Arg2)
			var cmp value.Value
			switch instr.Op {
			case "less":
				cmp = b.block.NewICmp(enum.IPredSLT, l, r)
			case "more":
				cmp = b.block.NewICmp(enum.IPredSGT, l, r)
			case "equal":
				cmp = b.block.NewICmp(enum.IPredEQ, l, r)
			}
			b.namedValues[instr.Res] = cmp

		case "iffalse":
			cond := b.getValue(instr.Arg1) // cond — это i1 (например, результат icmp)

			// Блок, если условие ложно (перейти на Res)
			falseBlock := labelBlocks[instr.Res]
			if falseBlock == nil {
				falseBlock = b.fnMain.NewBlock(instr.Res)
				labelBlocks[instr.Res] = falseBlock
			}

			// Продолжение, если условие истинно (continue блок)
			trueBlock := b.fnMain.NewBlock(b.uniqueLabel("continue"))
			b.block.NewCondBr(cond, trueBlock, falseBlock)

			// Продолжаем генерацию в trueBlock
			b.block = trueBlock

		case "goto":
			targetBlock := labelBlocks[instr.Res]
			if targetBlock == nil {
				targetBlock = b.fnMain.NewBlock(instr.Res)
				labelBlocks[instr.Res] = targetBlock
			}
			// Завершаем текущий блок переходом, если нет терминатора
			lastInsts := b.block.Insts
			if len(lastInsts) == 0 || (!isTerminator(lastInsts[len(lastInsts)-1])) {
				b.block.NewBr(targetBlock)
			}
			// Переключаемся на целевой блок
			b.block = targetBlock

		case "label":
			block := labelBlocks[instr.Res]
			if block == nil {
				block = b.fnMain.NewBlock(instr.Res)
				labelBlocks[instr.Res] = block
			}
			// Если текущий блок не закончен, завершаем его переходом на этот
			lastInsts := b.block.Insts
			if len(lastInsts) > 0 && !isTerminator(lastInsts[len(lastInsts)-1]) {
				b.block.NewBr(block)
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

		case "while_start":
			// Метка начала цикла
			startBlock := b.fnMain.NewBlock(instr.Res + "_start")
			endBlock := b.fnMain.NewBlock(instr.Res + "_end")
			labelBlocks[instr.Res+"_start"] = startBlock
			labelBlocks[instr.Res+"_end"] = endBlock

			// Переход из текущего блока в начало цикла
			if len(b.block.Insts) == 0 || !isTerminator(b.block.Insts[len(b.block.Insts)-1]) {
				b.block.NewBr(startBlock)
			}
			b.block = startBlock

		case "while_cond":
			// Проверяем условие цикла, ожидается, что Arg1 - условие, Res - имя while_start блока
			cond := b.getValue(instr.Arg1)
			zero := constant.NewInt(types.I32, 0)
			condVal := b.block.NewICmp(enum.IPredNE, cond, zero)

			bodyBlock := b.fnMain.NewBlock(instr.Res + "_body")
			endBlock := labelBlocks[instr.Res+"_end"]
			if endBlock == nil {
				endBlock = b.fnMain.NewBlock(instr.Res + "_end")
				labelBlocks[instr.Res+"_end"] = endBlock
			}
			labelBlocks[instr.Res+"_body"] = bodyBlock

			b.block.NewCondBr(condVal, bodyBlock, endBlock)

			b.block = bodyBlock

		case "while_end":
			// После тела цикла — переход к началу для проверки условия
			startBlock := labelBlocks[instr.Res+"_start"]
			endBlock := labelBlocks[instr.Res+"_end"]
			if startBlock == nil || endBlock == nil {
				panic("не найдены блоки while для " + instr.Res)
			}

			if len(b.block.Insts) == 0 || !isTerminator(b.block.Insts[len(b.block.Insts)-1]) {
				b.block.NewBr(startBlock)
			}
			b.block = endBlock
		}
	}

	// В конце ставим вызов system("pause") и ret
	pauseStr := b.ensurePauseString()
	pausePtr := b.block.NewGetElementPtr(
		pauseStr.ContentType,
		pauseStr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	systemFn := b.ensureSystemFunc()
	b.block.NewCall(systemFn, pausePtr)

	// Завершаем последний блок
	if len(b.block.Insts) == 0 || !isTerminator(b.block.Insts[len(b.block.Insts)-1]) {
		b.block.NewRet(constant.NewInt(types.I32, 0))
	}
}

// Вспомогательная функция для проверки терминатора
func isTerminator(inst ir.Instruction) bool {
	_, ok := inst.(ir.Terminator)
	return ok
}

func (b *LLVMBuilder) ensureVar(name string) value.Value {
	if ptr, ok := b.vars[name]; ok {
		return ptr
	}

	var varType types.Type = types.I32 // по умолчанию int

	if t, ok := b.varsTypes[name]; ok {
		varType = t
	}

	ptr := b.block.NewAlloca(varType)
	b.vars[name] = ptr
	return ptr
}

func (b *LLVMBuilder) getValue(name string) value.Value {
	// Временные значения
	if val, ok := b.namedValues[name]; ok {
		return val
	}

	// Булевы литералы
	if name == "true" {
		return constant.NewInt(types.I1, 1)
	}
	if name == "false" {
		return constant.NewInt(types.I1, 0)
	}

	// Целые числа
	if intVal, err := strconv.Atoi(name); err == nil {
		return constant.NewInt(types.I32, int64(intVal))
	}

	// Переменные
	if ptr, ok := b.vars[name]; ok {
		elemType := ptr.Type().(*types.PointerType).ElemType

		switch elemType {
		case types.I1:
			return b.block.NewLoad(types.I1, ptr)
		case types.I32:
			return b.block.NewLoad(types.I32, ptr)
		case types.I8Ptr: // i8* — указатель на строку
			// Для строки просто загрузим указатель i8*
			return b.block.NewLoad(types.I8Ptr, ptr)
		default:
			panic("неподдерживаемый тип переменной: " + name)
		}
	}

	panic("неизвестное значение: " + name)
}

func (b *LLVMBuilder) IR() *ir.Module {
	return b.mod
}

func (b *LLVMBuilder) ensureGlobalString(str, name string) *ir.Global {
	if g, ok := b.globalStrings[name]; ok {
		return g
	}
	global := b.mod.NewGlobalDef(name, constant.NewCharArrayFromString(str))
	b.globalStrings[name] = global
	return global
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
