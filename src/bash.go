//I can't beleive this works :O
//Terrible peformance... LOL
package main

import "errors"
import "flag"
import "strconv"
import "strings"

//Go flag.
var Bash bool
func init() {
	flag.BoolVar(&Bash, "bash", false, "Target Bash")
	
	RegisterAssembler(new(BashAssembler), &Bash, "bash", "#")
}

type BashAssembler struct {
	Indentation int
}


func (g *BashAssembler) indt(n ...int) string {
	if len(n) > 0 {
		return strings.Repeat("\t", int(g.Indentation)+n[0])
	} else {
		return strings.Repeat("\t", int(g.Indentation))
	}
}

func (g *BashAssembler) Header() []byte {
	return []byte(
	`#! /bin/bash

# A stack, using bash arrays.
# ---------------------------------------------------------------------------

# Create a new stack.
#
# Usage: stack_new name
#
# Example: stack_new x
function stack_new
{
    : ${1?'Missing stack name'}
    if stack_exists $1
    then
        echo "Stack already exists -- $1" >&2
        return 1
    fi

    eval "declare -ag _stack_$1"
    eval "declare -ig _stack_$1_i"
    eval "let _stack_$1_i=0"
    return 0
}

# Destroy a stack
#
# Usage: stack_destroy name
function stack_destroy
{
    : ${1?'Missing stack name'}
    eval "unset _stack_$1 _stack_$1_i"
    return 0
}

# Push one or more items onto a stack.
#
# Usage: stack_push stack item ...
function stack_push
{
    : ${1?'Missing stack name'}
    : ${2?'Missing item(s) to push'}

    if no_such_stack $1
    then
        echo "No such stack -- $1" >&2
        return 1
    fi

    stack=$1
    shift 1

    while (( $# > 0 ))
    do
        eval '_i=$'"_stack_${stack}_i"
        eval "_stack_${stack}[$_i]='$1'"
        eval "let _stack_${stack}_i+=1"
        shift 1
    done

    unset _i
    return 0
}

# Print a stack to stdout.
#
# Usage: stack_print name
function stack_print
{
    : ${1?'Missing stack name'}

    if no_such_stack $1
    then
        echo "No such stack -- $1" >&2
        return 1
    fi

    tmp=""
    eval 'let _i=$'_stack_$1_i
    while (( $_i > 0 ))
    do
        let _i=${_i}-1
        eval 'e=$'"{_stack_$1[$_i]}"
        tmp="$tmp $e"
    done
    echo "(" $tmp ")"
}

# Get the size of a stack
#
# Usage: stack_size name var
#
# Example:
#    stack_size mystack n
#    echo "Size is $n"
function stack_size
{
    : ${1?'Missing stack name'}
    : ${2?'Missing name of variable for stack size result'}
    if no_such_stack $1
    then
        echo "No such stack -- $1" >&2
        return 1
    fi
    eval "$2"='$'"{#_stack_$1[*]}"
}

# Pop the top element from the stack.
#
# Usage: stack_pop name var
#
# Example:
#    stack_pop mystack top
#    echo "Got $top"
function stack_pop
{
    : ${1?'Missing stack name'}
    : ${2?'Missing name of variable for popped result'}

    eval 'let _i=$'"_stack_$1_i"
    if no_such_stack $1
    then
        echo "No such stack -- $1" >&2
        return 1
    fi

    if [[ "$_i" -eq 0 ]]
    then
        echo "Empty stack -- $1" >&2
        return 1
    fi

    let _i-=1
    eval "$2"='$'"{_stack_$1[$_i]}"
    eval "unset _stack_$1[$_i]"
    eval "_stack_$1_i=$_i"
    unset _i
    return 0
}

function stack_len 
{
 	: ${1?'Missing stack name'}

    eval 'let _i=$'"_stack_$1_i"
    if no_such_stack $1
    then
        echo "No such stack -- $1" >&2
        return 1
    fi

    if [[ "$_i" -eq 0 ]]
    then
        echo 0
        return 0
    fi

   echo $_i
   return 0
}

function stack_index
{
 	: ${1?'Missing stack name'}
 	: ${2?'Missing index'}

    if no_such_stack $1
    then
        echo "No such stack -- $1" >&2
        return 1
    fi

    eval "echo "'$'"{_stack_$1[$2]}"
    return 0
}

function stack_set
{
 	: ${1?'Missing stack name'}
 	: ${2?'Missing index'}
 	: ${3?'Missing value'}

    if no_such_stack $1
    then
        echo "No such stack -- $1" >&2
        return 1
    fi

    eval "_stack_$1[$2]=$3"
    return 0
}

function stack_join {
	stack_new "$1_$S"
	local var
	local i
	for ((i=0;i<$(stack_len "$2"_"$3" );i++)); do
		stack_push "$1_$S" $(stack_index "$2"_"$3" $i) 
	done
	for ((i=0;i<$(stack_len "$4"_"$5" );i++)); do
		stack_push "$1_$S" $(stack_index "$4"_"$5" $i) 
	done
	S=$(expr $S + 1)
}

function no_such_stack
{
    : ${1?'Missing stack name'}
    stack_exists $1
    ret=$?
    declare -i x
    let x="1-$ret"
    return $x
}

function stack_exists
{
    : ${1?'Missing stack name'}

    eval '_i=$'"_stack_$1_i"
    if [[ -z "$_i" ]]
    then
        return 1
    else
        return 0
    fi
}

#Heap counter
S=0
stack_new N
stack_new N2
stack_new F
stack_new F2

ERROR=0

function push {
	stack_push N $1
}

function pushit {
	stack_push F2 $1
}

function popit {
	stack_pop F2 $1
}

function pop {
	stack_pop N $1
}

function pushstring {
	stack_push N2 $1
	stack_push N2 $2
}

function popstring {
	stack_pop N2 $2
	stack_pop N2 $1
}

function pushfunc {
	stack_push F $1
}

function popfunc {
	stack_pop F $1
}

function string {
	stack_new "$1_$S"
	S=$(expr $S + 1)
}

function stringpush {
	stack_push $1_$2 $3
}

function stringpop {
	stack_pop $1_$2 $3
}

function stringdata {
	stack_new "$1_$S"
	local var
	for var in "$@"
	do
		if [ $var != $1 ]; then
			stack_push "$1_$S" $var 
		fi
	done
	S=$(expr $S + 1)
}

function load {
	popstring text_1 text_2
	local name;
	local result_1=result; local result_2=$S; string result
	local variable=""
	
	local len=$(stack_len "$text_1"_"$text_2")
	if [ "$len" -gt 1 ] && [ $(stack_index "$text_1"_"$text_2" 0) = 36 ]; then
		local i
		for ((i=1;i<=$len;i++)); do
			local c=$(stack_index "$text_1"_"$text_2" $i)
			name="$name$(printf \\$(printf '%03o\t' "$c"))"
		done
		
		variable=$(eval echo -n "\$$name")
	else
		variable=$(eval echo -n "\$$(stack_index "$text_1"_"$text_2" 0)")
	fi
	for ((i=0;i<=${#variable};i++)); do
		stringpush $result_1 $result_2 $(printf %d\\n \'"${variable:$i:1}")
	done
	pushstring $result_1 $result_2
}

function open {
	popstring text_1 text_2
	local filename;
	
	local len=$(stack_len "$text_1"_"$text_2")
	local i
	for ((i=0;i<=$len;i++)); do
		local c=$(stack_index "$text_1"_"$text_2" $i)
		filename="$filename$(printf \\$(printf '%03o\t' "$c"))"
	done

	 eval "declare -ig $1_i=0"
	 eval 'declare -g $1_n="$filename"'
	  
	 if [ -e "$filename" ]; then
	 	push "0"
	 else
	 	push "-1"
	 fi
}

function stdout {
	local text_1
	local text_2
	popstring text_1 text_2
	
	local len=$(stack_len "$text_1"_"$text_2")
	local i
	for ((i=0;i<=$len;i++)); do
		local c=$(stack_index "$text_1"_"$text_2" $i)
		printf \\$(printf '%03o\t' "$c")
	done
}

function out {
	local text_1
	local text_2
	popstring text_1 text_2
	
	local n=$(eval "echo -n \$$1_n")
	if [ ${n: -1} = "/" ]; then
		if [ ! -d $(eval "echo -n \$$1_n") ]; then
			mkdir $(eval "echo -n \$$1_n") 2> /dev/null
		fi
		if [ "$?" = "1" ]; then
			push -1
			return
		fi
		push 0
		return
	elif [ ! -f $(eval "echo -n \$$1_n") ]; then
		touch $(eval "echo -n \$$1_n") 2> /dev/null
		if [ "$?" = "1" ]; then
			push -1
			return
		fi
	fi
	
	local len=$(stack_len "$text_1"_"$text_2")
	local i
	for ((i=0;i<=$len;i++)); do
		local c=$(stack_index "$text_1"_"$text_2" $i)
		if [ "$c" = "10" ]; then
			echo >> $(eval "echo -n \$$1_n")
		else
			echo -n $(printf \\$(printf '%03o\t' "$c")) >> $(eval "echo -n \$$1_n")
		fi
	done
	push 0
}

function inn {
	local len; pop len
	
	for ((i=0; i<=$len; i++)); do
		local v=$(tail -c +$(eval "echo -n \$$1_i") $(eval "echo -n \$$1_n") | head -c 1)
		if [ "$v" = "" ]; then
			push -1000
			return
		else
			push $(printf %d\\n \'"$v")
		fi
		declare $1_i=$(expr $(eval "echo -n \$$1_i") + 1)
	done
	
}

_stdi=0
function stdin {
	local len; pop len
	while true; do
		if [ $_stdi -eq ${#_stdb} ]; then
			if [ ${#_stdb} -gt 0 ]; then
				push 10
				len=$(bc <<< "$len - 1")
				if [ $(bc <<< "$len <= 0") -eq 1 ]; then
					return
				fi
			fi
		fi
		while [ $_stdi -lt ${#_stdb} ]; do
			push $(printf %d\\n \'"${_stdb:$_stdi:1}")
			
			_stdi=$(expr $_stdi + 1)
			
			len=$(bc <<< "$len - 1")
			if [ $(bc <<< "$len <= 0") -eq 1 ]; then
				return
			fi
		done
		read _stdb
	done
}

function div {
	local result=$(bc 2> /dev/null <<< "$1 / $2")
	if [ -z "$result" ]; then
		echo -n "0"
		exit 1
	else
		echo -n "$result"
	fi
}
`)
}

func (g *BashAssembler) SetFileName(s string) {
}

func (g *BashAssembler) Assemble(command string, args []string) ([]byte, error) {
	//Format argument expressions.
	//var expression bool
	for i, arg := range args {
		//RESERVED names in the language.
		switch arg {
			case "byte", "open", "load":
				args[i] = "u_"+args[i]
				arg = "u_"+arg
		}
		if arg == "=" {
			//expression = true
			continue
		}
		if _, err := strconv.Atoi(arg); err != nil {
			if arg[0] != '#' && arg[0] != '"' {
				args[i] = "$"+arg
				continue
			}
		}
		if arg[0] == '#' {
			args[i] = "$(stack_len \"$"+arg[1:]+"_1\"_\"$"+arg[1:]+"_2\")"
		}
		if arg[0] == '"' {
			var newarg string
			var j = i
			arg = arg[1:]
			
			stringloop:
			arg = strings.Replace(arg, "\\n", "\n", -1)
			for _, v := range arg {
				if v == '"' {
					goto end
				}
				newarg += strconv.Itoa(int(v))+" "
			}
			if len(arg) == 0 {
				goto end
			}
			newarg += strconv.Itoa(int(' '))+" "
			j++
			//println(arg)
			arg = args[j]
			goto stringloop
			end:
			//println(newarg)
			args[i] = newarg
		}
	} 
	switch command {
		case "ROUTINE":
			defer func() { g.Indentation++ }()
			return []byte("function main {\ntrue\n"), nil
		case "SUBROUTINE":
			defer func() { g.Indentation++ }()
			return []byte("function "+args[0][1:]+" {\n"), nil
		case "PUSH", "PUSHSTRING":
			if command == "PUSHSTRING" {
				return []byte(g.indt()+"pushstring "+args[0]+"_1 "+args[0]+"_2 \n"), nil
			} else {
				if len(args) == 1 {
					return []byte(g.indt()+"push "+args[0]+"\n"), nil
				} else {
					return []byte(g.indt()+"stringpush "+args[1]+"_1 "+args[1]+"_2 "+args[0]+" \n"), nil
				}
			}
		case "PUSHIT":
			return []byte(g.indt()+"pushit "+args[0][1:]+" \n"), nil
		case "PUSHFUNC":
			return []byte(g.indt()+"pushfunc "+args[0]+" \n"), nil
		case "POPFUNC":
			return []byte(g.indt()+"local "+args[0][1:]+"; popfunc "+args[0][1:]+" \n"), nil
		case "POPIT":
			return []byte(g.indt()+"local "+args[0][1:]+"; popit "+args[0][1:]+" \n"), nil
		case "FUNC":
			return []byte(g.indt()+args[0][1:]+"='"+args[1][1:]+"' \n"), nil
		
		//File stuff.
		case "OPEN":
			return []byte(g.indt()+"open "+args[0][1:]+"\n"), nil
		case "OUT":
			return []byte(g.indt()+"out "+args[0]+"\n"), nil
		case "IN":
			return []byte(g.indt()+"inn "+args[0]+"\n"), nil
		case "CLOSE":
			return []byte(""), nil
		
		case "ERROR":
			return []byte(g.indt()+"ERROR="+args[0]+"\n"), nil
			
		case "EXE":
			return []byte(g.indt()+"eval "+args[0]+"\n"), nil
		case "POP", "POPSTRING":
			if command == "POPSTRING" {
				return []byte(g.indt()+"local "+args[0][1:]+"_1; local "+args[0][1:]+"_2; popstring "+args[0][1:]+"_1 "+args[0][1:]+"_2 \n"), nil
			} else {
				if len(args) == 1 {
					return []byte(g.indt()+"local "+args[0][1:]+"; pop "+args[0][1:]+"\n"), nil
				} else {
					return []byte(g.indt()+"local "+args[0][1:]+" stringpop "+args[0]+" "+args[1][1:]+"_1 "+args[1][1:]+"_2 \n"), nil
				}
			}
		case "INDEX":
			return []byte(g.indt()+"local "+args[2][1:]+"=$(stack_index \""+args[0]+"_1\"_\""+args[0]+"_2\" "+args[1]+")\n"), nil
		case "SET":
			return []byte(g.indt()+"stack_set \""+args[0]+"_1\"_\""+args[0]+"_2\" "+args[1]+" "+args[2]+"\n"), nil
		case "VAR":
			if len(args) == 1 {
				return []byte(g.indt()+"local "+args[0][1:]+"=0 \n"), nil
			} else {
				return []byte(g.indt()+"local "+args[0][1:]+"="+args[1]+" \n"), nil
			}
		case "STRING":
			return []byte(g.indt()+args[0][1:]+"_1="+args[0][1:]+"; "+args[0][1:]+"_2=$S; string "+args[0][1:]+" \n"), nil
		case "STDOUT", "STDIN", "LOAD":
			return []byte(g.indt()+strings.ToLower(command)+"\n"), nil
		case "LOOP":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"while true; do\n"), nil
		case "END":
			if g.Indentation > 0 {
				g.Indentation--
				if g.Indentation == 0 {
					return []byte(g.indt()+"}\n"), nil
				}
			}
			return []byte(g.indt()+"fi\n"), nil
		case  "REPEAT":
			if g.Indentation > 0 {
				g.Indentation--
			}
			return []byte(g.indt()+"done\n"), nil
		case "DONE":
			if g.Indentation > 0 {
				g.Indentation--
			}
			return []byte(g.indt()+"}\n"), nil
		case "IF":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"if [ $(bc <<< \""+args[0]+"\") != '0' ]; then\n"), nil
		case "RUN":
			return []byte(g.indt()+args[0][1:]+"\n"), nil
		case "ELSE":
			return []byte(g.indt(-1)+"else\n"), nil
		case "ELSEIF":
			return []byte(g.indt(-1)+"elif [ $(bc <<< \""+args[0]+"\") != '0' ]; then\n"), nil
		case "BREAK":
			return []byte(g.indt()+"break\n"), nil
		case "RETURN":
			return []byte(g.indt()+"return\n"), nil
		case "STRINGDATA":
			return []byte(g.indt()+args[0][1:]+"_1="+args[0][1:]+"; "+args[0][1:]+"_2=$S; stringdata "+args[0][1:]+" "+args[1]+" \n"), nil
		case "JOIN":
			return []byte(g.indt()+"local "+args[0][1:]+"_1="+args[0][1:]+"; "+args[0][1:]+"_2=$S;"+
			" stack_join "+args[0][1:]+" "+args[1]+"_1 "+args[1]+"_2 "+args[2]+"_1 "+args[2]+"_2 \n"), nil
		
		//Maths.
		case "ADD":
			return []byte(g.indt()+args[0][1:]+"=$(bc <<< \""+args[1]+" + "+args[2]+"\")\n"), nil
		case "SUB":
			return []byte(g.indt()+args[0][1:]+"=$(bc <<< \""+args[1]+" - "+args[2]+"\")\n"), nil
		case "MUL":
			return []byte(g.indt()+args[0][1:]+"=$(bc <<< \""+args[1]+" * "+args[2]+"\")\n"), nil
		case "DIV":
			return []byte(g.indt()+args[0][1:]+"=$(div \""+args[1]+"\" \""+args[2]+"\")\n"), nil
		case "MOD":
			a, b := args[1], args[2]
			return []byte(g.indt()+args[0][1:]+"=$(bc <<< \""+a+" % "+b+" + ("+
			b+"*(("+a+"<0)*("+b+">0) + ("+a+">0)*("+b+"<0)) )\")\n"), nil
		case "POW":
			return []byte(g.indt()+args[0][1:]+"=$(bc <<< \""+args[1]+" ^ "+args[2]+"\")\n"), nil
			
		case "SLT":
			return []byte(g.indt()+args[0][1:]+"=$(bc <<< \""+args[1]+" < "+args[2]+"\")\n"), nil
		case "SEQ":
			return []byte(g.indt()+args[0][1:]+"=$(bc <<< \""+args[1]+" == "+args[2]+"\")\n"), nil
		case "SGE":
			return []byte(g.indt()+args[0][1:]+"=$(bc <<< \""+args[1]+" >= "+args[2]+"\")\n"), nil
		case "SGT":
			return []byte(g.indt()+args[0][1:]+"=$(bc <<< \""+args[1]+" > "+args[2]+"\")\n"), nil
		case "SNE":
			return []byte(g.indt()+args[0][1:]+"=$(bc <<< \""+args[1]+" != "+args[2]+"\")\n"), nil
		case "SLE":
			return []byte(g.indt()+args[0][1:]+"=$(bc <<< \""+args[1]+" <= "+args[2]+"\")\n"), nil
		default:
			/*if expression {
				return []byte(g.indt()+command+" "+strings.Join(args, " ")+"\n"), nil
			}*/
	}
	return nil, errors.New("Unrecognised expression: "+command+" "+strings.Join(args, " "))

}

func (g *BashAssembler) Footer() []byte {
	return []byte("main\n")
}
