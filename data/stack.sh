#! /bin/bash

stacknumbers=()
stackarray="0"
activearray="0"

#POP
function pop() {
	local length=${#stacknumbers[@]}
	echo ${stacknumbers[length-1]}
	unset stacknumbers[$length-1]
}

function tobytes() {
	local x=$1
	local i=0
	local y=()
	while [ $i -lt ${#x} ]; do y[$i]=${x:$i:1};  i=$((i+1));done
	echo $y
}

function stdout {
	stackarray=stackarray-1
	
	local array = $(eval '${stackarray_'$stackarray'[@]}')
	local len=${#array[@]}
	local i
	for ((i=0;i<=$len;i++)); do
		local c=${array[i]}
		printf \\$(printf '%03o\t' "$c")
	done
}

function make() {
	local array=()
	local i
	for ((i=0;i<$1;i++)); do
		array[i]=0
	done
	echo ${#array[@]}
}


toUpper() {
	local backupref
	if [ tmp = "$1" ]; then
		 declare -n backupref="$1";
		 local -n tmp=backupref
		  
		   for ((i=0; i<"${#tmp[@]}"; i++)); do
			 tmp[i]="${tmp[i]^^}"
			done;
		 
		 return
	fi
   declare -n tmp="$1"; 
   for ((i=0; i<"${#tmp[@]}"; i++)); do
      tmp[i]="${tmp[i]^^}"
    done;
   
}

defarr() {
	local arr=( "tests" "tests 1" "tests 3" )
}

function main() {
	#local arr=(t1 t2 t3)
	local tmp=( "test" "test 1" "test 3" )
	toUpper tmp
	printf "%s\n" "${tmp[@]}"
}
#main


function recursion() {
	declare -n tmp="$1";
	tmp+=(1)
	if [ ${#tmp[@]} -gt 4 ]; then 
		echo ${tmp[@]}
		return 
	fi
	tmp=$(recursion ${tmp[@]})
}

function sort() {
	local list=()
	recursion list
	echo ${list[@]}
}

sort
