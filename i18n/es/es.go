package es

import (
	"github.com/Oudwins/zog/zconst"
)

var Map zconst.LangMap = map[zconst.ZogType]map[zconst.ZogErrCode]string{
	zconst.TypeString: {
		zconst.ErrCodeRequired:        "Es obligatorio",
		zconst.ErrCodeMin:             "Cadena debe contener al menos {{min}} caracter(es)",
		zconst.ErrCodeMax:             "Cadena debe contener como máximo {{max}} caracter(es)",
		zconst.ErrCodeLen:             "Cadena debe tener exactamente {{len}} caracter(es)",
		zconst.ErrCodeEmail:           "Debe ser un correo electrónico válido",
		zconst.ErrCodeURL:             "Debe ser una URL válida",
		zconst.ErrCodeHasPrefix:       "Cadena debe comenzar con {{prefix}}",
		zconst.ErrCodeHasSuffix:       "Cadena debe terminar con {{suffix}}",
		zconst.ErrCodeContains:        "Cadena debe contener {{contained}}",
		zconst.ErrCodeContainsDigit:   "Cadena debe contener al menos un dígito",
		zconst.ErrCodeContainsUpper:   "Cadena debe contener al menos una letra mayúscula",
		zconst.ErrCodeContainsLower:   "Cadena debe contener al menos una letra minúscula",
		zconst.ErrCodeContainsSpecial: "Cadena debe contener al menos un carácter especial",
		zconst.ErrCodeOneOf:           "Cadena debe ser una de las siguientes: {{one_of_options}}",
		zconst.ErrCodeFallback:        "Cadena no es válida",
	},
	zconst.TypeBool: {
		zconst.ErrCodeRequired: "Es obligatorio",
		zconst.ErrCodeTrue:     "Debe ser verdadero",
		zconst.ErrCodeFalse:    "Debe ser falso",
		zconst.ErrCodeFallback: "Valor no es válido",
	},
	zconst.TypeNumber: {
		zconst.ErrCodeRequired: "Es obligatorio",
		zconst.ErrCodeLTE:      "Número debe ser menor o igual a {{lte}}",
		zconst.ErrCodeLT:       "Número debe ser menor que {{lt}}",
		zconst.ErrCodeGTE:      "Número debe ser mayor o igual a {{gte}}",
		zconst.ErrCodeGT:       "Número debe ser mayor que {{gt}}",
		zconst.ErrCodeEQ:       "Número debe ser igual a {{eq}}",
		zconst.ErrCodeOneOf:    "Número debe ser uno de los siguientes: {{options}}",
		zconst.ErrCodeFallback: "Número no es válido",
	},
	zconst.TypeTime: {
		zconst.ErrCodeRequired: "Es obligatorio",
		zconst.ErrCodeAfter:    "Fecha debe ser posterior a {{after}}",
		zconst.ErrCodeBefore:   "Fecha debe ser anterior a {{before}}",
		zconst.ErrCodeEQ:       "Fecha debe ser igual a {{eq}}",
		zconst.ErrCodeFallback: "Fecha no es válida",
	},
	zconst.TypeSlice: {
		zconst.ErrCodeRequired: "Es obligatorio",
		zconst.ErrCodeMin:      "Lista debe contener al menos {{min}} elementos",
		zconst.ErrCodeMax:      "Lista debe contener como máximo {{max}} elementos",
		zconst.ErrCodeLen:      "Lista debe contener exactamente {{len}} elementos",
		zconst.ErrCodeContains: "Lista debe contener {{contained}}",
		zconst.ErrCodeFallback: "Lista no es válida",
	},
	zconst.TypeStruct: {
		zconst.ErrCodeRequired: "Es obligatorio",
		zconst.ErrCodeFallback: "Estructura no es válida",
	},
}
