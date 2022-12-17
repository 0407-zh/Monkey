package vm

import (
	"Monkey/code"
	"Monkey/compiler"
	"Monkey/object"
	"fmt"
)

const StackSize = 2048 // 虚拟机栈大小

type VM struct {
	constants    []object.Object   // 编译器生成的常量
	instructions code.Instructions // 编译器生成的指令

	stack []object.Object // 虚拟机栈
	sp    int             // 栈指针， 始终指向栈中的下一位，即当前索引加一
}

// New 生成虚拟机实例
func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

// Run 启动虚拟机
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		// 取指令
		op := code.Opcode(vm.instructions[ip]) // 转换为操作码

		// 解码
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			// 执行指令
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value

			result := leftValue + rightValue
			vm.push(&object.Integer{Value: result})
		}
	}

	return nil
}

// StackTop 获取虚拟机栈顶对象
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

// 将object压入栈
func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize { // 栈溢出
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

// 弹栈
func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}
