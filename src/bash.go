package uct

import "flag"
import "fmt"
import "strconv"

var BashReserved = []string{
	

}

//This is the Java compiler for uct.
var Bash bool

func init() {
	flag.BoolVar(&Bash, "bash", false, "Target Bash")

	RegisterAssembler(BashAssembly, &Bash, "sh", "#")

	for _, word := range BashReserved {
		BashAssembly[word] = Reserved()
	}
	BashAssemblyPointer=&BashAssembly
}

var BashAssemblyPointer *Assemblable

var BashAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   `#! /bin/bash
ERROR=0
stacknumbers=()
stackarray="0"
activearray="0"
stackpipes=()

function _div() {
	if [ $(bc <<< "$2 == 0") -eq 1 ]; then
		if [ $(bc <<< "$1 == 0") -eq 1 ]; then
			od -A n -t d -N 1 /dev/urandom
		else
			echo 0
		fi
		return
	fi
	echo $(bc <<< "$1 / $2")
}

function _mod() {
	if [ $(bc <<< "$2 == 0") -eq 1 ]; then
		echo 0
		return
	fi
	echo $(bc <<< "$1 % $2 + ($2*(($1<0)*($2>0) + ($1>0)*($2<0)) )")
}

function _get() {
	local len=$(eval 'echo ${#'$activearray'[@]}')
	local index=$(_mod ${stacknumbers[${#stacknumbers[@]}-1]} $len)
	eval 'echo ${'$activearray"[$index]}"
}

function _set() {
	local len=$(eval 'echo ${#'$activearray'[@]}')
	local index=$(_mod ${stacknumbers[${#stacknumbers[@]}-1]} $len)
	eval $activearray"[$index]"=$1
}

function make() {
	local size=${stacknumbers[${#stacknumbers[@]}-1]}; unset stacknumbers[${#stacknumbers[@]}-1]
	local array=()
	local i
	for ((i=0;i<$size;i++)); do
		array[i]=0
	done
	eval "stackarray_$stackarray=array"; eval "stackarraydata_"$stackarray='${array[@]}'; stackarray=$((stackarray+1))
}

function load {
	stackarray=$((stackarray - 1)); declare -n text=$(eval 'echo ${stackarray_'$stackarray'[@]}')
	if [ -z "$text" ] || ! [[ "${text[@]}" == "$(eval 'echo ${stackarraydata_'$stackarray'[@]}')" ]]; then 
		text=($(eval 'echo ${stackarraydata_'$stackarray'[@]}'));
	fi
	local name;
	local result=()
	local variable=""
	
	local len=${#text[@]}
	if [ "$len" -gt 1 ] && [ ${text[0]} = 36 ]; then
		local i
		for ((i=1;i<=$len;i++)); do
			local c=${text[$i]}
			name="$name$(printf \\$(printf '%03o\t' "$c"))"
		done
		
		variable=$(eval echo -n "\$$name")
	else
		variable=$(eval echo -n "\$$(stack_index "$text_1"_"$text_2" 0)")
	fi
	for ((i=0;i<=${#variable};i++)); do
		result+=($(printf %d\\n \'"${variable:$i:1}"))
	done
	eval "stackarray_$stackarray=result"; eval "stackarraydata_"$stackarray='${result[@]}'; stackarray=$((stackarray+1))
}

function stdout {
	stackarray=$((stackarray - 1))

	declare -n array=$(eval 'echo ${stackarray_'$stackarray'[@]}');
	if [ -z "$array" ] || ! [[ "${array[@]}" == "$(eval 'echo ${stackarraydata_'$stackarray'[@]}')" ]]; then 
		array=($(eval 'echo ${stackarraydata_'$stackarray'[@]}'));
	fi
	local len=${#array[@]}
	local i
	for ((i=0;i<=$len;i++)); do
		local c=${array[i]}
		printf \\$(printf '%03o\t' "$c")
	done
}

_stdi=0
function stdin {
	local len=${stacknumbers[${#stacknumbers[@]}-1]}; unset stacknumbers[${#stacknumbers[@]}-1]
	
	if [ $len -eq 0 ]; then
		len=-10
	fi
	
	local text=()
	
	while true; do
		while [ $_stdi -lt ${#_stdb} ]; do
			char=$(printf %d\\n \'"${_stdb:$_stdi:1}")
			
			_stdi=$((_stdi + 1))
			
			if [ $(bc <<< "$len > 0") -eq 1 ]; then
				text+=($char)
				len=$(bc <<< "$len - 1")
				if [ $(bc <<< "$len <= 0") -eq 1 ]; then
					eval "stackarray_$stackarray=text"; eval "stackarraydata_"$stackarray='${text[@]}'; stackarray=$((stackarray+1))
					return
				fi
			else
				if [ $(bc <<< "$len*-1 == $char") -eq 1 ]; then
					eval "stackarray_$stackarray=text"; eval "stackarraydata_"$stackarray='${text[@]}'; stackarray=$((stackarray+1))
					return
				else
					text+=($char)
				fi
			fi
		done
		if [ $_stdi -eq ${#_stdb} ]; then
			if [ ${#_stdb} -gt 0 ]; then
				if [ $(bc <<< "$len > 0") -eq 1 ]; then
					text+=("10") #push 10
					len=$(bc <<< "$len - 1")
					if [ $(bc <<< "$len <= 0") -eq 1 ]; then
						eval "stackarray_$stackarray=text"; eval "stackarraydata_"$stackarray='${text[@]}'; stackarray=$stackarray+1
						return
					fi
				else
					if [ -10 -eq $len ]; then
						eval "stackarray_$stackarray=text"; eval "stackarraydata_"$stackarray='${text[@]}'; stackarray=$stackarray+1
						return
					else
						text+=("10") #push 10
					fi
				fi
			fi
		fi
		read _stdb
		if [ ${#_stdb} -eq 0 ]; then
			eval "stackarray_$stackarray=text"; eval "stackarraydata_"$stackarray='${text[@]}'; stackarray=$stackarray+1
			return
		fi
	done
}

slice() {
    local start end array_length length
    if [[ $1 == *:* ]]; then
        IFS=":"; read -r start end <<<"$1"
        shift
        array_length="$#"
        # defaults
        [ -z "$end" ] && end=$array_length
        [ -z "$start" ] && start=0
        (( start < 0 )) && let "start=(( array_length + start ))"
        (( end < 0 )) && let "end=(( array_length + end ))"
    else
        start="$1"
        shift
        array_length="$#"
        (( start < 0 )) && let "start=(( array_length + start ))"
        let "end=(( start + 1 ))"
    fi
    let "length=(( end - start ))"
    (( start < 0 )) && return 1
    # check bounds
    (( length < 0 )) && return 1
    (( start < 0 )) && return 1
    (( start >= array_length )) && return 1
    # parameters start with $1, so add 1 to $start
    let "start=(( start + 1 ))"
    echo "${@: $start:$length}"
}
`,
		Args: 1,
	},

	"FOOTER": Instruction{},

	"FILE": Instruction{
		Path: "",
	},
	
	"BASH": Instruction{All:true},

	"NUMBER": is("(echo '%s')", 1),
	"BIG": 	is("(echo '%s')", 1),
	"SIZE":   is("(echo ${#%s[@]})", 1),
	
	"STRING": Instruction{	
		Data: " ",
		Function: func(args []string) string {
			str, _ := strconv.Unquote(args[0])
			
			var result string = "("
			for _, char := range str {
				result += fmt.Sprint(int(char))+" "
			}
			result += ")"
			return result
		},
	},
	
	"ERRORS":  is("ERROR", 1),
	
	"LINK":  is(""),
	"CONNECT":  is(""),
	"SLICE":  is("slice ${stacknumbers[${#stacknumbers[@]}-1]}:${stacknumbers[${#stacknumbers[@]}-2]}; unset stacknumbers[${#stacknumbers[@]}-1]; unset stacknumbers[${#stacknumbers[@]}-1]"),

	"SOFTWARE": Instruction{
		Data:   "function main() {\n\techo -n",
		Indent:      1,
	},
	"EXIT": Instruction{
		Indented:    1,
		Data:        "}\nmain\n",
		Else: &Instruction{
			Data:   "exit $ERROR",
		},
	},
	
	"EVAL": is(""),

	"FUNCTION": is("function %s() {", 1, 1),
	"RETURN": Instruction{
		Indented:    1,
		Data:        "}",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "return",
		},
	},
	
	"SCOPE": is(`stackpipes+=(\"%s\")`, 1),
	
	"EXE": is("eval $%s", 1),

	"PUSH": is("stacknumbers+=($%s)", 1),
	"PULL": is("local %s=${stacknumbers[${#stacknumbers[@]}-1]}; unset stacknumbers[${#stacknumbers[@]}-1]", 1),

	"PUT":   is("eval $activearray'+=($%s)'", 1),
	"POP":   is("local %s=$(eval 'echo ${'$activearray'[${#'$activearray'[@]}-1]}'); eval 'unset '$activearray'[${#'$activearray'[@]}-1]'", 1),
	"PLACE": is("activearray=%s", 1),
	"ARRAY":  is("local %s=(); activearray=%s", 1),
	"RENAME": is("%s=($(eval 'echo ${'$activearray'[@]}'))", 1),
	"RELOAD": is("", 1),

	"SHARE": is(`eval "stackarray_$stackarray=%s"; eval "stackarraydata_"$stackarray='${%s[@]}'; stackarray=$((stackarray+1))`, 1),
	
	"TESTREF": Instruction{
		Args: 0,
	},
	
	"GRAB":  Instruction{
		Data: " ",
		Args: 1,
		Function: func(args []string) string {
			(*BashAssemblyPointer)["TESTREF"] = Instruction{Args:(*BashAssemblyPointer)["TESTREF"].Args+1}
			id := fmt.Sprint((*BashAssemblyPointer)["TESTREF"].Args)
			return 	`
	stackarray=$((stackarray-1))
	local -n _testref`+id+`=$(eval 'echo $stackarray_'$stackarray)
	local _evalstr='local -n `+args[0]+`=$stackarray_'$stackarray
	if [ -z "$_testref`+id+`" ] || ! [[ "${_testref`+id+`[@]}" == "$(eval 'echo ${stackarraydata_'$stackarray'[@]}')" ]]; then
		_evalstr='local `+args[0]+`=($(echo ${stackarraydata_'$stackarray'[@]}))'
	fi
	eval $_evalstr
	`
		},
	},

	"RELAY": is("stackpipes+=($%s)", 1),
	"TAKE":  is("local %s=${stackpipes[${#stackpipes[@]}-1]}; unset stackpipes[${#stackpipes[@]}-1]", 1),

	"GET": is("local %s=$(_get); unset stacknumbers[${#stacknumbers[@]}-1]", 1),
	"SET": is("_set $%s; unset stacknumbers[${#stacknumbers[@]}-1]", 1),

	"VAR": is("local %s=0", 1),

	"OPEN":   is("open"),
	"EXECUTE": is("execute"),
	"DELETE": is("delete"),
	"LOAD":   is("load"),
	"OUT":    is("out"),
	"STAT":   is("info"),
	"IN":     is("inn"),
	"STDOUT": is("stdout"),
	"STDIN":  is("stdin"),
	"HEAP":   is("heap"),
	"HEAPIT":   is("heapit"),
	
	"MAKE":   is("make"),

	"CLOSE": is("", 1),

	"LOOP":   is("while true; do", 0, 1),
	"BREAK":  is("break"),
	"REPEAT": is("done", 0, -1, -1),

	"IF":   is("if [ $(bc <<< \"$%s != 0\") -eq 1 ]; then", 1, 1),
	"ELSE": is("else", 0, 0, -1),
	"END":  is("fi", 0, -1, -1),

	"RUN":  is("%s", 1),
	"DATA": is_data("%s=(", "%s", ",", ")"),

	"FORK": is("%s &\n", 1),

	"ADD": is("local %s=$(bc <<< \"$%s + $%s\")", 3),
	"SUB": is("local %s=$(bc <<< \"$%s - $%s\")", 3),
	"MUL": is("local %s=$(bc <<< \"$%s * $%s\")", 3),
	"DIV": is("local %s=$(_div $%s $%s)", 3),
	"MOD": is("local %s=$(_mod $%s $%s)", 3), 

	"POW": is("local %s=$(bc <<< \"$%s ^ $%s\")", 3),

	"SLT": is("local %s=$(bc <<< \"$%s < $%s\")", 3),
	"SEQ": is("local %s=$(bc <<< \"$%s == $%s\")", 3),
	"SGE": is("local %s=$(bc <<< \"$%s >= $%s\")", 3),
	"SGT": is("local %s=$(bc <<< \"$%s > $%s\")", 3),
	"SNE": is("local %s=$(bc <<< \"$%s != $%s\")", 3),
	"SLE": is("local %s=$(bc <<< \"$%s <= $%s\")", 3),

	"JOIN": is("local %s=(${%s[@]} ${%s[@]})", 3),
	"ERROR": is("ERROR=$%s", 1),
}
