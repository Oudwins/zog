package zog

// TODO as per coerce in zod
// z.Coerce.string(). // returns a string validator but will first convert the "number" to string
// z.Coerce().String(). // returns a string validator but will first convert the "number" to string. This wont work for schemas

// either we duplicate the code for each type
// or we add a parse() method to the validator which returns the value and will call a getValue() method on the value
// Coerce will hook into the getValue() method and will return the coerced value

// best idea is:
// 1. all validators will have a getValue/convertValue/valueOf(originalValue any) method
// 2. the getValue() method will check for a function inside the vlaidator to see if its nil, "converter"
// 3. if its nil, it will return the originalValue
// 4. if its not nil, it will call the function "converter" and return the result
// 5. The coerce struct can have a method for each type that will return a validator of that type with the converter function set
