package vm

import (
	"Monkey/code"
	"Monkey/compiler"
	"Monkey/object"
	"fmt"
)

const StackSize = 2048 // 虚拟机栈大小

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

var Null = &object.Null{}

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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			if err := vm.executeBinaryOperation(op); err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpTrue:
			if err := vm.push(True); err != nil {
				return err
			}
		case code.OpFalse:
			if err := vm.push(False); err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			if err := vm.executeComparison(op); err != nil {
				return err
			}
		case code.OpBang:
			if err := vm.executeBangOperator(); err != nil {
				return err
			}
		case code.OpMinus:
			if err := vm.executeMinusOperator(); err != nil {
				return err
			}
		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}
		case code.OpNull:
			if err := vm.push(Null); err != nil {
				return err
			}
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

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

// 二元运算
func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left object.Object, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}
