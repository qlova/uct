extern crate num;
extern crate rand;

use num::bigint::BigInt;
use num::bigint::{ToBigInt, Sign};
use num::ToPrimitive;
use std::char;

use std::cell::Cell;
use std::ptr;
use std::io;
use std::io::prelude::*;
use std::env;
use std::collections::HashMap;

use rand::os::OsRng;
use rand::Rng;

fn EmptyFunction(stack: &mut Stack) {}

pub struct Pipe {
	pub Function: fn(&mut Stack)
}

impl Pipe {
	pub fn Exe(&mut self, stack: &mut Stack) {
		(self.Function)(stack)
	}
}

#[derive(Clone, Debug)]
pub struct Pointer {
	pub address: usize,
	pub counter: *mut usize,
	pub collect: bool,
}

impl Drop for Pointer {
    fn drop(&mut self) {
		if !self.counter.is_null() && self.collect {
			let mut rc: usize = unsafe { *self.counter };
			if rc > 0 {
				rc = rc - 1;
			} else if self.address > 0 {
				println!("Not enough reference counting! for {}", self.address)
			}
			unsafe { *self.counter = rc }
			//println!("{}", rc);
			//Collect garbage in the future...
		}
    }
}

//This is the stack.
pub struct Stack {
	pub Numbers: Vec<BigInt>,
	pub Arrays: Vec<Pointer>,
	pub Pipes: Vec<Pipe>,
	
	pub ERROR: BigInt,
	pub ActiveArray: Pointer,
	
	pub Borrow: Vec<Vec<BigInt>>,
	pub Library: Vec<Cell<usize>>,
	
	pub Map: HashMap<String, BigInt>,
	
	//Heaps.
	pub TheHeap: Vec<Pointer>,
	pub HeapRoom: Vec<usize>,
	
} 

pub fn Div(a: &BigInt, b: &BigInt) -> BigInt {
	if *b == 0.to_bigint().unwrap() {
		if *a == 0.to_bigint().unwrap() {
			let mut rng = match OsRng::new() {
				Ok(g) => g,
				Err(e) => panic!("Failed to obtain OS RNG: {}", e)
			};

			let num = (rng.next_u32()%256)+1;
			return num.to_bigint().unwrap()
		} else {
			return 0.to_bigint().unwrap()
		}
	} else {
		return a/b
	}
}

pub fn Mod(a: &BigInt, n: &BigInt) -> BigInt {
	if (a < n) && *a > 0.to_bigint().unwrap() {
		return a.clone();
	}
	let d = Div(a, n);
	if d == 0.to_bigint().unwrap() && *a < 0.to_bigint().unwrap() && *n > -a {
		return a- (-n);
	}
	a % n
}

//Methods for the stack.
impl Stack {
	pub fn new() -> Stack {
		Stack {
			Numbers: Vec::new(),
			Arrays: Vec::new(),
			Pipes: Vec::new(),
			ERROR: 0.to_bigint().unwrap(),
			ActiveArray: Pointer{ address: 0, counter: ptr::null_mut(), collect: false },
			
			//Get past rust's annoying borrow checker.
			Borrow: Vec::new(),
			Library: Vec::new(),
			Map: HashMap::new(),
			
			TheHeap: Vec::new(),
			HeapRoom: Vec::new(),
		}
	}
	
	pub fn Array(&mut self) -> Pointer {
		
		let mut found = false;
		let mut value = 0;
		
		for (index, counter) in self.Library.iter().enumerate() {
			if counter.get() == 0 && index > 0 {
				//println!("FREE SPACE!");
				self.Borrow[index] = Vec::new();
				counter.set(1);
				found = true;
				value = index;
			}
		}
		
		if !found {
			self.Borrow.push(Vec::new());
			self.Library.push(Cell::new(1));
			value = self.Borrow.len()-1;
		}
		
		self.ActiveArray = Pointer{ 
			address: value, 
			counter: self.Library[value].as_ptr(),
			collect: false,
		};
		Pointer {
			address: value, 
			counter: self.Library[value].as_ptr(),
			collect: true,
		}
	}
	
	pub fn NewStringArray(&mut self, string: &str) -> Pointer {
		let mut result: Vec<BigInt> = Vec::new();
		for letter in string.chars() {
			result.push((letter as u32).to_bigint().unwrap());
		}
		
		self.Borrow.push(result);
		self.Library.push(Cell::new(1));
		self.ActiveArray.address = self.Borrow.len()-1;
		Pointer {
			address: self.Borrow.len()-1, 
			counter:self.Library[self.Borrow.len()-1].as_ptr(),
			collect: true,
		}
	}
	
	pub fn Push(&mut self, value: BigInt) {
		self.Numbers.push(value)
	}
	pub fn Pull(&mut self) -> BigInt {
		self.Numbers.pop().unwrap()
	}
	
	pub fn Share(&mut self, value: Pointer) {
		let ref mut cell = self.Library[value.address];
		cell.set(cell.get()+1);
		self.Arrays.push(value)
	}
	pub fn Grab(&mut self) -> Pointer {
		self.Arrays.pop().unwrap()
	}
	
	pub fn Relay(&mut self, value: Pipe) {
		self.Pipes.push(value)
	}
	pub fn Take(&mut self) -> Pipe {
		self.Pipes.pop().unwrap()
	}
	
	pub fn Set(&mut self, value: BigInt) {
		let index = self.Pull();	
		let length = self.Borrow[self.ActiveArray.address].len();
		let ref mut array = self.Borrow[self.ActiveArray.address];	
		array[Mod(&index,&length.to_bigint().unwrap()).to_usize().unwrap()] = value
	}
	
	pub fn Get(&mut self) -> BigInt {
		let index = self.Pull();
		let ref array = self.Borrow[self.ActiveArray.address];
		array[Mod(&index, &array.len().to_bigint().unwrap()).to_usize().unwrap()].clone()
	}
	pub fn Put(&mut self, value: BigInt) {
		let ref mut array = self.Borrow[self.ActiveArray.address];
		array.push(value)
	}
	
	pub fn Pop(&mut self) -> BigInt {
		let ref mut array = self.Borrow[self.ActiveArray.address];
		array.pop().unwrap()
	}
	
	//Beware of memory leaks!
	//TODO optimise if a == destination?
	pub fn Join(&mut self, a: &Pointer, b: &Pointer) -> Pointer {
		let mut result : Vec<BigInt> = Vec::new();
		
		{ let ref mut cell = self.Library[a.address];
		cell.set(cell.get()+1); }
		{ let ref mut cell = self.Library[b.address];
		cell.set(cell.get()+1); }
		
		for letter in &self.Borrow[a.address] {
			result.push(letter.clone())
		}
		for letter in &self.Borrow[b.address] {
			result.push(letter.clone())
		}
		
		self.Borrow.push(result);
		self.Library.push(Cell::new(1));
		self.ActiveArray = Pointer{ 
			address: self.Borrow.len()-1, 
			counter: self.Library[self.Borrow.len()-1].as_ptr(),
			collect: false,
		};
		Pointer {
			address: self.Borrow.len()-1, 
			counter:self.Library[self.Borrow.len()-1].as_ptr(),
			collect: true,
		}
	}
	
	pub fn Slice(&mut self) {
		let s = self.Grab();
		let a = self.Pull().to_usize().unwrap();
		let b = self.Pull().to_usize().unwrap();
		
		let mut v = Vec::new();
		{
			let ref mut array = self.Borrow[s.address];
		
			let slice = &mut array[a..b];
		
			v = slice.to_vec();
		}
		self.Borrow.push(v);
		self.Library.push(Cell::new(1));
		self.ActiveArray = Pointer{ 
			address: self.Borrow.len()-1, 
			counter: self.Library[self.Borrow.len()-1].as_ptr(),
			collect: false,
		};
		let result = Pointer {
			address: self.Borrow.len()-1, 
			counter:self.Library[self.Borrow.len()-1].as_ptr(),
			collect: true,
		};
		self.Share(result);
	}
	
	pub fn Link(&mut self) {
		let s = self.Grab();
		let n = self.Pull();
		let mut name = String::new();
		for letter in &self.Borrow[s.address] {
			name.push(letter.to_u8().unwrap() as char);
		}
		self.Map.insert(name, n);
	}
	
	pub fn Connect(&mut self) {
		let s = self.Grab();
		let mut name = String::new();
		for letter in &self.Borrow[s.address] {
			name.push(letter.to_u8().unwrap() as char);
		}
		let mut value = 0.to_bigint().unwrap();
		let mut problem = false;
		// Takes a reference and returns Option<&V>
    	match self.Map.get(&name) {
		    Some(number) => {
		    	value = number.to_bigint().unwrap();
		    },
		    _ => {
		    	problem = true
		    },
		}
		if problem {
			self.ERROR = 1.to_bigint().unwrap();
		}
		self.Push(value)
	}
	
	pub fn Make(&mut self) {
		let size = self.Pull().to_usize().unwrap();
		let mut vec: Vec<BigInt> = Vec::new();
		vec.resize(size, 0.to_bigint().unwrap());
		
		self.Borrow.push(vec);
		self.Library.push(Cell::new(1));
		self.ActiveArray = Pointer{ 
			address: self.Borrow.len()-1, 
			counter: self.Library[self.Borrow.len()-1].as_ptr(),
			collect: false,
		};
		let result = Pointer {
			address: self.Borrow.len()-1, 
			counter:self.Library[self.Borrow.len()-1].as_ptr(),
			collect: true,
		};
		self.Share(result);
	}
	
	//Not working properly!
	pub fn Heap(&mut self) {
	
		let original = self.Pull();
		let usign = original.to_isize().unwrap();
		
	
		let mut address: usize = 0;
		
		if usign < 0 {
			address = (-usign).to_usize().unwrap();
		} else {
			address = usign.to_usize().unwrap();
		} 
		
		//Retrieve item.
		if usign > 0 {
			let ln = self.TheHeap.len();
			let index = Mod(&original, &(ln+1).to_bigint().unwrap()) - 1.to_bigint().unwrap();
			let obj = self.TheHeap[index.to_usize().unwrap()].clone();
			self.Share(obj);
		
		//Delete item.
		} else if usign < 0 {
			let ln = self.TheHeap.len();
			let index = Mod(&address.to_bigint().unwrap(), &(ln+1).to_bigint().unwrap()) - 1.to_bigint().unwrap();
			
			self.TheHeap[index.to_usize().unwrap()] = Pointer{ address: 0, counter: ptr::null_mut(), collect: false };
			self.HeapRoom.push(address);
		
		//Add item.
		} else {
			if self.HeapRoom.len() > 0 {
				let address = self.HeapRoom.pop().unwrap();
				let ln = self.TheHeap.len();
				let index = Mod(&address.to_bigint().unwrap(), &(ln+1).to_bigint().unwrap()) - 1.to_bigint().unwrap();
				
				let mut value = self.Grab();
				value.collect = false;
				{
					let ref mut cell = self.Library[value.address];
					cell.set(cell.get()+2);
				}
				
				self.TheHeap[index.to_usize().unwrap()] = value;
				self.Push(address.to_bigint().unwrap());
			} else {
				let l = self.TheHeap.len();
				let mut value = self.Grab();
				
				{
					let ref mut cell = self.Library[value.address];
					cell.set(cell.get()+2);
				}
				value.collect = false;
				
				self.TheHeap.push(value);
				self.Push((l+1).to_bigint().unwrap());
			}
		}
	}
	
	pub fn Stdin(&mut self) {
		let mut length = self.Pull();
		let mut text: Vec<BigInt> = Vec::new();
		let stdin = io::stdin();		
		let mut buffer = String::new();
		
		if length == 0.to_bigint().unwrap() {
			stdin.read_line(&mut buffer);
			
			for letter in buffer.chars() {
				text.push((letter as u32).to_bigint().unwrap());
			}
		} else if length > 0.to_bigint().unwrap() {
		
			let mut buffer: Vec<u8> = Vec::new();
		
			stdin.take(length.to_u64().unwrap()).read_to_end(&mut buffer);
			
			for letter in buffer {
				text.push(letter.to_bigint().unwrap());
			}
		} else if length < 0.to_bigint().unwrap() {
			
			let mut buffer: Vec<u8> = Vec::new();
			
			//Could be optimsed.
			let c = (-length.clone()).to_u32();
			
			stdin.lock().read_until((-length).to_u8().unwrap(),&mut buffer);
			
			for letter in buffer {
				
				if (letter as u32) == c.unwrap() {
					break;
				}
				text.push((letter as u32).to_bigint().unwrap());
			}
		}	
		
		self.Borrow.push(text);
		self.Library.push(Cell::new(1));
		self.ActiveArray = Pointer{ 
			address: self.Borrow.len()-1, 
			counter: self.Library[self.Borrow.len()-1].as_ptr(),
			collect: false,
		};
		let result = Pointer {
			address: self.Borrow.len()-1, 
			counter:self.Library[self.Borrow.len()-1].as_ptr(),
			collect: true,
		};
		self.Share(result);
	}
	
	pub fn Stdout(&mut self) {
		let text = self.Grab();
		
		for letter in &self.Borrow[text.address] {
			let c = letter.to_u8().unwrap() as char;
			print!("{}", c)
		}
	}
	
	pub fn Load(&mut self) {
		let text = self.Grab();
		let mut result: Vec<BigInt> = Vec::new();
		let mut name = String::new();
		
		//Enviromental variables.
		if self.Borrow[text.address].len() > 0 && self.Borrow[text.address][0] == 36.to_bigint().unwrap() {
			for n in 1..self.Borrow[text.address].len() {
				name.push(self.Borrow[text.address][n].to_u8().unwrap() as char)
			}
			
			match env::var(name) {
				Ok(val) => name = val,
				Err(e) => name = String::new(),
			}
		
		
		} else {
			if self.Borrow[text.address].len() > 1{
				name = String::new();
				self.ERROR = 1.to_bigint().unwrap();
				
			//System arguments.
			} else {
				let index = self.Borrow[text.address][0].to_usize().unwrap();
				let mut args = env::args();
				
				if index < args.len() { 
					name = args.nth(index).unwrap();
				} else {
					self.ERROR = 1.to_bigint().unwrap();
				}
			}
		}
		
		for letter in name.chars() {
			result.push((letter as u32).to_bigint().unwrap());
		}
		
		self.Borrow.push(result);
		self.Library.push(Cell::new(1));
		self.ActiveArray = Pointer{ 
			address: self.Borrow.len()-1, 
			counter: self.Library[self.Borrow.len()-1].as_ptr(),
			collect: false,
		};
		let result = Pointer {
			address: self.Borrow.len()-1, 
			counter:self.Library[self.Borrow.len()-1].as_ptr(),
			collect: true,
		};
		self.Share(result);
	}
}
