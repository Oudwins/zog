package zes

import (
	"github.com/Oudwins/zog/conf"
	"github.com/Oudwins/zog/zconst"
)

var Map conf.LangMap = map[zconst.ZogType]map[zconst.ZogErrCode]string{
	zconst.TypeString: {
		zconst.ErrCodeRequired:        "Es obligatorio",
		zconst.ErrCodeMin:             "La cadena debe contener al menos {{min}} caracter(es)",
		zconst.ErrCodeMax:             "La cadena debe contener como máximo {{min}} caracter(es)",
		zconst.ErrCodeLen:             "La cadena debe tener exactamente {{len}} caracter(es)",
		zconst.ErrCodeEmail:           "Debe ser un correo electrónico válido",
		zconst.ErrCodeURL:             "Debe ser una URL válida",
		zconst.ErrCodeHasPrefix:       "La cadena debe comenzar con {{prefix}}",
		zconst.ErrCodeHasSuffix:       "La cadena debe terminar con {{suffix}}",
		zconst.ErrCodeContains:        "La cadena debe contener {{contained}}",
		zconst.ErrCodeContainsDigit:   "La cadena debe contener al menos un dígito",
		zconst.ErrCodeContainsUpper:   "La cadena debe contener al menos una letra mayúscula",
		zconst.ErrCodeContainsLower:   "La cadena debe contener al menos una letra minúscula",
		zconst.ErrCodeContainsSpecial: "La cadena debe contener al menos un carácter especial",
		zconst.ErrCodeOneOf:           "La cadena debe ser uno de los siguientes: {{one_of_options}}",
		zconst.ErrCodeFallback:        "La cadena no es válida",
	},
	zconst.TypeBool: {
		zconst.ErrCodeRequired: "Es obligatorio",
		zconst.ErrCodeTrue:     "Debe ser verdadero",
		zconst.ErrCodeFalse:    "Debe ser falso",
		zconst.ErrCodeFallback: "El valor no es válido",
	},
	zconst.TypeNumber: {
		zconst.ErrCodeRequired: "Es obligatorio",
		zconst.ErrCodeLTE:      "El número debe ser menor o igual a {{lte}}",
		zconst.ErrCodeLT:       "El número debe ser menor que {{lt}}",
		zconst.ErrCodeGTE:      "El número debe ser mayor o igual a {{gte}}",
		zconst.ErrCodeGT:       "El número debe ser mayor que {{gt}}",
		zconst.ErrCodeEQ:       "El número debe ser igual a {{eq}}",
		zconst.ErrCodeOneOf:    "El número debe ser uno de los siguientes: {{options}}",
		zconst.ErrCodeFallback: "El número no es válido",
	},
	zconst.TypeTime: {
		zconst.ErrCodeRequired: "Es obligatorio",
		zconst.ErrCodeAfter:    "La fecha debe ser posterior a {{after}}",
		zconst.ErrCodeBefore:   "La fecha debe ser anterior a {{before}}",
		zconst.ErrCodeEQ:       "La fecha debe ser igual a {{eq}}",
		zconst.ErrCodeFallback: "La fecha no es válida",
	},
	zconst.TypeSlice: {
		zconst.ErrCodeRequired: "Es obligatorio",
		zconst.ErrCodeMin:      "La colección debe contener al menos {{min}} elementos",
		zconst.ErrCodeMax:      "La colección debe contener como máximo {{max}} elementos",
		zconst.ErrCodeLen:      "La colección debe contener exactamente {{len}} elementos",
		zconst.ErrCodeContains: "La colección debe contener {{contained}}",
		zconst.ErrCodeFallback: "La colección no es válida",
	},
	zconst.TypeStruct: {
		zconst.ErrCodeRequired: "Es obligatorio",
		zconst.ErrCodeFallback: "La estructura no es válida",
	},
}
