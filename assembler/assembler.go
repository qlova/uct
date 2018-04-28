package assembler

import "errors"
import "io"
import "math/big"
 
type Assembler struct {
	Instruction byte
	String string
	Int big.Int
	Data []big.Int
	
	indentation int
	
	Linker func(name string) io.ReadCloser
	
	Input []io.ReadCloser
	Output io.Writer
	Runtime io.Writer
}

type Target struct {
	Name, FileExtension, Header, Footer string
	Runtime Runtime
	
	Assembler func(*Assembler)
}

type Runtime struct {
	Name, Data string
}

var Targets = make(map[string]Target)

func RegisterTarget(t Target) {
	Targets[t.FileExtension] = t
}

func (assembler *Assembler) Assemble(ext string) error {
	target, ok := Targets[ext]
	if !ok {
		return errors.New("Invalid Target: "+ext)
	}
	
	assembler.WriteLine(target.Header)
	var RuntimeLoaded bool
	
	for {
		var b [1]byte
		_, err := assembler.Input[len(assembler.Input)-1].Read(b[:])
		if err != nil {
			if err != io.EOF {
				return err
			}
			
			if len(assembler.Input) == 1 {
				break
			} else {
				assembler.Input[len(assembler.Input)-1].Close()
				assembler.Input = assembler.Input[:len(assembler.Input)-1]
				continue
			}
		}	
		
		if b[0] != 1 {
			return errors.New("Invalid .u file!")
		}
		
		assembler.Input[len(assembler.Input)-1].Read(b[:])
		
		assembler.Instruction = b[0]
		
		if assembler.Instruction != NATIVE && assembler.Instruction != LINK && !RuntimeLoaded {
			assembler.WriteLine(target.Runtime.Data)
			RuntimeLoaded = true
		}
		
		switch assembler.Instruction {
			case LINK, CALL, PUSH, PUSH_LIST, PUSH_PIPE, PULL, PULL_LIST, PULL_PIPE, NAME, NAME_LIST, NAME_PIPE, WRAP, CODE:
				assembler.String, err = assembler.ReadString()
				if err != nil {
					return err
				}
			case INT:
				assembler.Int, err = assembler.ReadInt()
				if err != nil {
					return err
				}
			case NATIVE:
				t, err := assembler.ReadString()
				if err != nil {
					return err
				}
				
				assembler.String, err = assembler.ReadString()
				if err != nil {
					return err
				}
				
				if t != target.FileExtension {
					continue
				}
				
			case DATA:
				assembler.String, err = assembler.ReadString()
				if err != nil {
					return err
				}
				
				assembler.Data, err = assembler.ReadData()
				if err != nil {
					return err
				}
		}
		
		if assembler.Instruction == LINK {
			if assembler.Linker == nil {
				panic("Linker not set!")
			}
			
			assembler.Input = append(assembler.Input, assembler.Linker(assembler.String))
			continue
		} else {
			
			target.Assembler(assembler)
		}
	}
	
	assembler.WriteLine(target.Footer)
	
	return nil
}

func (assembler *Assembler) Indentation() int {
	return assembler.indentation
}

func (assembler *Assembler) IncreaseIndentation() {
	assembler.indentation++
}

func (assembler *Assembler) DecreaseIndentation() {
	assembler.indentation--
	if assembler.indentation < 0 {
		assembler.indentation = 0
	}
}

func (assembler *Assembler) ReadData() ([]big.Int, error) {
	var length, err = assembler.ReadInt()
	if err != nil {
		return nil, err
	}
	
	var result []big.Int
	
	for i:=0; i < int(length.Int64()); i++ {
		var value, err =  assembler.ReadInt()
		if err != nil {
			return nil, err
		}
		
		result = append(result, value)
	}
	
	return result, nil
}

func (assembler *Assembler) ReadString() (string, error) {
	var length, err = assembler.ReadInt()
	if err != nil {
		return "", err
	}
	
	var buffer = make([]byte, length.Int64())
	_, err = assembler.Input[len(assembler.Input)-1].Read(buffer[:])
	if err != nil {
		return "", err
	}
	
	return string(buffer), nil
}

func (assembler *Assembler) ReadInt() (big.Int, error) {
	var b [1]byte
	_, err := assembler.Input[len(assembler.Input)-1].Read(b[:])
	if err != nil {
		return big.Int{}, err
	}
	var length = b[0]
	
	var buffer = make([]byte, length)
	_, err = assembler.Input[len(assembler.Input)-1].Read(buffer[:])
	if err != nil {
		return big.Int{}, err
	}
	
	var result big.Int
	
	result.SetBytes(buffer)
	
	return result, nil
}

func (assembler *Assembler) WriteString(data string) error {
	for i := 0; i < assembler.indentation; i++ {
		_, err := assembler.Output.Write([]byte{'\t'})
		if err != nil {
			return err
		}
	}
	_, err := assembler.Output.Write([]byte(data))
	if err != nil {
		return err
	}
	
	return nil
}


func (assembler *Assembler) WriteLine(data string) error {
	if err := assembler.WriteString(data); err != nil {
		return err
	}
	
	_, err := assembler.Output.Write([]byte{'\n'})
	return err
}
