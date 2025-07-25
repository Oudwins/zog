package es

import (
	"github.com/Oudwins/zog/zconst"
)

var Map zconst.LangMap = map[zconst.ZogType]map[zconst.ZogIssueCode]string{
	zconst.TypeString: {
		zconst.NotIssueCode(zconst.IssueCodeLen):             "Cadena no debe tener exactamente {{len}} caracter(es)",
		zconst.NotIssueCode(zconst.IssueCodeEmail):           "No debe ser un correo electrónico válido",
		zconst.NotIssueCode(zconst.IssueCodeUUID):            "No debe ser un UUID válido",
		zconst.NotIssueCode(zconst.IssueCodeMatch):           "Cadena no es válida",
		zconst.NotIssueCode(zconst.IssueCodeURL):             "No debe ser una URL válida",
		zconst.NotIssueCode(zconst.IssueCodeHasPrefix):       "Cadena no debe comenzar con {{prefix}}",
		zconst.NotIssueCode(zconst.IssueCodeHasSuffix):       "Cadena no debe terminar con {{suffix}}",
		zconst.NotIssueCode(zconst.IssueCodeContains):        "Cadena no debe contener {{contained}}",
		zconst.NotIssueCode(zconst.IssueCodeContainsDigit):   "Cadena no debe contener ningún dígito",
		zconst.NotIssueCode(zconst.IssueCodeContainsUpper):   "Cadena no debe contener ninguna letra mayúscula",
		zconst.NotIssueCode(zconst.IssueCodeContainsLower):   "Cadena no debe contener ninguna letra minúscula",
		zconst.NotIssueCode(zconst.IssueCodeContainsSpecial): "Cadena no debe contener ningún carácter especial",
		zconst.NotIssueCode(zconst.IssueCodeOneOf):           "Cadena no debe ser una de las siguientes: {{one_of_options}}",
		zconst.IssueCodeRequired:                             "Es obligatorio",
		zconst.IssueCodeNotNil:                               "No debe estar vacio",
		zconst.IssueCodeMin:                                  "Cadena debe contener al menos {{min}} caracter(es)",
		zconst.IssueCodeMax:                                  "Cadena debe contener como máximo {{max}} caracter(es)",
		zconst.IssueCodeLen:                                  "Cadena debe tener exactamente {{len}} caracter(es)",
		zconst.IssueCodeEmail:                                "Debe ser un correo electrónico válido",
		zconst.IssueCodeUUID:                                 "Debe ser un UUID válido",
		zconst.IssueCodeMatch:                                "Cadena no es válida",
		zconst.IssueCodeURL:                                  "Debe ser una URL válida",
		zconst.IssueCodeHasPrefix:                            "Cadena debe comenzar con {{prefix}}",
		zconst.IssueCodeHasSuffix:                            "Cadena debe terminar con {{suffix}}",
		zconst.IssueCodeContains:                             "Cadena debe contener {{contained}}",
		zconst.IssueCodeContainsDigit:                        "Cadena debe contener al menos un dígito",
		zconst.IssueCodeContainsUpper:                        "Cadena debe contener al menos una letra mayúscula",
		zconst.IssueCodeContainsLower:                        "Cadena debe contener al menos una letra minúscula",
		zconst.IssueCodeContainsSpecial:                      "Cadena debe contener al menos un carácter especial",
		zconst.IssueCodeOneOf:                                "Cadena debe ser una de las siguientes: {{one_of_options}}",
		zconst.IssueCodeFallback:                             "Cadena no es válida",
	},
	zconst.TypeBool: {
		zconst.IssueCodeRequired: "Es obligatorio",
		zconst.IssueCodeNotNil:   "No debe estar vacio",
		zconst.IssueCodeTrue:     "Debe ser verdadero",
		zconst.IssueCodeFalse:    "Debe ser falso",
		zconst.IssueCodeFallback: "Valor no es válido",
	},
	zconst.TypeNumber: {
		zconst.IssueCodeRequired:                   "Es obligatorio",
		zconst.IssueCodeNotNil:                     "No debe estar vacio",
		zconst.IssueCodeLTE:                        "Número debe ser menor o igual a {{lte}}",
		zconst.IssueCodeLT:                         "Número debe ser menor que {{lt}}",
		zconst.IssueCodeGTE:                        "Número debe ser mayor o igual a {{gte}}",
		zconst.IssueCodeGT:                         "Número debe ser mayor que {{gt}}",
		zconst.IssueCodeEQ:                         "Número debe ser igual a {{eq}}",
		zconst.NotIssueCode(zconst.IssueCodeEQ):    "Número no debe ser igual a {{eq}}",
		zconst.IssueCodeOneOf:                      "Número debe ser uno de los siguientes: {{one_of_options}}",
		zconst.NotIssueCode(zconst.IssueCodeOneOf): "Número no debe ser uno de los siguientes: {{one_of_options}}",
		zconst.IssueCodeFallback:                   "Número no es válido",
	},
	zconst.TypeTime: {
		zconst.IssueCodeRequired: "Es obligatorio",
		zconst.IssueCodeNotNil:   "No debe estar vacio",
		zconst.IssueCodeAfter:    "Fecha debe ser posterior a {{after}}",
		zconst.IssueCodeBefore:   "Fecha debe ser anterior a {{before}}",
		zconst.IssueCodeEQ:       "Fecha debe ser igual a {{eq}}",
		zconst.IssueCodeFallback: "Fecha no es válida",
	},
	zconst.TypeSlice: {
		zconst.IssueCodeRequired:                      "Es obligatorio",
		zconst.IssueCodeNotNil:                        "No debe estar vacio",
		zconst.IssueCodeMin:                           "Lista debe contener al menos {{min}} elementos",
		zconst.IssueCodeMax:                           "Lista debe contener como máximo {{max}} elementos",
		zconst.IssueCodeLen:                           "Lista debe contener exactamente {{len}} elementos",
		zconst.NotIssueCode(zconst.IssueCodeLen):      "Lista no debe contener exactamente {{len}} elementos",
		zconst.IssueCodeContains:                      "Lista debe contener {{contained}}",
		zconst.NotIssueCode(zconst.IssueCodeContains): "Lista no debe contener {{contained}}",
		zconst.IssueCodeFallback:                      "Lista no es válida",
	},
	zconst.TypeStruct: {
		zconst.IssueCodeRequired: "Es obligatorio",
		zconst.IssueCodeNotNil:   "No debe estar vacio",
		zconst.IssueCodeFallback: "Estructura no es válida",
		// JSON
		zconst.IssueCodeInvalidJSON: "JSON no válido",
		// ZHTTP ISSUES
		zconst.IssueCodeZHTTPInvalidForm:  "Formulario no válido",
		zconst.IssueCodeZHTTPInvalidQuery: "Parámetros de consulta no válidos",
	},
}
