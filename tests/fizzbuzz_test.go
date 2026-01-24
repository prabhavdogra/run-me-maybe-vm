package tests

var FizzBuzzTest = ProgramTestCase{
	name: "fizzbuzz",
	program: `@imp "stddefs.wm"
			@def N 100 
			jmp main
			push_str "Fizz"
			push_str "Buzz"
			push_str "\n"

			handle_three:
			get_str 0
			push STDOUT
			write
			inswap 0
			push 1
			add
			inswap 0
			jmp three_continue

			handle_five:
			get_str 1
			push STDOUT
			write
			inswap 0
			push 1
			add
			inswap 0
			jmp five_continue

			handle_number:
			dup
			int_to_str
			push STDOUT
			write
			get_str 2
			push STDOUT
			write
			jmp number_continue

			main:
			push 1

			loop:
			push 0
			inswap 0

			dup
			push 3
			mod
			zjmp handle_three

			three_continue:
			dup
			push 5
			mod
			zjmp handle_five
			five_continue:

			inswap 0
			zjmp handle_number 
			get_str 2
			push STDOUT
			write

			number_continue:
			push 1
			add

			push N 
			cmpl
			nzjmp end 
			pop


			jmp loop

			end:
			push 0
			exit`,
	expected: []string{
		"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz",
		"11", "Fizz", "13", "14", "FizzBuzz", "16", "17", "Fizz", "19", "Buzz",
		"Fizz", "22", "23", "Fizz", "Buzz", "26", "Fizz", "28", "29", "FizzBuzz",
		"31", "32", "Fizz", "34", "Buzz", "Fizz", "37", "38", "Fizz", "Buzz",
		"41", "Fizz", "43", "44", "FizzBuzz", "46", "47", "Fizz", "49", "Buzz",
		"Fizz", "52", "53", "Fizz", "Buzz", "56", "Fizz", "58", "59", "FizzBuzz",
		"61", "62", "Fizz", "64", "Buzz", "Fizz", "67", "68", "Fizz", "Buzz",
		"71", "Fizz", "73", "74", "FizzBuzz", "76", "77", "Fizz", "79", "Buzz",
		"Fizz", "82", "83", "Fizz", "Buzz", "86", "Fizz", "88", "89", "FizzBuzz",
		"91", "92", "Fizz", "94", "Buzz", "Fizz", "97", "98", "Fizz", "Buzz",
	},
	additionalFiles: StdDefs,
}
