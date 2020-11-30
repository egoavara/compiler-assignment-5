package cdtgo

import "unicode"

type Rule int

func (r Rule) IsNaming() bool {
	return unicode.IsUpper([]rune(r.String())[0])
}
func (r Rule) String() string {
	return ruleName[r-1]
}

var ruleName = [NO_RULES]string{
	// 0
	"PROGRAM",
	"translation_unit[0]",
	"translation_unit[1]",
	"external_dcl[0]",
	"external_dcl[1]",
	"FUNC_DEF",
	"FUNC_HEAD",
	"DCL_SPEC",
	"dcl_specifiers[0]",
	"dcl_specifiers[1]",
	// 10
	"dcl_specifier[0]",
	"dcl_specifier[1]",
	"CONST_TYPE",
	"INT_TYPE",
	"VOID_TYPE",
	"function_name",
	"FORMAL_PARA",
	"opt_formal_param[0]",
	"opt_formal_param[1]",
	"formal_param_list[0]",
	// 20
	"formal_param_list[1]",
	"PARAM_DCL",
	"COMPOUND_ST",
	"DCL_LIST",
	"DCL_LIST",
	"declaration_list[0]",
	"declaration_list[1]",
	"DCL",
	"init_dcl_list[0]",
	"init_dcl_list[1]",
	// 30
	"DCL_ITEM",
	"DCL_ITEM",
	"SIMPLE_VAR",
	"ARRAY_VAR",
	"opt_number[0]",
	"opt_number[1]",
	"STAT_LIST",
	"opt_stat_list[1]",
	"statement_list[0]",
	"statement_list[1]",
	// 40
	"statement[0]",
	"statement[1]",
	"statement[2]",
	"statement[3]",
	"statement[4]",
	"statement[5]",
	"statement[6]",
	"statement[7]",
	"statement[8]",
	"statement[9]",
	// 50
	"EXP_ST",
	"opt_expression[0]",
	"opt_expression[1]",
	"CASE_ST",
	"DEFAULT_ST",
	"CONTINUE_ST",
	"BREAK_ST",
	"IF_ST",
	"IF_ELSE_ST",
	"WHILE_ST",
	// 60
	"DO_WHILE_ST",
	"SWITCH_ST",
	"FOR_ST",
	"INIT_PART",
	"CONDITION_PART",

	"POST_PART",
	"RETURN_ST",
	"expression",
	"assignment_exp[0]",
	"ASSIGN_OP",
	// 70
	"ADD_ASSIGN",
	"SUB_ASSIGN",
	"MUL_ASSIGN",
	"DIV_ASSIGN",
	"MOD_ASSIGN",

	"logical_or_exp[0]",
	"LOGICAL_OR",
	"logical_and_exp[0]",
	"LOGICAL_AND",
	"equality_exp[0]",
	// 80
	"EQ",
	"NE",
	"relational_exp[0]",
	"GT",
	"LT",

	"GE",
	"LE",
	"additive_exp[0]",
	"ADD",
	"SUB",
	// 90
	"multiplicative_exp[0]",
	"MUL",
	"DIV",
	"REMAINDER",
	"unary_exp[0]",
	"UNARY_MINUS",
	"LOGICAL_NOT",
	"PRE_INC",
	"PRE_DEC",
	"postfix_exp[0]",
	// 100
	"INDEX",
	"CALL",
	"POST_INC",
	"POST_DEC",
	"opt_actual_param[0]",
	"opt_actual_param[1]",
	"ACTUAL_PARAM",
	"actual_param_list[0]",
	"actual_param_list[1]",
	"primary_exp[0]",
	// 110
	"primary_exp[1]",
	"primary_exp[2]",
}
